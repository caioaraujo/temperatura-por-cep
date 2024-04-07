package main

import (
	"net/http"

	"github.com/caioaraujo/temperatura-por-cep/configs"
	"github.com/caioaraujo/temperatura-por-cep/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	weatherHandler := handlers.NewWeatherHandler(configs.WeatherApiKey)
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Route("/temperatura", func(r chi.Router) {
		r.Get("/{cep}", weatherHandler.GetWeather)
	})

	http.ListenAndServe(":8000", r)
}
