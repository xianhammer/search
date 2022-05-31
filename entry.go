package search

// Entry contain a table keyword and other data relevate for searching.
type Entry struct {
	common, id int
	value      []byte
	next       *Entry
}

// ID of this keyword.
func (e *Entry) ID() int {
	return e.id
}

// Len return keyword length.
func (e *Entry) Len() int {
	return len(e.value)
}

// Bytes return keyword as bytes.
func (e *Entry) Bytes() []byte {
	return e.value[:]
}

func (e *Entry) accepted(x int, b byte) bool {
	return x+1 == len(e.value) && e.value[x] == b
}

func (e *Entry) forward(x int, b byte) *Entry {
	// fmt.Printf("   forward(%d, %x) on %v\n", x, b, e)
	for e != nil && x <= e.common && (x >= len(e.value) || e.value[x] < b) {
		e = e.next
	}

	// fmt.Printf("   - new e = %v (%v<%v, %v==%v)\n", e, x, len(e.value), e.value[x], b)
	// if e != nil && x < len(e.value) && e.value[x] >= b {
	if e != nil && x < len(e.value) && e.value[x] == b {
		return e
	}

	return nil
}
