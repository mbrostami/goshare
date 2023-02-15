package main

import (
	"context"
	"github.com/mbrostami/goshare/internal/config"
	"github.com/mbrostami/goshare/pkg/tracer"
	"os"

	"github.com/mbrostami/goshare/api/cmd"
	"github.com/mbrostami/goshare/internal/services/client"
	"github.com/mbrostami/goshare/internal/services/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	cfg := config.Load()
	ctx := context.Background()
	if cfg.Tracing {
		closeFunc, err := tracer.InitTracerWithJaegerExporter(cfg.Jaeger, "goshare", "1.0.0")
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		defer closeFunc(ctx)
	}

	cli := cmd.NewCli(
		server.NewService(),
		client.NewService(),
	)
	if err := cli.Run(ctx); err != nil {
		log.Error().Err(err).Send()
	}
}
