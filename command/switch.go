package command

import (
	"fmt"
	"github.com/Songmu/prompter"
	"github.com/taku-k/swiro/aws"
	"github.com/urfave/cli"
	"os"
	"strings"
)

func CmdSwitch(c *cli.Context) error {
	key := c.String("route-table")
	vip := c.String("vip")

	var instanceKey string
	if instanceKey = c.String("instance"); instanceKey == "" {
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
	fmt.Fprintln(os.Stdout, "Success!!")

	return nil
}
