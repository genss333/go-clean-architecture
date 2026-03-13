package main

import (
	"context"
	"log"

	postgres "github.com/genss333/go-clean-architecture/infrastructure/connection"
	db_repositories "github.com/genss333/go-clean-architecture/infrastructure/database/repositories"
	"github.com/genss333/go-clean-architecture/infrastructure/logger"
	"github.com/genss333/go-clean-architecture/internal/config"
	"github.com/genss333/go-clean-architecture/internal/delivery/http/handler"
	"github.com/genss333/go-clean-architecture/internal/delivery/http/router"
	depart_usecases "github.com/genss333/go-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
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

	gin.SetMode(cfg.Server.Mode)

	pool, err := postgres.NewPgxPool(context.Background(), cfg.Postgres)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()

	departmentRepo := db_repositories.NewDepartmentSQLCRepository(pool)
	departmentUC := depart_usecases.NewDepartmentUsecase(departmentRepo)
	departmentHandler := handler.NewDepartmentHandler(departmentUC)

	r := router.New(departmentHandler)

	log.Printf("Server starting on :%s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
