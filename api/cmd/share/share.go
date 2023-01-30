package share

import "github.com/rs/zerolog/log"

type Handler struct {
	opts *Options
}

func New(opts *Options) *Handler {
	return &Handler{
		opts: opts,
	}
}

func (h *Handler) Run() error {
	log.Debug().Msg("debug message!")
	log.Info().Msg("info message!")
	log.Warn().Msg("wanr message!")
	log.Error().Msg("err message!")
	return nil
}
