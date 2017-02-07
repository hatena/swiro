package aws

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"strings"
	"time"
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

	// TODO: check whether the route has actually replaced

	return nil
}

func (t *RouteTable) ListPossibleVips() []string {
	//for _, r := range t.table.Routes {
	//	if *r.DestinationCidrBlock ==
	//}
	return nil
}

func (t *RouteTable) GetSrcInstanceByVip(vip string) (*Ec2Meta, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	id, err := t.getSrcInstanceIdByVip(vip)
	if err != nil {
		return nil, err
	}
	name, err := t.cli.getInstanceNameById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Ec2Meta{Name: name, Id: id}, nil
}

func (t *RouteTable) getSrcInstanceIdByVip(vip string) (string, error) {
	vipCidrBlock := vip
	if !strings.HasSuffix(vipCidrBlock, "/32") {
		vipCidrBlock = fmt.Sprintf("%s/32", vipCidrBlock)
	}
	for _, route := range t.table.Routes {
		if *route.DestinationCidrBlock == vipCidrBlock && *route.InstanceId != "" {
			return *route.InstanceId, nil
		}
	}
	return "", errors.New("Given vip is not found")
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
