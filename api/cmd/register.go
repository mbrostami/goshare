package cmd

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/mbrostami/goshare/internal/services/client"
)

type registerOptions struct {
	Server   string `short:"s" long:"server" description:"address of the server <ip>:<port>" required:"true"`
	Username string `short:"u" long:"username" description:"username" required:"true"`
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
	user, err := h.clientService.Register(h.opts.Username, h.opts.Server)
	if err != nil {
		return err
	}
	fmt.Printf("Public Key: %x\n", user.PubKey)
	return nil
}
