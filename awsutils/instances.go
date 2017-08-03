package awsutils

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// AWSInstance holds aws instance information
type AWSInstance struct {
	ID           string
	InstanceType string
	Name         string
	PrivateDNS   string
	PrivateIP    string
	PublicDNS    string
	PublicIP     string
	State        string
	Tags         map[string]string
}

// NewAWSInstance creates a new instance from an ec2 instance
func NewAWSInstance(
	id string,
	instanceType string,
	name string,
	privateDNS string,
	privateIP string,
	publicDNS string,
	publicIP string,
	state string,
	tags map[string]string,
) *AWSInstance {
	return &AWSInstance{
		ID:           id,
		InstanceType: instanceType,
		Name:         name,
		PrivateDNS:   privateDNS,
		PrivateIP:    privateIP,
		PublicDNS:    publicDNS,
		PublicIP:     publicIP,
		State:        state,
		Tags:         tags,
	}
}

// String returns a string representation of the instance
func (i *AWSInstance) String() string {
	return fmt.Sprintf("<instance id=%s instancetype=%s name=%s privateDNS=%s privateIP=%s publicDNS=%s publicIP=%s state=%s tags=%s",
		i.ID,
		i.InstanceType,
		i.Name,
		i.PrivateDNS,
		i.PrivateIP,
		i.PublicDNS,
		i.PublicIP,
		i.State,
		i.Tags,
	)
}

// DiscoverInstances discovers all AWS instances
func DiscoverInstances(region string) ([]*AWSInstance, error) {
	config := &aws.Config{
		Region: aws.String(region),
		CredentialsChainVerboseErrors: aws.Bool(true),
	}

	s := ec2.New(session.New(config))
	input := &ec2.DescribeInstancesInput{}

	return discoverInstances(s, input)
}

func discoverInstances(s ec2iface.EC2API, input *ec2.DescribeInstancesInput) ([]*AWSInstance, error) {
	result, err := s.DescribeInstances(input)

	if err != nil {
		return nil, err
	}

	var instances []*AWSInstance
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {

			// Organize tags
			tags := convertTags(instance.Tags)

			// Convert instance
			i := NewAWSInstance(
				*instance.InstanceId,
				*instance.InstanceType,
				*instance.KeyName,
				*instance.PrivateDnsName,
				*instance.PrivateIpAddress,
				*instance.PublicDnsName,
				*instance.PublicIpAddress,
				*instance.State.Name,
				tags,
			)

			// TODO: Get the port to the instance. They are stored in the Security Groups

			instances = append(instances, i)
		}
	}
	return instances, nil
}

func convertTags(tags []*ec2.Tag) map[string]string {
	ts := make(map[string]string, len(tags))
	for _, kv := range tags {
		ts[*kv.Key] = *kv.Value
	}
	return ts
}
