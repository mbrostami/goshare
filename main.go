package main

import (
	"os"

	"github.com/mbrostami/goshare/api/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	cli := cmd.NewCli()
	if err := cli.Process(); err != nil {
		log.Error().Err(err).Send()
	}
}
