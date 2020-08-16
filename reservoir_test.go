package reservoir

import (
	"fmt"
	"testing"
)

const (
	TestRuns         = 100000
	Delta    float32 = 10.0
)

func key(s []int) string {
	x := ""
	for _, e := range s {
		x = x + fmt.Sprintf("%v,", e)
	}
	return x
}

type ProbableOutput struct {
	Probability float32
	Output      []int
}

func ExampleSample() {
	c := make(chan int)

	go func(c chan int) {
		for _, i := range []int{1, 2, 3} {
			c <- i
		}
		close(c)
	}(c)

	Sample(1, c)
	// will be one of [1], [2], or [3]
}

func TestSample(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		num         int
		stream      []int
		outputs     []ProbableOutput
	}{
		{
			description: "0 elements from a stream with 0 elements",
			num:         0,
			stream:      []int{},
			outputs: []ProbableOutput{
				{
					Probability: 100,
					Output:      []int{},
				},
			},
		},
		{
			description: "0 elements from a stream with 1 elements",
			num:         0,
			stream:      []int{1},
			outputs: []ProbableOutput{
				{
					Probability: 100,
					Output:      []int{},
				},
			},
		},
		{
			description: "1 element from a stream with 1 element",
			num:         1,
			stream:      []int{1},
			outputs: []ProbableOutput{
				{
					Probability: 100,
					Output:      []int{1},
				},
			},
		},
		{
			description: "1 element from a stream with 2 elements",
			num:         1,
			stream:      []int{1, 2},
			outputs: []ProbableOutput{
				{
					Probability: 50,
					Output:      []int{1},
				},
				{
					Probability: 50,
					Output:      []int{2},
				},
			},
		},
		{
			description: "2 elements from a stream with 2 elements",
			num:         2,
			stream:      []int{1, 2},
			outputs: []ProbableOutput{
				{
					Probability: 100,
					Output:      []int{1, 2},
				},
			},
		},
		{
			description: "2 elements from a stream with 3 elements",
			num:         2,
			stream:      []int{1, 2, 3},
			outputs: []ProbableOutput{
				{
					Probability: 33.3,
					Output:      []int{1, 2},
				},
				{
					Probability: 33.3,
					Output:      []int{3, 2},
				},
				{
					Probability: 33.3,
					Output:      []int{1, 3},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			results := map[string]int{}

			for i := 0; i < TestRuns; i++ {
				c := make(chan int)

				go func(s []int, c chan int) {
					for _, e := range s {
						c <- e
					}
					close(c)
				}(tt.stream, c)

				samples := Sample(tt.num, c)

				results[key(samples)] += 1
			}

			for _, output := range tt.outputs {
				n := results[key(output.Output)]
				p := float32(n) / float32(TestRuns) * 100

				if p < output.Probability-Delta || p > output.Probability+Delta {
					fmt.Println(results)

					t.Fatalf(
						"taking %v samples from %v, expected samples %v with probability %v, got probability %v",
						tt.num,
						tt.stream,
						output.Output,
						output.Probability,
						p,
					)
				}
			}
		})
	}
}
