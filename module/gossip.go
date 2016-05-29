package main

import (
	"flag"
	"fmt"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/gossip"
)

var filename string

func init() {

	const (
		defaultFilename = "default.ini"
		usageFilename   = "the path to the file containing the configurations"
	)
	flag.StringVar(&filename, "f", defaultFilename, usageFilename)

}

func main() {

	flag.Parse()

	config, err := cfg.New(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	module, err := gossip.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	module.Run()
}
