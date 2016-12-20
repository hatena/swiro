package command

import (
	"github.com/urfave/cli"
	"errors"
	"github.com/Songmu/prompter"
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/taku-k/swiro/aws"
)

func CmdSwitch(c *cli.Context) error {
	// ルートテーブル, vip, instance id は metadata から取る
	if c.NArg() < 2 {
		cli.ShowCommandHelp(c, "switch")
		return errors.New("Route table ID and VIP are required")
	}

	routeTableId := c.Args()[0]
	vip := c.Args()[1]

	if !prompter.YN("Switch the following VIP in the route table.\n  " + vip + "(" + routeTableId + ")\nAre you sure?", true) {
		logging.INFO("Switching is canceled")
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

	// 
	return nil
}
