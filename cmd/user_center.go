package main

import (
	"github.com/iodsp/user_center/commands"
	"github.com/urfave/cli"
	"os"
)

var version = "development"

func main() {
	app := cli.NewApp()
	app.Name = "iodsp/user_center"
	app.Usage = "User management system for iodsp"
	app.Version = version
	app.Copyright = "(c) 2019 The iodsp/user_center contributors <fionawp@126.com>"
	app.EnableBashCompletion = true
	app.Flags = commands.GlobalFlags

	app.Commands = []cli.Command{
		commands.ConfigCommand,
		commands.StartCommand,
	}

	app.Run(os.Args)
}
