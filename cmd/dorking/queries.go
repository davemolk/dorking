package main

import (
	"fmt"
	"strings"
)

// queryData is a struct containing the possible components
// to be used in constructing query strings.
type queryData struct {
	base     string
	contains string
	feed     string
	filetype string
	host string
	inbody   string
	intitle  string
	inurl    string
	name     string
	not      string
	notsite  string
	or       string
	site     string
	spacer   string
}

// getQueryData returns a slice of queryData, containing
// advanced query information for bing, brave, duck duck go,
// and yahoo search engines.
func (d *dorking) getQueryData() []queryData {
	qdSlice := make([]queryData, 0, 4)
	bing := queryData{
		base:     "https://bing.com/search?q=",
		contains: "contains%3A",
		feed:     "feed%3A",
		filetype: "filetype%3A",
		host: "www.bing.com",
		inbody:   "inbody%3A",
		intitle:  "intitle%3A",
		inurl:    "inanchor%3A",
		name:     "bing",
		not:      "-",
		notsite:  "-site%3A",
		or:       "OR",
		site:     "site%3A",
		spacer:   "+",
	}

	// doesn't publish query info, so this is
	// assembled from poking around...
	brave := queryData{
		base:    "https://search.brave.com/search?q=",
		feed:     "feed%3A",
		filetype: "filetype%3A",
		host: "search.brave.com",
		inbody:   "inbody%3A",
		intitle:  "intitle%3A",
		inurl:   "inurl%3A",
		name:    "brave",
		not:     "-",
		notsite:  "-site%3A",
		or:      "OR",
		site:    "site%3A",
		spacer:  "+",
	}

	// quack quack quack mr ducksworth
	ddg := queryData{
		base:     "https://html.duckduckgo.com/html?q=",
		feed:     "feed%3A",
		filetype: "filetype%3A",
		host: "duckduckgo.com",
		inbody:   "inbody%3A",
		intitle:  "intitle%3A",
		inurl:    "inurl%3A",
		name:     "duckduckgo",
		not:      "-",
		notsite:  "-site%3A",
		or:       "OR",
		site:     "site%3A",
		spacer:   "+",
	}

	// has a bunch of weird query params but it seems like
	// just specifying within p works...
	yahoo := queryData{
		base:     "https://search.yahoo.com/search?p=",
		feed:     "feed%3A",
		filetype: "filetype%3A",
		host: "search.yahoo.com",
		inbody:   "inbody%3A",
		intitle:  "intitle%3A",
		inurl:    "inanchor%3A",
		name:     "yahoo",
		not:      "-",
		notsite:  "-site%3A",
		or:       "OR",
		site:     "site%3A",
		spacer:   "+",
	}

	qdSlice = append(qdSlice, bing, brave, ddg, yahoo)

	return qdSlice
}

// makeQueryStrings returns a slice of query strings
// specified for each search engine.
func (d *dorking) makeQueryStrings() []string {
	qdSlice := d.getQueryData()
	urls := make([]string, 0, len(qdSlice))
	for _, qd := range qdSlice {
		var url string
		var components []string
		var cleanedQuery string
		switch {
		case d.config.exact:
			cleanedQuery = strings.Replace(d.config.query, " ", qd.spacer, -1)
			cleanedQuery = fmt.Sprintf("\"%s\"", cleanedQuery)
			components = append(components, cleanedQuery)
		case d.config.query != "":
			cleanedQuery = strings.Replace(d.config.query, " ", qd.spacer, -1)
			components = append(components, cleanedQuery)
		}
		if d.config.contains != "" && qd.contains != "" {
			contains := fmt.Sprintf("%s%s", qd.contains, d.config.contains)
			components = append(components, contains)
		}
		if d.config.feed != "" && qd.feed != "" {
			feed := fmt.Sprintf("%s%s", qd.feed, d.config.feed)
			components = append(components, feed)
		}
		if d.config.filetype != "" && qd.filetype != "" {
			filetype := fmt.Sprintf("%s%s", qd.filetype, d.config.filetype)
			components = append(components, filetype)
		}
		if d.config.inbody != "" && qd.inbody != "" {
			inbody := fmt.Sprintf("%s%s", qd.inbody, d.config.inbody)
			components = append(components, inbody)
		}
		if d.config.intitle != "" && qd.intitle != "" {
			intitle := fmt.Sprintf("%s%s", qd.intitle, d.config.intitle)
			components = append(components, intitle)
		}
		// ddg needs the url enclosed in double quotes to work properly, and the search
		// results from the other three search engines aren't adversely impacted by the quotes.
		if d.config.inurl != "" && qd.inurl != "" {
			inurl := fmt.Sprintf("%s%s%s%s", qd.inurl, `"`, d.config.inurl, `"`)
			components = append(components, inurl)
		}
		if d.config.not != "" && qd.not != "" {
			cleanedQuery = strings.Replace(d.config.not, " ", qd.spacer, -1)
			cleanedQuery = fmt.Sprintf("%s%s", qd.not, cleanedQuery)
			components = append(components, cleanedQuery)
		}
		if d.config.notsite != "" && qd.notsite != "" {
			notsite := fmt.Sprintf("%s%s", qd.notsite, d.config.notsite)
			components = append(components, notsite)
		}
		if d.config.or != "" && qd.or != "" {
			cleanedQuery = strings.Replace(d.config.or, " ", qd.spacer, -1)
			cleanedQuery = fmt.Sprintf("%s%s", qd.or, cleanedQuery)
			components = append(components, cleanedQuery)
		}
		if d.config.site != "" && qd.site != "" {
			site := fmt.Sprintf("%s%s", qd.site, d.config.site)
			components = append(components, site)
		}
		params := strings.Join(components, "+")
		url = fmt.Sprintf("%s%s", qd.base, params)
		urls = append(urls, url)
	}

	return urls
}

// selectorSlice adds the search query URL to each search engine
// selector and returns the results as a selector slice.
func (d *dorking) selectorSlice() []selectors {
	urls := d.makeQueryStrings()
	selectorSlice := d.getSelectors()
	if len(urls) != len(selectorSlice) {
		d.errorLog.Fatal("mismatch between query urls and query data")
	}
	for i, u := range urls {
		selectorSlice[i].url = u
	}
	return selectorSlice
}
