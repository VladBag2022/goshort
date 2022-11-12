package main

import (
	"context"
	"fmt"
	_ "net/http/pprof" //nolint:gosec // enable debug handler for education
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/VladBag2022/goshort/internal/misc"
	"github.com/VladBag2022/goshort/internal/server"
	"github.com/VladBag2022/goshort/internal/storage"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

const NA string = "N/A"

func main() {
	if len(buildVersion) == 0 {
		buildVersion = NA
	}
	if len(buildDate) == 0 {
		buildDate = NA
	}
	if len(buildCommit) == 0 {
		buildCommit = NA
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	log.Trace("Read configuration from environment variables...")
	config, err := server.NewConfig()
	if err != nil {
		log.Error(fmt.Sprintf("Unable to read configuration from environment variables: %s", err))
		return
	}

	log.Trace("Read configuration from command line arguments...")
	serverAddressPtr := flag.StringP("address", "a", "", "server address: host:port")
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

	app, postgresRepository, memoryRepository, err := newApp(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		app.ListenAndServe()
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

func newApp(cfg *server.Config) (server.Server, *storage.PostgresRepository, *storage.MemoryRepository, error) {
	if len(cfg.DatabaseDSN) == 0 {
		if len(cfg.FileStoragePath) == 0 {
			mem := storage.NewMemoryRepository(
				misc.Shorten,
				misc.UUID,
			)
			return server.NewServer(mem, nil, cfg), nil, mem, nil
		}

		coolStorage, _ := storage.NewCoolStorage(cfg.FileStoragePath)
		mem := storage.NewMemoryRepositoryWithCoolStorage(
			misc.Shorten,
			misc.UUID,
			coolStorage,
		)
		if err := mem.Load(context.Background()); err != nil {
			fmt.Println(err)
		}
		return server.NewServer(mem, nil, cfg), nil, mem, nil
	}

	pg, err := storage.NewPostgresRepository(
		context.Background(),
		cfg.DatabaseDSN,
		misc.Shorten,
		misc.UUID,
	)
	return server.NewServer(pg, pg, cfg), pg, nil, err
}
