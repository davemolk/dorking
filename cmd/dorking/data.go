package main

import "sync"

type searchMap struct {
	mu sync.Mutex
	searches map[string]string
}

func newSearchMap() *searchMap {
	return &searchMap{
		searches: make(map[string]string),
	}
}

func (s *searchMap) store(url, blurb string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.searches[url]; !ok {
		s.searches[url] = blurb
	}
}