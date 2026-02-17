package forecast

import (
	"math/rand"

	"github.com/rodrigocitadin/monte-markov-weather-simulation/internal/model"
)

func GenerateSyntheticData(days int) []model.WeatherState {
	history := make([]model.WeatherState, days)
	current := model.WeatherState{MinTemp: 20, MaxTemp: 26, Rain: 0}

	for i := range days {
		history[i] = current

		deltaMin := (rand.Intn(3) - 1) * model.BucketSize
		deltaMax := (rand.Intn(3) - 1) * model.BucketSize

		nextMin := current.MinTemp + deltaMin
		nextMax := current.MaxTemp + deltaMax

		nextRain := 0
		rainProb := 0.2

		if current.Rain > 0 {
			rainProb = 0.6
			nextMax -= model.BucketSize
		}

		if rand.Float64() < rainProb {
			if rand.Float64() < 0.3 {
				nextRain = 2
			} else {
				nextRain = 1
			}
		}

		if nextMin < -10 {
			nextMin = -10
		}
		if nextMax > 45 {
			nextMax = 45
		}

		if nextMax <= nextMin {
			nextMax = nextMin + model.BucketSize
		}

		current = model.WeatherState{
			MinTemp: nextMin,
			MaxTemp: nextMax,
			Rain:    nextRain,
		}
	}

	return history
}
