package server

type Handler struct {
	opts *Options
}

func New(opts *Options) *Handler {
	return &Handler{
		opts: opts,
	}
}

func (h *Handler) Run() error {
	return nil
}
