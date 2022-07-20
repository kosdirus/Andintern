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

	r.Route("/api", func(r chi.Router) {
		r.Get("/", srv.getCar)       // для вывода всех машин или одной машины по бренду или айди
		r.Post("/", srv.createCar)   // для создания записи о машине
		r.Put("/", srv.updateCar)    // для обновления записи о машине по айди (может меняться бренд и/или цена)
		r.Delete("/", srv.deleteCar) // для удаления записи по бренду/айди/цене ниже указанной

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
