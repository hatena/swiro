package main

import (
	"fmt"
	"os"

	"github.com/taku-k/swiro/command"
	"github.com/urfave/cli"
)

var GlobalFlags = []cli.Flag{}

var Commands = []cli.Command{
	{
		Name:   "switch",
		Usage:  "Switch the route based on given arguments",
		Action: command.CmdSwitch,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "r, route-table", Usage: "route table id or name"},
			cli.StringFlag{Name: "v, vip", Usage: "Virtual IP address"},
			cli.StringFlag{Name: "I, instance", Usage: "instance id or name"},
			cli.BoolFlag{Name: "f, force", Usage: "force switching (default: false"},
		},
	},
	{
		Name:   "list",
		Usage:  "List all route tables associated with Virtual IP",
		Action: command.CmdList,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
