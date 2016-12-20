package aws

import (
	"fmt"
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

func (c *RouteTableClient) DescribeRouteTables() {
	res, err := c.ec2Svc.DescribeRouteTables(nil)
	if err != nil || len(res.RouteTables) < 1 {
		//return nil, err
		fmt.Println(err)
		return
	}

	fmt.Println(res.RouteTables)
}
