# Vulpes [![GoDoc][godoc-badge]][godoc] [![Build Status][travis-ci-badge]][travis-ci] [![Report Card][report-card-badge]][report-card]
A Negamax AI with Alpha-Beta pruning implementation. For zero-sum, turn-based, two-player games, written in pure Go.

Vulpes is designed to be both performant and easy to use, and comes built-in with an example Connect Four AI capable of searching to depths 15-20 in a matter of seconds, and with modifications could probably solve the game in a few hours.

### How?

Just supply 2 functions to fulfill the Game interface:
```go
// Game describes a two-player, zero-sum, turn-based game.
type Game interface {
	// Children returns the child nodes from this one. If the game is not ended, this must return at least 1 child.
	Children() []Game
	// Evaluate returns an evaluation of the current game state from the perspective of the current player. 'ending' must be one of {LOSS, TIE, WIN, UNFINISHED}. 'heuristic' is only required when ending is UNFINISHED.
	Evaluate() (ending int, heuristic float64)
}
```

[travis-ci-badge]:   https://api.travis-ci.org/argusdusty/vulpes.svg?branch=master
[travis-ci]:         https://api.travis-ci.org/argusdusty/vulpes
[godoc-badge]:       https://godoc.org/github.com/argusdusty/vulpes?status.svg
[godoc]:             https://godoc.org/github.com/argusdusty/vulpes
[report-card-badge]: https://goreportcard.com/badge/github.com/argusdusty/vulpes
[report-card]:       https://goreportcard.com/report/github.com/argusdusty/vulpes