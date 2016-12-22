package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type RouteTable struct {
	table *ec2.RouteTable
	cli   *Ec2Client
}

func NewRouteTables(ctx context.Context) ([]*RouteTable, error) {
	cli := newEc2Client()
	ec2Tables, err := cli.getRouteTables(ctx)
	if err != nil {
		return nil, err
	}
	tables := make([]*RouteTable, 0)
	for _, t := range ec2Tables {
		tables = append(tables, &RouteTable{table: t, cli: cli})
	}
	return tables, nil
}

func NewRouteTable(ctx context.Context, routeTableKey string) (*RouteTable, error) {
	cli := newEc2Client()
	ec2Table, err := cli.getRouteTableByKey(ctx, routeTableKey)
	if err != nil {
		return nil, err
	}
	return &RouteTable{table: ec2Table, cli: cli}, nil
}

func (t *RouteTable) ReplaceRoute(ctx context.Context, vip, instance string) error {
	routeTableId := *t.table.RouteTableId
	destinationCidrBlock := fmt.Sprintf("%s/32", vip)
	instanceId, err := t.cli.getInstanceId(ctx, instance)
	if err != nil {
		return err
	}
	if err = t.cli.replaceRoute(ctx, routeTableId, destinationCidrBlock, instanceId); err != nil {
		return err
	}
	return nil
}

func (t *RouteTable) String() string {
	return t.table.String()
}
