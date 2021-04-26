// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
//
// Crawl3 adds support for depth limiting.
//
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gopl.io/ch5/links"
)

type Worker struct {
	urlList []string
	depth   int
}

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

func crawl(url string, w *bufio.Writer) []string {
	// write into the file
	tokens <- struct{}{} // acquire a token
	_, err := fmt.Fprintf(w, "%v\n", url)
	if err != nil {
		fmt.Println("could not write into file :(")
	}
	list, err := links.Extract(url)
	<-tokens // release the token

	if err != nil {
		log.Print(err)
	}
	return list
}

//!-sema

//!+
func main() {

	// handle args
	args := os.Args[1:]
	if len(args) < 3 {
		fmt.Printf("not enough arguments, please do: \n" +
			"./web-crawler -depth=<n> -results=<filename> <url> \n")
		return
	}

	var depth int
	var filename string
	var url string

	// without last element because its the url
	url = args[len(args)-1]
	for _, arg := range args[:len(args)-1] {
		param := strings.Split(arg, "=")
		switch param[0] {
		case "-depth":
			num, err := strconv.Atoi(param[1])
			depth = num
			if err != nil {
				fmt.Printf("bad arguments, please use: \n" +
					"./web-crawler -depth=<n> -results=<filename> <url> \n")
				return
			}
		case "-results":
			filename = param[1]
		default:
			fmt.Printf("bad arguments, please use: \n" +
				"./web-crawler -depth=<n> -results=<filename> <url> \n")
			return
		}
	}

	// create file in which we will print
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("could not create resutls file :(")
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()

	// worklist
	var first Worker
	var tmp []string
	tmp = append(tmp, url)

	first.depth = 0
	first.urlList = tmp

	worklist := make(chan Worker)
	go func() { worklist <- first }()

	var n int // number of pending sends to worklist
	n++

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- {

		workers := <-worklist
		actualDepth := workers.depth

		for _, link := range workers.urlList {
			if !seen[link] {
				seen[link] = true
				if actualDepth >= depth {
					continue
				}
				n++
				go func(link string) {
					var tmp Worker
					tmp.depth = actualDepth + 1
					tmp.urlList = crawl(link, w)
					worklist <- tmp
				}(link)
			}
		}
	}
}

//!-
