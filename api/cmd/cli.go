package cmd

import (
	"errors"
	"os"
	"sync"

	"github.com/mbrostami/goshare/api/cmd/receive"
	"github.com/rs/zerolog"

	"github.com/mbrostami/goshare/api/cmd/share"

	"github.com/jessevdk/go-flags"
	"github.com/mbrostami/goshare/api/cmd/server"
)

type Options struct {
	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	Verbose []bool           `short:"v" long:"verbose" description:"Show verbose debug information"`
	Server  *server.Options  `command:"server"`
	Share   *share.Options   `command:"share"`
	Receive *receive.Options `command:"receive"`
}

type Handler interface {
	Run() error
}

type Cli struct {
	handlers sync.Map
	opts     *Options
}

func NewCli() *Cli {
	var cli Cli
	cli.opts = &Options{
		Server:  &server.Options{},
		Share:   &share.Options{},
		Receive: &receive.Options{},
	}
	cli.addHandler("server", server.New(cli.opts.Server))
	cli.addHandler("share", share.New(cli.opts.Share))
	cli.addHandler("receive", receive.New(cli.opts.Receive))
	return &cli
}

func (c *Cli) addHandler(name string, h Handler) {
	c.handlers.Store(name, h)
}

func (c *Cli) Process() error {
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

	return command.Run()
}
