package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/limoges/p2pnet/client"
)

func main() {

	var filename string
	var err error
	var c *client.Client

	flag.Parse()

	filename = flag.Arg(0)
	if _, err = os.Stat(filename); err != nil {
		fmt.Printf("Path '%v' does not exist.\n", filename)
		return
	}

	if c, err = client.New(filename); err != nil {
		fmt.Printf("Could not create client.\n")
		return
	}

	c.Run()
}
