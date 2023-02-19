package cmd

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"

	"github.com/jessevdk/go-flags"
	"github.com/mbrostami/goshare/internal/services/sharing"
)

type shareOptions struct {
	File       string   `short:"f" long:"file" description:"file path you want to share" required:"true"`
	Servers    []string `short:"s" long:"server" description:"address of the servers <ip>:<port>" required:"true"`
	CAPath     string   `long:"ca-path" description:"CA's certificate path to verify server's certificate'"`
	WithTLS    bool     `long:"with-tls" description:"connect with tls encryption"`
	SkipVerify bool     `long:"skip-verify" description:"skip tls certificate verification"`
}

type shareHandler struct {
	opts           *shareOptions
	sharingService *sharing.Service
}

func newShareHandler(sharingService *sharing.Service) *shareHandler {
	return &shareHandler{
		opts:           &shareOptions{},
		sharingService: sharingService,
	}
}

func (h *shareHandler) Run(ctx context.Context, command *flags.Command) error {
	log.Debug().Msg("checking servers...")
	if err := h.sharingService.VerifyServers(ctx, h.opts.Servers, h.opts.CAPath, h.opts.WithTLS, h.opts.SkipVerify); err != nil {
		return err
	}

	log.Debug().Msg("generating the key...")
	key, uid := h.sharingService.GenerateKey(h.opts.Servers)
	fmt.Printf("share this key -> %s\n", key)

	log.Debug().Msg("starting the share...")
	return h.sharingService.Share(ctx, h.opts.File, uid, h.opts.Servers, h.opts.CAPath, h.opts.WithTLS, h.opts.SkipVerify)
}
