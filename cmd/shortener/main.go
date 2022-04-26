package main

import (
	"context"
	"fmt"
	"github.com/VladBag2022/goshort/internal/server"
	"github.com/VladBag2022/goshort/internal/shortener"
	"github.com/VladBag2022/goshort/internal/storage"
	"os"
	"os/signal"
	"syscall"
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
		if err := memoryRepository.Load(context.Background()); err != nil {
			fmt.Println(err)
		}
	} else {
		memoryRepository = storage.NewMemoryRepository(shortener.Shorten)
	}
	app := server.NewServer(memoryRepository, config)

	go func() {
		app.ListenAndServer()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sigChan

	if err := memoryRepository.Dump(context.Background()); err != nil {
		fmt.Println(err)
	}
}