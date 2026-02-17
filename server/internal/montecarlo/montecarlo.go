package montecarlo

import "sync"

func SimulateMultiStep[T comparable](
	initial T,
	days int,
	runs int,
	next func(T) T,
) []map[T]float64 {
	results := make([]map[T]int, days)
	for i := range results {
		results[i] = make(map[T]int)
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	workers := 8
	chunk := runs / workers

	for range workers {
		wg.Go(func() {
			local := make([]map[T]int, days)
			for i := range local {
				local[i] = make(map[T]int)
			}

			for range chunk {
				state := initial

				for d := range days {
					state = next(state)
					local[d][state]++
				}
			}

			mu.Lock()
			for d := range days {
				for k, v := range local[d] {
					results[d][k] += v
				}
			}
			mu.Unlock()
		})
	}

	wg.Wait()

	final := make([]map[T]float64, days)

	for d := range days {
		final[d] = make(map[T]float64)
		for k, v := range results[d] {
			final[d][k] = float64(v) / float64(runs)
		}
	}

	return final
}
