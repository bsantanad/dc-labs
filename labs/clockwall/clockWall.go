/* Netcatl is a read-only TCP client */
package main

import (
	"os"
	"strings"
	"io"
	"net"
	"log"
)

func getTime(conn io.Reader, stdout io.Writer, host string) {
	log.Print(host)
	log.Print(conn)
	var response io.Reader
	//response = conn
	log.Print(response)
	_, err := io.Copy(stdout, conn)
	if err != nil {
		log.Fatal(err)
	}
}


func main() {

	args := os.Args[1:]
	hosts := make([]string, len(args))
	//time := make(chan io.Writer)
	//time := make(chan io.Reader)

	for _, arg := range args {
		indexEq := strings.Index(arg, "=") + 1
		hosts = append(hosts, arg[indexEq:])
	}
	log.Print(hosts[1])
	conn, err := net.Dial("tcp", hosts[1])
	if err != nil {
		log.Fatal(err)
	}
	log.Print("im here 1")
	log.Print(conn)
	defer conn.Close()
	getTime(conn, os.Stdout, hosts[1])
}
