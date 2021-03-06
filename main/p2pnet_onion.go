package main

import (
	"flag"
	"fmt"
	"github.com/limoges/p2pnet"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/onion"
)

var (
	filename string
)

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

	module, err := onion.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := p2pnet.Run(module); err != nil {
		fmt.Println(err)
		return
	}
}
