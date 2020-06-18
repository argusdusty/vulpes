package connect4

import (
	"github.com/argusdusty/vulpes"
)

type connect4 struct {
	currentPlayer bitboard
	taken         bitboard
}

func (c connect4) canPlay(col int) bool {
	return !c.taken.index(5, col)
}

func (c connect4) play(col int) connect4 {
	// Swap the current player to the next player
	currentPlayer := c.currentPlayer ^ c.taken
	// Since taken columns are of the form 0...01...1, adding 1 to them turns them into 0...10...0, so you can | them to produce a new taken spots with the new move
	taken := c.taken | (c.taken + (1 << (7 * col)))
	// Now that taken has been updated, the last player (now currentPlayer ^ taken) has their move set.
	return connect4{currentPlayer, taken}
}

func (c connect4) Children() []vulpes.Game {
	children := make([]vulpes.Game, 0, 7)
	for j := 0; j < 7; j++ {
		if !c.canPlay(j) {
			continue
		}
		children = append(children, c.play(j))
	}
	return children
}

func (c connect4) Evaluate() (ending int, heuristic float64) {
	//fmt.Println("1a", (c.currentPlayer ^ c.taken).isWin(), c.taken.filled(), c.score())
	// If there is a win, it must be from the previous player.
	if (c.currentPlayer ^ c.taken).isWin() {
		return vulpes.LOSS, 0.0
	}
	if c.taken.filled() {
		return vulpes.TIE, 0
	}
	// Game's not ovr yet, so compute the heuristic scores as the number of free spaces available for a set of 4
	return vulpes.UNFINISHED, float64(c.currentPlayer.heur(c.taken))
}

func (c connect4) String(turn bool) string {
	out := ""
	for i := 5; i >= 0; i-- {
		for j := 0; j < 7; j++ {
			x1 := c.currentPlayer.index(i, j)
			x2 := c.taken.index(i, j)
			if !x2 {
				out += "_"
			} else if x1 == turn {
				out += "X"
			} else {
				out += "O"
			}
		}
		if i > 0 {
			out += "\n"
		}
	}
	return out
}

// AI uses vulpes to play Connect Four
type AI struct {
	State connect4
	Turn  bool
}

// NewEmptyAI returns a Connect4 AI from an empty board
func NewEmptyAI() *AI {
	return NewAI([6][7]int{})
}

// NewAI returns a Connect4 AI from a given board
func NewAI(board [6][7]int) *AI {
	var currentPlayer bitboard
	var nextPlayer bitboard
	var sum int
	for i := 0; i < 6; i++ {
		for j := 0; j < 7; j++ {
			sum += board[i][j]
			if board[i][j] == 1 {
				currentPlayer.set(i, j)
			} else if board[i][j] == -1 {
				nextPlayer.set(i, j)
			}
		}
	}
	if sum == 1 {
		// Player 1 has had an extra move, so it's the other player's turn
		currentPlayer, nextPlayer = nextPlayer, currentPlayer
	} else if sum != 0 {
		panic("Invalid board")
	}
	return &AI{State: connect4{currentPlayer: currentPlayer, taken: currentPlayer ^ nextPlayer}, Turn: sum != 1}
}

// MakeMove takes the best move (searching to the given depth) and plays it, updating State. If the game is over, it returns the ending state, and makes no changes to State.
func (C *AI) MakeMove(depth uint) float64 {
	best, score := vulpes.SolveGame(C.State, depth)
	C.State = best.(connect4)
	C.Turn = !C.Turn
	return score
}

// String returns a string representation of the game board
func (C *AI) String() string {
	return C.State.String(C.Turn)
}
