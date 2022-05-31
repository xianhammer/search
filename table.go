package search

import (
	"bytes"
	"sort"
)

// Table for keyword searchings.
type Table struct {
	entries         []Entry
	dirty           bool
	caseInsensitive bool
}

// NewTable create a new tabæle
func NewTable(caseInsensitiveSearch bool) (t *Table) {
	t = new(Table)
	t.caseInsensitive = caseInsensitiveSearch
	return
}

// Len method from Sort interface.
func (t *Table) Len() int { return len(t.entries) }

// Swap method from Sort interface.
func (t *Table) Swap(i, j int) { t.entries[i], t.entries[j] = t.entries[j], t.entries[i] }

// Less method from Sort interface.
func (t *Table) Less(i, j int) bool { return compare(t.entries[i].value, t.entries[j].value) > 0 }

// Add a new keyword to the table.
func (t *Table) Add(w []byte, id int) {
	if t.caseInsensitive {
		w = bytes.ToLower(w)
	}

	t.entries = append(t.entries, Entry{0, id, w, nil})
	t.dirty = true
}

// MaxPrefixLength return the longest sequence of (prefix) bytes common to the keywords.
// The returned length can be used to create buffer for reparsing.
func (t *Table) MaxPrefixLength() (l int) {
	for i := 0; i < len(t.entries); i++ {
		if l < t.entries[i].common {
			l = t.entries[i].common
		}
	}
	return
}

// NewSearcher prepare the table for searching and return a searcher object.
// Runtime: If table is dirty, O(n*log(n)) does to sorting of arrays. Otherwise O(n), n is number of keywords.
func (t *Table) NewSearcher() *Searcher {
	if t.dirty {
		t.prepareTable()
	}

	s := new(Searcher)
	s.caseInsensitive = t.caseInsensitive
	for i := 0; i < len(t.entries); i++ {
		key := t.entries[i].value[0]
		if s.first[key] == nil {
			s.first[key] = &t.entries[i]
		}
	}

	s.Clear()
	s.firstPush = true
	return s
}

/*
// ParseBytes will read from byte slice b and call the accept callback for each found instance of a keyword.
// Longest matching keyword is returned - so a table of "a" and "abb", will return "abb" and "a" for input "abba".
func (t *Table) ParseBytes(b []byte, accept Acceptor) {
	t.NewSearcher().ParseBytes(b, accept)
}

// Parse will read from reader r and call the accept callback for each found instance of a keyword.
// Longest matching keyword is returned - so a table of "a" and "abb", will return "abb" and "a" for input "abba".
func (t *Table) Parse(r io.Reader, accept Acceptor) {
	t.NewSearcher().Parse(r, accept)
}
*/

func (t *Table) prepareTable() {
	t.dirty = false
	sort.Sort(t)

	for i := 1; i < len(t.entries); i++ {
		t.entries[i-1].common = commonPrefix(t.entries[i-1].value, t.entries[i].value)
		if t.entries[i-1].common > 0 {
			t.entries[i-1].next = &t.entries[i]
		}
	}
}
