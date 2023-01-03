package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
)

type config struct {
	contains string
	exact    bool
	feed     string
	filetype string
	inbody   string
	intitle  string
	inurl    string
	json     bool
	not      string
	notsite  string
	or       string
	os string
	query    string
	site     string
	timeout  int
	verbose  bool
	write    bool
}

type dorking struct {
	config   config
	errorLog *log.Logger
	noBlank  *regexp.Regexp
	searches *searchMap
}

func main() {
	var config config
	flag.StringVar(&config.contains, "contains", "", "return sites with links to specified file types")
	flag.BoolVar(&config.exact, "exact", false, "match exact words")
	flag.StringVar(&config.feed, "feed", "", "return RSS or Atom feeds for search term(s)")
	flag.StringVar(&config.filetype, "filetype", "", "file type")
	flag.StringVar(&config.inbody, "inbody", "", "return sites with search term(s) in body")
	flag.StringVar(&config.intitle, "intitle", "", "return sites with search term(s) in site title")
	flag.StringVar(&config.inurl, "inurl", "", "return sites with search term(s) in site URL")
	flag.BoolVar(&config.json, "j", false, "print json results to stdout")
	flag.StringVar(&config.notsite, "notsite", "", "site/domain to exclude")
	flag.StringVar(&config.not, "not", "", "term(s) to exclude")
	flag.StringVar(&config.or, "or", "", "OR term(s)")
	flag.StringVar(&config.os, "os", "w", "operating system (w or m)")
	flag.StringVar(&config.query, "q", "", "search query")
	flag.StringVar(&config.site, "site", "", "site/domain to search")
	flag.IntVar(&config.timeout, "t", 5000, "timeout for request")
	flag.BoolVar(&config.verbose, "v", false, "chatty mode")
	flag.BoolVar(&config.write, "w", false, "write results to file")
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ltime|log.Lshortfile)
	noBlank := regexp.MustCompile(`\s{2,}`)
	searches := newSearchMap()

	d := &dorking{
		config:   config,
		errorLog: errorLog,
		noBlank:  noBlank,
		searches: searches,
	}

	selectorSlice := d.selectorSlice()

	var wg sync.WaitGroup
	wg.Add(len(selectorSlice))
	for _, s := range selectorSlice {
		go func(s selectors) {
			defer wg.Done()
			b, err := d.makeRequest(s.url)
			if err != nil {
				if d.config.verbose {
					errorLog.Printf("unable to make request for %s: %v\n", s.name, err)
				}
				return
			}
			d.parseData(b, s)
		}(s)
	}
	wg.Wait()

	if config.json || config.write {
		b, err := d.json(d.searches.searches)
		if err != nil {
			errorLog.Fatal("unable to get json:", err)
		}
		if config.json {
			fmt.Println(string(b))
		}
		if config.write {
			if err := os.WriteFile("results.json", b, 0644); err != nil {
				errorLog.Printf("writing error: %v\n", err)
			}
		}
	}
}
