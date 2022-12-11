package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type searchMap struct {
	mu       sync.Mutex
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

func (d *dorking) encode(data map[string]string) ([]byte, error) {
	buf := &bytes.Buffer{}
	e := json.NewEncoder(buf)
	e.SetEscapeHTML(false)
	e.SetIndent("", "    ")
	err := e.Encode(data)
	return bytes.TrimRight(buf.Bytes(), "\n"), err
}

func (d *dorking) json() ([]byte, error) {
	data, err := d.encode(d.searches.searches)
	if err != nil {
		return nil, fmt.Errorf("encoding error: %w", err)
	}
	return data, nil
}

func (d *dorking) write() error {
	data, err := d.json()
	if err != nil {
		return fmt.Errorf("unable to get json: %w", err)
	}
	err = os.WriteFile("results.json", data, 0644)
	if err != nil {
		return fmt.Errorf("writing error: %w", err)
	}
	return nil
}
