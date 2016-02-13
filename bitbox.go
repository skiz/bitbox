package bitbox

// BitBox is a dynamically sized bit container which allows
// bits to be set and read quickly and somewhat efficiently.
type BitBox struct {
	max   int
	Bytes []byte
}

// NewBitBox returns a new BitBox configured for the given
// number of bits.
func NewBitBox(bits int) *BitBox {
	b := &BitBox{}
	b.Resize(bits)
	return b
}

// Resize will reallocate the BitBox, used both to optimize
// future insertions, and to shrink the BitBox as needed.
// The actual allocation will reflect a suitable byte size.
func (b *BitBox) Resize(bits int) int {
	if bits > b.max {
		add := (bits >> 3) - len(b.Bytes)
		if bits&7 > 0 {
			add++
		}
		b.Bytes = append(b.Bytes, make([]byte, add)...)
	} else {
		bytes := (bits >> 3)
		if bits&7 > 0 {
			bytes++
		}
		b.Bytes = append([]byte(nil), b.Bytes[:bytes]...)
	}
	b.max = len(b.Bytes) << 3
	return b.max
}

// position returns the byte index and a bitmask for a given bit.
func (b *BitBox) Position(n int) (int, uint8) {
	byteNum := n >> 3
	n -= (byteNum << 3)
	var bitPos uint8 = 1
	bitPos <<= uint(7 - n)
	return byteNum, bitPos
}

// Get the given bit position, returning true if it is set.
func (b *BitBox) Get(n int) bool {
	if n >= b.max {
		return false
	}
	byteNum, bitPos := b.Position(n)
	return !((b.Bytes[byteNum] & bitPos) == 0)
}

// Set the given bit position to true.
func (b *BitBox) Set(n int) {
	byteNum, bitPos := b.Position(n)
	if n >= b.max {
		if n == 0 {
			b.Resize(1)
		} else {
			b.Resize(n)
		}
	}
	b.Bytes[byteNum] |= bitPos
}

// Unset sets the given bit position to false.
func (b *BitBox) Unset(n int) {
	if n < b.max {
		byteNum, bitPos := b.Position(n)
		b.Bytes[byteNum] &^= bitPos
	}
}

// Toggle flips the given bit position.
func (b *BitBox) Toggle(n int) {
	if n >= b.max {
		b.Resize(n)
	}
	byteNum, bitPos := b.Position(n)
	b.Bytes[byteNum] ^= bitPos
}

// GetByte returns a byte from the BitBox
func (b *BitBox) GetByte(n int) byte {
	if len(b.Bytes) < n {
		return 0x0
	}
	return b.Bytes[n]
}

// Size() returns the bit size of the BitBox
func (b *BitBox) Size() int {
	return b.max //len(b.Bytes) * 8
}

// Clear zeros all bits, effectively setting all bits to false.
func (b *BitBox) Clear() {
	for i := range b.Bytes {
		b.Bytes[i] = 0x00
	}
}

// And executes a bitwise AND operation on the given bit positions.
func (b *BitBox) And(bs []int) bool {
	for _, n := range bs {
		if !b.Get(n) {
			return false
		}
	}
	return true
}

// Or executes a bitwise OR operation on the given bit positions.
func (b *BitBox) Or(bs []int) bool {
	for _, n := range bs {
		if b.Get(n) {
			return true
		}
	}
	return false
}

// Xor executes a bitwise XOR operation on the given bit positions.
func (b *BitBox) Xor(bs []int) bool {
	var c int
	for _, n := range bs {
		if b.Get(n) {
			c++
		}
		if c > 1 {
			return false
		}
	}
	if c == 1 {
		return true
	}
	return false
}
