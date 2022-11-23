package main

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
	related  string
	site     string
	spacer   string
}

func (d *dorking) makeQueryData() []queryData {
	var qdSlice []queryData

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
		site:     "site%3A",
		spacer:   "+",
	}
	qdSlice = append(qdSlice, bing)

	ddg := queryData{
		base:     "https://html.duckduckgo.com/html?q=",
		filetype: "filetype%3A",
		intitle:  "intitle%3A",
		inurl:    "inurl%3A",
		nosite:   "-site%3A",
		site:     "site%3A",
		spacer:   "+",
	}
	qdSlice = append(qdSlice, ddg)

	return qdSlice
}