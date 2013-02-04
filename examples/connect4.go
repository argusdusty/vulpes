package main

import (
	"fmt"
	"github.com/argusdusty/vulpes"
	"time"
)

const DEPTH = uint32(12)

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
	// Yes, I know this isn't as good as EndState's
	a, b1, c, d, e, f, g := b[0][0], b[0][1], b[0][2], b[0][3], b[0][4], b[0][5], b[0][6]
	h, i, j, k, l, m, n := b[1][0], b[1][1], b[1][2], b[1][3], b[1][4], b[1][5], b[1][6]
	o, p, q, r, s, t, u := b[2][0], b[2][1], b[2][2], b[2][3], b[2][4], b[2][5], b[2][6]
	v, w, x, y, z, aa, ab := b[3][0], b[3][1], b[3][2], b[3][3], b[3][4], b[3][5], b[3][6]
	ac, ad, ae, af, ag, ah, ai := b[4][0], b[4][1], b[4][2], b[4][3], b[4][4], b[4][5], b[4][6]
	aj, ak, al, am, an, ao, ap := b[5][0], b[5][1], b[5][2], b[5][3], b[5][4], b[5][5], b[5][6]

	score += count7[729*a+243*b1+81*c+27*d+9*e+3*f+g+1093]
	score += count7[729*h+243*i+81*j+27*k+9*l+3*m+n+1093]
	score += count7[729*o+243*p+81*q+27*r+9*s+3*t+u+1093]
	score += count7[729*v+243*w+81*x+27*y+9*z+3*aa+ab+1093]
	score += count7[729*ac+243*ad+81*ae+27*af+9*ag+3*ah+ai+1093]
	score += count7[729*aj+243*ak+81*al+27*am+9*an+3*ao+ap+1093]

	score += count6[243*a+81*h+27*o+9*v+3*ac+aj+364]
	score += count6[243*b1+81*i+27*p+9*w+3*ad+ak+364]
	score += count6[243*c+81*j+27*q+9*x+3*ae+al+364]
	score += count6[243*d+81*k+27*r+9*y+3*af+am+364]
	score += count6[243*e+81*l+27*s+9*z+3*ag+an+364]
	score += count6[243*f+81*m+27*t+9*aa+3*ah+ao+364]
	score += count6[243*g+81*n+27*u+9*ab+3*ai+ap+364]

	score += count6[243*a+81*i+27*q+9*y+3*ag+ao+364]
	score += count6[243*aj+81*ad+27*x+9*r+3*l+f+364]
	score += count6[243*b1+81*j+27*r+9*z+3*ah+ap+364]
	score += count6[243*ak+81*ae+27*y+9*s+3*m+g+364]

	score += count5[81*c+27*k+9*s+3*aa+ai+121]
	score += count5[81*al+27*af+9*z+3*t+n+121]
	score += count5[81*h+27*p+9*x+3*af+an+121]
	score += count5[81*ac+27*w+9*q+3*k+e+121]

	score += count4[27*d+9*l+3*t+ab+40]
	score += count4[27*am+9*ag+3*aa+u+40]
	score += count4[27*o+9*w+3*ae+am+40]
	score += count4[27*v+9*p+3*j+d+40]

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
			if b03 == 1 {
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
