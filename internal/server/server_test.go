package server

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	
	"github.com/VladBag2022/goshort/internal/misc"
	"github.com/VladBag2022/goshort/internal/storage"
)

func Example() {
	const url string = "http://localhost:8080"

	cfg, err := NewConfig()
	if err != nil {
		log.Error(fmt.Sprintf("Unable to read configuration from environment variables: %s", err))
		return
	}
	
	mem := storage.NewMemoryRepository(
		misc.Shorten,
		misc.UUID,
	)

	app := NewServer(mem, nil, cfg)
	go func() {
		app.ListenAndServe()
	}()

	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		Get(url + "/ping")

	if err != nil {
		log.Errorf("failed to run ping request. %s", err)
		return
	}
	fmt.Println("Status Code:", resp.StatusCode())

	// Output:
	// Status Code: 500
}