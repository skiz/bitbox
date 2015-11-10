package bitbox

import (
	"testing"
)

func TestNewBitBox(t *testing.T) {
	b := NewBitBox(129)
	if b.Size() != 136 {
		t.Errorf("Expected bitbox to have size of 136, but was %d", b.Size())
	}
}

func TestBitPosition(t *testing.T) {
	b := &BitBox{}

	by, bi := b.position(0)
	if by != 0 || bi != 0x80 {
		t.Errorf("Expected position(0) to be 0, 0x80, but was %d, %0b", by, bi)
	}

	by, bi = b.position(7)
	if by != 0 || bi != 0x01 {
		t.Errorf("Expected position(7) to be 0, 0x01, but was %d, %0b", by, bi)

	}

	by, bi = b.position(8)
	if by != 1 || bi != 0x80 {
		t.Errorf("Expected position(8) to be 1, 0x80, but was %d, %0b", by, bi)
	}

	by, bi = b.position(36)
	if by != 4 || bi != 0x08 {
		t.Errorf("Expected position(37) to be 4, 0x08, but was %d, %0b", by, bi)
	}
}

func TestBasicBitManipulation(t *testing.T) {
	b := &BitBox{}
	b.Set(5)
	if !b.Get(5) {
		t.Errorf("Expected bit 5 to be true, but got false")
	}

	if b.bytes[0] != uint8(4) {
		t.Errorf("Expected byte uint8(4), but got uint8(%d)", b.bytes[0])
	}

	b.Set(15)
	if !b.Get(15) {
		t.Errorf("Expected bit 15 to be set, but got %v", b.Get(15))
	}

	b.Unset(15)
	if b.Get(15) {
		t.Errorf("Expected bit 15 to be unset, but to %v", b.Get(15))
	}

	if b.max != 16 {
		t.Errorf("Expected Unset to not change max, but changed to %d", b.max)
	}

	if len(b.bytes) != 2 {
		t.Errorf("Expected Unset to not change bytes, but is %d", len(b.bytes))
	}
}

func TestByteExpansion(t *testing.T) {
	b := &BitBox{}
	b.Set(0)
	if b.max != 8 {
		t.Errorf("Expected 8 bit max, but got %d", b.max)
	}
	if len(b.bytes) != 1 {
		t.Errorf("Expected 1 byte, but was %d", len(b.bytes))
	}

	b.Set(15)
	if b.max != 16 {
		t.Errorf("Expected 16 bit max, but got %d", b.max)
	}
	if len(b.bytes) != 2 {
		t.Errorf("Expected 2 bytes, but was %d", len(b.bytes))
	}

	b.Set(641)
	if b.max != 648 {
		t.Errorf("Expected 648 bit max, but got %d", b.max)
	}
	if len(b.bytes) != 81 {
		t.Errorf("Expected 81 bytes, but was %d", len(b.bytes))
	}

	b.Get(4500)
	if b.max != 648 {
		t.Errorf("Get should not expand max, but expanded to %d", b.max)
	}
	if len(b.bytes) != 81 {
		t.Errorf("Expected 81 bytes, but was %d", len(b.bytes))
	}
}

func TestGetByte(t *testing.T) {
	b := &BitBox{}
	b.Set(3)
	if b.GetByte(0) != 0x10 {
		t.Errorf("Expected byte 0 to by 0x10, but was %x", b.GetByte(1))
	}

	if b.GetByte(32) != 0x0 {
		t.Errorf("Expected byte 32 to by 0x00, but was %x", b.GetByte(1))
	}
}

func TestResizing(t *testing.T) {
	b := &BitBox{}
	b.Set(90)

	b.Clear()
	if len(b.bytes) != 12 {
		t.Errorf("Clear expected byte length of 0, but was %d", len(b.bytes))
	}
	for i := 0; i < 12; i++ {
		if b.bytes[i] != 0x00 {
			t.Errorf("Expected all bits to be cleared, but byte %x wasnt", i)
		}
	}

	b.Resize(15)
	if b.max != 16 {
		t.Errorf("Expected Resize to set max to 16, but was %d", b.max)
	}
	if len(b.bytes) != 2 {
		t.Errorf("Expected Resize byte length of 2, but was %d", len(b.bytes))
	}

	b.Set(15)
	b.Resize(15)
	if !b.Get(15) {
		t.Errorf("Expected growing with resize to not clear bit position 15")
	}

	b.Resize(4)
	if b.Get(15) {
		t.Errorf("Expected resize to clear bit position 15")
	}
	if len(b.bytes) != 1 {
		t.Errorf("Expected Resize byte length of 1, but was %d", len(b.bytes))
	}

	if b.Resize(16) != 16 {
		t.Errorf("Expected returned size of 16, but was %d", b.Size())
	}

	b = NewBitBox(63)
	if b.Resize(64) != 64 {
		t.Errorf("Expected returned size of 64, but was %d", b.Size())
	}
}

func TestToggleBits(t *testing.T) {
	b := &BitBox{}

	b.Toggle(23)
	if !b.Get(23) {
		t.Errorf("Expected Toggle to set bit 23, but was %v", b.Get(23))
	}

	b.Toggle(23)
	if b.Get(23) {
		t.Errorf("Expected Toggle to unset bit 23, but was %v", b.Get(23))
	}
}

