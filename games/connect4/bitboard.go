package connect4

import (
	"math/bits"
)

const (
	filledBoard = bitboard(0b0111111011111101111110111111011111101111110111111)
	colMask     = filledBoard & (filledBoard >> 1) & (filledBoard >> 2) & (filledBoard >> 3)
	rowMask     = filledBoard & (filledBoard >> 7) & (filledBoard >> 14) & (filledBoard >> 21)
	ldiagMask   = filledBoard & (filledBoard >> 8) & (filledBoard >> 16) & (filledBoard >> 24)
	rdiagMask   = filledBoard & (filledBoard >> 6) & (filledBoard >> 12) & (filledBoard >> 18)
)

type bitboard uint64

func (b bitboard) index(row, col int) bool {
	return b&(1<<(row+7*col)) != 0
}

func (b bitboard) set(row, col int) bitboard {
	return b | 1<<(row+7*col)
}

func (b bitboard) filled() bool {
	return b == filledBoard
}

func (b bitboard) isWin() bool {
	return b&(b>>1)&(b>>2)&(b>>3)|b&(b>>6)&(b>>12)&(b>>18)|b&(b>>7)&(b>>14)&(b>>21)|b&(b>>8)&(b>>16)&(b>>24) != 0
}

func (b bitboard) heur(taken bitboard) int {
	var heur int = 0
	ob := taken ^ b

	x1 := b >> 1
	x2 := b >> 2
	x3 := b >> 3
	y1 := ob >> 1
	y2 := ob >> 2
	y3 := ob >> 3
	allow := ob | y1 | y2 | y3
	oallow := b | x1 | x2 | x3
	heur += bits.OnesCount64(uint64((b&^(x1&(x2^x3)^x2&x3) ^ x1&^(x2&x3) ^ x2 ^ x3) & colMask &^ allow))
	heur -= bits.OnesCount64(uint64((ob&^(y1&(y2^y3)^y2&y3) ^ y1&^(y2&y3) ^ y2 ^ y3) & colMask &^ oallow))
	heur += bits.OnesCount64(uint64((b&(x1|(x2^x3)^x2&x3)^x1&(x2|x3)^x2&x3)&colMask&^allow)) << 4
	heur -= bits.OnesCount64(uint64((ob&(y1|(y2^y3)^y2&y3)^y1&(y2|y3)^y2&y3)&colMask&^oallow)) << 4
	heur += bits.OnesCount64(uint64((b&(x1&(x2^x3)^x2&x3)^x1&x2&x3)&colMask&^allow)) << 8
	heur -= bits.OnesCount64(uint64((ob&(y1&(y2^y3)^y2&y3)^y1&y2&y3)&colMask&^oallow)) << 8

	x1 = b >> 6
	x2 = b >> 12
	x3 = b >> 18
	y1 = ob >> 6
	y2 = ob >> 12
	y3 = ob >> 18
	allow = ob | y1 | y2 | y3
	oallow = b | x1 | x2 | x3
	heur += bits.OnesCount64(uint64((b&^(x1&(x2^x3)^x2&x3) ^ x1&^(x2&x3) ^ x2 ^ x3) & rdiagMask &^ allow))
	heur -= bits.OnesCount64(uint64((ob&^(y1&(y2^y3)^y2&y3) ^ y1&^(y2&y3) ^ y2 ^ y3) & rdiagMask &^ oallow))
	heur += bits.OnesCount64(uint64((b&(x1|(x2^x3)^x2&x3)^x1&(x2|x3)^x2&x3)&rdiagMask&^allow)) << 4
	heur -= bits.OnesCount64(uint64((ob&(y1|(y2^y3)^y2&y3)^y1&(y2|y3)^y2&y3)&rdiagMask&^oallow)) << 4
	heur += bits.OnesCount64(uint64((b&(x1&(x2^x3)^x2&x3)^x1&x2&x3)&rdiagMask&^allow)) << 8
	heur -= bits.OnesCount64(uint64((ob&(y1&(y2^y3)^y2&y3)^y1&y2&y3)&rdiagMask&^oallow)) << 8

	x1 = b >> 7
	x2 = b >> 14
	x3 = b >> 21
	y1 = ob >> 7
	y2 = ob >> 14
	y3 = ob >> 21
	allow = ob | y1 | y2 | y3
	oallow = b | x1 | x2 | x3
	heur += bits.OnesCount64(uint64((b&^(x1&(x2^x3)^x2&x3) ^ x1&^(x2&x3) ^ x2 ^ x3) & rowMask &^ allow))
	heur -= bits.OnesCount64(uint64((ob&^(y1&(y2^y3)^y2&y3) ^ y1&^(y2&y3) ^ y2 ^ y3) & rowMask &^ oallow))
	heur += bits.OnesCount64(uint64((b&(x1|(x2^x3)^x2&x3)^x1&(x2|x3)^x2&x3)&rowMask&^allow)) << 4
	heur -= bits.OnesCount64(uint64((ob&(y1|(y2^y3)^y2&y3)^y1&(y2|y3)^y2&y3)&rowMask&^oallow)) << 4
	heur += bits.OnesCount64(uint64((b&(x1&(x2^x3)^x2&x3)^x1&x2&x3)&rowMask&^allow)) << 8
	heur -= bits.OnesCount64(uint64((ob&(y1&(y2^y3)^y2&y3)^y1&y2&y3)&rowMask&^oallow)) << 8

	x1 = b >> 8
	x2 = b >> 16
	x3 = b >> 24
	y1 = ob >> 8
	y2 = ob >> 16
	y3 = ob >> 24
	allow = ob | y1 | y2 | y3
	oallow = b | x1 | x2 | x3
	heur += bits.OnesCount64(uint64((b&^(x1&(x2^x3)^x2&x3) ^ x1&^(x2&x3) ^ x2 ^ x3) & ldiagMask &^ allow))
	heur -= bits.OnesCount64(uint64((ob&^(y1&(y2^y3)^y2&y3) ^ y1&^(y2&y3) ^ y2 ^ y3) & ldiagMask &^ oallow))
	heur += bits.OnesCount64(uint64((b&(x1|(x2^x3)^x2&x3)^x1&(x2|x3)^x2&x3)&ldiagMask&^allow)) << 4
	heur -= bits.OnesCount64(uint64((ob&(y1|(y2^y3)^y2&y3)^y1&(y2|y3)^y2&y3)&ldiagMask&^oallow)) << 4
	heur += bits.OnesCount64(uint64((b&(x1&(x2^x3)^x2&x3)^x1&x2&x3)&ldiagMask&^allow)) << 8
	heur -= bits.OnesCount64(uint64((ob&(y1&(y2^y3)^y2&y3)^y1&y2&y3)&ldiagMask&^oallow)) << 8
	return heur
}
