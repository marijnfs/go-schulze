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

func rank_to_vote(rank []int, N int) []bool {
	vote := make([]bool, N*N)
	for i, n := range rank {
		for _, j := range rank[i:] {
			vote[N*n+j] = true
		}
	}
	return vote
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

	N := 5
	V := 10000

	voters := make([][]int, 5)
	voters[0] = []int{3,4,0,2,1}
	voters[1] = []int{4,3,2,1,0}
	voters[2] = []int{0,1,2,3,4}
	voters[3] = []int{1,2,0,4,3}
	voters[4] = []int{0,2,1,3,4}
	var votes [][]bool

	fmt.Println("voting")
	for i := 0; i < V; i++ {
		r := int(rand.Float64() * 2.0999)
		if r == 2 {
			r += int(rand.Float64() * 2.999)
		}
		
		vote := rank_to_vote(voters[r], N)
		votes = append(votes, vote)
	}

	fmt.Println("counting")

	max_val := 0
	counter := make(map[int]int, 0)
	for i := 0; i < 100000; i++ {
		C := 3
		rank := rand.Perm(N)[:C]
		var indices []int
		for n, v := range rank_to_vote(rank, N) {
			if v { indices = append(indices, n) }
		}
		
		c, _ := count(votes, indices)
		if c >= max_val {
			fmt.Println(c, rank)
			max_val = c
			q := rank[0] * 100 + rank[1] * 10 + rank[2]
			if _, ok := counter[q]; !ok { counter[q] = 0 }
			counter[q]++

		}
	}

	max_val = 0
	for i, c := range counter {
		if c >= max_val {
			fmt.Println(c, i)
			max_val = c
		}
	}
}
