package aws

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
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

func (c *Ec2Client) getRouteTables(ctx context.Context, retry int) ([]*ec2.RouteTable, error) {
	req, resp := c.ec2Svc.DescribeRouteTablesRequest(nil)
	req.HTTPRequest = req.HTTPRequest.WithContext(ctx)
	var err error
	for i := 0; i < retry; i++ {
		err = req.Send()
		if err == nil && len(resp.RouteTables) > 0 {
			break
		}
	}
	if err != nil || len(resp.RouteTables) < 1 {
		return nil, err
	}

	return resp.RouteTables, nil
}

func (c *Ec2Client) getRouteTableByKey(ctx context.Context, retry int, key string) (*ec2.RouteTable, error) {
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
	var err error
	for i := 0; i < retry; i++ {
		err := req.Send()
		if err == nil && len(resp.RouteTables) > 0 {
			break
		}
	}
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

func (c *Ec2Client) replaceRoute(ctx context.Context, retry int, routeTableId, destinationCidrBlock, instanceId string) error {
	req, _ := c.ec2Svc.ReplaceRouteRequest(&ec2.ReplaceRouteInput{
		RouteTableId:         aws.String(routeTableId),
		InstanceId:           aws.String(instanceId),
		DestinationCidrBlock: aws.String(destinationCidrBlock),
	})
	req.HTTPRequest = req.HTTPRequest.WithContext(ctx)
	var err error
	for i := 0; i < retry; i++ {
		err := req.Send()
		if err == nil {
			break
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *Ec2Client) getInstanceIdByDest(ctx context.Context, retry int, routeTableId, dest string) (string, error) {
	t, err := c.getRouteTableByKey(ctx, retry, routeTableId)
	if err != nil {
		return "", err
	}
	for _, r := range t.Routes {
		if *r.DestinationCidrBlock == dest {
			return *r.InstanceId, nil
		}
	}
	return "", errors.New("Not found")
}

func (c *Ec2Client) getInstanceByKey(ctx context.Context, retry int, key string) (*ec2.Instance, error) {
	var input *ec2.DescribeInstancesInput
	if strings.HasPrefix(key, "i-") {
		input = &ec2.DescribeInstancesInput{
			InstanceIds: []*string{
				aws.String(key),
			},
		}
	} else {
		input = &ec2.DescribeInstancesInput{
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

	var err error
	req, resp := c.ec2Svc.DescribeInstancesRequest(input)
	req.HTTPRequest = req.HTTPRequest.WithContext(ctx)
	for i := 0; i < retry; i++ {
		err := req.Send()
		if err == nil && len(resp.Reservations) > 0 {
			break
		}
	}
	switch {
	case err != nil:
		return nil, err
	case len(resp.Reservations) == 0:
		return nil, errors.New("Given instance is not found")
	case len(resp.Reservations[0].Instances) != 1:
		return nil, errors.New("Too much instances are fetched")
	}
	return resp.Reservations[0].Instances[0], nil
}

func (c *Ec2Client) getInstanceId(ctx context.Context, retry int, key string) (string, error) {
	instance, err := c.getInstanceByKey(ctx, retry, key)
	if err != nil {
		return "", err
	}
	return *instance.InstanceId, nil
}

func (c *Ec2Client) getInstanceNameById(ctx context.Context, retry int, instanceId string) (string, error) {
	instance, err := c.getInstanceByKey(ctx, retry, instanceId)
	if err != nil {
		return "", err
	}
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			return *tag.Value, nil
		}
	}
	return "", nil
}

func (c *Ec2Client) getENINameById(ctx context.Context, retry int, ENIId string) (string, error) {
	input := &ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("network-interface-id"),
				Values: []*string{
					aws.String(ENIId),
				},
			},
		},
	}

	req, resp := c.ec2Svc.DescribeNetworkInterfacesRequest(input)
	req.HTTPRequest = req.HTTPRequest.WithContext(ctx)
	var err error
	for i := 0; i < retry; i++ {
		err := req.Send()
		if err == nil && len(resp.NetworkInterfaces) > 0 {
			break
		}
	}
	switch {
	case err != nil:
		return "", err
	case len(resp.NetworkInterfaces) == 0:
		return "", errors.New("Given interface is not found")
	case len(resp.NetworkInterfaces) != 1:
		return "", errors.New("Too much instances are fetched")
	}
	eni := resp.NetworkInterfaces[0]
	for _, tag := range eni.TagSet {
		if *tag.Key == "Name" {
			return *tag.Value, nil
		}
	}
	return "", nil
}
