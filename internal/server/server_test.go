package server

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"

	"github.com/VladBag2022/goshort/internal/misc"
	"github.com/VladBag2022/goshort/internal/storage"
)

func ExampleServer() {
	cfg, err := NewConfig()
	if err != nil {
		log.Error(fmt.Sprintf("Unable to read configuration from environment variables: %s", err))
		return
	}
	cfg.Address = "localhost:51515"

	mem := storage.NewMemoryRepository(
		misc.Shorten,
		misc.UUID,
	)

	app := NewServer(mem, nil, cfg)
	go func() {
		app.ListenAndServe()
	}()

	time.Sleep(time.Second) // let server start

	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		Get("http://" + cfg.Address + "/ping")
	if err != nil {
		log.Errorf("failed to run ping request. %s", err)
		return
	}
	fmt.Println("Status Code:", resp.StatusCode())

	// Output:
	// Status Code: 500
}
