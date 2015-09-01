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
		for _, j := range rank[i+1:] {
			vote[N*n+j] = true
		}
	}
	return vote
}

func disperse(votes [][]bool, lr float64) (strength []float64) {
	strength = make([]float64, len(votes[0]))
	weights := make([][]float64, len(votes))
	N := make([]int, len(votes))

	//initialise weights
	for i, v := range votes {
		weights[i] = make([]float64, len(v))
		for _, b := range v {
			if b { N[i]++ }
		}
		
		budget := 1.0 / float64(N[i])
		for j, b := range v {
			if b { weights[i][j] = budget }
		}
	}
	
	//loop this + remove till ranked
	for epoch := 0; epoch < 200; epoch++ {
		//zero
		for i := range strength {
			strength[i] = 0
		}
		
		//sum weights
		for _, w := range weights {
			for i, v := range w {
				strength[i] += math.Sqrt(v)
			}
		}

		//avg
		for i := range strength {
			strength[i] /= float64(len(weights))
		}

		//deriv step / step size
		var total_delta float64
		for i, v := range votes {
			//fmt.Println(weights[i])
			delta := make([]float64, len(v))
			for j, b := range v {
				if b {
					delta[j] = -(1.0 - strength[j]) * math.Sqrt(weights[i][j]) // * 2
				}
			}
			if i == 1 {
				//fmt.Println("DB", delta)
			}
			
			//get norm
			var sum float64
			for j, b := range v {
				if b { sum += delta[j] }
			}

			//normalize
			for j, b := range v {
				if b {
					delta[j] -= sum / float64(N[i])
					weights[i][j] += delta[j] * lr
					total_delta += delta[j] * delta[j]

					//norm
					if weights[i][j] < 0 { weights[i][j] = 0 }
					
				}
			}
			if i == 1 {
				//fmt.Println("W", weights[i])
				//fmt.Println("D", delta)
			}
			
			sum = 0
			for _, w := range weights[i] { sum += w }
			for j := range weights[i] { weights[i][j] /= sum }
		}
		//for _, w := range weights {
		//	pp(w, 5)
		//}
		//fmt.Println(strength)
		fmt.Println("total change:", total_delta)
	}
		
	return
}

func count(votes [][]bool, mask []int) (minbin int, bins []int) {
	N := len(votes[0])
	bins = make([]int, N)
	assignments := make([]int, len(votes))

	//initial assignment
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
	/*

	//improve weakest
	for {
		//find weakest
		var c, minbin, maxbin int
		minbin, maxbin = math.MaxInt64, 0
		for _, m := range mask {
			if bins[m] < minbin {
				minbin = bins[m]
				c = m
			}
			if bins[m] > maxbin {
				maxbin = bins[m]
			}
		}

		if maxbin - minbin < 2 {
			break
		}

		//find path
		
		
	}
	
	
	//find weakest link
	minbin = math.MaxInt64
	for _, m := range mask {
		if bins[m] < minbin {
			minbin = bins[m]
		}
	}
*/
	return
}

func pp(vals []float64, N int) {
	for c := 0; c < N; c++ {
		for r := 0; r < N; r++ {
			fmt.Printf("%.3f ", vals[c*N+r])
		}
		fmt.Println()
	}
}

func main() {
	fmt.Println("Stol Vote")

	
	N := 5
	V := 10000

	voters := make([][]int, 5)
	voters[0] = []int{3,4,0,1,2}
	voters[1] = []int{4,3,2,1,0}
	voters[2] = []int{0,1,2,3,4}
	voters[3] = []int{0,2,1,4,3}
	voters[4] = []int{0,2,1,3,4}
	
	var votes [][]bool

	fmt.Println("voting")
	for i := 0; i < V; i++ {
		r := i % 3
		//r := int(rand.Float64() * 2.999)
		if r == 2 {
			r += int(rand.Float64() * 2.999)
		}
		
		vote := rank_to_vote(voters[r], N)
		votes = append(votes, vote)
	}
	

	//N := 3
	/*votes := make([][]bool, 3)
	votes[0] = rank_to_vote([]int{0,1,2}, N)
	votes[1] = rank_to_vote([]int{0,1,2}, N)
	votes[2] = rank_to_vote([]int{1,0,2}, N)
 */
	fmt.Println("counting")
	strengths := disperse(votes, .02)
	pp(strengths, N)
	

	/*
	max_val := 0
	counter := make(map[int]int, 0)
	for i := 0; i < 10000; i++ {
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
*/
}
