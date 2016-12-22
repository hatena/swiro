package aws

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"log"
	"os"
	"strings"
)

type RouteTable struct {
	table *ec2.RouteTable
	cli   *RouteTableClient
}

type RouteTableClient struct {
	ec2Svc ec2iface.EC2API
}

func NewRouteTables(ctx context.Context) ([]*RouteTable, error) {
	cli := newRouteTableClient()
	ec2Tables, err := cli.describeRouteTables(ctx)
	if err != nil {
		return nil, err
	}
	tables := make([]*RouteTable, 0)
	for _, t := range ec2Tables {
		tables = append(tables, &RouteTable{table: t, cli: cli})
	}
	return tables, nil
}

func NewRouteTable(ctx context.Context, routeTableId string) (*RouteTable, error) {
	cli := newRouteTableClient()
	ec2Table, err := cli.describeRouteTableById(ctx, routeTableId)
	if err != nil {
		return nil, err
	}
	return &RouteTable{table: ec2Table, cli: cli}, nil
}

func (t *RouteTable) ReplaceRoute(vip, instance string) error {

	return nil
}

func (t *RouteTable) String() string {
	return t.table.String()
}

func newRouteTableClient() *RouteTableClient {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal("Creating session is failed")
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region, _ = NewMetaDataClientFromSession(sess).GetRegion()
	}
	ec2Svc := ec2.New(sess, aws.NewConfig().WithRegion("ap-northeast-1"))

	return &RouteTableClient{ec2Svc: ec2Svc}
}

func (c *RouteTableClient) describeRouteTables(ctx context.Context) ([]*ec2.RouteTable, error) {
	req, resp := c.ec2Svc.DescribeRouteTablesRequest(nil)
	req.HTTPRequest = req.HTTPRequest.WithContext(ctx)
	if err := req.Send(); err != nil || len(resp.RouteTables) < 1 {
		return nil, err
	}

	return resp.RouteTables, nil
}

func (c *RouteTableClient) describeRouteTableById(ctx context.Context, key string) (*ec2.RouteTable, error) {
	var input *ec2.DescribeRouteTablesInput
	if strings.HasPrefix(key, "rtb-") {
		input = &ec2.DescribeRouteTablesInput{
			RouteTableIds: []*string{
				aws.String(key),
			},
		}
	} else {
		input = &ec2.DescribeRouteTablesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("tag-key"),
					Values: []*string{
						aws.String("Name"),
					},
				},
				{
					Name: aws.String("tag-value"),
					Values: []*string{
						aws.String(key),
					},
				},
			},
		}
	}
	
	req, resp := c.ec2Svc.DescribeRouteTablesRequest(input)
	req.HTTPRequest = req.HTTPRequest.WithContext(ctx)
	if err := req.Send(); err != nil || len(resp.RouteTables) < 1 {
		return nil, err
	}

	return resp.RouteTables[0], nil
}
