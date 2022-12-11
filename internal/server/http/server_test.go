package http

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"

	"github.com/VladBag2022/goshort/internal/misc"
	"github.com/VladBag2022/goshort/internal/server"
	"github.com/VladBag2022/goshort/internal/storage"
)

func ExampleServer() {
	cfg := server.NewConfig()
	cfg.Address = "localhost:51515"

	mem := storage.NewMemoryRepository(
		misc.Shorten,
		misc.UUID,
	)

	abstract := server.NewServer(mem, nil, cfg)
	app := NewServer(&abstract)
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
