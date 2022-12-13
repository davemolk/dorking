package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
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

// write saves the search results to a json file.
func (d *dorking) write(data map[string]string) error {
	js, err := d.json(data)
	if err != nil {
		return fmt.Errorf("unable to encode to json: %w", err)
	}
	err = os.WriteFile("results.json", js, 0644)
	if err != nil {
		return fmt.Errorf("writing error: %w", err)
	}
	return nil
}
