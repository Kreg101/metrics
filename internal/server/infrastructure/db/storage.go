package db

import (
	"context"
	"database/sql"
	"github.com/Kreg101/metrics/internal/server/entity"
	"github.com/Kreg101/metrics/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"time"
)

// Storage структура для работы с базой данных. Содержит в себе соединение и логер.
// Реализует интерфейс transport.Repository
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
func (s Storage) Add(ctx context.Context, m entity.Metric) {

	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		s.log.Errorf("can't use transaction: %v", err)
		return
	}
	defer tx.Rollback()

	// проверяем существование метрики в хранилище
	inStore, err := sqlElementExist(ctx, tx, m)
	if err != nil {
		s.log.Errorf("can't check element's existing: %v", err)
		return
	}

	if inStore {

		switch m.MType {

		case "counter":
			*m.Delta, err = sqlUpdateDelta(ctx, m, tx)
			if err != nil {
				s.log.Errorf("can't update delta entity: %v", err)
				return
			}

		case "gauge":
			err = sqlUpdateValue(ctx, tx, m)
			if err != nil {
				s.log.Errorf("can't update value entity: %v", err)
			}
		}
	} else {
		err = sqlInsert(ctx, tx, m)
		if err != nil {
			s.log.Errorf("can't insert entity into storage: %v", err)
		}
	}

	// завершаем транзакцию
	err = tx.Commit()
	if err != nil {
		s.log.Errorf("can't commit transaction: %v", err)
	}
}

// Get возвращает метрику из хранилища по имени и true, если она есть,
// либо пустую метрику и false, если ее нет
func (s Storage) Get(ctx context.Context, name string) (entity.Metric, bool) {
	m, err := s.sqlGetMetric(ctx, name)
	if err != nil {
		if err != sql.ErrNoRows {
			s.log.Errorf("can't get existing value from data base: %v", err)
		}
		return entity.Metric{}, false
	}

	return m, true
}

// GetAll получает все метрики из базы данных и пытается иx преобразовать к entity.Metrics
func (s Storage) GetAll(ctx context.Context) entity.Metrics {
	metrics := make(entity.Metrics, 0)
	rows, err := s.conn.QueryContext(ctx, `SELECT * FROM metrics`)

	if err != nil {
		s.log.Errorf("can't get all elements from data base: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {

		m := entity.Metric{}

		err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
		if err != nil {
			s.log.Errorf("can't get entity %s from data base: %v", m, err)
			return nil
		}

		metrics[m.ID] = m
	}

	err = rows.Err()
	if err != nil {
		s.log.Errorf("can't parse metrics from rows: %v", err)
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

func sqlElementExist(ctx context.Context, tx *sql.Tx, m entity.Metric) (bool, error) {
	row := tx.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT * FROM metrics WHERE id = $1 AND mtype = $2)`,
		m.ID, m.MType)

	var inStore bool
	err := row.Scan(&inStore)

	if err != nil {
		return false, err
	}

	return inStore, nil
}

//func sqlGetDelta(ctx context.Context, m entity.Metric, tx *sql.Tx) (int64, error) {
//	var prev int64
//	row := tx.QueryRowContext(ctx,
//		`SELECT delta FROM metrics WHERE $1 = id AND $2 = mtype`,
//		m.ID, m.MType)
//
//	err := row.Scan(&prev)
//	if err != nil {
//		return 0, err
//	}
//	return prev, nil
//}

func sqlUpdateDelta(ctx context.Context, m entity.Metric, tx *sql.Tx) (int64, error) {
	row := tx.QueryRowContext(ctx,
		`UPDATE metrics SET delta = delta + $1 WHERE id = $2 AND mtype = $3 RETURNING delta`,
		m.Delta, m.ID, m.MType)

	var newDelta int64
	err := row.Scan(&newDelta)
	if err != nil {
		return 0, err
	}

	return newDelta, nil
}

func sqlUpdateValue(ctx context.Context, tx *sql.Tx, m entity.Metric) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE metrics SET value = $1 WHERE id = $2 AND mtype = $3`,
		m.Value, m.ID, m.MType)

	return err
}

func sqlInsert(ctx context.Context, tx *sql.Tx, m entity.Metric) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO metrics (id, mtype, delta, value) VALUES ($1, $2, $3, $4)`,
		m.ID, m.MType, m.Delta, m.Value)

	return err
}

func (s Storage) sqlGetMetric(ctx context.Context, id string) (entity.Metric, error) {
	m := entity.Metric{}
	row := s.conn.QueryRowContext(ctx,
		`SELECT * FROM metrics WHERE id = $1`, id)

	err := row.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
	if err != nil {
		return entity.Metric{}, err
	}
	return m, nil
}
