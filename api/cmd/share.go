package cmd

import (
	"github.com/jessevdk/go-flags"
	"github.com/mbrostami/goshare/internal/services/client"
	"github.com/rs/zerolog/log"
)

type shareOptions struct {
	File string `short:"f" long:"file" description:"file path you want to share" required:"true"`
}

type shareHandler struct {
	opts          *shareOptions
	clientService *client.Service
}

func newShareHandler(clientService *client.Service) *shareHandler {
	return &shareHandler{
		opts:          &shareOptions{},
		clientService: clientService,
	}
}

func (h *shareHandler) Run(command *flags.Command) error {
	log.Debug().Msg("debug message!")
	log.Info().Msg("info message!")
	log.Warn().Msg("wanr message!")
	log.Error().Msg("err message!")

	return nil
}
