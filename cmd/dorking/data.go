package main

import (
	"bytes"
	"encoding/json"
	"sync"
)

// searchMap contains a mutex-protected map for storing search
// results (formatted as URL:blurb).
type searchMap struct {
	mu       sync.Mutex
	searches map[string]string
}

// newSearchMap initializes and returns a new searchMap.
func newSearchMap() *searchMap {
	return &searchMap{
		searches: make(map[string]string),
	}
}

// store checks if a URL has already been stored. If it hasn't,
// both the URL and the associated blurb will be stored.
func (s *searchMap) store(url, blurb string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.searches[url]; !ok {
		s.searches[url] = blurb
	}
}

// json encodes the search results to json and returns
// a byte slice and any errors.
func (d *dorking) json(data map[string]string) ([]byte, error) {
	buf := &bytes.Buffer{}
	e := json.NewEncoder(buf)
	e.SetEscapeHTML(false)
	e.SetIndent("", "    ")
	err := e.Encode(data)
	return bytes.TrimRight(buf.Bytes(), "\n"), err
}
