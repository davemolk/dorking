package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// seed random number generator to get random user agents
func init() {
	rand.Seed(time.Now().UnixNano())
}

// makeRequest takes in a URL, makes a GET request, and returns
// the results as a *bytes.Buffer (along with any errors).
func (d *dorking) makeRequest(url string) (*bytes.Buffer, error) {
	if d.config.verbose {
		log.Printf("requesting %s\n", url)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(d.config.timeout)*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't make request for %s: %v", url, err)
	}

	req = d.headers(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't get response for %s: %v", url, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP response: %d for %s", resp.StatusCode, url)
	}

	g, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to decode gzip for %s: %w", url, err)
	}

	b, err := io.ReadAll(g)
	if err != nil {
		return nil, fmt.Errorf("unable to read body for %s: %v", url, err)
	}
	
	buf := bytes.NewBuffer(b)
	return buf, nil
}

// headers randomly adds firefox or chrome headers to the request.
func (d *dorking) headers(r *http.Request) *http.Request {
	if rand.Intn(2) == 1 {
		return d.ff(r)
	}
	return d.chrome(r)
}

// ff returns a request with firefox headers added.
func (d *dorking) ff(r *http.Request) *http.Request {
	uAgent := d.ffUA()
	r.Header.Set("Host", r.URL.Host)
	r.Header.Set("User-Agent", uAgent)
	r.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	r.Header.Set("Accept-Language", "en-US,en;q=0.5")
	r.Header.Set("Accept-Encoding", "gzip")
	r.Header.Set("DNT", "1")
	r.Header.Set("Connection", "keep-alive")
	r.Header.Set("Upgrade-Insecure-Requests", "1")
	r.Header.Set("Sec-Fetch-Dest", "document")
	r.Header.Set("Sec-Fetch-Mode", "navigate")
	r.Header.Set("Sec-Fetch-Site", "none")
	r.Header.Set("Sec-Fetch-User", "?1")
	r.Header.Set("Sec-GCP", "1")
	return r
}

// chrome returns a request with chrome headers added.
func (d *dorking) chrome(r *http.Request) *http.Request {
	uAgent := d.chromeUA()
	r.Header.Set("Host", r.URL.Host)
	r.Header.Set("Connection", "keep-alive")
	r.Header.Set("Cache-Control", "max-age=0")
	r.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="99", "Google Chrome";v="99"`)
	r.Header.Set("sec-ch-ua-mobile", "?0")
	switch d.config.os {
	case "m":
		r.Header.Set("sec-ch-ua-platform", "Macintosh")
	default:
		r.Header.Set("sec-ch-ua-platform", "Windows")
	}
	r.Header.Set("Upgrade-Insecure-Requests", "1")
	r.Header.Set("User-Agent", uAgent)
	r.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	r.Header.Set("Sec-Fetch-Site", "none")
	r.Header.Set("Sec-Fetch-Mode", "navigate")
	r.Header.Set("Sec-Fetch-User", "?1")
	r.Header.Set("Sec-Fetch-Dest", "document")
	r.Header.Set("Accept-Encoding", "gzip")
	r.Header.Set("Accept-Language", "en-US,en;q=0.5")
	return r
}

// ffUA returns a randomly selected firefox user agent.
func (d *dorking) ffUA() string {
	var userAgents []string
	switch d.config.os {
	case "m":
		userAgents = []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:107.0) Gecko/20100101 Firefox/107.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:106.0) Gecko/20100101 Firefox/106.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:104.0) Gecko/20100101 Firefox/104.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:102.0) Gecko/20100101 Firefox/102.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:101.0) Gecko/20100101 Firefox/101.0",
		}
	default:
		userAgents = []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:99.0) Gecko/20100101 Firefox/99.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:102.0) Gecko/20100101 Firefox/102.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:106.0) Gecko/20100101 Firefox/106.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:104.0) Gecko/20100101 Firefox/104.0",
		}
	}
	random := rand.Intn(len(userAgents))
	return userAgents[random]
}

// chromeUA returns a randomly selected chrome user agent.
func (d *dorking) chromeUA() string {
	var userAgents []string
	switch d.config.os {
	case "m":
		userAgents = []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4692.56 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4889.0 Safari/537.36",
		}
	default:
		userAgents = []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36",
		}
	}
	random := rand.Intn(len(userAgents))
	return userAgents[random]
}
