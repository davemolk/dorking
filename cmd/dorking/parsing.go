package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type selectors struct {
	blurbSelector string
	itemSelector  string
	linkSelector  string
	name          string
	url string
}

func (d *dorking) getSelectors() []selectors {
	var s []selectors

	bing := selectors{
		blurbSelector: "div.b_caption p",
		itemSelector:  "li.b_algo",
		linkSelector:  "h2 a",
		name:          "bing",
	}
	s = append(s, bing)

	ddg := selectors{
		blurbSelector: "div.links_main > a",
		itemSelector:  "div.web-result",
		linkSelector:  "div.links_main > a",
		name:          "duck",
	}
	s = append(s, ddg)

	return s
}

func (d *dorking) parseData(ctx context.Context, s selectors) {
	b, err := d.makeRequest(ctx, s.url)
	if err != nil {
		fmt.Println("here", err)
		return
	}
	defer b.Close()
	doc, err := goquery.NewDocumentFromReader(b)
	if err != nil {
		fmt.Println("no here", err)
		return
	}
	doc.Find(s.itemSelector).Each(func(_ int, g *goquery.Selection) {
		if link, ok := g.Find(s.linkSelector).Attr("href"); !ok {
			fmt.Printf("no link found: %s\n", s.url)
			return
		} else {
			cleanLink := d.cleanLinks(link)
			fmt.Println(cleanLink)
			blurb := g.Find(s.blurbSelector).Text()
			if blurb == "" {
				fmt.Printf("can't get blurb for %s\n", s.name)
			}
			cleanBlurb := d.cleanBlurb(blurb)
			fmt.Println(cleanBlurb)
			fmt.Println()
		}
	})
}

func (d *dorking) cleanBlurb(s string) string {
	cleanB := d.noBlank.ReplaceAllString(s, " ")
	cleanB = strings.TrimSpace(cleanB)
	cleanB = strings.ReplaceAll(cleanB, "\n", "")
	return cleanB
}

func (d *dorking) cleanLinks(s string) string {
	u, err := url.QueryUnescape(s)
	if err != nil {
		fmt.Println(err)
		return s
	}
	if strings.HasPrefix(u, "//duck") {
		removePrefix := strings.Split(u, "=")
		u = removePrefix[1]
		removeSuffix := strings.Split(u, "&rut")
		u = removeSuffix[0]
	}
	return u
}
