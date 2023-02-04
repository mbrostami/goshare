package cmd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/jessevdk/go-flags"
	"github.com/mbrostami/goshare/internal/services/client"
)

type receiveOptions struct {
	Key string `short:"k" long:"key" description:"download key" required:"true"`
}

type receiveHandler struct {
	opts          *receiveOptions
	clientService *client.Service
}

func newReceiveHandler(clientService *client.Service) *receiveHandler {
	return &receiveHandler{
		opts:          &receiveOptions{},
		clientService: clientService,
	}
}

func (h *receiveHandler) Run(command *flags.Command) error {
	ctx := context.TODO()
	servers, id, err := h.clientService.ParseKey(h.opts.Key)
	if err != nil {
		return err
	}

	log.Debug().Msg("checking servers...")
	if err := h.clientService.VerifyServers(ctx, servers); err != nil {
		return err
	}

	fileName, err := h.clientService.Receive(ctx, id, servers)
	if err != nil {
		return err
	}
	fmt.Printf("Downloaded file: %s\n", fileName)

	return nil
}
