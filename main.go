package main

import (
	"os"

	"github.com/mbrostami/goshare/api/cmd"
	"github.com/mbrostami/goshare/internal/services/client"
	"github.com/mbrostami/goshare/internal/services/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	//cfg := config.Load()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	cli := cmd.NewCli(
		server.NewService(),
		client.NewService(),
	)
	if err := cli.Run(); err != nil {
		log.Error().Err(err).Send()
	}
}
