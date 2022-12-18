# dorking

advanced searching with bing, brave, duck duck go, and yahoo

## Overview
Google dorking is great, but Google's recaptchas etc are less great. Use dorking to search bing, brave, duck duck go, and yahoo instead, all off a single advanced search query. By default, dorking prints the results (url and blurb) to stdout, but you can change to json with a flag and/or save the results to a file. I'm hoping to add more search options when I have time. Great!

## Examples
```
dorking -filetype pdf -q goroutines
(json results truncated)
{
    "http://mesl.ucsd.edu/pubs/zhou_SIGBED16.pdf": "Go language’s concurrency is enabled through goroutines and invoked with keyword go. A user creates a goroutine and associates it with a program using go func(arg). After creation, go runtime scheduler automatically allocates gor-outines to run on OS threads.",
    "http://www.cs.uky.edu/~raphael/grad/keepingCurrent/roberts-concurrency.pdf": "Goroutines Goroutines enable concurrency Like threads, but lighter Spawn one by prefixing a function call with the go keyword Similar to backgrounding a process in Linux Scheduled onto OS threads by the Go runtime Goroutines share an address space, but sharing data structures is not idiomatic.",
    "http://www1.cs.columbia.edu/~aho/cs6998/reports/12-12-11_DeshpandeSponslerWeiss_GO.pdf": "scheduling goroutines onto threads e ectively is crucial for the e cient performance of Go programs. The idea behind goroutines is that they are capable of running concurrently, like threads, but are also extremely lightweight in compar-ison. So, while there might be multiple threads created ...",
    "https://arxiv.org/pdf/2204.00764.pdf": "Goroutines are considered “lightweight”, and the Go runtime context switches them on the operating-system (OS) threads. Go programmers use goroutines liberally both for symmetric and asymmetric tasks. Two or more goroutines can commu-nicate data via message passing (channels [13]) or shared ..."
}
```

## Install
First, you'll need to [install go](https://golang.org/doc/install). Then, run the following command:

```
go install github.com/davemolk/dorking/cmd/dorking@latest
```

# Additional Notes
* Some of the operators work better than others. 

* Brave doesn't publish advanced query info (at least that I found), so what's there is from me poking around.

* Bing's reported operators can be unreliable (ext, hasfeed, ip, and info don't seem to work, so I've excluded them).

* While Yahoo has a special query system (v*_vt, for instance), just using p seems to work, so I stuck with that.


## Flags
```
Usage of dorking:
  -contains string
    	return sites with links to specified file types
  -exact bool
    	match exact words
  -feed string
    	return RSS or Atom feeds for search term(s)
  -filetype string
    	filetypes
  -inbody string
    	return sites with search term(s) in body
  -intitle string
    	return sites with search term(s) in title
  -inurl string
    	return sites with search term(s) in url
  -j bool
    	print json results to stdout
  -nostite string
    	site/domain to exclude
  -or string
    	OR term(s)
  -q string
    	search query
  -site string
    	site/domain to search
  -t int
    	timeout for request (in ms)
  -v bool
    	verbose mode
  -w bool
    	write results to file
```
