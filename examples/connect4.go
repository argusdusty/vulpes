package main

import (
	"fmt"
	"github.com/argusdusty/vulpes"
	"time"
)

const DEPTH = uint32(13)

// Computes a heuristic score for a row of 4
// if all pieces are for the same player, return +-16^n, else 0
func countd(a, b, c, d int) int {
	r := 0
	hr := uint8(0)
	hb := uint8(0)
	if a == 1 {
		hr += 4
	} else if a == -1 {
		hb += 4
	}
	if b == 1 {
		hr += 4
	} else if b == -1 {
		hb += 4
	}
	if c == 1 {
		hr += 4
	} else if c == -1 {
		hb += 4
	}
	if d == 1 {
		hr += 4
	} else if d == -1 {
		hb += 4
	}

	if hr == 0 {
		if hb != 0 {
			r -= (1 << hb)
		}
	} else if hb == 0 {
		r += (1 << hr)
	}
	return r
}

// Produces a heuristic score for a row of 5
func countc(a, b, c, d, e int) int {
	return countd(a, b, c, d) + countd(b, c, d, e)
}

// Produces a heuristic score for a row of 6
func countb(a, b, c, d, e, f int) int {
	return countc(a, b, c, d, e) + countd(c, d, e, f)
}

// Produces a heuristic score for a row of 7
func counta(a, b, c, d, e, f, g int) int {
	return countb(a, b, c, d, e, f) + countd(d, e, f, g)
}

// Precompute the possible values
var count4 [81]int
var count5 [243]int
var count6 [729]int
var count7 [2187]int

func init_counts() {
	for a := 0; a < 3; a++ {
		for b := 0; b < 3; b++ {
			for c := 0; c < 3; c++ {
				for d := 0; d < 3; d++ {
					count4[27*a+9*b+3*c+d] = countd(a-1, b-1, c-1, d-1)
					for e := 0; e < 3; e++ {
						count5[81*a+27*b+9*c+3*d+e] = countc(a-1, b-1, c-1, d-1, e-1)
						for f := 0; f < 3; f++ {
							count6[243*a+81*b+27*c+9*d+3*e+f] = countb(a-1, b-1, c-1, d-1, e-1, f-1)
							for g := 0; g < 3; g++ {
								count7[729*a+243*b+81*c+27*d+9*e+3*f+g] = counta(a-1, b-1, c-1, d-1, e-1, f-1, g-1)
							}
						}
					}
				}
			}
		}
	}
}

type Connect4 struct {
	Board [6][7]int
}

func (C Connect4) Children(Turn bool) []vulpes.Game {
	Children := make([]vulpes.Game, 0)
	if Turn {
		for j := 0; j < 7; j++ {
			i := 0
			for {
				if C.Board[i][j] != 0 || i == 5 {
					break
				}
				i++
			}
			if C.Board[i][j] != 0 {
				if i == 0 {
					continue
				}
				i--
			}
			C.Board[i][j] = 1
			Children = append(Children, Connect4{C.Board})
			C.Board[i][j] = 0
		}
	} else {
		for j := 0; j < 7; j++ {
			i := 0
			for {
				if C.Board[i][j] != 0 || i == 5 {
					break
				}
				i++
			}
			if C.Board[i][j] != 0 {
				if i == 0 {
					continue
				}
				i--
			}
			C.Board[i][j] = -1
			Children = append(Children, Connect4{C.Board})
			C.Board[i][j] = 0
		}
	}
	return Children
}

