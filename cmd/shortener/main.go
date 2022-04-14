package main

import (
	"github.com/VladBag2022/goshort/internal/server"
	"github.com/VladBag2022/goshort/internal/shortener"
	"github.com/VladBag2022/goshort/internal/storage"
)

func main() {
	r := storage.NewMemoryRepository(shortener.Shorten)
	s := server.New(r, "localhost", 8080)
	s.ListenAndServer()
}