package search

func commonPrefix(a, b []byte) (count int) {
	l := len(a)
	if l > len(b) {
		l = len(b)
	}

	for ; count < l && a[count] == b[count]; count++ {
	}
	return
}

// a < b  ==> >0
// a == b ==> 0
// a > b  ==> <0
func compare(a, b []byte) int {
	shortest, longest := len(a), len(b)
	if shortest > longest {
		shortest, longest = longest, shortest
	}

	for i := 0; i < shortest; i++ {
		if a[i] != b[i] {
			return int(b[i]) - int(a[i])
		}
	}

	return len(b) - len(a)
}
