package main

import (
	"fmt"
	"time"
	//"math"
	"math/rand"
	"bytes"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Table struct {
	votes []int
	schulze_votes []int
	C int
}

func MakeTable(C int) *Table {
	return &Table{make([]int, C * C), make([]int, C * C), C}
}

func (t *Table) Prefer(i, j int) {
	*t.Vote(i, j)++
}

func (t *Table) Vote(i, j int) *int {
	return &t.votes[j * t.C + i]
}

func (t *Table) SchulzeVote(i, j int) *int {
	return &t.schulze_votes[j * t.C + i]
}

func (t *Table) AddVote(ranks []int) {
	for i, x := range ranks {
		for j, y := range ranks {
			if x < y {
				t.Prefer(i, j)
			}
		}
	}
}

func (t *Table) Schulze() {
	for i := range t.schulze_votes {
		t.schulze_votes[i] = 0
	}

	
	for i := 0; i < t.C; i++ {
		for j := 0; j < t.C; j++ {
			if i != j { //&& *t.Vote(i, j) >= *t.Vote(j, i) { //experimenting
				*t.SchulzeVote(i, j) = *t.Vote(i, j)
			}
		}
	}

	for i := 0; i < t.C; i++ {
		for j := 0; j < t.C; j++ {
			if i != j {
				for k := 0; k < t.C; k++ {
					if i != k && j != k {
						*t.SchulzeVote(j, k) = max(*t.SchulzeVote(j, k), min(*t.SchulzeVote(j, i), *t.SchulzeVote(i, k)))
					}
				}
			}
		}
	}
}

func (t *Table) String() string {
	var b bytes.Buffer
	for i := 0; i < t.C; i++ {
		for j := 0; j < t.C; j++ {
			if *t.Vote(i, j) > *t.Vote(j, i) {
				fmt.Fprint(&b, BLUE)
			} else if *t.Vote(i, j) < *t.Vote(j, i) {
				fmt.Fprint(&b, RED)
			} else {
				fmt.Fprint(&b, COL_RESET)
			}
			fmt.Fprint(&b, *t.Vote(i, j), " ")
		}
		fmt.Fprintln(&b)
	}
	fmt.Fprint(&b, COL_RESET)
	return b.String()
}

func (t *Table) SchulzeString() string {
	var b bytes.Buffer
	for i := 0; i < t.C; i++ {
		for j := 0; j < t.C; j++ {
			if *t.SchulzeVote(i, j) > *t.SchulzeVote(j, i) {
				fmt.Fprint(&b, BLUE)
			} else if *t.SchulzeVote(i, j) < *t.SchulzeVote(j, i) {
				fmt.Fprint(&b, RED)
			} else {
				fmt.Fprint(&b, COL_RESET)
			}
			fmt.Fprint(&b, *t.SchulzeVote(i, j), " ")
		}
		fmt.Fprintln(&b)
	}
	fmt.Fprint(&b, COL_RESET)
	return b.String()
}

func (t *Table) SchulzeRankString() string {
	var b bytes.Buffer
	done := make([]bool, t.C)

	ranks := make([]int, 0)

	rank := 1
a:
	for {
		winners := make([]int, 0)

		for i := 0; i < t.C; i++ {
			if !done[i] {
				win := true
				for j := 0; j < t.C; j++ {
					if i != j && !done[j] && *t.SchulzeVote(i, j) < *t.SchulzeVote(j, i) {
						win = false
						break
					}
				}
				if win {
					fmt.Fprintln(&b, rank, ":", i)
					winners = append(winners, i)
					ranks = append(ranks, i)
				}
			}
		}

		for _, v := range winners {
			done[v] = true
		}
		rank++
		for _, v := range done {
			if !v { continue a}
		}
		break
	}

	
	for _, i := range ranks {
		for _, j := range ranks {
			if *t.Vote(i, j) > *t.Vote(j, i) {
				fmt.Fprint(&b, BLUE)
			} else if *t.Vote(i, j) < *t.Vote(j, i) {
				fmt.Fprint(&b, RED)
			} else {
				fmt.Fprint(&b, COL_RESET)
			}
			fmt.Fprint(&b, *t.Vote(i, j), " ")
		}
		fmt.Fprintln(&b)
	}
	fmt.Fprint(&b, COL_RESET)

	for _, i := range ranks {
		for _, j := range ranks {
			if *t.SchulzeVote(i, j) > *t.SchulzeVote(j, i) {
				fmt.Fprint(&b, BLUE)
			} else if *t.SchulzeVote(i, j) < *t.SchulzeVote(j, i) {
				fmt.Fprint(&b, RED)
			} else {
				fmt.Fprint(&b, COL_RESET)
			}
			fmt.Fprint(&b, *t.SchulzeVote(i, j), " ")
		}
		fmt.Fprintln(&b)
	}
	fmt.Fprint(&b, COL_RESET)

	return b.String()
}

func main() {
	t := MakeTable(6)
	fmt.Print("Voting...")

	users := make([][]int, 3)
	users[0] = []int{1,1,1,5,5,5}
	users[1] = []int{5,5,5,1,5,5}
	users[2] = []int{5,5,5,5,1,5}
	for i := 0; i < 100000; i++ {
		//t.AddVote(rand.Perm(t.C))
		user := i % 3 //rand.Intn(3)
		vote := make([]int, 6)
		for i, v := range users[user] {
			vote[i] = v + rand.Intn(5) - 2
		}
		t.AddVote(vote)
	}
	fmt.Println("Done")
	
	fmt.Println(t)
	
	fmt.Print("Counting...")
	t.Schulze()
	fmt.Println("Done")
	fmt.Println(t.SchulzeString())
	fmt.Println(t.SchulzeRankString())

}
