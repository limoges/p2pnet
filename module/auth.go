package main

import (
	"flag"
	"fmt"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/client"
)

const (
	// DefaultConfigurationsPath is the default search path for the program's
	// configrations.
	DefaultConfigurationsPath = "./default.ini"
	usageConfiguration        = "path of the ini-formatted configuration file"
)

var (
	// configurationsPath allows the use of the specified .ini configurations
	// file found at the given path through -c [path].
	configurationsPath = flag.String("c", DefaultConfigurationsPath,
		usageConfiguration)
)

func main() {

	// Must be ran to allow the usage of the flag package functions (init, ...)
	flag.Parse()

	conf, err := cfg.New(*configurationsPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	auth, err := auth.New(conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	auth.Run()
}
