package cfg

import (
	"fmt"
	"github.com/go-ini/ini"
)

// Configurations holds the configurations file which contains all the key-values used to
// specify the runtime behaviour of the application.
type Configurations struct {
	File        *ini.File
	HostKeyPath string
}

func (c *Configurations) Load(path string) error {

	const (
		hostKeyToken = "HOSTKEY"
	)
	// Open the file at the given path.
	fmt.Printf("Parsing configuration found at '%v'\n", path)

	// We load the configuration with disregard for the lower/upper-caseyness.
	cfg, err := ini.InsensitiveLoad(path)
	// c.HostKeyPath, _ = cfg.GetKey(hostKeyToken).String()
	c.File = cfg

	return err
}
