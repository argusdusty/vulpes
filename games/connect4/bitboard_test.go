package connect4

import (
	"math/rand"
	"testing"
	"time"
)

const (
	N = 100
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type bruteForceBitboard [6][7]bool

func (bf bruteForceBitboard) index(row, col int) bool {
	return bf[row][col]
}

func (bf bruteForceBitboard) set(row, col int) bruteForceBitboard {
	bf[row][col] = true
	return bf
}

func (bf bruteForceBitboard) apply(f func(i, j int)) {
	for i := 0; i < 6; i++ {
		for j := 0; j < 7; j++ {
			f(i, j)
		}
	}
}

func (b bitboard) String() string {
	return frombitboard(b).String()
}

func (bf bruteForceBitboard) String() string {
	var out string
	bf.apply(func(i, j int) {
		if bf.index(i, j) {
			out += "X"
		} else {
			out += "O"
		}
		if j == 6 {
			out += "\n"
		}
	})
	return out
}

func (bf bruteForceBitboard) filled() bool {
	var f bool = true
	bf.apply(func(i, j int) {
		if !bf.index(i, j) {
			f = false
		}
	})
	return f
}

func (bf bruteForceBitboard) cnt4(i0, j0, mi, mj int) int {
	c := 0
	if bf.index(i0, j0) {
		c++
	}
	if bf.index(i0+mi, j0+mj) {
		c++
	}
	if bf.index(i0+2*mi, j0+2*mj) {
		c++
	}
	if bf.index(i0+3*mi, j0+3*mj) {
		c++
	}
	return c
}

func (bf bruteForceBitboard) set4(i0, j0, mi, mj int) bruteForceBitboard {
	bf = bf.set(i0, j0)
	bf = bf.set(i0+mi, j0+mj)
	bf = bf.set(i0+2*mi, j0+2*mj)
	bf = bf.set(i0+3*mi, j0+3*mj)
	return bf
}

func (bf bruteForceBitboard) apply4(f func(i0, j0, mi, mj int)) {
	bf.apply(func(i, j int) {
		if j+3 < 7 {
			f(i, j, 0, 1)
		}
		if i+3 < 6 {
			f(i, j, 1, 0)
		}
		if i+3 < 6 && j+3 < 7 {
			f(i, j, 1, 1)
		}
		if i-3 >= 0 && j+3 < 7 {
			f(i, j, -1, 1)
		}
	})
}

func (bf bruteForceBitboard) highlightRowAdds(cnt int) bruteForceBitboard {
	var r bruteForceBitboard
	bf.apply4(func(i0, j0, mi, mj int) {
		if mi == 0 && mj == 1 && bf.cnt4(i0, j0, mi, mj) == cnt {
			r = r.set(i0, j0)
		}
	})
	return r
}

func (bf bruteForceBitboard) highlightColAdds(cnt int) bruteForceBitboard {
	var r bruteForceBitboard
	bf.apply4(func(i0, j0, mi, mj int) {
		if mi == 1 && mj == 0 && bf.cnt4(i0, j0, mi, mj) == cnt {
			r = r.set(i0, j0)
		}
	})
	return r
}

func (bf bruteForceBitboard) highlightLeftDiagAdds(cnt int) bruteForceBitboard {
	var r bruteForceBitboard
	bf.apply4(func(i0, j0, mi, mj int) {
		if mi == 1 && mj == 1 && bf.cnt4(i0, j0, mi, mj) == cnt {
			r = r.set(i0, j0)
		}
	})
	return r
}

func (bf bruteForceBitboard) highlightRightDiagAdds(cnt int) bruteForceBitboard {
	var r bruteForceBitboard
	bf.apply4(func(i0, j0, mi, mj int) {
		if mi == -1 && mj == 1 && bf.cnt4(i0, j0, mi, mj) == cnt {
			r = r.set(i0, j0)
		}
	})
	return r
}

func (bf bruteForceBitboard) expandAll2s() bruteForceBitboard {
	var r bruteForceBitboard
	bf.apply4(func(i0, j0, mi, mj int) {
		if bf.cnt4(i0, j0, mi, mj) == 2 {
			r = r.set4(i0, j0, mi, mj)
		}
	})
	return r
}

func (bf bruteForceBitboard) expandAll3s() bruteForceBitboard {
	var r bruteForceBitboard
	bf.apply4(func(i0, j0, mi, mj int) {
		if bf.cnt4(i0, j0, mi, mj) == 3 {
			r = r.set4(i0, j0, mi, mj)
		}
	})
	return r
}

func (bf bruteForceBitboard) highlightAll4s() bruteForceBitboard {
	var r bruteForceBitboard
	bf.apply4(func(i0, j0, mi, mj int) {
		if bf.cnt4(i0, j0, mi, mj) == 4 {
			r = r.set4(i0, j0, mi, mj)
		}
	})
	return r
}

func (bf bruteForceBitboard) isWin() bool {
	var f bool
	bf.apply4(func(i0, j0, mi, mj int) {
		if bf.cnt4(i0, j0, mi, mj) == 4 {
			f = true
		}
	})
	return f
}

func (bf bruteForceBitboard) invert() bruteForceBitboard {
	var r bruteForceBitboard
	bf.apply(func(i, j int) {
		if !bf.index(i, j) {
			r = r.set(i, j)
		}
	})
	return r
}

func (bf bruteForceBitboard) popcount() int {
	var c int = 0
	bf.apply(func(i, j int) {
		if bf.index(i, j) {
			c++
		}
	})
	return c
}

func (bf bruteForceBitboard) heur(taken bruteForceBitboard) int {
	var obf bruteForceBitboard
	bf.apply(func(i, j int) {
		if taken.index(i, j) && !bf.index(i, j) {
			obf = obf.set(i, j)
		}
	})
	var heur int = 0
	bf.apply4(func(i, j, mi, mj int) {
		c := bf.cnt4(i, j, mi, mj)
		oc := obf.cnt4(i, j, mi, mj)
		if c > 0 && oc == 0 && c < 4 {
			heur += (1 << (4 * (c - 1)))
		} else if oc > 0 && c == 0 && oc < 4 {
			heur -= (1 << (4 * (oc - 1)))
		}
	})
	return heur
}

func (bf bruteForceBitboard) tobitboard() bitboard {
	var b bitboard
	bf.apply(func(i, j int) {
		if bf.index(i, j) {
			b = b.set(i, j)
		}
	})
	return b
}

func frombitboard(b bitboard) bruteForceBitboard {
	var bf bruteForceBitboard
	bf.apply(func(i, j int) {
		if b.index(i, j) {
			bf = bf.set(i, j)
		}
	})
	return bf
}

func randboard() bruteForceBitboard {
	var bf bruteForceBitboard
	bf.apply(func(i, j int) {
		if rand.Intn(2) == 1 {
			bf = bf.set(i, j)
		}
	})
	return bf
}

func TestBitboardConversion(t *testing.T) {
	for i := 0; i < N; i++ {
		bf := randboard()
		tbf := frombitboard(bf.tobitboard())
		if tbf != bf {
			t.Errorf("Bitboard conversion failed: %v != %v", bf, tbf)
		}
		if tbf.tobitboard() != bf.tobitboard() {
			t.Errorf("Bitboard conversion failed: %v != %v", bf.tobitboard(), tbf.tobitboard())
		}
	}
}

func TestBitboardWin(t *testing.T) {
	c := 0
	for i := 0; i < N; i++ {
		bf := randboard()
		b := bf.tobitboard()
		if bf.isWin() != b.isWin() {
			t.Errorf("Bitboard win check failed: bf: %v (%v), b: %v (%v)", bf, bf.isWin(), b, b.isWin())
		}
		if bf.isWin() {
			c++
		}
	}
	t.Logf("Number of wins: %v/%v", c, N)
}

func TestBitboardFilled(t *testing.T) {
	for i := 0; i < N; i++ {
		bf := randboard()
		b := bf.tobitboard()
		if bf.filled() != b.filled() {
			t.Errorf("Bitboard filled check failed: bf: %v (%v), b: %v (%v)", bf, bf.filled(), b, b.filled())
		}
	}
	b := filledBoard
	if !b.filled() {
		t.Errorf("Filled board not filled?, b: %v", b)
	}
	bf := frombitboard(b)
	if !bf.filled() {
		t.Errorf("Filled bf board not filled?, bf: %v", bf)
	}
}

func TestBitboardHeur(t *testing.T) {
	for i := 0; i < N; i++ {
		bf := randboard()
		obf := randboard()
		b := bf.tobitboard()
		ob := obf.tobitboard()
		taken := ob | b
		takenf := frombitboard(taken)
		if bf.heur(takenf) != b.heur(taken) {
			t.Errorf("Bitboard heur comparison failed: b: %v, taken: %v, ob: %v, (%v != %v)", b, taken, taken^b, bf.heur(takenf), b.heur(taken))
		}
	}
	b := filledBoard
	if b.heur(b) != 0 {
		t.Errorf("Filled board heur fail, b: %v (%v)", b, b.heur(b))
	}
	if bitboard(0).heur(0) != 0 {
		t.Errorf("Empty board heur fail, b: %v (%v)", 0, bitboard(0).heur(0))
	}
}
