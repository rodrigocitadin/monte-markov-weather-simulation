package markov

import (
	"math/rand"
	"sync"
)

type Chain[T comparable] struct {
	mu    sync.RWMutex
	trans map[T]map[T]float64
}

func NewChain[T comparable]() *Chain[T] {
	return &Chain[T]{
		trans: make(map[T]map[T]float64),
	}
}

func (c *Chain[T]) Train(states []T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	counts := make(map[T]map[T]float64)

	for i := 0; i < len(states)-1; i++ {
		from := states[i]
		to := states[i+1]

		if counts[from] == nil {
			counts[from] = make(map[T]float64)
		}
		counts[from][to]++
	}

	for from, m := range counts {
		total := 0.0
		for _, v := range m {
			total += v
		}

		if c.trans[from] == nil {
			c.trans[from] = make(map[T]float64)
		}

		for to, v := range m {
			c.trans[from][to] = v / total
		}
	}
}

func (c *Chain[T]) Next(state T) T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	transitions := c.trans[state]
	r := rand.Float64()
	acc := 0.0

	for s, p := range transitions {
		acc += p
		if r <= acc {
			return s
		}
	}

	return state
}
