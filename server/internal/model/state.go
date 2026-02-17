package model

import "math"

type WeatherState struct {
	MinTemp int
	MaxTemp int
	Rain    int
}

const BucketSize = 2

func ToState(min, max, rainSum float64) WeatherState {
	return WeatherState{
		MinTemp: roundToBucket(min),
		MaxTemp: roundToBucket(max),
		Rain:    bucketRain(rainSum),
	}
}

func roundToBucket(val float64) int {
	return int(math.Round(val/float64(BucketSize)) * BucketSize)
}

func bucketRain(r float64) int {
	switch {
	case r < 0.1:
		return 0
	case r < 5.0:
		return 1
	default:
		return 2
	}
}
