package cmd

import (
	"github.com/jessevdk/go-flags"
	"github.com/mbrostami/goshare/internal/services/client"
)

type registerOptions struct {
	Server   string `short:"s" long:"server" description:"address of the server <ip>:<port>" required:"true"`
	Username string `short:"u" long:"username" description:"username" required:"true"`
	Key      string `short:"k" long:"key" description:"path to the public key" required:"true" default:"~/.ssh/id_rsa.pub"`
}

type registerHandler struct {
	opts          *registerOptions
	clientService *client.Service
}

func newRegisterHandler(clientService *client.Service) *registerHandler {
	return &registerHandler{
		opts:          &registerOptions{},
		clientService: clientService,
	}
}

func (h *registerHandler) Run(command *flags.Command) error {
	//  TODO read the pubKey and send as string
	return h.clientService.Register(h.opts.Username, h.opts.Server, h.opts.Key)
}