func TestAnd(t *testing.T) {
	b := &BitBox{}
	b.Set(0)
	b.Set(19)
	b.Set(43)

	and := b.And([]int{0, 19, 43})
	if !and {
		t.Errorf("Expected And to be true but was %v", and)
	}

	and = b.And([]int{0, 19, 99})
	if and {
		t.Errorf("Expected And to be false but was %v", and)
	}

	and = b.And([]int{0})
	if !and {
		t.Errorf("Expected And to be true but was %v", and)
	}

	and = b.And([]int{42})
	if and {
		t.Errorf("Expected And to be true but was %v", and)
	}

}

func TestOr(t *testing.T) {
	b := &BitBox{}
	b.Set(7)
	b.Set(12)

	or := b.Or([]int{0, 7})
	if !or {
		t.Errorf("Expected Or to be true but was %v", or)
	}

	or = b.Or([]int{0, 33})
	if or {
		t.Errorf("Expected Or to be false but was %v", or)
	}

	or = b.Or([]int{1, 2, 3, 4, 7})
	if !or {
		t.Errorf("Expected Or to be true but was %v", or)
	}

	or = b.Or([]int{7, 12})
	if !or {
		t.Errorf("Expected Or to be true but was %v", or)
	}
}

func TestXOr(t *testing.T) {
	b := &BitBox{}
	b.Set(7)
	b.Set(12)

	xor := b.Xor([]int{0, 7})
	if !xor {
		t.Errorf("Expected Xor to be true but was %v", xor)
	}

	xor = b.Xor([]int{0, 33})
	if xor {
		t.Errorf("Expected Xor to be false but was %v", xor)
	}

	xor = b.Xor([]int{1, 2, 3, 4, 7})
	if !xor {
		t.Errorf("Expected Xor to be true but was %v", xor)
	}

	xor = b.Xor([]int{7, 12})
	if xor {
		t.Errorf("Expected Xor to be false but was %v", xor)
	}

	xor = b.Xor([]int{7, 0, 12})
	if xor {
		t.Errorf("Expected Xor to be false but was %v", xor)
	}
}

func BenchmarkToggle(b *testing.B) {
	bb := NewBitBox(10000)
	for n := 0; n < b.N; n++ {
		bb.Toggle(500)
	}
}

func BenchmarkGet(b *testing.B) {
	bb := NewBitBox(10000)
	for n := 0; n < b.N; n++ {
		bb.Get(500)
	}
}

func BenchmarkSet(b *testing.B) {
	bb := NewBitBox(10000)
	for n := 0; n < b.N; n++ {
		bb.Set(500)
	}

}

func BenchmarkUnset(b *testing.B) {
	bb := NewBitBox(10000)
	for n := 0; n < b.N; n++ {
		bb.Unset(500)
	}
}

func BenchmarkClear(b *testing.B) {
	bb := NewBitBox(10000)
	for n := 0; n < b.N; n++ {
		bb.Clear()
	}
}

func BenchmarkTwoAnd(b *testing.B) {
	bb := NewBitBox(10000)
	bb.Set(300)
	bb.Set(500)
	for n := 0; n < b.N; n++ {
		bb.And([]int{300, 500})
	}
}

func BenchmarkThreeAnd(b *testing.B) {
	bb := NewBitBox(10000)
	bb.Set(100)
	bb.Set(300)
	bb.Set(500)
	for n := 0; n < b.N; n++ {
		bb.And([]int{100, 300, 500})
	}
}

func BenchmarkFourAndWorstCase(b *testing.B) {
	bb := NewBitBox(10000)
	bb.Set(100)
	bb.Set(300)
	bb.Set(500)
	bb.Set(800)
	for n := 0; n < b.N; n++ {
		bb.And([]int{100, 300, 500, 800})
	}
}

func BenchmarkFourAndBestCase(b *testing.B) {
	bb := NewBitBox(10000)
	bb.Set(100)
	bb.Set(300)
	bb.Set(500)
	bb.Set(800)
	for n := 0; n < b.N; n++ {
		bb.And([]int{0, 300, 500, 800})
	}
}

func BenchmarkTwoOrFirst(b *testing.B) {
	bb := NewBitBox(10000)
	bb.Set(1700)
	for n := 0; n < b.N; n++ {
		bb.Or([]int{1700, 800})
	}
}

func BenchmarkTwoOrSecond(b *testing.B) {
	bb := NewBitBox(10000)
	bb.Set(1700)
	for n := 0; n < b.N; n++ {
		bb.Or([]int{800, 1700})
	}
}

func BenchmarkThreeOrNone(b *testing.B) {
	bb := NewBitBox(10000)
	for n := 0; n < b.N; n++ {
		bb.Or([]int{800, 1700, 2500})
	}
}

func BenchmarkXorTwo(b *testing.B) {
	bb := NewBitBox(10000)
	bb.Set(1776)
	for n := 0; n < b.N; n++ {
		bb.Or([]int{800, 1776})
	}
}

func BenchmarkXorTwoNone(b *testing.B) {
	bb := NewBitBox(10000)
	bb.Set(1776)
	for n := 0; n < b.N; n++ {
		bb.Or([]int{800, 43})
	}
}

func BenchmarkXorThreeWorstCase(b *testing.B) {
	bb := NewBitBox(10000)
	bb.Set(1776)
	for n := 0; n < b.N; n++ {
		bb.Or([]int{1906, 9999, 1776})
	}
}

func BenchmarkXorThreeBestCase(b *testing.B) {
	bb := NewBitBox(10000)
	bb.Set(1776)
	for n := 0; n < b.N; n++ {
		bb.Or([]int{1776, 1906, 9999})
	}
}
