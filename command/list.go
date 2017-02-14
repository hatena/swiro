package command

import (
	"github.com/urfave/cli"
	"github.com/taku-k/swiro/aws"
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
