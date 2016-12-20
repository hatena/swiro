package aws

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"github.com/ktakuya/aptly/_vendor/src/github.com/mitchellh/goamz/ec2"
	"github.com/aws/aws-sdk-go/aws"
)

type RouteTableClient struct {
	svc ec2iface.EC2API
}

func NewRouteTableClient() *RouteTableClient {
	session := session.New()
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region, _ = NewMetaDataClientFromSession(session).GetRegion()
	}
	svc := ec2.New(session, &aws.Config{Region: aws.String(region)})

	return &RouteTableClient{svc: svc}
}

