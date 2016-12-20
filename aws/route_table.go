package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"log"
	"os"
	"context"
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

func (c *RouteTableClient) DescribeRouteTables(ctx context.Context) {
	req, resp := c.ec2Svc.DescribeRouteTablesRequest(nil)
	req.HTTPRequest = req.HTTPRequest.WithContext(ctx)
	if err := req.Send(); err != nil || len(resp.RouteTables) < 1 {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Println(resp)
}
