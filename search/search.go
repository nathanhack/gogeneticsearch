package search

import (
	"context"
	"fmt"
	"github.com/cheggaaa/pb/v3"
)

type Random func() string
type RandomTopHistory func() string
type Mutate func(s1, s2 string) string
type Test func(s string) string
type Store func(s, result string)

func Run(ctx context.Context, iterations, randPerIter, mutatedPerIter int, random Random, history RandomTopHistory, mutate Mutate, test Test, store Store, showProgress bool) error {
	if random == nil {
		return fmt.Errorf("random")
	}

	var bar *pb.ProgressBar
	if showProgress {
		bar = pb.StartNew(iterations)
	}

	for i := 0; i < iterations; i++ {
		samples := make([]string, 0, randPerIter+mutatedPerIter)

		for j := 0; j < mutatedPerIter; j++ {
			h1 := history()
			h2 := history()
			if h1 == "" {
				break
			}
			samples = append(samples, mutate(h1, h2))
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
		if showProgress {
			bar.Increment()
		}
	}
	if showProgress {
		bar.Finish()
	}

	return nil
}
