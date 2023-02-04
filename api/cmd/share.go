package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/jessevdk/go-flags"
	"github.com/mbrostami/goshare/internal/services/client"
)

type shareOptions struct {
	File    string   `short:"f" long:"file" description:"file path you want to share" required:"true"`
	Servers []string `short:"s" long:"server" description:"address of the servers <ip>:<port>" required:"true"`
}

type shareHandler struct {
	opts          *shareOptions
	clientService *client.Service
}

func newShareHandler(clientService *client.Service) *shareHandler {
	return &shareHandler{
		opts:          &shareOptions{},
		clientService: clientService,
	}
}

func (h *shareHandler) Run(command *flags.Command) error {
	ctx := context.TODO()
	log.Debug().Msg("checking servers...")
	if err := h.clientService.VerifyServers(ctx, h.opts.Servers); err != nil {
		return err
	}

	log.Debug().Msg("generating the key...")
	key, uid := h.clientService.GenerateKey(h.opts.Servers)
	fmt.Printf("share this key -> %s\n", key)

	//if !askForConfirmation("receiver started receiving with the above key?") {
	//	log.Debug().Msg("aborted")
	//	return nil
	//}

	log.Debug().Msg("starting the share...")
	return h.clientService.Share(ctx, h.opts.File, uid, h.opts.Servers)
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [Y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "" || response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
