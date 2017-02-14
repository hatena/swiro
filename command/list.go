package command

import (
	"github.com/taku-k/swiro/aws"
	"github.com/urfave/cli"
)

func CmdList(c *cli.Context) error {
	ts, err := aws.NewRouteTables()
	if err != nil {
		return err
	}
	for _, t := range ts {
		t.ListPossibleVips().Output()
	}
	return nil
}
