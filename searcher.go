package search

import "io"

type SearcherCallback func(start, end, id int, b []byte)

// Searcher contain data for repeated calls to the Match method.
type Searcher struct {
	first           [256]*Entry
	current         *Entry
	x, pos          int
	firstPush       bool
	prevAccepted    bool
	callback        SearcherCallback
	caseInsensitive bool
}

func (s *Searcher) ReadFrom(r io.Reader) (n int64, err error) {
	var buffer [1024]byte
	n_, err := r.Read(buffer[:])
	for err != nil {
		for i, b := range buffer[:n_] {
			s.Push(s.pos+i, b) // if used -> ?
		}
		n += int64(n_)
		s.pos += n_
		n_, err = r.Read(buffer[:])
	}
	return
}

func (s *Searcher) Callback(callback SearcherCallback) {
	s.callback = callback
}

func (s *Searcher) Flush(pos int) {
	if s.prevAccepted && s.current != nil && s.callback != nil {
		s.callback(s.pos, pos, s.current.ID(), s.current.Bytes())
	}
	s.current = nil
	s.prevAccepted = false
	s.firstPush = true
}

func (s *Searcher) Clear() {
	s.firstPush = false
	s.x = 0
	s.pos = 0
	s.current = nil
	s.prevAccepted = false
}

func (s *Searcher) Push(pos int, b byte) (used bool) {
	if s.caseInsensitive && 'A' <= b && b <= 'Z' {
		b += 'a' - 'A'
	}

	if s.current != nil {
		s.x++
		current := s.current.forward(s.x, b)
		accepted := s.current.accepted(s.x, b)
		if s.prevAccepted && current == nil && s.callback != nil {
			s.callback(s.pos, pos-1, s.current.ID(), s.current.Bytes())
		}
		s.prevAccepted = accepted
		s.firstPush = accepted
		s.current = current
		return accepted
	}

	entry := s.first[b]
	if s.firstPush && entry != nil {
		s.Clear()
		s.pos = pos
		s.current = entry
	}

	return s.current != nil
}
