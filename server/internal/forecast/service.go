package forecast

import (
	"errors"
	"math"

	"github.com/rodrigocitadin/monte-markov-weather-simulation/internal/markov"
	"github.com/rodrigocitadin/monte-markov-weather-simulation/internal/model"
	"github.com/rodrigocitadin/monte-markov-weather-simulation/internal/montecarlo"
	"github.com/rodrigocitadin/monte-markov-weather-simulation/internal/weatherapi"
)

type Service struct{}

type DayPrediction struct {
	Day      int     `json:"day"`
	MinTemp  float64 `json:"min"`
	MaxTemp  float64 `json:"max"`
	RainProb float64 `json:"rain_prob"`
	RainInt  int     `json:"rain_intensity"`
}

func (s *Service) Forecast(lat, lon float64, daysToPredict int) ([]DayPrediction, error) {
	resp, err := weatherapi.Fetch(lat, lon)
	if err != nil {
		return nil, err
	}

	if len(resp.Daily.Time) < 1 {
		return nil, errors.New("insufficient data")
	}

	var realHistory []model.WeatherState
	for i := range resp.Daily.Time {
		state := model.ToState(
			resp.Daily.TempMin[i],
			resp.Daily.TempMax[i],
			resp.Daily.RainSum[i],
		)
		realHistory = append(realHistory, state)
	}

	chain := markov.NewChain[model.WeatherState]()

	syntheticData := GenerateSyntheticData(5000)
	chain.Train(syntheticData)

	for range 20 {
		chain.Train(realHistory)
	}

	lastState := realHistory[len(realHistory)-1]
	simulationRuns := 10000

	distributions := montecarlo.SimulateMultiStep(
		lastState,
		daysToPredict,
		simulationRuns,
		chain.Next,
	)

	var predictions []DayPrediction
	for d, dist := range distributions {
		pred := calculateExpectedValues(dist)
		pred.Day = d + 1
		predictions = append(predictions, pred)
	}

	return predictions, nil
}

func calculateExpectedValues(dist map[model.WeatherState]float64) DayPrediction {
	var sumMin, sumMax, probRain float64
	var maxProb float64
	var mostLikelyRain int

	for state, prob := range dist {
		sumMin += float64(state.MinTemp) * prob
		sumMax += float64(state.MaxTemp) * prob

		if state.Rain > 0 {
			probRain += prob
		}

		if prob > maxProb {
			maxProb = prob
			mostLikelyRain = state.Rain
		}
	}

	return DayPrediction{
		MinTemp:  math.Round(sumMin*10) / 10,
		MaxTemp:  math.Round(sumMax*10) / 10,
		RainProb: math.Round(probRain*100) / 100,
		RainInt:  mostLikelyRain,
	}
}
