package main

import (
	"context"
	"errors"
	"moneytor/api"
	db "moneytor/database/sqlc"
	"moneytor/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := utils.LoadConfig("config.json")
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to load config")
	}

	if config.Env == "DEV" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx := context.Background()
	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to database")
	}
	defer connPool.Close()

	store := db.NewStore(connPool)

	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create server")
	}

	httpServer := &http.Server{
		Addr:    config.HttpServerAddress,
		Handler: server.Router(),
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Unable to start server")
		}
	}()

	log.Info().Msgf("server started at %s", config.HttpServerAddress)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("server exited")
}
