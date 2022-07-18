package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type Server struct {
	*http.Server
}

func NewServer() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	//registerSwagger(r)

	r.Route("/api", func(r chi.Router) {
		r.Get("/", getCar)       // для вывода всех машин или одной машины по бренду или айди
		r.Post("/", createCar)   // для создания записи о машине
		r.Put("/", updateCar)    // для обновления записи о машине по айди (может меняться бренд и/или цена)
		r.Delete("/", deleteCar) // для удаления записи по бренду/айди/цене

	})

	http.ListenAndServe(":3000", r)
}

/*func registerSwagger(r *chi.Mux) {
	r.HandleFunc("/internal/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/internal/swagger/", http.StatusFound)
	})

	swaggerHandler := http.StripPrefix("/internal/", http.FileServer(http.FS(assets.SwaggerFiles)))
	r.Get("/internal/swagger/*", swaggerHandler.ServeHTTP)

}*/
