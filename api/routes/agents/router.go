package agents

import (
	"github.com/aacebo/orca/api/common"
	"github.com/go-chi/chi/v5"
)

func New(r chi.Router, ctx common.Context) {
	r.HandleFunc(
		"/sockets",
		Handler(ctx),
	)
}
