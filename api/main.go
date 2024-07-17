package main

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aacebo/orca/api/amqp"
	"github.com/aacebo/orca/api/common"
	"github.com/aacebo/orca/api/postgres"
	"github.com/aacebo/orca/api/routes"
	"github.com/aacebo/orca/api/sockets"
	"github.com/aacebo/orca/api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-chi/render"
)

func main() {
	startedAt := time.Now()
	os.Setenv("TZ", "") // UTC
	gob.Register(map[string]any{})

	pg := postgres.New()
	defer pg.Close()

	amqp := amqp.New()
	defer amqp.Close()

	ctx := common.Context{
		"amqp":    amqp,
		"pg":      pg,
		"sockets": sockets.New(),
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.NoCache)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(httprate.LimitByRealIP(600, 1*time.Minute))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.AllowAll().Handler)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]any{
			"started_at": startedAt.UnixNano(),
		})
	})

	r.Mount("/v1", routes.New(ctx))

	http.ListenAndServe(
		fmt.Sprintf(":%s", utils.GetEnv("PORT", "3000")),
		r,
	)
}
