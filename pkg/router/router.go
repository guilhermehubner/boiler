package router

import (
	"compress/flate"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rafaelsq/boiler/pkg/graphql"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/boiler/pkg/rest"
	"github.com/rafaelsq/boiler/pkg/website"
)

func ApplyMiddlewares(r chi.Router) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Compress(flate.BestCompression))
	r.Use(middleware.Timeout(2 * time.Second))
}

func ApplyRoute(r chi.Router, us iface.UserService, es iface.EmailService) {
	// website
	r.Get("/", website.Handle)
	r.Get("/favicon.ico", http.NotFound)
	r.Handle("/static/*", http.FileServer(http.Dir("pkg/website")))

	// graphql
	r.Route("/graphql", func(g chi.Router) {
		g.Get("/play", graphql.PlayHandle())
		g.HandleFunc("/query", graphql.QueryHandleFunc(us, es))
	})

	// rest
	r.Route("/rest", func(r chi.Router) {
		r.Route("/users", func(u chi.Router) {
			r.Get("/", rest.ListUsersHandle(us))
			r.Post("/", rest.AddUserHandle(us))
			r.Route("/{userID:[0-9]+}", func(u chi.Router) {
				r.Get("/", rest.GetUserHandle(us))
			})
		})

		r.Route("/emails", func(u chi.Router) {
			r.Post("/", rest.AddEmailHandle(es))
			r.Get("/{userID:[0-9]+}", rest.ListUserEmailsHandle(es))
		})
	})
}
