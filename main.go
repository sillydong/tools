package main

import (
	"os"
	"github.com/urfave/cli"
)

var (
	wd      string
)

var commands []cli.Command

func init() {
	wd, _ = os.Getwd()
	os.Chdir(wd)
}

func main() {
	app := cli.NewApp()
	app.Name = "Tool"
	app.Usage = "tools for daily use"
	app.Commands = commands
	app.Run(os.Args)
}
