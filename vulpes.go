package vulpes

import "sort"

type Game interface {
	Children(Turn bool) []Game   // Determines the children nodes from this one. Must return an empty list
	Heuristic(Turn bool) float64 // Determines the heuristic score of the current node, from the perspective of the first player
	EndState(Turn bool) uint8    // 0 for tie, 1 for first player win, 2 for loss, 3 for not ended
}

// Just returns the score of a given state.
func Search(State Game, Depth uint32, Turn bool, Alpha, Beta, MinScore, MaxScore float64) float64 {
	End := State.EndState(Turn)
	if End == 1 {
		return MaxScore
	} else if End == 2 {
		return MinScore
	} else if End == 0 {
		return 0
	}
	if Depth <= 0 {
		return State.Heuristic(Turn)
	}
	if Turn {
		if Depth <= 3 {
			for _, Child := range State.Children(true) {
				TempScore := Search(Child, Depth-1, false, Alpha, Beta, MinScore, MaxScore)
				if TempScore >= Alpha {
					Alpha = TempScore
					if Beta <= Alpha {
						return Alpha
					}
				}
			}
		} else {
			Children := State.Children(true)
			s := len(Children)
			Moves := make([]int, s)
			Scores := make([]float64, s)
			for i, Child := range Children {
				Score := Child.Heuristic(true)
				low := 0
				high := i
				for low < high {
					mid := (low + high) >> 1
					if Scores[mid] > Score {
						low = mid + 1
					} else {
						high = mid
					}
				}
				Moves[i] = i
				Scores[i] = Score
				for j := low; j < i; j++ {
					Moves[j], Moves[i] = Moves[i], Moves[j]
					Scores[j], Scores[i] = Scores[i], Scores[j]
				}
			}
			for _, i := range Moves {
				Child := Children[i]
				TempScore := Search(Child, Depth-1, false, Alpha, Beta, MinScore, MaxScore)
				if TempScore >= Alpha {
					Alpha = TempScore
					if Beta <= Alpha {
						return Alpha
					}
				}
			}
		}
		return Alpha
	}
	if Depth <= 3 {
		for _, Child := range State.Children(false) {
			TempScore := Search(Child, Depth-1, true, Alpha, Beta, MinScore, MaxScore)
			if TempScore <= Beta {
				Beta = TempScore
				if Beta <= Alpha {
					return Beta
				}
			}
		}
	} else {
		Children := State.Children(false)
		s := len(Children)
		Moves := make([]int, s)
		Scores := make([]float64, s)
		for i, Child := range Children {
			Score := Child.Heuristic(false)
			low := 0
			high := i
			for low < high {
				mid := (low + high) >> 1
				if Scores[mid] < Score {
					low = mid + 1
				} else {
					high = mid
				}
			}
			Moves[i] = i
			Scores[i] = Score
			for j := low; j < i; j++ {
				Moves[j], Moves[i] = Moves[i], Moves[j]
				Scores[j], Scores[i] = Scores[i], Scores[j]
			}
		}
		for _, i := range Moves {
			Child := Children[i]
			TempScore := Search(Child, Depth-1, true, Alpha, Beta, MinScore, MaxScore)
			if TempScore <= Beta {
				Beta = TempScore
				if Beta <= Alpha {
					return Beta
				}
			}
		}
	}
	return Beta
}

// Takes a starting node for the game, and returns the best child node and it's score
func SolveGame(State Game, Depth uint32, Turn bool, MinScore, MaxScore float64) (Game, float64) {
	End := State.EndState(Turn)
	if End == 1 {
		return State, MaxScore
	} else if End == 2 {
		return State, MinScore
	} else if End == 0 {
		return State, 0
	}
	if Depth <= 0 {
		return State, State.Heuristic(Turn)
	}
	Best := State
	if Turn {
		for _, Child := range State.Children(true) {
			TempScore := Search(Child, Depth-1, false, MinScore, MaxScore, MinScore, MaxScore)
			if TempScore > MinScore {
				Best = Child
				MinScore = TempScore
			}
		}
		return Best, MinScore
	}
	for _, Child := range State.Children(false) {
		TempScore := Search(Child, Depth-1, true, MinScore, MaxScore, MinScore, MaxScore)
		if TempScore < MaxScore {
			Best = Child
			MaxScore = TempScore
		}
	}
	return Best, MaxScore
}
