package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func count(votes [][]bool, mask []int) (minbin int, bins []int) {
	N := len(votes[0])
	bins = make([]int, N)
	assignments := make([]int, len(votes))
	for v, vote := range votes {
		c, minbin := -1, math.MaxInt64
		for _, m := range rand.Perm(len(mask)) {
			p := mask[m]
			if vote[p] && bins[p] < minbin {
				c, minbin = p, bins[p]
			}
		}

		if c >= 0 {
			bins[c]++
			assignments[v] = c
		}
	}

	minbin = math.MaxInt64
	for _, m := range mask {
		if bins[m] < minbin {
			minbin = bins[m]
		}
	}
	return
}


func main() {
	fmt.Println("Stol Vote")

	N := 300 * 300
	V := 100000

	var votes [][]bool

	fmt.Println("voting")
	for i := 0; i < V; i++ {
		vote := make([]bool, N)
		C := rand.Intn(N)
		for _, v := range rand.Perm(N)[:C] {
			vote[v] = true
		}
		votes = append(votes, vote)
	}

	fmt.Println("counting")
	for i := 0; i < 10; i++ {
		C := 5
		mask := rand.Perm(N)[:C]
		{
			c, _ := count(votes, mask)
			fmt.Println(c, mask)
		}
		c, _ := count(votes, mask)
		fmt.Println(c, mask)
	}
}
