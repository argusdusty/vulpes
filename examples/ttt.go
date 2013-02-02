package main

import (
	"fmt"
	"github.com/argusdusty/vulpes"
	"time"
)

type TicTacToe struct {
	Board [9]int
}

// If Turn, it's the first player's move
func (T TicTacToe) Children(Turn bool) []Game {
	Children := make([]Game, 0)
	if Turn {
		for i := 0; i < 9; i++ {
			if T.Board[i] == 0 {
				T.Board[i] = 1
				Children = append(Children, TicTacToe{T.Board})
				T.Board[i] = 0
			}
		}
	} else {
		for i := 0; i < 9; i++ {
			if T.Board[i] == 0 {
				T.Board[i] = -1
				Children = append(Children, TicTacToe{T.Board})
				T.Board[i] = 0
			}
		}
	}
	return Children
}

// We want to solve TTT, not approximate it
func (T TicTacToe) Heuristic(Turn bool) float64 {
	return 0.0
}

// Turn is not used
// Returns 1 if the first player ('X') wins, 2 if the second player ('O') wins, and 0 for a tie
func (T TicTacToe) EndState(Turn bool) uint8 {
	a := T.Board[0]
	b := T.Board[1]
	c := T.Board[2]
	d := T.Board[3]
	e := T.Board[4]
	f := T.Board[5]
	g := T.Board[6]
	h := T.Board[7]
	i := T.Board[8]
	if e != 0 {
		if (a == e && e == i) || (b == e && e == h) || (c == e && e == g) || (d == e && e == f) {
			if e == 1 {
				return 1
			}
			return 2
		}
	}
	if a != 0 {
		if (a == b && a == c) || (a == d && a == g) {
			if a == 1 {
				return 1
			}
			return 2
		}
	}
	if i != 0 {
		if (c == i && f == i) || (g == i && h == i) {
			if i == 1 {
				return 1
			}
			return 2
		}
	}
	if a == 0 || b == 0 || c == 0 || d == 0 || e == 0 || f == 0 || g == 0 || h == 0 || i == 0 {
		return 3
	}
	return 0
}

func (T TicTacToe) ToString() string {
	out := ""
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			x := T.Board[3*i+j]
			if x == 0 {
				out += "_"
			} else if x == 1 {
				out += "X"
			} else {
				out += "O"
			}
		}
		out += "\n"
	}
	return out
}

func main() {
	State := TicTacToe{[9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}}
	Depth := uint32(9)
	Turn := true
	for State.EndState() == 3 {
		t := time.Now()
		Best, Score := vulpes.SolveGame(State, Depth, Turn, -100, 100)
		State = Best.(TicTacToe)
		fmt.Println("Board:")
		fmt.Print(State.ToString)
		fmt.Println("Score:", Score)
		fmt.Println("Time Taken:", time.Now().Sub(t))
		Turn = !Turn
	}
}
