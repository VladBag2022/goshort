package main

import (
	"github.com/VladBag2022/goshort/internal/server"
	"github.com/VladBag2022/goshort/internal/shortener"
	"github.com/VladBag2022/goshort/internal/storage"
	"log"
)

func main() {
	r := storage.NewMemoryRepository(shortener.Shorten)
	c, err := server.NewConfig()
	if err != nil{
		log.Fatal(err)
		return
	}
	s := server.NewServer(r, c)
	s.ListenAndServer()
}