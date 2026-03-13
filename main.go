package main

import (
	"log"

	"github.com/genss333/go-clean-architecture/infrastructure/logger"
	"github.com/genss333/go-clean-architecture/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	if err := logger.Init(cfg.Server.Mode); err != nil {
		log.Fatal("Failed to init logger:", err)
	}
	defer logger.Sync()

}
