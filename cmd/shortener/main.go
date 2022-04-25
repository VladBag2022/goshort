package main

import (
	"context"
	"fmt"
	"github.com/VladBag2022/goshort/internal/server"
	"github.com/VladBag2022/goshort/internal/shortener"
	"github.com/VladBag2022/goshort/internal/storage"
)

func main() {
	config, err := server.NewConfig()
	if err != nil{
		fmt.Println(err)
		return
	}
	var memoryRepository *storage.MemoryRepository
	if len(config.FileStoragePath) != 0 {
		coolStorage, _ := storage.NewCoolStorage(config.FileStoragePath)
		memoryRepository = storage.NewMemoryRepositoryWithCoolStorage(shortener.Shorten, coolStorage)
		memoryRepository.Load(context.Background())
	} else {
		memoryRepository = storage.NewMemoryRepository(shortener.Shorten)
	}
	app := server.NewServer(memoryRepository, config)
	app.ListenAndServer()
	memoryRepository.Dump(context.Background())
}