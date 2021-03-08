/* Netcatl is a read-only TCP client */
package main

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func printTime(conn io.Reader, stdout io.Writer) {
	_, err := io.Copy(stdout, conn)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	/* manage input */
	args := os.Args[1:]
	hosts := make([]string, 0)
	for _, arg := range args {
		indexEq := strings.Index(arg, "=") + 1
		hosts = append(hosts, arg[indexEq:])
	}

	/* unbuffered chan */
	response := make(chan io.Reader)
	defer close(response)

	/* make request */
	for {
		for _, host := range hosts {
			go getTime(response, host)
			time := <-response
			go printTime(time, os.Stdout)
		}
	}
}

func getTime(time chan io.Reader, host string) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	time <- conn
}
