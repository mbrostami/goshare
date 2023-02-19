package cmd

import (
	"context"
	"github.com/jessevdk/go-flags"
	"github.com/mbrostami/gcert"
)

type certOptions struct {
	Host string `short:"h" long:"host" required:"true" description:"Comma-separated hostnames and IPs to generate a certificate for"`
	Dst  string `long:"dst" default:"." description:"Directory path to store and override key.pem and cert.pem files"`
}

type certHandler struct {
	opts *certOptions
}

func newCertHandler() *certHandler {
	return &certHandler{
		opts: &certOptions{},
	}
}

func (h *certHandler) Run(ctx context.Context, command *flags.Command) error {
	return gcert.Generate(
		h.opts.Host,
		h.opts.Dst,
	)
}
