package cmd

import (
	"errors"
	"os"
	"sync"

	"github.com/mbrostami/goshare/internal/services/client"

	"github.com/mbrostami/goshare/internal/services/server"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
)

type Handler interface {
	Run(command *flags.Command) error
}

type Cli struct {
	handlers sync.Map
	opts     *Options
}

func NewCli(serverService *server.Service, clientService *client.Service) *Cli {
	var cli Cli
	cli.opts = &Options{}

	srvHandler := newServerHandler(serverService)
	cli.opts.Server = srvHandler.opts
	cli.handlers.Store("server", srvHandler)

	shrHandler := newShareHandler(clientService)
	cli.opts.Share = shrHandler.opts
	cli.handlers.Store("share", shrHandler)

	rcvHandler := newReceiveHandler(clientService)
	cli.opts.Receive = rcvHandler.opts
	cli.handlers.Store("receive", rcvHandler)

	return &cli
}

func (c *Cli) Run() error {
	parser := flags.NewParser(c.opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
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

	return command.Run(parser.Active)
}