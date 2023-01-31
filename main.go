package main

import (
	"os"

	"github.com/mbrostami/goshare/internal/services/client"

	"github.com/mbrostami/goshare/internal/repositories/database"

	"github.com/boltdb/bolt"
	"github.com/mbrostami/goshare/api/cmd"
	"github.com/mbrostami/goshare/internal/config"
	"github.com/mbrostami/goshare/internal/services/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.Load()

	bdb, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		log.Fatal().Msgf("database connect: %s", err)
	}
	defer bdb.Close()

	repo, err := database.NewRepository(bdb)
	if err != nil {
		log.Fatal().Msgf("database initialize: %s", err)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	cli := cmd.NewCli(
		server.NewService(repo),
		client.NewService(repo),
	)
	if err := cli.Run(); err != nil {
		log.Error().Err(err).Send()
	}
}
