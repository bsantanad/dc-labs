/* Netcatl is a read-only TCP client */
package main

import (
	"os"
	"strings"
	"io"
	"net"
	"log"
)

func printTime(conn io.Reader, stdout io.Writer) {
	_, err := io.Copy(stdout, conn)
	if err != nil {
		log.Fatal(err)
	}
}


func main() {

	args := os.Args[1:]
	hosts := make([]string, 0)

	for _, arg := range args {
		indexEq := strings.Index(arg, "=") + 1
		hosts = append(hosts, arg[indexEq:])
	}

	response := make(chan io.Reader, 10)
	defer close(response)

	for _, host := range hosts {
		go getTime(response, host)
		time := <-response
		go printTime(time, os.Stdout)
	}
	//getTime(<-response, os.Stdout, hosts[1])
}

func getTime(c chan io.Reader,/* conn net.Conn,*/ host string){
	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	c <- conn
}
