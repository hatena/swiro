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
		return errors.New("Route table ID or Name and VIP are required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := c.Args()[0]
	vip := c.Args()[1]

	if !prompter.YN("Switch the following VIP in the route table.\n  "+vip+"("+key+")\nAre you sure?", true) {
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

	var routeTable *aws.RouteTable
	var err error
	if routeTable, err = aws.NewRouteTable(ctx, key); err != nil {
		return err
	}
	if err = routeTable.ReplaceRoute(ctx, vip, instanceId); err != nil {
		return err
	}
	// vip の存在の確認, バリデーション
	// instance_id の存在確認
	// route table の
	return nil
}
