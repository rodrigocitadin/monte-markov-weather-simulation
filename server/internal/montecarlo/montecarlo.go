package montecarlo

import (
	"sync"
)

func Simulate[T comparable](
	initial T,
	days int,
	runs int,
	next func(T) T,
) map[T]float64 {
	results := make(map[T]int)
	var mu sync.Mutex
	var wg sync.WaitGroup

	workers := 8
	chunk := runs / workers

	for range workers {
		wg.Go(func() {
			local := make(map[T]int)

			for range chunk {
				state := initial

				for range days {
					state = next(state)
				}

				local[state]++
			}

			mu.Lock()
			for k, v := range local {
				results[k] += v
			}
			mu.Unlock()
		})
	}

	wg.Wait()

	final := make(map[T]float64)
	for k, v := range results {
		final[k] = float64(v) / float64(runs)
	}

	return final
}
