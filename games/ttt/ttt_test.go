package ttt

import (
	"fmt"
	"testing"

	"github.com/argusdusty/vulpes"
)

func TestWin(t *testing.T) {
	for _, i := range [2]int{-1, 1} {
		for _, turn := range [2]bool{false, true} {
			for _, board := range [...][9]int{{i, i, i, 0, 0, 0, 0, 0, 0}, {0, 0, 0, i, i, i, 0, 0, 0}, {0, 0, 0, 0, 0, 0, i, i, i}, {i, 0, 0, i, 0, 0, i, 0, 0}, {0, i, 0, 0, i, 0, 0, i, 0}, {0, 0, i, 0, 0, i, 0, 0, i}, {i, 0, 0, 0, i, 0, 0, 0, i}, {0, 0, i, 0, i, 0, i, 0, 0}} {
				ending, _ := ttt{board: board, turn: turn}.Evaluate()
				expectation := vulpes.WIN
				if turn != (i == 1) {
					expectation = vulpes.LOSS
				}
				if ending != expectation {
					t.Errorf("Incorrect ending state, board: %v, move: %v, turn: %v, ending: %v != %v", board, i, turn, ending, expectation)
				}
			}
		}
	}
}

func TestTie(t *testing.T) {
	for bits := 0; bits < 1024; bits++ {
		board := [9]int{}
		for i := 0; i < 9; i++ {
			board[i] = 2*((bits>>uint(i))&1) - 1
		}
		ending, _ := ttt{board: board, turn: (bits >> 9) == 1}.Evaluate()
		if ending == vulpes.UNFINISHED {
			t.Errorf("Incorrect ending state, board: %v, ending: %v != %v", board, ending, vulpes.TIE)
		}
	}
}

func TestNewAI(t *testing.T) {
	c := NewAI([9]int{1, 0, 0, 0, 0, 0, 0, 0, 0})
	if c.Turn != false {
		t.Errorf("Wrong turn for single X board: %v", c)
	}
}

func TestNewInvalidAI(t *testing.T) {
	defer func() {
		if r := recover(); r != "Invalid board" {
			panic(r)
		}
	}()
	NewAI([9]int{-1, 0, 0, 0, 0, 0, 0, 0, 0})
	t.Errorf("Invalid board didn't panic")
}

func TestAI(t *testing.T) {
	c := NewEmptyAI()
	score := c.MakeMove(9)
	target := `X__
___
___`
	if c.String() != target {
		t.Errorf("Bad opening X move: %s != %s", c.String(), target)
	}
	if score != 0 {
		t.Errorf("Non-zero TTT score: %v", score)
	}
	score = c.MakeMove(9)
	target = `X__
_O_
___`
	if c.String() != target {
		t.Errorf("Bad opening O move: %s != %s", c.String(), target)
	}
	if score != 0 {
		t.Errorf("Non-zero TTT score: %v", score)
	}
}

func BenchmarkAI(b *testing.B) {
	for depth := uint(0); depth < 10; depth++ {
		b.Run(fmt.Sprintf("Depth %d", depth), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c := NewEmptyAI()
				c.MakeMove(depth)
			}
		})
	}
}
