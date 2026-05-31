package main

import (
	"context"
	"moneytor/api"
	db "moneytor/database/sqlc"
	"moneytor/utils"
	"os"

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

	if err := server.Start(config.HttpServerAddress); err != nil {
		log.Fatal().Err(err).Msg("Unable to start server")
	}
}
