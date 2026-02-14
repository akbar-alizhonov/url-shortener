package main

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/http/handlers"
	"awesomeProject/internal/http/middlewares"
	"awesomeProject/internal/repositiries"
	"awesomeProject/internal/service"
	"awesomeProject/pkg/logger"
	"awesomeProject/pkg/postgres"
	"context"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	cfg := config.MustLoad()
	log := logger.NewLogger(cfg.Env)
	log.Info("app initialized")
	log.Debug("debug enabled")

	pool, err := postgres.NewPostgres(cfg.Postgres)
	if err != nil {
		log.Error("failed to connect to postgres", err)
	}

	repo := repositiries.NewUrlRepository(pool)
	generator := service.NewAliasGenerator()
	serv := service.NewUrlService(repo, generator, log, cfg.BaseUrl())
	urlHandler := handlers.NewUrlHandler(serv)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// echo
	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middlewares.RequestContext)
	e.Use(middlewares.RequestLogger(log))
	e.Use(middleware.Recover())

	// routes
	e.POST("/url", urlHandler.SaveUrl)
	e.GET("/list", urlHandler.ListUrls)
	e.GET("/url/:id", urlHandler.Redirect)
	e.PUT("/url", urlHandler.Update)
	e.DELETE("/url/:id", urlHandler.Delete)

	if err = e.Start(":8080"); err != nil {
		log.Error("failed to start server", err)
	}
}
