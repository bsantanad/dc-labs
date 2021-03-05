/* Netcatl is a read-only TCP client */
package main

import (
	"os"
	"strings"
	"fmt"
)

/*func main(dst io.Writer, src io.Reader) {


}*/

func main() {

	args := os.Args[1:]
	hosts := make([]string, len(args))

	for _, arg := range args {
		indexEq := strings.Index(arg, "=") + 1
		fmt.Printf("%s\n", arg[indexEq:])
		hosts = append(hosts, arg[indexEq:])
	}
	//conn, err := net.Dial("tcp", "")


}
