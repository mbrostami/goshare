package cmd

import (
	"context"
	"fmt"
	"github.com/mbrostami/goshare/pkg/tracer"

	"github.com/rs/zerolog/log"

	"github.com/jessevdk/go-flags"
	"github.com/mbrostami/goshare/internal/services/client"
)

type receiveOptions struct {
	Key        string `short:"k" long:"key" description:"download key" required:"true"`
	WithTLS    bool   `long:"with-tls" description:"connect with tls encryption"`
	SkipVerify bool   `long:"skip-verify" description:"skip tls certificate verification"`
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

func (h *receiveHandler) Run(ctx context.Context, command *flags.Command) error {
	ctx, span := tracer.NewSpan(ctx, "receive-run")
	defer span.End()

	servers, id, err := h.clientService.ParseKey(h.opts.Key)
	if err != nil {
		return err
	}

	log.Debug().Msg("checking servers...")
	if err := h.clientService.VerifyServers(ctx, servers, h.opts.WithTLS, h.opts.SkipVerify); err != nil {
		return err
	}

	fileName, err := h.clientService.Receive(ctx, id, servers, h.opts.WithTLS, h.opts.SkipVerify)
	if err != nil {
		return err
	}
	fmt.Printf("Downloaded file: %s\n", fileName)

	return nil
}
