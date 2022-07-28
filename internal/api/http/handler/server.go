package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kosdirus/andintern/assets"
	"github.com/kosdirus/andintern/internal/andintern"
	"github.com/kosdirus/andintern/internal/config"
	"net/http"
)

type Server struct {
	*http.Server
	andintern *andintern.Core
	cfg       *config.Config
}

func NewServer(cfg *config.Config, andintern *andintern.Core) (*Server, error) {
	srv := &Server{
		Server: &http.Server{
			Addr:         cfg.API.Address,
			ReadTimeout:  cfg.API.ReadTimeout,
			WriteTimeout: cfg.API.WriteTimeout,
		},
		andintern: andintern,
		cfg:       cfg,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	if cfg.API.ServeSwagger {
		registerSwagger(r)
	}

	r.Route("/api/car", func(r chi.Router) {
		r.Get("/", srv.getCar)       // to display all cars or one car by brand or ID
		r.Post("/", srv.createCar)   // to create a car record in database
		r.Put("/", srv.updateCar)    // to update a car record by ID (brand and/or price may change)
		r.Delete("/", srv.deleteCar) // to delete car entry by brand/ID/price below the specified one

	})

	srv.Handler = r

	return srv, nil
}

func registerSwagger(r *chi.Mux) {
	r.HandleFunc("/internal/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/internal/swagger/", http.StatusFound)
	})

	swaggerHandler := http.StripPrefix("/internal/", http.FileServer(http.FS(assets.SwaggerFiles)))
	r.Get("/internal/swagger/*", swaggerHandler.ServeHTTP)

}
