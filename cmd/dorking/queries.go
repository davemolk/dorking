package main

import (
	"fmt"
	"strings"
)

type queryData struct {
	base     string
	contains string
	ext      string
	feed     string
	filetype string
	hasfeed  string
	inbody   string
	info     string
	intitle  string
	inurl    string
	ip       string
	nosite   string
	not      string
	related  string
	site     string
	spacer   string
}

func (d *dorking) getQueryData() []queryData {
	qdSlice := make([]queryData, 0, 3)

	bing := queryData{
		base:     "https://bing.com/search?q=",
		contains: "contains%3A",
		ext:      "ext%3A",
		feed:     "feed%3A",
		filetype: "filetype%3A",
		hasfeed:  "hasfeed%3A",
		inbody:   "inbody%3A",
		info:     "info%3A",
		intitle:  "intitle%3A",
		inurl:    "inanchor%3A",
		ip:       "ip%3A",
		nosite:   "-site%3A",
		not:      "-",
		site:     "site%3A",
		spacer:   "+",
	}
	qdSlice = append(qdSlice, bing)

	// doesn't publish query info, so this is
	// assembled from poking around...
	brave := queryData{
		base:   "https://search.brave.com/search?q=",
		inurl:  "inurl%3A",
		nosite: "-site:%3A",
		not:    "-",
		site:   "site:%3A",
		spacer: "+",
	}
	qdSlice = append(qdSlice, brave)

	ddg := queryData{
		base:     "https://html.duckduckgo.com/html?q=",
		filetype: "filetype%3A",
		intitle:  "intitle%3A",
		inurl:    "inurl%3A",
		nosite:   "-site%3A",
		not:      "-",
		site:     "site%3A",
		spacer:   "+",
	}
	qdSlice = append(qdSlice, ddg)

	return qdSlice
}

func (d *dorking) makeQueryStrings() []string {
	qdSlice := d.getQueryData()
	var urls []string
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

		if d.config.ext != "" && qd.ext != "" {
			ext := fmt.Sprintf("%s%s", qd.ext, d.config.ext)
			components = append(components, ext)
		}

		if d.config.feed != "" && qd.feed != "" {
			feed := fmt.Sprintf("%s%s", qd.feed, d.config.feed)
			components = append(components, feed)
		}

		if d.config.filetype != "" && qd.filetype != "" {
			filetype := fmt.Sprintf("%s%s", qd.filetype, d.config.filetype)
			components = append(components, filetype)
		}

		if d.config.hasfeed != "" && qd.hasfeed != "" {
			hasfeed := fmt.Sprintf("%s%s", qd.hasfeed, d.config.hasfeed)
			components = append(components, hasfeed)
		}

		if d.config.inbody != "" && qd.inbody != "" {
			inbody := fmt.Sprintf("%s%s", qd.inbody, d.config.inbody)
			components = append(components, inbody)
		}

		if d.config.info != "" && qd.info != "" {
			info := fmt.Sprintf("%s%s", qd.info, d.config.info)
			components = append(components, info)
		}

		if d.config.intitle != "" && qd.intitle != "" {
			intitle := fmt.Sprintf("%s%s", qd.intitle, d.config.intitle)
			components = append(components, intitle)
		}

		if d.config.inurl != "" && qd.inurl != "" {
			inurl := fmt.Sprintf("%s%s", qd.inurl, d.config.inurl)
			components = append(components, inurl)
		}

		if d.config.ip != "" && qd.ip != "" {
			ip := fmt.Sprintf("%s%s", qd.ip, d.config.ip)
			components = append(components, ip)
		}

		if d.config.notsite != "" && qd.nosite != "" {
			nosite := fmt.Sprintf("%s%s", qd.nosite, d.config.notsite)
			components = append(components, nosite)
		}

		if d.config.not != "" && qd.not != "" {
			cleanedQuery = strings.Replace(d.config.not, " ", qd.spacer, -1)
			cleanedQuery = fmt.Sprintf("%s%s", "-", cleanedQuery)
			components = append(components, cleanedQuery)
		}

		if d.config.related != "" && qd.related != "" {
			related := fmt.Sprintf("%s%s", qd.related, d.config.related)
			components = append(components, related)
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
