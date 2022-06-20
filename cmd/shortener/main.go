package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/VladBag2022/goshort/internal/misc"
	"github.com/VladBag2022/goshort/internal/server"
	"github.com/VladBag2022/goshort/internal/storage"
)

func main() {
	log.SetLevel(log.TraceLevel)

	log.Trace("Read configuration from environment variables...")
	config, err := server.NewConfig()
	if err != nil{
		log.Error(fmt.Sprintf("Unable to read configuration from environment variables: %s", err))
		return
	}

	log.Trace("Read configuration from command line arguments...")
	serverAddressPtr := flag.StringP("address", "a", "","server address: host:port")
	baseURLPtr := flag.StringP("base", "b", "", "base url for URL misc")
	fileStoragePathPtr := flag.StringP("file", "f", "", "file storage path")
	databasePtr := flag.StringP("database", "d", "", "database DSN")
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
	if len(*databasePtr) != 0 {
		config.DatabaseDSN = *databasePtr
	}

	var app server.Server
	var postgresRepository *storage.PostgresRepository
	var memoryRepository *storage.MemoryRepository

	if len(config.DatabaseDSN) != 0 {
		postgresRepository, err = storage.NewPostgresRepository(
			context.Background(),
			config.DatabaseDSN,
			misc.Shorten,
			misc.Register,
		)
		if err != nil {
			fmt.Println(err)
			return
		}

		app = server.NewServer(postgresRepository, postgresRepository, config)
	} else {
		if len(config.FileStoragePath) != 0 {
			coolStorage, _ := storage.NewCoolStorage(config.FileStoragePath)
			memoryRepository = storage.NewMemoryRepositoryWithCoolStorage(
				misc.Shorten,
				misc.Register,
				coolStorage,
			)
			if err = memoryRepository.Load(context.Background()); err != nil {
				fmt.Println(err)
			}
		} else {
			memoryRepository = storage.NewMemoryRepository(
				misc.Shorten,
				misc.Register,
			)
		}
		app = server.NewServer(memoryRepository, postgresRepository, config)
	}

	go func() {
		app.ListenAndServer()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sigChan

	if memoryRepository != nil {
		if err = memoryRepository.Dump(context.Background()); err != nil {
			fmt.Println(err)
		}
		memoryRepository.Close()
	}

	if postgresRepository != nil {
		postgresRepository.Close()
	}
}