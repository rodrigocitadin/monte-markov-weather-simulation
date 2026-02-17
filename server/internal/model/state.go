package model

type TempState int
type RainState int

const (
	TempCold TempState = iota
	TempMild
	TempHot
)

const (
	RainDry RainState = iota
	RainLight
	RainHeavy
)

func BucketTemp(t float64) TempState {
	switch {
	case t < 5:
		return TempCold
	case t < 25:
		return TempMild
	default:
		return TempHot
	}
}

func BucketRain(r float64) RainState {
	switch {
	case r == 0:
		return RainDry
	case r < 10:
		return RainLight
	default:
		return RainHeavy
	}
}
