package main

import (
	"BlockchainCurrency/config"
	"BlockchainCurrency/internal/adapter/blockchain"
	"BlockchainCurrency/internal/adapter/repository"
	"BlockchainCurrency/internal/service"
	"BlockchainCurrency/internal/transport/http/handlers"
	"BlockchainCurrency/internal/transport/http/server"
	"BlockchainCurrency/pkg/logger"
	"BlockchainCurrency/pkg/migrations"
	"BlockchainCurrency/pkg/pgdb"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Failed app start:%v", err)
	}
}

func run() error {
	logger, err := logger.New()
	if err != nil {
		return fmt.Errorf("init logger: %v", err)
	}
	defer logger.Sync()

	cfg, err := config.New()
	if err != nil {
		logger.Errorf("init config", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	db, err := pgdb.ConnectToDB(ctx, cfg.Postgres.ToDSN())
	if err != nil {
		logger.Errorf("connect to db", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.Errorf("failed close db", err)
		}
	}()

	if err := migrations.Run(db, logger); err != nil {
		logger.Errorf("failed up migrations", err)
	}
	quotesCache := repository.NewQuotesCache()
	quotesPostgres := repository.NewQuotesRepository(db)
	blockchainAPI := bc.NewBclockchain(cfg)
	quotesService := service.NewQuotesService(logger, cfg, quotesPostgres, blockchainAPI, quotesCache)
	quotesService.Start(ctx)

	quotesHandlers := handlers.NewQuotesHandlers(quotesService, logger, &cfg.Credentials)

	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	v1 := router.Group("/api")
	{
		v1.GET("/tickers", quotesHandlers.GetQuotes)
		v1.POST("/quit", quotesHandlers.BasicAuthMiddleware(), quotesHandlers.Quit)
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		<-sigCh

		if err := quotesService.Shutdown(ctx); err != nil {
			logger.Error("Failed to shutdown quotes service", err)
		}

		logger.Info("Received shutdown signal (SIGINT/SIGTERM)")
	}()
	srv := server.NewServer(&cfg.HTTP, logger, router)

	go srv.Run()
	quotesService.ListenOnClose(cancel)

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Shutdown Server ...", err)
	}

	logger.Info("Server stopped.")

	return nil
}
