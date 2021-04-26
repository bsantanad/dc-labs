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
	"fmt"
	"log"
	"os"
    "strings"
    "strconv"
    "bufio"

	"gopl.io/ch5/links"
)

type Workelement struct {
    channel chan []string
    myDepth int
}

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

func crawl(url string, w *bufio.Writer) []string {
	_, err := fmt.Fprintf(w, "%v\n", url)
    if err != nil {
        fmt.Println("could not write into file :(")
    }
	tokens <- struct{}{} // acquire a token
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
    //worklist := make(chan []string, int)
	var n int // number of pending sends to worklist

	// Start with the command-line arguments.
	n++
	//go func() { worklist <- os.Args[1:] }()

    // handle args
	args := os.Args[1:]
	if len(args) < 3 {
		fmt.Printf("not enough arguments, please do: \n" +
			"./web-crawler -depth=<n> -results=<filename> <url> \n")
        return
	}

    // fill the worklist
    //var we Workelement
    //we.channel I:= 
    worklist := make(chan []string)
    var tmp []string
    tmp = append(tmp, args[2])
    go func() { worklist <- tmp }()

    var depth int
    var filename string
    for _, arg := range args[:2] {
        param := strings.Split(arg, "=")
        switch param[0] {
            case "-depth":
                num, err := strconv.Atoi(param[1])
                depth = num
                if err != nil{
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
    fmt.Printf("%v\n", depth)

    f, err := os.Create(filename)
    if err != nil {
        fmt.Println("could not create resutls file :(");
        return
    }
    defer f.Close()
    w := bufio.NewWriter(f)
    defer w.Flush()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
        fmt.Printf("%v\n", n)
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklist <- crawl(link, w)
				}(link)
			}
		}
	}
}

//!-
