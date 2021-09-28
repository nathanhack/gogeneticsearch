package search

import (
	"context"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/nathanhack/threadpool"
)

// Random should return a sample that was randomly generated
type Random func() string

// RandomTopHistory should return a sample pick at random from the "top" performing samples
type RandomTopHistory func() string

// Mutate should take two samples to create and returns a third sample, the goal is the third sample is randomly created based on the two samples s1 and s2.
type Mutate func(s1, s2 string) string

// Test takes a sample to test and returns a results string.
type Test func(s string) string

// Store take in the sample and results strings, with the goal to store the two for later use (in the RandomTopHistory and possibly the Random and Mutate function to dedup).
type Store func(s, result string)

func Run(ctx context.Context, iterations, randPerIter, mutatedPerIter int, random Random, history RandomTopHistory, mutate Mutate, test Test, store Store, showProgress bool, parallelThreads int) error {
	if random == nil {
		return fmt.Errorf("random")
	}

	var bar *pb.ProgressBar
	if showProgress {
		bar = pb.StartNew(iterations)
	}

	pool := threadpool.New(ctx, parallelThreads, iterations)

	for i := 0; i < iterations; i++ {
		pool.Add(func() {
			samples := make([]string, 0, randPerIter+mutatedPerIter)

			for j := 0; j < mutatedPerIter; j++ {
				h1 := history()
				h2 := history()
				if h1 == "" || h2 == "" {
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
			if showProgress {
				bar.Increment()
			}
		})

		// pool.Add() will check the ctx, but if
		// iterations is large it could lead to
		// a long time from context cancel to end
		// of loop, so we check up on the context
		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}

	pool.Wait()
	if showProgress {
		bar.Finish()
	}

	return nil
}
