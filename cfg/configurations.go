package cfg

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/vaughan0/go-ini"
)

const (
	// hostKeyToken holds the configuration key by which one can retrieve
	// the path of the client's public key in PEM-format.
	hostKeyToken = "HOSTKEY"
)

// Configurations holds the configurations file which contains all the
// key-values used to specify the runtime behaviour of the application.
type Configurations struct {
	file ini.File
}

var (
	ErrConfigurationFileNotFound = errors.New("configurations could not be found")
)

func New(path string) (conf *Configurations, err error) {

	conf = &Configurations{}

	fmt.Printf("Attempting to load configuration from '%v'\n", path)
	// Load the configuration, disregarding case.
	file, err := ini.LoadFile(path)

	if err != nil {
		fmt.Println(err)
		return nil, ErrConfigurationFileNotFound
	}

	conf.file = file
	fmt.Printf("Loaded configuration from: %v\n", path)
	return conf, nil
}

func (c Configurations) Get(data interface{}, section string, key string) bool {

	value, ok := c.file.Get(section, key)
	if !ok {
		fmt.Printf("Could not find section.key = %v.%v\n", section, key)
		return ok
	}

	switch data := data.(type) {
	case *string:
		*data = value
	case *int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println(err)
			ok = false
		}
		*data = intValue
	default:
		panic(fmt.Sprintf("handling for %v not implemented", reflect.TypeOf(data)))
	}
	return ok
}

func (c Configurations) Init(data interface{}, section string, key string,
	defaultValue interface{}) {

	if ok := c.Get(data, section, key); ok {
		return
	}

	switch data := data.(type) {
	case *int:
		*data = defaultValue.(int)
	case *string:
		*data = defaultValue.(string)
	default:
		panic(fmt.Sprintf("handling for %v not implemented", data))
	}
}
