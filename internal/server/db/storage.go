package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Kreg101/metrics/internal/metric"
	"github.com/Kreg101/metrics/internal/server/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"time"
)

// -d="host=localhost user=postgres password=Kravchenko01 dbname=really sslmode=disable" - параметры командной строки

// Storage структура для работы с базой данных. Содержит в себе соединение и логер.
// Реализует интерфейс handler.Repository
type Storage struct {
	conn *sql.DB
	log  *zap.SugaredLogger
}

// NewStorage в соединении создает нужную таблицу и инициализирует внутренний логер
func NewStorage(conn *sql.DB, log *zap.SugaredLogger) (Storage, error) {
	fmt.Println("here")
	s := Storage{conn: conn}

	if log != nil {
		s.log = log
	} else {
		s.log = logger.Default()
	}

	_, err := s.conn.ExecContext(context.Background(), `
        CREATE TABLE IF NOT EXISTS metrics (
            id VARCHAR(128) PRIMARY KEY,
            mtype VARCHAR(30) NOT NULL,
            delta INTEGER,
            value DOUBLE PRECISION         
        )
    `)

	if err != nil {
		return Storage{}, err
	}

	return s, nil
}

// Add добавляет метрику в бд. Если она там уже есть, то обновляет ее значение в соответствие с типом метрики
// Гарантируется, что сюда поступают правильные метрики
func (s Storage) Add(ctx context.Context, m metric.Metric) {

	// для того чтобы не рассматривать много случаев, если данной метрики еще нет в бд
	if m.Delta == nil {
		m.Delta = new(int64)
	} else {
		m.Value = new(float64)
	}

	// проверяем существование метрики в хранилище
	row := s.conn.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT * FROM metrics WHERE $1 = id AND $2 = mtype)`,
		m.ID, m.MType)

	var inStore bool
	err := row.Scan(&inStore)
	if err != nil {
		panic(err)
	}

	if inStore {
		switch m.MType {
		case "counter":

			// по ТЗ нам нужно вернуть обновленное значение метрики, поэтому после обновления
			// вытаскиваем второй раз ее из хранилища. Чтобы эти операции происходили слитно и если что
			// откатились, используем транзакции
			tx, err := s.conn.BeginTx(ctx, nil)

			if err != nil {
				s.log.Errorf("can't use transaction: %e", err)
				return
			}

			defer tx.Rollback()

			_, err = s.conn.ExecContext(ctx,
				`UPDATE metrics SET delta = delta + $1 WHERE $2 = id AND $3 = mtype`,
				*m.Delta, m.ID, m.MType)

			if err != nil {
				s.log.Errorf("can't update counter metric: %e", err)
				return
			}

			row := s.conn.QueryRowContext(ctx,
				`SELECT delta FROM metrics WHERE $1 = id AND $2 = mtype`,
				m.ID, m.MType)

			err = row.Scan(m.Delta)
			if err != nil {
				s.log.Errorf("can't get metric value after update: %e", err)
				return
			}

			tx.Commit()

		case "gauge":

			fmt.Println(m.ID, m.MType, *m.Value)
			_, err = s.conn.ExecContext(ctx,
				`UPDATE metrics SET value = $1 WHERE $2 = id AND $3 = mtype`,
				*m.Value, m.ID, m.MType)
		}

	} else {
		_, err = s.conn.ExecContext(ctx,
			`INSERT INTO metrics (id, mtype, delta, value) VALUES ($1, $2, $3, $4)`,
			m.ID, m.MType, *m.Delta, *m.Value)

		//normal(m)
	}

	if err != nil {
		s.log.Errorf("can't add metric %s to storage: %e", m, err)
	}
}

// Get возвращает метрику из хранилища по имени и true, если она есть,
// либо пустую метрику и false, если ее нет
func (s Storage) Get(ctx context.Context, name string) (metric.Metric, bool) {
	m := metric.Metric{Delta: new(int64), Value: new(float64)}
	row := s.conn.QueryRowContext(ctx,
		`SELECT * FROM metrics WHERE id = $1`, name)

	err := row.Scan(&m.ID, &m.MType, m.Delta, m.Value)
	if err != nil {
		if err != sql.ErrNoRows {
			s.log.Errorf("can't get existing value from data base: %e", err)
		}
		return metric.Metric{}, false
	}

	return normal(m), true
}

// Normal приводит метрику к каноническому виду, после того, как ее
// достали из хранилища
func normal(m metric.Metric) metric.Metric {
	res := m
	if res.MType == "gauge" {
		res.Delta = nil
	} else {
		res.Value = nil
	}
	return res
}

// GetAll получает все метрики из базы данных и пытается иx преобразовать к metric.Metrics
func (s Storage) GetAll(ctx context.Context) metric.Metrics {
	metrics := make(metric.Metrics, 0)
	rows, err := s.conn.QueryContext(ctx,
		`SELECT * FROM metrics`)

	if err != nil {
		s.log.Errorf("can't get all elements from data base: %e", err)
		return nil
	}

	defer rows.Close()

	for rows.Next() {
		m := metric.Metric{Delta: new(int64), Value: new(float64)}

		err = rows.Scan(&m.ID, &m.MType, m.Delta, m.Value)
		if err != nil {
			s.log.Errorf("can't get metric %s from data base: %e", m, err)
			return nil
		}

		metrics[m.ID] = normal(m)
	}

	err = rows.Err()
	if err != nil {
		s.log.Errorf("can't parse metrics from rows: %e", err)
		return nil
	}

	return metrics
}

// Close закрывает соединение с базой данных
func (s Storage) Close() {
	s.conn.Close()
}

// Ping проверяет соединение с базой данных
func (s Storage) Ping(pctx context.Context) error {
	ctx, cancel := context.WithTimeout(pctx, 1*time.Second)
	defer cancel()
	if err := s.conn.PingContext(ctx); err != nil {
		return err
	}

	return nil
}