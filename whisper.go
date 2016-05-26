package main

import (
	"flag"
	"fmt"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/gossip"
)

const (
	// DefaultConfigurationsPath is the default search path for the program's configrations.
	DefaultConfigurationsPath = "./default.ini"
)

var (
	// configurationsPath allows the use of the specified .ini configurations file
	// found at the given path through -c [path].
	configurationsPath string
)

func init() {

	const (
		usageConfiguration = "path of the ini-formatted configuration file"
	)
	flag.StringVar(&configurationsPath, "c", DefaultConfigurationsPath, usageConfiguration)
}

func main() {

	// Must be ran to allow the usage of the flag package functions (init, ...)
	flag.Parse()

	c := cfg.Configurations{}
	err := c.Load(configurationsPath)

	if err != nil {
		// We'll just use the default configuration
		fmt.Printf("Problem parsing configuration. Using default.\n")
	}

	g := gossip.Gossip{}
	g.LoadConfiguration(c)
	g.Run()
}
