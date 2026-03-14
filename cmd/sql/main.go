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
	// FIX 1: Catch errors if the folder doesn't exist or path is wrong
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}

		fullPath := filepath.Join(folderPath, entry.Name())
		fmt.Printf("Executing script: %s...\n", entry.Name())

		// FIX 2: Check for errors when reading the file
		content, err := os.ReadFile(fullPath)
		if err != nil {
			log.Fatalf("Failed to read file %s: %v", entry.Name(), err)
		}

		// FIX 3: Check for SQL execution errors and log exactly which file failed
		_, err = pool.Exec(ctx, string(content))
		if err != nil {
			log.Fatalf("❌ SQL ERROR in file %s:\n%v", entry.Name(), err)
		}

		fmt.Printf("✅ Successfully executed: %s\n", entry.Name())
	}

	fmt.Println("All migrations completed successfully!")
}
