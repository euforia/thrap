package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	_version   string
	_buildtime string
)

// Version returns the full version of the app. This is populated at
// build time
func Version() string {
	return _version + " " + _buildtime
}

var (
	showVersion = flag.Bool("version", false, "show version")
)

func init() {

	if len(_version) == 0 {
		_version = "unknown"
	}
	if len(_buildtime) == 0 {
		_buildtime = time.Now().UTC().String()
	}

}

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Println(Version())
		os.Exit(0)
	}

	//
	// TODO: Add code here
	//
}
