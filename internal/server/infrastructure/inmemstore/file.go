package inmemstore

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Kreg101/metrics/internal/entity"
	"io"
	"os"
)

// lineCounter - вспомогательная функция для подсчета количества метрик в файле
// так как количество строк совпадает с количеством метрик(каждая метрика на отдельной строке)
// то достаточно посчитать строки
func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// Filer - структура для записи и чтения хранилища из файла. Переодичность записи, нужно ли загружать
// прошле метрики из файла при старте, это определяется конфигурацией
type Filer struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

// writeMetric записывает струтуру в формате json
func (f Filer) writeMetric(m *entity.Metric) error {
	return f.encoder.Encode(&m)
}

// readMetric считывавет в формате json одну единственную метрику
// и мапит ее в структуру
func (f Filer) readMetric() (*entity.Metric, error) {
	event := &entity.Metric{}
	if err := f.decoder.Decode(&event); err != nil {
		return nil, err
	}

	return event, nil
}

// LoadFile - предварительная загрузка из файла
func (f Filer) load() (entity.Metrics, error) {
	s := entity.Metrics{}
	help, err := os.Open(f.file.Name())

	if err != nil {
		return nil, err
	}

	count, err := lineCounter(help)
	if err != nil {
		return nil, err
	}

	for i := 0; i < count; i++ {
		m, err := f.readMetric()
		if err != nil {
			return nil, err
		}
		s[m.ID] = *m
	}

	return s, nil
}

// Write переносит в файл все содержимое хранилища, предварительно очистив файл
func (s *InMemStorage) Write() {
	if err := os.Truncate(s.filer.file.Name(), 0); err != nil {
		s.log.Errorf("failed to truncate: %v", err)
		return
	}

	for _, m := range s.GetAll(context.Background()) {
		err := s.filer.writeMetric(&m)
		if err != nil {
			s.log.Errorf("can't add entity %v to file: %s", m, err)
		}
	}
}
