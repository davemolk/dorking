package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// selectors is a struct containing the necessary information
// to parse the search engine results for each URL and blurb.
type selectors struct {
	blurbSelector string
	itemSelector  string
	linkSelector  string
	name          string
	url           string
}

// getSelectors returns a slice of selectors consisting of one
// selector struct for each of the search engines.
func (d *dorking) getSelectors() []selectors {
	s := make([]selectors, 0, 4)

	bing := selectors{
		blurbSelector: "div.b_caption p",
		itemSelector:  "li.b_algo",
		linkSelector:  "h2 a",
		name:          "bing",
	}

	brave := selectors{
		blurbSelector: "div.snippet-content p.snippet-description",
		itemSelector:  "div.fdb",
		linkSelector:  "div.fdb > a.result-header",
		name:          "brave",
	}

	ddg := selectors{
		blurbSelector: "div.links_main > a",
		itemSelector:  "div.web-result",
		linkSelector:  "div.links_main > a",
		name:          "duck",
	}

	yahoo := selectors{
		blurbSelector: "div.compText",
		itemSelector:  "div.algo",
		linkSelector:  "h3 > a",
		name:          "yahoo",
	}

	s = append(s, bing, brave, ddg, yahoo)

	return s
}

// parseData creates and parses a goquery document for the
// result URLs and blurbs. These are stored in the searchMap
// and printed to stdout (unless the -j flag is true).
func (d *dorking) parseData(b *bytes.Buffer, s selectors) {
	doc, err := goquery.NewDocumentFromReader(b)
	if err != nil {
		if d.config.verbose {
			d.errorLog.Printf("unable to generate goquery doc for %s: %v\n", s.name, err)
		}
		return
	}
	doc.Find(s.itemSelector).Each(func(_ int, g *goquery.Selection) {
		link, ok := g.Find(s.linkSelector).Attr("href")
		if !ok {
			return
		}
		cleanedLink := d.cleanLinks(link)
		blurb := g.Find(s.blurbSelector).Text()
		cleanedBlurb := d.cleanBlurb(blurb)
		if !d.config.json {
			d.printStdout(cleanedLink, cleanedBlurb)
		}
		d.searches.store(cleanedLink, cleanedBlurb)
	})
}

// cleanBlurb removes any extraneous whitespace and \n from a string.
func (d *dorking) cleanBlurb(s string) string {
	cleanB := d.noBlank.ReplaceAllString(s, " ")
	cleanB = strings.TrimSpace(cleanB)
	cleanB = strings.ReplaceAll(cleanB, "\n", "")
	return cleanB
}

// cleanLinks strips any unnecessary information added to the result
// links by duck duck go and yahoo.
func (d *dorking) cleanLinks(s string) string {
	u, err := url.QueryUnescape(s)
	if err != nil {
		if d.config.verbose {
			d.errorLog.Printf("unable to clean %s: %v\n", s, err)
		}
		return s
	}
	switch {
	case strings.HasPrefix(u, "//duck"):
		// ddg links will sometimes take the following format:
		// //duckduckgo.com/l/?uddg=actualURLHere/&rut=otherStuff
		noPre := strings.Split(u, "=")
		u = noPre[1]
		noSuf := strings.Split(u, "&rut")
		u = noSuf[0]
	case strings.HasPrefix(u, "https://r.search.yahoo.com/"):
		// gotta clean these too
		noPre := strings.Split(u, "/RU=")
		u = noPre[1]
		noSuf := strings.Split(u, "/RK=")
		u = noSuf[0]
	}
	return u
}

// printStdout truncates any blurb with a length longer
// than 200 and prints to stdout.
func (d *dorking) printStdout(cleanLink, cleanBlurb string) {
	fmt.Println(cleanLink)
	if len(cleanBlurb) > 200 {
		cleanBlurb = cleanBlurb[:200]
	}
	fmt.Println(cleanBlurb)
	fmt.Println()
}
