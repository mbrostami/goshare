package cmd

import (
	"context"
	"github.com/jessevdk/go-flags"
	"github.com/mbrostami/goshare/api/grpc"
)

type serverOptions struct {
	CertPath string `long:"cert-path" description:"path to cert.pem and key.pem" default:"."`
	IP       string `long:"ip" description:"ip address to listen on"`
	Port     string `short:"p" long:"port" description:"port number to listen on" default:"8080"`
	WithTLS  bool   `long:"with-tls" description:"enable tls encryption"`
}

type serverHandler struct {
	opts *serverOptions
}

func newServerHandler() *serverHandler {
	return &serverHandler{
		opts: &serverOptions{},
	}
}

func (h *serverHandler) Run(ctx context.Context, command *flags.Command) error {
	return grpc.ListenAndServe(h.opts.WithTLS, h.opts.CertPath, h.opts.IP+":"+h.opts.Port)
}
