package db

import (
	"context"
	"database/sql"
	"github.com/Kreg101/metrics/internal/metric"
	"github.com/Kreg101/metrics/internal/server/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"time"
)

// Storage структура для работы с базой данных. Содержит в себе соединение и логер.
// Реализует интерфейс handler.Repository
type Storage struct {
	conn *sql.DB
	log  *zap.SugaredLogger
}

// NewStorage в соединении создает нужную таблицу и инициализирует внутренний логер
func NewStorage(conn *sql.DB, log *zap.SugaredLogger) (Storage, error) {
	s := Storage{conn: conn}

	if log != nil {
		s.log = log
	} else {
		s.log = logger.Default()
	}

	err := s.sqlCreateTable()
	if err != nil {
		return Storage{}, err
	}

	return s, nil
}

// Add добавляет метрику в бд. Если она там уже есть, то обновляет ее значение в соответствие с типом метрики
// Гарантируется, что сюда поступают правильные метрики
func (s Storage) Add(ctx context.Context, m metric.Metric) {

	// проверяем существование метрики в хранилище
	inStore, err := s.sqlElementExist(ctx, m)
	if err != nil {
		s.log.Errorf("can't check element's existing: %e", err)
		return
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

			// считываем предыдущее значение метрики
			prev, err := s.sqlGetDelta(ctx, m)
			if err != nil {
				s.log.Errorf("can't get delta metric from storage: %e", err)
				return
			}

			// обновляем текущую
			*m.Delta += prev

			// обновляем ее в бд
			err = s.sqlUpdateDelta(ctx, m)
			if err != nil {
				s.log.Errorf("can't update delta metric: %e", err)
				return
			}

			// завершаем транзакцию
			err = tx.Commit()
			if err != nil {
				s.log.Errorf("can't commit transaction: %s", err.Error())
			}

		case "gauge":

			err = s.sqlUpdateValue(ctx, m)
			if err != nil {
				s.log.Errorf("can't update value metric: %e", err)
			}
		}
	} else {
		err = s.sqlInsert(ctx, m)
		if err != nil {
			s.log.Errorf("can't insert metric into storage: %e", err)
		}
	}
}

// Get возвращает метрику из хранилища по имени и true, если она есть,
// либо пустую метрику и false, если ее нет
func (s Storage) Get(ctx context.Context, name string) (metric.Metric, bool) {
	m, err := s.sqlGetMetric(ctx, name)
	if err != nil {
		if err != sql.ErrNoRows {
			s.log.Errorf("can't get existing value from data base: %e", err)
		}
		return metric.Metric{}, false
	}

	return m, true
}

// GetAll получает все метрики из базы данных и пытается иx преобразовать к metric.Metrics
func (s Storage) GetAll(ctx context.Context) metric.Metrics {
	metrics := make(metric.Metrics, 0)
	rows, err := s.conn.QueryContext(ctx, `SELECT * FROM metrics`)

	if err != nil {
		s.log.Errorf("can't get all elements from data base: %e", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {

		m := metric.Metric{}

		err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
		if err != nil {
			s.log.Errorf("can't get metric %s from data base: %e", m, err)
			return nil
		}

		metrics[m.ID] = m
	}

	err = rows.Err()
	if err != nil {
		s.log.Errorf("can't parse metrics from rows: %e", err)
		return nil
	}

	return metrics
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

func (s Storage) sqlCreateTable() error {
	_, err := s.conn.ExecContext(context.Background(), `
        CREATE TABLE IF NOT EXISTS metrics (
            id VARCHAR(128) PRIMARY KEY,
            mtype VARCHAR(30) NOT NULL,
            delta BIGINT,
            value DOUBLE PRECISION         
        )
    `)
	return err
}

func (s Storage) sqlElementExist(ctx context.Context, m metric.Metric) (bool, error) {
	row := s.conn.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT * FROM metrics WHERE id = $1 AND mtype = $2)`,
		m.ID, m.MType)

	var inStore bool
	err := row.Scan(&inStore)

	if err != nil {
		return false, err
	}

	return inStore, nil
}

func (s Storage) sqlGetDelta(ctx context.Context, m metric.Metric) (int64, error) {
	var prev int64
	row := s.conn.QueryRowContext(ctx,
		`SELECT delta FROM metrics WHERE $1 = id AND $2 = mtype`,
		m.ID, m.MType)

	err := row.Scan(&prev)
	if err != nil {
		return 0, err
	}
	return prev, nil
}

func (s Storage) sqlUpdateDelta(ctx context.Context, m metric.Metric) error {
	_, err := s.conn.ExecContext(ctx,
		`UPDATE metrics SET delta = $1 WHERE id = $2 AND mtype = $3`,
		m.Delta, m.ID, m.MType)

	return err
}

func (s Storage) sqlUpdateValue(ctx context.Context, m metric.Metric) error {
	_, err := s.conn.ExecContext(ctx,
		`UPDATE metrics SET value = $1 WHERE id = $2 AND mtype = $3`,
		m.Value, m.ID, m.MType)

	return err
}

func (s Storage) sqlInsert(ctx context.Context, m metric.Metric) error {
	_, err := s.conn.ExecContext(ctx,
		`INSERT INTO metrics (id, mtype, delta, value) VALUES ($1, $2, $3, $4)`,
		m.ID, m.MType, m.Delta, m.Value)

	return err
}

func (s Storage) sqlGetMetric(ctx context.Context, id string) (metric.Metric, error) {
	m := metric.Metric{}
	row := s.conn.QueryRowContext(ctx,
		`SELECT * FROM metrics WHERE id = $1`, id)

	err := row.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
	if err != nil {
		return metric.Metric{}, err 
	}
	return m, nil 
}
