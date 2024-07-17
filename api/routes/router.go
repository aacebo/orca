package routes

import (
	"github.com/aacebo/orca/api/common"
	"github.com/aacebo/orca/api/routes/agents"
	"github.com/go-chi/chi/v5"
)

func New(ctx common.Context) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		agents.New(r, ctx)
	})

	return r
}
