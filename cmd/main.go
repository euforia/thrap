package main

import (
	"fmt"
	"os"
	"time"

	"github.com/euforia/thrap/cli"
)

var (
	_version   string
	_buildtime string
)

func init() {
	if _version == "" {
		_version = "unknown"
	}
	if _buildtime == "" {
		_buildtime = time.Now().UTC().String()
	}
}

func version() string {
	return _version + " " + _buildtime
}

func main() {
	app := cli.NewCLI(version())
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
