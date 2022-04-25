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

type coolStorageDumper struct {
	file    *os.File
	encoder *json.Encoder
}

type coolStorageLoader struct {
	file    *os.File
	decoder *json.Decoder
}

type CoolStorage struct {
	loader 	*coolStorageLoader
	dumper  *coolStorageDumper
}

func NewCoolStorage(fileName string) (*CoolStorage, error) {
	loader, err := NewCoolStorageLoader(fileName)
	if err != nil {
		return nil, err
	}
	dumper, err := NewCoolStorageDumper(fileName)
	if err != nil {
		return nil, err
	}
	return &CoolStorage{
		loader: loader,
		dumper: dumper,
	}, nil
}

func (c *CoolStorage) Dump(records []*CoolStorageRecord) error {
	for _, record := range records {
		if err := c.dumper.Dump(record); err != nil {
			return err
		}
	}
	return nil
}

func (c *CoolStorage) Load() ([]*CoolStorageRecord, error) {
	var records []*CoolStorageRecord
	for {
		record, err := c.loader.Load()
		if err != nil {
			fmt.Println(err)
			return records, nil
		}
		records = append(records, record)
	}
}

func NewCoolStorageDumper(fileName string) (*coolStorageDumper, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &coolStorageDumper{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func NewCoolStorageLoader(fileName string) (*coolStorageLoader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &coolStorageLoader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *coolStorageDumper) Dump(record *CoolStorageRecord) error {
	return c.encoder.Encode(&record)
}

func (c *coolStorageLoader) Load() (*CoolStorageRecord, error) {
	record := &CoolStorageRecord{}
	if err := c.decoder.Decode(&record); err != nil {
		return nil, err
	}
	return record, nil
}

func (c *coolStorageDumper) Close() error {
	return c.file.Close()
}

func (c *coolStorageLoader) Close() error {
	return c.file.Close()
}