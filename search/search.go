package search

import (
	"context"
	"fmt"
)

type Random func() string
type RandomTopHistory func() string
type Mutate func(s1, s2 string) string
type Test func(s string) string
type Store func(s, result string)

func Run(ctx context.Context, iterations, randPerIter, mutatedPerIter int, random Random, history RandomTopHistory, mutate Mutate, test Test, store Store) error {
	if random == nil {
		return fmt.Errorf("random")
	}

	for i := 0; i < iterations; i++ {
		samples := make([]string, 0, randPerIter+mutatedPerIter)

		for j := 0; j < mutatedPerIter; j++ {
			samples = append(samples, mutate(history(), history()))
		}

		for len(samples) < randPerIter+mutatedPerIter {
			samples = append(samples, random())
		}

		//now go through all the samples
		for _, sample := range samples {
			store(sample, test(sample))
		}

		//check up on the context
		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}

	return nil
}
