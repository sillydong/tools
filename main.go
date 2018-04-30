package main

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

var (
	appName string
	wd      string
)

var commands []cli.Command

func init() {
	wd, _ = os.Getwd()
	appName = filepath.Base(wd)
	os.Chdir(wd)
}

func main() {
	app := cli.NewApp()
	app.Name = "Tool"
	app.Usage = "tools for daily use"
	app.Commands = commands
	app.Run(os.Args)
}
