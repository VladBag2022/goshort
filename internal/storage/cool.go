package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type CoolStorageRecord struct {
	Origin string	`json:"origin"`
	ID     string	`json:"id"`
}

type coolStorageWriter struct {
	file    *os.File
	encoder *json.Encoder
}

type coolStorageReader struct {
	file    *os.File
	decoder *json.Decoder
}

type CoolStorage struct {
	reader *coolStorageReader
	writer *coolStorageWriter
}

func NewCoolStorage(fileName string) (*CoolStorage, error) {
	reader, err := NewCoolStorageReader(fileName)
	if err != nil {
		return nil, err
	}
	writer, err := NewCoolStorageWriter(fileName)
	if err != nil {
		return nil, err
	}
	return &CoolStorage{
		reader: reader,
		writer: writer,
	}, nil
}

func (c *CoolStorage) PutRecords(records []*CoolStorageRecord) error {
	for _, record := range records {
		if err := c.writer.PutRecord(record); err != nil {
			return err
		}
	}
	return nil
}

func (c *CoolStorage) FetchRecords() ([]*CoolStorageRecord, error) {
	var records []*CoolStorageRecord
	for {
		record, err := c.reader.FetchRecord()
		if err != nil {
			fmt.Println(err)
			return records, nil
		}
		records = append(records, record)
	}
}

func NewCoolStorageWriter(fileName string) (*coolStorageWriter, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return &coolStorageWriter{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func NewCoolStorageReader(fileName string) (*coolStorageReader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return &coolStorageReader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *coolStorageWriter) PutRecord(record *CoolStorageRecord) error {
	return c.encoder.Encode(&record)
}

func (c *coolStorageReader) FetchRecord() (*CoolStorageRecord, error) {
	record := &CoolStorageRecord{}
	if err := c.decoder.Decode(&record); err != nil {
		return nil, err
	}
	return record, nil
}

func (c *coolStorageWriter) Close() error {
	return c.file.Close()
}

func (c *coolStorageReader) Close() error {
	return c.file.Close()
}