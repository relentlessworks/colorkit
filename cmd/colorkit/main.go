package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/relentlessworks/colorkit/internal/api"
	"github.com/relentlessworks/colorkit/internal/config"
)

func main() {
	cfg := config.Load()
	handler := api.New(cfg.Secret)
	mux := handler.Routes()

	log.Printf("colorkit listening on %s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		fmt.Printf("error: server failed: %v\n", err)
	}
}
