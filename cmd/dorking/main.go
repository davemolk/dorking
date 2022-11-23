package main

import (
	"flag"
)

type config struct {
	contains   string
	ext        string
	feed       string
	filetype   string
	hasfeed    string
	inbody     string
	info       string
	intitle    string
	inurl      string
	ip         string
	query      string
	queryExact string
	nosite     string
	related    string
	site       string
	timeout    int
}

type dorking struct {
	config config
}


func main() {
	var config config
	flag.StringVar(&config.contains, "contains", "", "return sites with links to specified file types")
	flag.StringVar(&config.ext, "ext", "", "return sites with specified file name extension")
	flag.StringVar(&config.feed, "feed", "", "return RSS or Atom feeds for search term(s)")
	flag.StringVar(&config.filetype, "filetype", "", "file type")
	flag.StringVar(&config.hasfeed, "hasfeed", "", "return sites with RSS or Atom feeds for search term(s)")
	flag.StringVar(&config.inbody, "inbody", "", "return sites with search term(s) in body")
	flag.StringVar(&config.info, "info", "", "return information that Bing has about a site")
	flag.StringVar(&config.intitle, "intitle", "", "return sites with search term(s) in site title")
	flag.StringVar(&config.inurl, "inurl", "", "return sites with search term(s) in site URL")
	flag.StringVar(&config.ip, "ip", "", "return sites hosted by specific ip")
	flag.StringVar(&config.nosite, "nosite", "", "site/domain to exclude")
	flag.StringVar(&config.query, "q", "", "search query")
	flag.StringVar(&config.queryExact, "qe", "", "search query (exact matching)")
	flag.StringVar(&config.related, "related", "", "return sites similar to input site")
	flag.StringVar(&config.site, "site", "", "site/domain to search")
	flag.IntVar(&config.timeout, "t", 5000, "timeout for request")
	flag.Parse()

	d := &dorking{
		config: config,
	}
	
	qdSlice := d.makeQueryData()
	_ = qdSlice
}