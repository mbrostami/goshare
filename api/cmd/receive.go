package cmd

import (
	"github.com/jessevdk/go-flags"
	"github.com/mbrostami/goshare/internal/services/client"
)

type receiveOptions struct {
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
	return nil
}
