package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Kreg101/metrics/internal/metric"
	"github.com/Kreg101/metrics/internal/server/logger"
	"io"
	"os"
)

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

//func NewProducer(fileName string) (*Producer, error) {
//	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
//	if err != nil {
//		return nil, err
//	}
//
//	return &Producer{
//		file:    file,
//		encoder: json.NewEncoder(file),
//	}, nil
//}

func (p *Producer) WriteMetric(m *metric.Metric) error {
	return p.encoder.Encode(&m)
}

func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *Consumer) ReadMetric() (*metric.Metric, error) {
	event := &metric.Metric{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

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

func (c *Consumer) LoadFile() (Metrics, error) {
	s := Metrics{}
	help, err := os.Open(c.file.Name())

	if err != nil {
		return nil, err
	}

	count, err := lineCounter(help)
	if err != nil {
		return nil, err
	}

	for i := 0; i < count; i++ {
		m, err := c.ReadMetric()
		if err != nil {
			return nil, err
		}
		s[m.ID] = *m
	}

	return s, nil

}

func (s *Storage) Write() {
	log := logger.Default()

	fmt.Println("here")

	if err := os.Truncate(s.producer.file.Name(), 0); err != nil {
		log.Errorf("Failed to truncate: %v", err)
		fmt.Println("lox")
		return
	}

	for _, m := range s.GetAll() {
		fmt.Println(m)
		err := s.producer.WriteMetric(&m)
		if err != nil {
			log.Errorf("can't add metric %v to file: %s", m, err)
		}
	}
}
