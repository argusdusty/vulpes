package connect4

import (
	"fmt"
	"testing"
)

func TestAI(t *testing.T) {
	c := NewEmptyAI()
	c.MakeMove(13)
	target := `_______
_______
_______
_______
_______
___X___`
	if c.String() != target {
		t.Errorf("Bad opening move: %s != %s", c.String(), target)
	}
}

func BenchmarkAI(b *testing.B) {
	for depth := uint(0); depth < 15; depth++ {
		b.Run(fmt.Sprintf("Depth %d", depth), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c := NewEmptyAI()
				c.MakeMove(depth)
			}
		})
	}
}

/*
func BenchmarkAISolve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := NewEmptyAI()
		c.MakeMove(42)
	}
}
*/
