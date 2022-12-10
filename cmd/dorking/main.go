package main

import (
	"flag"
	"log"
	"regexp"
	"sync"
)

type config struct {
	contains string
	exact    bool
	ext      string
	feed     string
	filetype string
	hasfeed  string
	inbody   string
	info     string
	intitle  string
	inurl    string
	ip       string
	json     bool
	not      string
	notsite  string
	or       string
	query    string
	site     string
	timeout  int
	verbose  bool
}

type dorking struct {
	config   config
	noBlank  *regexp.Regexp
	searches *searchMap
}

func main() {
	var config config
	flag.StringVar(&config.contains, "contains", "", "return sites with links to specified file types")
	flag.BoolVar(&config.exact, "exact", false, "match exact words")
	flag.StringVar(&config.ext, "ext", "", "return sites with specified file name extension")
	flag.StringVar(&config.feed, "feed", "", "return RSS or Atom feeds for search term(s)")
	flag.StringVar(&config.filetype, "filetype", "", "file type")
	flag.StringVar(&config.hasfeed, "hasfeed", "", "return sites with RSS or Atom feeds for search term(s)")
	flag.StringVar(&config.inbody, "inbody", "", "return sites with search term(s) in body")
	flag.StringVar(&config.info, "info", "", "return information that Bing has about a site")
	flag.StringVar(&config.intitle, "intitle", "", "return sites with search term(s) in site title")
	flag.StringVar(&config.inurl, "inurl", "", "return sites with search term(s) in site URL")
	flag.StringVar(&config.ip, "ip", "", "return sites hosted by specific ip")
	flag.BoolVar(&config.json, "json", false, "write results to json")
	flag.StringVar(&config.notsite, "notsite", "", "site/domain to exclude")
	flag.StringVar(&config.not, "not", "", "term(s) to exclude")
	flag.StringVar(&config.or, "or", "", "OR term(s)")
	flag.StringVar(&config.query, "q", "", "search query")
	flag.StringVar(&config.site, "site", "", "site/domain to search")
	flag.IntVar(&config.timeout, "t", 5000, "timeout for request")
	flag.BoolVar(&config.verbose, "v", false, "chatty mode")
	flag.Parse()

	noBlank := regexp.MustCompile(`\s{2,}`)
	searches := newSearchMap()

	d := &dorking{
		config:   config,
		noBlank:  noBlank,
		searches: searches,
	}

	urls := d.makeQueryStrings()
	selectors := d.getSelectors()
	if len(urls) != len(selectors) {
		log.Fatal("mismatch between query urls and query data")
	}
	for i, u := range urls {
		selectors[i].url = u
	}
	
	var wg sync.WaitGroup
	wg.Add(len(selectors))
	for _, s := range selectors {
		go d.parseData(&wg, s)
	}
	wg.Wait()

	if config.json {
		if err := d.write(); err != nil {
			log.Fatalf("unable to write file: %v", err)
		}
	}
}
