package aws

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"log"
	"os"
	"strings"
)

type Ec2Client struct {
	ec2Svc ec2iface.EC2API
}

func newEc2Client() *Ec2Client {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal("Creating session is failed")
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region, _ = NewMetaDataClientFromSession(sess).GetRegion()
	}
	ec2Svc := ec2.New(sess, aws.NewConfig().WithRegion("ap-northeast-1"))

	return &Ec2Client{ec2Svc: ec2Svc}
}

func (c *Ec2Client) getRouteTables(ctx context.Context) ([]*ec2.RouteTable, error) {
	req, resp := c.ec2Svc.DescribeRouteTablesRequest(nil)
	req.HTTPRequest = req.HTTPRequest.WithContext(ctx)
	if err := req.Send(); err != nil || len(resp.RouteTables) < 1 {
		return nil, err
	}

	return resp.RouteTables, nil
}

func (c *Ec2Client) getRouteTableByKey(ctx context.Context, key string) (*ec2.RouteTable, error) {
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
	err := req.Send()
	switch {
	case err != nil:
		return nil, err
	case len(resp.RouteTables) == 0:
		return nil, errors.New("Route table is not found")
	case len(resp.RouteTables) > 1:
		return nil, errors.New("Too much tables are found")
	}
	return resp.RouteTables[0], nil
}

func (c *Ec2Client) replaceRoute(ctx context.Context, routeTableId, destinationCidrBlock, instanceId string) error {
	req, _ := c.ec2Svc.ReplaceRouteRequest(&ec2.ReplaceRouteInput{
		RouteTableId:         aws.String(routeTableId),
		InstanceId:           aws.String(instanceId),
		DestinationCidrBlock: aws.String(destinationCidrBlock),
	})
	req.HTTPRequest = req.HTTPRequest.WithContext(ctx)
	if err := req.Send(); err != nil {
		return err
	}
	return nil
}

func (c *Ec2Client) getInstanceId(ctx context.Context, instance string) (string, error) {
	if strings.HasPrefix(instance, "i-") {
		return instance, nil
	}

	req, resp := c.ec2Svc.DescribeInstancesRequest(&ec2.DescribeInstancesInput{
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
					aws.String(instance),
				},
			},
		},
	})
	req.HTTPRequest = req.HTTPRequest.WithContext(ctx)
	err := req.Send()
	switch {
	case err != nil:
		return "", err
	case len(resp.Reservations) == 0:
		return "", errors.New("Given instance is not found")
	case len(resp.Reservations[0].Instances) != 1:
		return "", errors.New("Too much instances are fetched")
	}
	return *resp.Reservations[0].Instances[0].InstanceId, nil
}
