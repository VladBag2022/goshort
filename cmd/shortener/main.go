package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	flag "github.com/spf13/pflag"

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
	serverAddressPtr := flag.StringP("address", "a", "","server address: host:port")
	baseURLPtr := flag.StringP("base", "b", "", "base url for URL shortener")
	fileStoragePathPtr := flag.StringP("file", "f", "", "file storage path")
	flag.Parse()

	if len(*serverAddressPtr) != 0 {
		config.Address = *serverAddressPtr
	}
	if len(*baseURLPtr) != 0 {
		config.BaseURL = *baseURLPtr
	}
	if len(*fileStoragePathPtr) != 0 {
		config.FileStoragePath = *fileStoragePathPtr
	}

	var memoryRepository *storage.MemoryRepository
	if len(config.FileStoragePath) != 0 {
		coolStorage, _ := storage.NewCoolStorage(config.FileStoragePath)
		memoryRepository = storage.NewMemoryRepositoryWithCoolStorage(
			shortener.Shorten,
			shortener.Register,
			coolStorage,
		)
		if err = memoryRepository.Load(context.Background()); err != nil {
			fmt.Println(err)
		}
	} else {
		memoryRepository = storage.NewMemoryRepository(
			shortener.Shorten,
			shortener.Register,
		)
	}
	defer memoryRepository.Close()
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