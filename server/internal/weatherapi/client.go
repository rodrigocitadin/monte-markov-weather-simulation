package weatherapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Daily struct {
	Time    []string  `json:"time"`
	TempMax []float64 `json:"temperature_2m_max"`
	TempMin []float64 `json:"temperature_2m_min"`
	RainSum []float64 `json:"rain_sum"`
}

type Response struct {
	Daily Daily `json:"daily"`
}

func Fetch(lat, lon float64) (*Response, error) {
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&daily=temperature_2m_max,temperature_2m_min,rain_sum&past_days=30&timezone=auto",
		lat, lon,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}
