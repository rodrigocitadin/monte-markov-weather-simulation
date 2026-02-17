package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rodrigocitadin/monte-markov-weather-simulation/internal/forecast"
)

type Handler struct {
	Service *forecast.Service
}

func (h *Handler) Forecast(w http.ResponseWriter, r *http.Request) {
	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lon, _ := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)

	result, err := h.Service.Forecast(lat, lon)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(result)
}
