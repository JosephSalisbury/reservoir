// Package reservoir provides primitives for sampling streams, using reservoir sampling.
package reservoir

import (
	"math/rand"
	"time"
)

var (
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// Sample samples at most k values from the channel c.
func Sample(k int, c chan int) []int {
	i := 0
	x := make([]int, k)

	for e := range c {
		i++

		if i-1 < k {
			x[i-1] = e
			continue
		}

		if r.Float32() < float32(k)/float32(i) {
			x[r.Intn(len(x))] = e
		}
	}

	return x
}
