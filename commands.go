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
		Usage:  "",
		Action: command.CmdSwitch,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "I, instance-id", Usage: "instance id"},
			cli.IntFlag{Name: "n, max-attempts", Value: 10, Usage: "the maximum number of attempts to poll replacing route (default: 10)"},
			cli.IntFlag{Name: "i, interval", Value: 2, Usage: "the interval in seconds to poll replacing route (default: 2)"},
		},
	},
	{
		Name:   "list",
		Usage:  "",
		Action: command.CmdList,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
