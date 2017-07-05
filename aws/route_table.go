package aws

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
)

const timeOut = 3 * time.Second

type RouteTable struct {
	table *ec2.RouteTable
	cli   *Ec2Client
}

type Ec2Meta struct {
	Name string
	Id   string
}

func NewRouteTables() ([]*RouteTable, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	cli := newEc2Client()
	ec2Tables, err := cli.getRouteTables(ctx)
	if err != nil {
		return nil, err
	}
	tables := make([]*RouteTable, 0, len(ec2Tables))
	for _, t := range ec2Tables {
		tables = append(tables, &RouteTable{table: t, cli: cli})
	}
	return tables, nil
}

func NewRouteTable(routeTableKey string) (*RouteTable, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	cli := newEc2Client()
	ec2Table, err := cli.getRouteTableByKey(ctx, routeTableKey)
	if err != nil {
		return nil, err
	}
	return &RouteTable{table: ec2Table, cli: cli}, nil
}

func (t *RouteTable) ReplaceRoute(vip, instance string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	routeTableId := *t.table.RouteTableId
	destinationCidrBlock := fmt.Sprintf("%s/32", vip)
	instanceId, err := t.cli.getInstanceId(ctx, instance)
	if err != nil {
		return err
	}
	if err = t.cli.replaceRoute(ctx, routeTableId, destinationCidrBlock, instanceId); err != nil {
		return err
	}

	changed, err := t.cli.getInstanceIdByDest(ctx, routeTableId, destinationCidrBlock)
	if err != nil {
		return err
	}
	if changed != instanceId {
		return errors.New("Route has not been replaced yet")
	}

	return nil
}

func (t *RouteTable) ListPossibleVips() *MaybeVips {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	ids := make([]string, 0, len(t.table.Routes))
	vips := make([]string, 0, len(t.table.Routes))
	names := make([]string, 0, len(t.table.Routes))
	for _, r := range t.table.Routes {
		if strings.HasSuffix(*r.DestinationCidrBlock, "/32") {
			if r.InstanceId != nil {
				ids = append(ids, *r.InstanceId)
				vips = append(vips, *r.DestinationCidrBlock)
				name, err := t.cli.getInstanceNameById(ctx, *r.InstanceId)
				if err != nil {
					name = "unknown"
				}
				names = append(names, name)
			}
		}
	}
	return &MaybeVips{t, ids, vips, names}
}

func (t *RouteTable) GetSrcByVip(vip string) (*Ec2Meta, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	id, state, err := t.getSrcByVip(vip)
	if err != nil {
		return nil, err
	}

	var name string
	switch {
	case strings.HasPrefix(id, "i-") && state == "active":
		name, err = t.cli.getInstanceNameById(ctx, id)
		if err != nil {
			return nil, err
		}
	case strings.HasPrefix(id, "eni-") && state == "active":
		name, err = t.cli.getENINameById(ctx, id)
		if err != nil {
			return nil, err
		}
	case state == "blackhole":
		name = "blackhole"
	default:
		return nil, errors.New("Not support to switch from neither instance nor ENI destination")
	}
	return &Ec2Meta{Name: name, Id: id}, nil
}

func (t *RouteTable) getSrcByVip(vip string) (string, string, error) {
	vipCidrBlock := vip
	if !strings.HasSuffix(vipCidrBlock, "/32") {
		vipCidrBlock = fmt.Sprintf("%s/32", vipCidrBlock)
	}
	for _, route := range t.table.Routes {
		if *route.DestinationCidrBlock == vipCidrBlock {
			switch {
			case route.InstanceId != nil && *route.InstanceId != "":
				return *route.InstanceId, *route.State, nil
			case route.NetworkInterfaceId != nil && *route.NetworkInterfaceId != "":
				return *route.NetworkInterfaceId, *route.State, nil
			}
		}
	}
	return "", "", errors.New("Given vip is not found")
}

func (t *RouteTable) GetRouteTableId() string {
	return *t.table.RouteTableId
}

func (t *RouteTable) GetRouteTableName() string {
	ret := "-"
	if len(t.table.Tags) != 0 {
		for _, tag := range t.table.Tags {
			if *tag.Key == "Name" {
				ret = *tag.Value
			}
		}
	}
	return ret
}

func (t *RouteTable) String() string {
	return t.table.String()
}
