package forecast

import (
	"math/rand"

	"github.com/rodrigocitadin/monte-markov-weather-simulation/internal/model"
)

func GenerateSyntheticData(totalSteps int) []model.WeatherState {
	var history []model.WeatherState

	runs := 50
	stepsPerRun := max(totalSteps/runs, 10)

	for range runs {
		startTemp := -10 + rand.Intn(51)

		current := model.WeatherState{
			MinTemp: model.RoundToBucket(float64(startTemp)),
			MaxTemp: model.RoundToBucket(float64(startTemp + 5)),
			Rain:    0,
		}

		if rand.Float64() < 0.3 {
			current.Rain = 1
		}

		for range stepsPerRun {
			history = append(history, current)

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

			if nextMin < -15 {
				nextMin = -15
			}
			if nextMax > 50 {
				nextMax = 50
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
	}

	return history
}
