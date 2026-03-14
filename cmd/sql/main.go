package cmd_sql

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	postgres "github.com/genss333/go-clean-architecture/infrastructure/connection"
	"github.com/genss333/go-clean-architecture/internal/config"
)

func main() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	pool, err := postgres.NewPgxPool(context.Background(), cfg.Postgres)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()

	folderPath := "../../db/migrations/"
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", folderPath, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}

		fullPath := filepath.Join(folderPath, entry.Name())
		fmt.Printf("Executing script: %s...\n", entry.Name())

		content, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("failed to read SQL file %s: %w", fullPath, err)
		}

		_, err = pool.Exec(ctx, string(content))
		if err != nil {
			return fmt.Errorf("failed to execute SQL file %s: %w", fullPath, err)
		}
	}

	return nil
}
