package cmd

import (
	"context"
	"errors"
	"github.com/mbrostami/goshare/pkg/tracer"
	"os"
	"sync"

	"github.com/mbrostami/goshare/internal/services/sharing"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
)

type Handler interface {
	Run(ctx context.Context, command *flags.Command) error
}

type Cli struct {
	handlers sync.Map
	opts     *Options
}

func NewCli(sharingService *sharing.Service) *Cli {
	var cli Cli
	cli.opts = &Options{}

	srvHandler := newServerHandler()
	cli.opts.Server = srvHandler.opts
	cli.handlers.Store("server", srvHandler)

	shrHandler := newShareHandler(sharingService)
	cli.opts.Share = shrHandler.opts
	cli.handlers.Store("share", shrHandler)

	rcvHandler := newReceiveHandler(sharingService)
	cli.opts.Receive = rcvHandler.opts
	cli.handlers.Store("receive", rcvHandler)

	crtHandler := newCertHandler()
	cli.opts.Cert = crtHandler.opts
	cli.handlers.Store("cert", crtHandler)

	return &cli
}

func (c *Cli) Run(ctx context.Context) error {
	ctx, span := tracer.NewSpan(ctx, "cli")
	defer span.End()

	parser := flags.NewParser(c.opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			return nil
		}
		parser.WriteHelp(os.Stdout)
		return err
	}

	if parser.Active == nil {
		parser.WriteHelp(os.Stdout)
		return errors.New("expected a valid command")
	}

	cmd, ok := c.handlers.Load(parser.Active.Name)
	if !ok {
		parser.WriteHelp(os.Stdout)
		return errors.New("command is not valid")
	}

	command, ok := cmd.(Handler)
	if !ok {
		parser.WriteHelp(os.Stdout)
		return errors.New("command is not valid")
	}

	switch len(c.opts.Verbose) {
	case 3:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case 1:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}

	return command.Run(ctx, parser.Active)
}
