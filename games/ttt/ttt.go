package ttt

import (
	"github.com/argusdusty/vulpes"
)

type ttt struct {
	board [9]int
	turn  bool
}

func (t ttt) Children() []vulpes.Game {
	children := make([]vulpes.Game, 0, 9)
	for i := 0; i < 9; i++ {
		if t.board[i] == 0 {
			if t.turn {
				t.board[i] = 1
			} else {
				t.board[i] = -1
			}
			children = append(children, ttt{t.board, !t.turn})
			t.board[i] = 0
		}
	}
	return children
}

func midpointEval(a, b, c, d, e, f, g, h, i int, turn bool) int {
	if e == 0 {
		return vulpes.UNFINISHED
	}
	if (a == e && e == i) || (b == e && e == h) || (c == e && e == g) || (d == e && e == f) {
		if (e == 1) == turn {
			return vulpes.WIN
		}
		return vulpes.LOSS
	}
	return vulpes.UNFINISHED
}

func cornerEval(a, b, c, d, g int, turn bool) int {
	if a == 0 {
		return vulpes.UNFINISHED
	}
	if (a == b && a == c) || (a == d && a == g) {
		if (a == 1) == turn {
			return vulpes.WIN
		}
		return vulpes.LOSS
	}
	return vulpes.UNFINISHED
}

func tieEval(a, b, c, d, e, f, g, h, i int) int {
	if a == 0 || b == 0 || c == 0 || d == 0 || e == 0 || f == 0 || g == 0 || h == 0 || i == 0 {
		return vulpes.UNFINISHED
	}
	return vulpes.TIE
}

// We want to solve TTT, not approximate it, so heuristic is always 0
func (t ttt) Evaluate() (ending int, heuristic float64) {
	a := t.board[0]
	b := t.board[1]
	c := t.board[2]
	d := t.board[3]
	e := t.board[4]
	f := t.board[5]
	g := t.board[6]
	h := t.board[7]
	i := t.board[8]
	v := midpointEval(a, b, c, d, e, f, g, h, i, t.turn)
	if v != vulpes.UNFINISHED {
		return v, 0
	}
	v = cornerEval(a, b, c, d, g, t.turn)
	if v != vulpes.UNFINISHED {
		return v, 0
	}
	v = cornerEval(i, h, g, f, c, t.turn)
	if v != vulpes.UNFINISHED {
		return v, 0
	}
	return tieEval(a, b, c, d, e, f, g, h, i), 0
}

func (t ttt) String() string {
	out := ""
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			x := t.board[3*i+j]
			if x == 0 {
				out += "_"
			} else if x == 1 {
				out += "X"
			} else {
				out += "O"
			}
		}
		if i != 2 {
			out += "\n"
		}
	}
	return out
}

// AI uses vulpes to play perfect Tic-Tac-Toe
type AI struct {
	State ttt
	Turn  bool
}

// NewEmptyAI returns a Tic-Tac-Toe AI from an empty board
func NewEmptyAI() *AI {
	return NewAI([9]int{})
}

// NewAI returns a Tic-Tac-Toe AI from a given board
func NewAI(board [9]int) *AI {
	var sum int
	for i := 0; i < 9; i++ {
		sum += board[i]
	}
	var turn bool = true
	if sum == 1 {
		// Player 1 has had an extra move, so it's the other player's turn
		turn = false
	} else if sum != 0 {
		panic("Invalid board")
	}
	return &AI{State: ttt{board: board, turn: turn}}
}

// MakeMove takes the best move (searching to the given depth) and plays it, updating State. If the game is over, it returns the ending state, and makes no changes to State.
func (C *AI) MakeMove(depth uint) float64 {
	best, score := vulpes.SolveGame(C.State, depth)
	C.State = best.(ttt)
	C.Turn = !C.Turn
	return score
}

// String returns a string representation of the game board
func (C *AI) String() string {
	return C.State.String()
}
