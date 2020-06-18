package vulpes

import (
	"math"
	"sort"
)

const (
	// LOSS means the game is over and the current player has lost the game.
	LOSS = iota
	// TIE means the game is over and ended in a tie.
	TIE
	// WIN means the game is over and the current player has won the game.
	WIN
	// UNFINISHED means the game is not yet over.
	UNFINISHED
)

// Game describes a two-player, zero-sum, turn-based game.
type Game interface {
	// Children returns the child nodes from this one. If the game is not ended, this must return at least 1 child.
	Children() []Game
	// Evaluate returns an evaluation of the current game state from the perspective of the current player. 'ending' must be one of {LOSS, TIE, WIN, UNFINISHED}. 'heuristic' is only required when ending is UNFINISHED.
	Evaluate() (ending int, heuristic float64)
}

type moveScore struct {
	moveIndex int
	moveScore float64
}

type moveScores []moveScore

func (s moveScores) Len() int           { return len(s) }
func (s moveScores) Less(i, j int) bool { return s[i].moveScore > s[j].moveScore }
func (s moveScores) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Search returns the computed score of a given state.
func Search(state Game, depth uint, alpha, beta float64) (Game, float64) {
	ending, heuristic := state.Evaluate()
	switch ending {
	case LOSS:
		return state, math.Inf(-1)
	case TIE:
		return state, 0
	case WIN:
		return state, math.Inf(1)
	}
	if depth == 0 {
		return state, heuristic
	}
	children := state.Children()
	moveScores := make(moveScores, len(children))
	for i := range children {
		moveScores[i] = moveScore{i, 0.0}
	}
	var tmpScore float64
	if depth > 1 {
		// Pre-sort the possible moves by their score to speed up the pruning
		for i, child := range children {
			// Depth-0 search to force heuristic scoring
			_, tmpScore = Search(child, 0, -beta, -alpha)
			moveScores[i].moveScore = -tmpScore
		}
		sort.Sort(moveScores)
	}
	var bestChild Game
	for _, moveScore := range moveScores {
		child := children[moveScore.moveIndex]
		_, tmpScore = Search(child, depth-1, -beta, -alpha)
		tmpScore = -tmpScore
		if tmpScore > alpha {
			alpha = tmpScore
			bestChild = child
			if beta <= alpha {
				return bestChild, beta
			}
		}
		if bestChild == nil {
			// Take the first child, in case all the children are terrible.
			bestChild = child
		}
	}
	if bestChild == nil {
		// No possible moves, so return the current state.
		bestChild = state
	}
	return bestChild, alpha
}

// SolveGame takes a starting node for the game, and returns the best child node and its score, after searching to the specified depth
func SolveGame(state Game, depth uint) (Game, float64) {
	return Search(state, depth, math.Inf(-1), math.Inf(1))
}
