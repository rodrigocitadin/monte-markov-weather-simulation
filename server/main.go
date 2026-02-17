package main

import (
	"log"
	"net/http"

	"github.com/rodrigocitadin/monte-markov-weather-simulation/internal/forecast"
	apphttp "github.com/rodrigocitadin/monte-markov-weather-simulation/internal/http"
)

func main() {
	service := &forecast.Service{}
	handler := &apphttp.Handler{
		Service: service,
	}

	http.HandleFunc("/forecast", handler.Forecast)

	log.Println("running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
