package command

import (
	"github.com/urfave/cli"
	"fmt"
	"os"
)

func CmdList(c *cli.Context) error {
	fmt.Fprintln(os.Stderr, "Not implemented")
	return nil
}
