package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	postgres "github.com/genss333/go-clean-architecture/infrastructure/connection"
	"github.com/genss333/go-clean-architecture/infrastructure/logger"
	"github.com/genss333/go-clean-architecture/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	if err := logger.Init(cfg.Server.Mode); err != nil {
		log.Fatal("Failed to init logger:", err)
	}
	defer logger.Sync()

	pool, err := postgres.NewPgxPool(context.Background(), cfg.Postgres)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()

	folderPath := "../../db/migrations/"
	entries, err := os.ReadDir(folderPath)

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}

		fullPath := filepath.Join(folderPath, entry.Name())
		fmt.Printf("Executing script: %s...\n", entry.Name())

		content, _ := os.ReadFile(fullPath)

		_, err = pool.Exec(ctx, string(content))

	}

}
