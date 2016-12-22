package command

import (
	"errors"
	"fmt"
	"github.com/Songmu/prompter"
	"github.com/taku-k/swiro/aws"
	"github.com/urfave/cli"
	"os"
	"strings"
)

func CmdSwitch(c *cli.Context) error {
	if c.NArg() < 2 {
		cli.ShowCommandHelp(c, "switch")
		return errors.New("Route table ID or Name and VIP are required")
	}

	key := c.Args()[0]
	vip := c.Args()[1]

	var instanceKey string
	if instanceKey = c.String("instance-id"); instanceKey == "" {
		var err error
		if instanceKey, err = aws.NewMetaDataClient().GetInstanceID(); err != nil {
			return err
		}
	}

	routeTable, err := aws.NewRouteTable(key)
	if err != nil {
		return err
	}

	promptStr := `Switch the route below setting:
============================================
Route Table: %s (%s)
Virtual IP:  %s -------- Src:  %s (%s)
             %s \\
             %s  ======> Dest: %s
============================================
Are you sure?`
	routeTableName := routeTable.GetRouteTableName()
	routeTableId := routeTable.GetRouteTableId()
	srcInstance, err := routeTable.GetSrcInstanceByVip(vip)
	if err != nil {
		return err
	}
	ws := strings.Repeat(" ", len(vip))
	if !prompter.YN(fmt.Sprintf(promptStr, routeTableName, routeTableId, vip, srcInstance.Name, srcInstance.Id, ws, ws, instanceKey), true) {
		fmt.Fprintln(os.Stderr, "Switching is canceled")
		return nil
	}

	err = routeTable.ReplaceRoute(vip, instanceKey)
	if err != nil {
		return err
	}

	return nil
}
