package aws

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"log"
	"os"
)

type RouteTableClient struct {
	ec2Svc ec2iface.EC2API
}

func NewRouteTableClient() *RouteTableClient {
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

func (c *RouteTableClient) DescribeRouteTables(ctx context.Context) ([]*ec2.RouteTable, error) {
	req, resp := c.ec2Svc.DescribeRouteTablesRequest(nil)
	req.HTTPRequest = req.HTTPRequest.WithContext(ctx)
	if err := req.Send(); err != nil || len(resp.RouteTables) < 1 {
		return nil, err
	}

	return resp.RouteTables, nil
}

func (c *RouteTableClient) DescribeRouteTableById(ctx context.Context, routeTableId string) (*ec2.RouteTable, error) {
	req, resp := c.ec2Svc.DescribeRouteTablesRequest(&ec2.DescribeRouteTablesInput{
		RouteTableIds: []*string{
			aws.String(routeTableId),
		},
	})
	req.HTTPRequest = req.HTTPRequest.WithContext(ctx)
	if err := req.Send(); err != nil || len(resp.RouteTables) < 1 {
		return nil, err
	}

	return resp.RouteTables[0], nil
}
