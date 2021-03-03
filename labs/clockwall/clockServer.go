// Clock Server is a concurrent TCP server that periodically writes the time.
package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"time"
	"fmt"
)

type Input struct {
	port     int
	timezone string
}

func handleConn(c net.Conn, tz string) {
	defer c.Close()
	loc, _ := time.LoadLocation(tz)
	for {
		/*print time in loc*/
		_, err := io.WriteString(c,
				time.Now().In(loc).Format("15:04:05\n"))

		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {

	var input Input
	var host string

	input = manageInput()
	host = "localhost:" + fmt.Sprintf("%d", input.port)

	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn, input.timezone) // handle connections concurrently
	}
}

func manageInput() Input {

	var input Input
	var tmpPort *int

	input.timezone = os.Getenv("TZ")

	tmpPort = flag.Int("port", 9000, "port number.")
	flag.Parse()
	input.port = *tmpPort

	return input

}