func (C Connect4) Heuristic(Turn bool) float64 {
	score := 0

	b := C.Board
	b0, b1, b2, b3, b4, b5 := b[0], b[1], b[2], b[3], b[4], b[5]
	b00, b01, b02, b03, b04, b05, b06 := b0[0], b0[1], b0[2], b0[3], b0[4], b0[5], b0[6]
	b10, b11, b12, b13, b14, b15, b16 := b1[0], b1[1], b1[2], b1[3], b1[4], b1[5], b1[6]
	b20, b21, b22, b23, b24, b25, b26 := b2[0], b2[1], b2[2], b2[3], b2[4], b2[5], b2[6]
	b30, b31, b32, b33, b34, b35, b36 := b3[0], b3[1], b3[2], b3[3], b3[4], b3[5], b3[6]
	b40, b41, b42, b43, b44, b45, b46 := b4[0], b4[1], b4[2], b4[3], b4[4], b4[5], b4[6]
	b50, b51, b52, b53, b54, b55, b56 := b5[0], b5[1], b5[2], b5[3], b5[4], b5[5], b5[6]

	// Rows
	score += count7[729*b00+243*b01+81*b02+27*b03+9*b04+3*b05+b06+1093]
	score += count7[729*b10+243*b11+81*b12+27*b13+9*b14+3*b15+b16+1093]
	score += count7[729*b20+243*b21+81*b22+27*b23+9*b24+3*b25+b26+1093]
	score += count7[729*b30+243*b31+81*b32+27*b33+9*b34+3*b35+b36+1093]
	score += count7[729*b40+243*b41+81*b42+27*b43+9*b44+3*b45+b46+1093]
	score += count7[729*b50+243*b51+81*b52+27*b53+9*b54+3*b55+b56+1093]

	// Columns
	score += count6[243*b00+81*b10+27*b20+9*b30+3*b40+b50+364]
	score += count6[243*b01+81*b11+27*b21+9*b31+3*b41+b51+364]
	score += count6[243*b02+81*b12+27*b22+9*b32+3*b42+b52+364]
	score += count6[243*b03+81*b13+27*b23+9*b33+3*b43+b53+364]
	score += count6[243*b04+81*b14+27*b24+9*b34+3*b44+b54+364]
	score += count6[243*b05+81*b15+27*b25+9*b35+3*b45+b55+364]
	score += count6[243*b06+81*b16+27*b26+9*b36+3*b46+b56+364]

	// Length-6 diagonals
	score += count6[243*b00+81*b11+27*b22+9*b33+3*b44+b55+364]
	score += count6[243*b50+81*b41+27*b32+9*b23+3*b14+b05+364]
	score += count6[243*b01+81*b12+27*b23+9*b34+3*b45+b56+364]
	score += count6[243*b51+81*b42+27*b33+9*b24+3*b15+b06+364]

	// Length 5 diagonals
	score += count5[81*b02+27*b13+9*b24+3*b35+b46+121]
	score += count5[81*b52+27*b43+9*b34+3*b25+b16+121]
	score += count5[81*b10+27*b21+9*b32+3*b43+b54+121]
	score += count5[81*b40+27*b31+9*b22+3*b13+b04+121]

	// Length 4 diagonals
	score += count4[27*b03+9*b14+3*b25+b36+40]
	score += count4[27*b53+9*b44+3*b35+b26+40]
	score += count4[27*b20+9*b31+3*b42+b53+40]
	score += count4[27*b30+9*b21+3*b12+b03+40]

	return float64(score)
}

