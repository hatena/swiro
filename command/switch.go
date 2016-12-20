package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/Songmu/prompter"
	"github.com/taku-k/swiro/aws"
	"github.com/urfave/cli"
	"os"
	"time"
)

func CmdSwitch(c *cli.Context) error {
	if c.NArg() < 2 {
		cli.ShowCommandHelp(c, "switch")
		return errors.New("Route table ID and VIP are required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	routeTableId := c.Args()[0]
	vip := c.Args()[1]

	if !prompter.YN("Switch the following VIP in the route table.\n  "+vip+"("+routeTableId+")\nAre you sure?", true) {
		fmt.Fprintln(os.Stderr, "Switching is canceled")
		return nil
	}

	var instanceId string
	if instanceId = c.String("instance-id"); instanceId == "" {
		var err error
		if instanceId, err = aws.NewMetaDataClient().GetInstanceID(); err != nil {
			return err
		}
	}

	routeTableCli := aws.NewRouteTableClient()
	routeTableCli.DescribeRouteTables()
	// vip の存在の確認, バリデーション
	// instance_id の存在確認
	// route table の
	return nil
}
