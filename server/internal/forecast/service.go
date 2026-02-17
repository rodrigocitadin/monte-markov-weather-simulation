package forecast

import (
	"errors"

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
	RainProb float64 `json:"chuva"`
}

type TempStats struct {
	Sum   float64
	Count int
}

type FinalPrediction struct {
	MinTemp  float64 `json:"min"`
	MaxTemp  float64 `json:"max"`
	RainProb float64 `json:"chuva"`
}

func (s *Service) Forecast(lat, lon float64) ([]DayPrediction, error) {
	resp, err := weatherapi.Fetch(lat, lon)
	if err != nil {
		return nil, err
	}

	if len(resp.Daily.Time) < 3 {
		return nil, errors.New("not enough data")
	}

	var minStates []model.TempState
	var maxStates []model.TempState
	var rainStates []model.RainState

	minStats := make(map[model.TempState]*TempStats)
	maxStats := make(map[model.TempState]*TempStats)

	for i := range resp.Daily.Time {
		minTemp := resp.Daily.TempMin[i]
		maxTemp := resp.Daily.TempMax[i]
		rain := resp.Daily.RainSum[i]

		minState := model.BucketTemp(minTemp)
		maxState := model.BucketTemp(maxTemp)
		rainState := model.BucketRain(rain)

		minStates = append(minStates, minState)
		maxStates = append(maxStates, maxState)
		rainStates = append(rainStates, rainState)

		if minStats[minState] == nil {
			minStats[minState] = &TempStats{}
		}
		minStats[minState].Sum += minTemp
		minStats[minState].Count++

		if maxStats[maxState] == nil {
			maxStats[maxState] = &TempStats{}
		}
		maxStats[maxState].Sum += maxTemp
		maxStats[maxState].Count++
	}

	minChain := markov.NewChain[model.TempState]()
	maxChain := markov.NewChain[model.TempState]()
	rainChain := markov.NewChain[model.RainState]()

	minChain.Train(minStates)
	maxChain.Train(maxStates)
	rainChain.Train(rainStates)

	runs := 20000
	days := 3

	minDist := montecarlo.SimulateMultiStep(
		minStates[len(minStates)-1],
		days,
		runs,
		minChain.Next,
	)

	maxDist := montecarlo.SimulateMultiStep(
		maxStates[len(maxStates)-1],
		days,
		runs,
		maxChain.Next,
	)

	rainDist := montecarlo.SimulateMultiStep(
		rainStates[len(rainStates)-1],
		days,
		runs,
		rainChain.Next,
	)

	var predictions []DayPrediction

	for d := 0; d < 3; d++ {
		predictions = append(predictions, DayPrediction{
			Day:      d + 1,
			MinTemp:  expectedTemp(minDist[d], minStats),
			MaxTemp:  expectedTemp(maxDist[d], maxStats),
			RainProb: expectedRain(rainDist[d]),
		})
	}

	return predictions, nil
}

func expectedTemp(
	dist map[model.TempState]float64,
	stats map[model.TempState]*TempStats,
) float64 {
	total := 0.0

	for state, prob := range dist {
		s := stats[state]
		if s == nil || s.Count == 0 {
			continue
		}
		mean := s.Sum / float64(s.Count)
		total += mean * prob
	}

	return total
}

func expectedRain(dist map[model.RainState]float64) float64 {
	total := 0.0

	for state, prob := range dist {

		var rainVal float64

		switch state {
		case model.RainDry:
			rainVal = 0.0
		case model.RainLight:
			rainVal = 0.4
		case model.RainHeavy:
			rainVal = 0.9
		}

		total += rainVal * prob
	}

	return total
}
