package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Kreg101/metrics/internal/metric"
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

func (f Filer) writeMetric(m *metric.Metric) error {
	return f.encoder.Encode(&m)
}

func (f Filer) readMetric() (*metric.Metric, error) {
	event := &metric.Metric{}
	if err := f.decoder.Decode(&event); err != nil {
		return nil, err
	}

	return event, nil
}

// LoadFile - предварительная загрузка из файла
func (f Filer) load() (metric.Metrics, error) {
	s := metric.Metrics{}
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

// Write записывает в файл содержимое хранилища, предварительно очистив файл
func (s *Storage) Write() {
	if err := os.Truncate(s.filer.file.Name(), 0); err != nil {
		s.log.Errorf("failed to truncate: %v", err)
		return
	}

	for _, m := range s.GetAll() {
		fmt.Println(m)
		err := s.filer.writeMetric(&m)
		if err != nil {
			s.log.Errorf("can't add metric %v to file: %s", m, err)
		}
	}
}
