package main

import (
	"fmt"
	"os"

	"github.com/hatena/swiro/build"
	"github.com/urfave/cli"
)

func main() {
	info := build.GetInfo()

	app := cli.NewApp()
	app.Name = info.Name
	app.Version = info.Version
	app.Author = "taku-k"
	app.Email = "taakuu19@gmail.com"
	app.Usage = ""

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
