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

func (t *Table) Prefer(x, y int) {
	t.votes[y * t.C + x]++
}

func (t *Table) Vote(i, ii int) *int {
	return &t.votes[ii * t.C + i]
}

func (t *Table) SchulzeVote(i, ii int) *int {
	return &t.schulze_votes[ii * t.C + i]
}

func (t *Table) AddVote(ranks []int) {
	for i, x := range ranks {
		for ii, y := range ranks {
			if x < y {
				t.Prefer(i, ii)
			}
		}
	}
}

func (t *Table) Schulze() {
	for i, _ := range t.schulze_votes {
		t.schulze_votes[i] = 0
	}

	
	for i := 0; i < t.C; i++ {
		for j := 0; j < t.C; j++ {
			if i != j && *t.Vote(i, j) >= *t.Vote(j, i) {
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
	for j := 0; j < t.C; j++ {
		for i := 0; i < t.C; i++ {
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
	for j := 0; j < t.C; j++ {
		for i := 0; i < t.C; i++ {
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
	t := MakeTable(100)
	fmt.Print("Voting...")
	for i := 0; i < 1000000; i++ {
		t.AddVote(rand.Perm(t.C))
	}
	fmt.Println("Done")
	
	fmt.Println(t)
	
	fmt.Print("Counting...")
	t.Schulze()
	fmt.Println("Done")
	fmt.Println(t.SchulzeString())

}
