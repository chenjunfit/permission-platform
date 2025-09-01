package bitring

import "sync"

const (
	bitsPerWord = 64
	bitsMask    = bitsPerWord - 1
	bitsShift   = 6

	defaultSize        = 128
	defaultConsecutive = 3
)

type BitRing struct {
	words       []uint64
	filled      bool
	size        int
	pos         int
	threshold   float64
	consecutive int
	eventCount  int
	mu          sync.RWMutex
}

func NewBitRing(size, consecutive int, threshold float64) *BitRing {
	if size <= 0 {
		size = defaultSize
	}
	if consecutive <= 0 {
		consecutive = defaultConsecutive
	}
	if consecutive > size {
		consecutive = size
	}
	if threshold < 0 {
		threshold = 0
	}
	if threshold > 1 {
		threshold = 1
	}
	return &BitRing{
		words:       make([]uint64, (size+bitsMask)/bitsPerWord),
		size:        size,
		threshold:   threshold,
		consecutive: consecutive,
	}
}
func (b *BitRing) bitAt(index int) bool {
	word := index >> bitsShift
	off := int(index & bitsMask)
	return (b.words[word] >> off & 1) == 1
}
func (b *BitRing) setBit(index int, v bool) {
	word := index >> bitsShift
	off := int(index & bitsMask)
	if v {
		b.words[word] |= 1 << off
	} else {
		b.words[word] &^= 1 << off
	}
}
func (b *BitRing) WindowSize() int {
	if b.filled {
		return b.size
	}
	return b.pos
}

func (b *BitRing) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	// 重置所有状态
	for i := range b.words {
		b.words[i] = 0
	}
	b.pos = 0
	b.filled = false
	b.eventCount = 0
}

func (b *BitRing) IsConditionMet() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	window := b.WindowSize()
	if window == 0 {
		return false
	}
	if window >= b.consecutive {
		allEvents := true
		for i := 1; i <= b.consecutive; i++ {
			pos := (b.pos - i + b.size) % b.size
			if !b.bitAt(pos) {
				allEvents = false
				break
			}
		}
		if allEvents {
			return true
		}
	}
	minWindow := max(b.consecutive, b.size/2)
	if window >= minWindow && float64(b.eventCount)/float64(window) > b.threshold {
		return true
	}
	//if float64(b.eventCount)/float64(window) > b.threshold {
	//	return true
	//}
	return false
}

func (b *BitRing) Add(eventHappend bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	oldBit := b.bitAt(b.pos)
	if b.filled && oldBit {
		b.eventCount--
	}
	b.setBit(b.pos, eventHappend)
	if eventHappend {
		b.eventCount++
	}
	b.pos++
	if b.pos == b.size {
		b.pos = 0
		b.filled = true
	}
}
