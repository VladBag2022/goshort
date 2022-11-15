package main

import (
	"context"
	"fmt"
	_ "net/http/pprof" //nolint:gosec // enable debug handler for education
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/VladBag2022/goshort/internal/misc"
	"github.com/VladBag2022/goshort/internal/server"
	"github.com/VladBag2022/goshort/internal/storage"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string

	configFile string
	rootCmd    = &cobra.Command{
		Use: "shortener",
		Run: rootRun,
	}
)

const NA string = "N/A"

func init() {
	cobra.OnInitialize(initConfig)

	// add flags.
	rootCmd.PersistentFlags().StringP("address", "a", "", "server address: host:port")
	rootCmd.PersistentFlags().StringP("base", "b", "", "base url for URL misc")
	rootCmd.PersistentFlags().StringP("file", "f", "", "file storage path")
	rootCmd.PersistentFlags().StringP("database", "d", "", "database DSN")
	rootCmd.PersistentFlags().BoolP("https", "s", false, "enable HTTPS")
	rootCmd.PersistentFlags().StringP("cert", "e", "", "cert PEM file for HTTPS")
	rootCmd.PersistentFlags().StringP("key", "p", "", "key PEM file for HTTPS")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")

	// bind flags.
	for key, f := range map[string]string{
		"Address":         "address",
		"BaseURL":         "base",
		"FileStoragePath": "file",
		"DatabaseDSN":     "database",
		"EnableHTTPS":     "https",
		"CertPEMFile":     "cert",
		"KeyPEMFile":      "key",
	} {
		if err := viper.BindPFlag(key, rootCmd.PersistentFlags().Lookup(f)); err != nil {
			log.Errorf("failed to bind flag %s. %s", f, err)
		}
	}

	// bind ENV variables.
	for key, env := range map[string]string{
		"Address":         "SERVER_ADDRESS",
		"BaseURL":         "BASE_URL",
		"FileStoragePath": "FILE_STORAGE_PATH",
		"AuthCookieName":  "AUTH_COOKIE",
		"AuthCookieKey":   "AUTH_KEY",
		"DatabaseDSN":     "DATABASE_DSN",
		"EnableHTTPS":     "ENABLE_HTTPS",
		"CertPEMFile":     "CERT_PEM",
		"KeyPEMFile":      "KEY_PEM",
	} {
		if err := viper.BindEnv(key, env); err != nil {
			log.Errorf("failed to bind ENV variable %s. %s", env, err)
		}
	}

	// set default values.
	viper.SetDefault("Address", "localhost:8080")
	viper.SetDefault("AuthCookieName", "X-AUTH")
	viper.SetDefault("AuthCookieKey", "gopher")
	viper.SetDefault("CertPEMFile", "cert.pem")
	viper.SetDefault("KeyPEMFile", "key.pem")
}

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

	if err := rootCmd.Execute(); err != nil {
		log.Errorf("failed to execute root command. %s", err)
	}
}

func rootRun(_ *cobra.Command, _ []string) {
	config := server.NewConfig()
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

	if err = app.Shutdown(); err != nil {
		log.Errorf("failed to shutdown HTTP server gracefully. %s", err)
	}

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

func initConfig() {
	if configFile == "" {
		configFile = os.Getenv("CONFIG")
	}

	// read config from file.
	if configFile != "" {
		viper.SetConfigName(configFile)
		viper.SetConfigType("json")
		viper.AddConfigPath(".")

		err := viper.ReadInConfig()
		if err != nil {
			log.Errorf("failed to read config file: %s", err)
		}
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