func (C Connect4) EndState(Turn bool) uint8 {
	found := false
	for i := 0; i < 6; i++ {
		for j := 0; j < 7; j++ {
			if C.Board[i][j] == 0 {
				found = true
			}
		}
	}
	if !found {
		return 0
	}

	var x1, x2 bool

	b := C.Board
	b0, b1, b2, b3, b4, b5 := b[0], b[1], b[2], b[3], b[4], b[5]
	b00, b01, b02, b03, b04, b05, b06 := b0[0], b0[1], b0[2], b0[3], b0[4], b0[5], b0[6]
	b10, b11, b12, b13, b14, b15, b16 := b1[0], b1[1], b1[2], b1[3], b1[4], b1[5], b1[6]
	b20, b21, b22, b23, b24, b25, b26 := b2[0], b2[1], b2[2], b2[3], b2[4], b2[5], b2[6]
	b30, b31, b32, b33, b34, b35, b36 := b3[0], b3[1], b3[2], b3[3], b3[4], b3[5], b3[6]
	b40, b41, b42, b43, b44, b45, b46 := b4[0], b4[1], b4[2], b4[3], b4[4], b4[5], b4[6]
	b50, b51, b52, b53, b54, b55, b56 := b5[0], b5[1], b5[2], b5[3], b5[4], b5[5], b5[6]

	// Rows
	if b03 != 0 {
		x1 = (b04 == b03)
		x2 = (b02 == b03)
		if (((b01 == b03) && x2) && ((b00 == b03) || x1)) || ((x2 || (b06 == b03)) && ((b05 == b03) && x1)) {
			if b03 == 1 {
				return 1
			}
			return 2
		}
	}
	if b13 != 0 {
		x1 = (b14 == b13)
		x2 = (b12 == b13)
		if (((b11 == b13) && x2) && ((b10 == b13) || x1)) || ((x2 || (b16 == b13)) && ((b15 == b13) && x1)) {
			if b13 == 1 {
				return 1
			}
			return 2
		}
	}
	if b23 != 0 {
		x1 = (b24 == b23)
		x2 = (b22 == b23)
		if (((b21 == b23) && x2) && ((b20 == b23) || x1)) || ((x2 || (b26 == b23)) && ((b25 == b23) && x1)) {
			if b23 == 1 {
				return 1
			}
			return 2
		}
	}
	if b33 != 0 {
		x1 = (b34 == b33)
		x2 = (b32 == b33)
		if (((b31 == b33) && x2) && ((b30 == b33) || x1)) || ((x2 || (b36 == b33)) && ((b35 == b33) && x1)) {
			if b33 == 1 {
				return 1
			}
			return 2
		}
	}
	if b43 != 0 {
		x1 = (b44 == b43)
		x2 = (b42 == b43)
		if (((b41 == b43) && x2) && ((b40 == b43) || x1)) || ((x2 || (b46 == b43)) && ((b45 == b43) && x1)) {
			if b43 == 1 {
				return 1
			}
			return 2
		}
	}
	if b53 != 0 {
		x1 = (b54 == b53)
		x2 = (b52 == b53)
		if (((b51 == b53) && x2) && ((b50 == b53) || x1)) || ((x2 || (b56 == b53)) && ((b55 == b53) && x1)) {
			if b53 == 1 {
				return 1
			}
			return 2
		}
	}

	// Columns
	if (b20 != 0) && (b20 == b30) {
		x1 = (b10 == b20)
		if ((b40 == b20) && (x1 || (b50 == b20))) || (x1 && (b00 == b20)) {
			if b20 == 1 {
				return 1
			}
			return 2
		}
	}
	if (b21 != 0) && (b21 == b31) {
		x1 = (b11 == b21)
		if ((b41 == b21) && (x1 || (b51 == b21))) || (x1 && (b01 == b21)) {
			if b21 == 1 {
				return 1
			}
			return 2
		}
	}
	if (b22 != 0) && (b22 == b32) {
		x1 = (b12 == b22)
		if ((b42 == b22) && (x1 || (b52 == b22))) || (x1 && (b02 == b22)) {
			if b22 == 1 {
				return 1
			}
			return 2
		}
	}
	if (b23 != 0) && (b23 == b33) {
		x1 = (b13 == b23)
		if ((b43 == b23) && (x1 || (b53 == b23))) || (x1 && (b03 == b23)) {
			if b23 == 1 {
				return 1
			}
			return 2
		}
	}
	if (b24 != 0) && (b24 == b34) {
		x1 = (b14 == b24)
		if ((b44 == b24) && (x1 || (b54 == b24))) || (x1 && (b04 == b24)) {
			if b24 == 1 {
				return 1
			}
			return 2
		}
	}
	if (b25 != 0) && (b25 == b35) {
		x1 = (b15 == b25)
		if ((b45 == b25) && (x1 || (b55 == b25))) || (x1 && (b05 == b25)) {
			if b25 == 1 {
				return 1
			}
			return 2
		}
	}
	if (b26 != 0) && (b26 == b36) {
		x1 = (b16 == b26)
		if ((b46 == b26) && (x1 || (b56 == b26))) || (x1 && (b06 == b26)) {
			if b26 == 1 {
				return 1
			}
			return 2
		}
	}

	// Length 6 diagonals
	if (b22 != 0) && (b22 == b33) {
		x1 = (b11 == b22)
		if ((b44 == b22) && (x1 || (b55 == b22))) || (x1 && (b00 == b22)) {
			if b22 == 1 {
				return 1
			}
			return 2
		}
	}
	if (b32 != 0) && (b32 == b23) {
		x1 = (b41 == b32)
		if ((b14 == b32) && (x1 || (b05 == b32))) || (x1 && (b50 == b32)) {
			if b32 == 1 {
				return 1
			}
			return 2
		}
	}
	if (b23 != 0) && (b23 == b34) {
		x1 = (b12 == b23)
		if ((b45 == b23) && (x1 || (b56 == b23))) || (x1 && (b01 == b23)) {
			if b23 == 1 {
				return 1
			}
			return 2
		}
	}
	if (b33 != 0) && (b33 == b24) {
		x1 = (b42 == b33)
		if ((b15 == b33) && (x1 || (b06 == b33))) || (x1 && (b51 == b33)) {
			if b33 == 1 {
				return 1
			}
			return 2
		}
	}

	// Length 5 diagonals
	if (b34 != 0) && (b43 == b34) && (b25 == b34) && ((b52 == b34) || (b16 == b34)) {
		if b34 == 1 {
			return 1
		}
		return 2
	}
	if (b24 != 0) && (b14 == b24) && (b35 == b24) && ((b02 == b24) || (b46 == b24)) {
		if b24 == 1 {
			return 1
		}
		return 2
	}
	if (b32 != 0) && (b21 == b32) && (b43 == b32) && ((b10 == b32) || (b54 == b32)) {
		if b32 == 1 {
			return 1
		}
		return 2
	}
	if (b22 != 0) && (b31 == b22) && (b13 == b22) && ((b40 == b22) || (b04 == b22)) {
		if b22 == 1 {
			return 1
		}
		return 2
	}

	// Length 4 diagonals
	if (b03 != 0) && (b14 == b03) && (b25 == b03) && (b36 == b03) {
		if b03 == 1 {
			return 1
		}
		return 2
	}
	if (b53 != 0) && (b44 == b53) && (b35 == b53) && (b26 == b53) {
		if b53 == 1 {
			return 1
		}
		return 2
	}
	if (b20 != 0) && (b31 == b20) && (b42 == b20) && (b53 == b20) {
		if b20 == 1 {
			return 1
		}
		return 2
	}
	if (b30 != 0) && (b21 == b30) && (b12 == b30) && (b03 == b30) {
		if b30 == 1 {
			return 1
		}
		return 2
	}
	return 3
}

func (C Connect4) ToString() string {
	out := ""
	for i := 0; i < 6; i++ {
		for j := 0; j < 7; j++ {
			x := C.Board[i][j]
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
	init_counts()
	var Board [6][7]int
	for i := 0; i < 6; i++ {
		for j := 0; j < 7; j++ {
			Board[i][j] = 0
		}
	}
	State := Connect4{Board}
	Turn := true
	for State.EndState(Turn) == 3 {
		t := time.Now()
		Best, Score := vulpes.SolveGame(State, DEPTH, Turn, -0x100000, 0x100000)
		State = Best.(Connect4)
		fmt.Println("Board:")
		fmt.Print(State.ToString())
		fmt.Println("Score:", Score)
		fmt.Println("Time Taken:", time.Now().Sub(t))
		Turn = !Turn
	}
	fmt.Println(State.EndState(Turn))
}
