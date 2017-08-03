package awsutils

import (
	"fmt"
	"sort"
	"strconv"

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
	Ports        []string
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
	ports []string,
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
		Ports:        ports,
	}
}

// String returns a string representation of the instance
func (i *AWSInstance) String() string {
	return fmt.Sprintf("<instance id=%s instancetype=%s name=%s privateDNS=%s privateIP=%s publicDNS=%s publicIP=%s state=%s tags=%s ports=%s",
		i.ID,
		i.InstanceType,
		i.Name,
		i.PrivateDNS,
		i.PrivateIP,
		i.PublicDNS,
		i.PublicIP,
		i.State,
		i.Tags,
		i.Ports,
	)
}

// DiscoverInstances discovers all AWS instances
func DiscoverInstances(region string) ([]*AWSInstance, error) {
	config := &aws.Config{
		Region: aws.String(region),
		CredentialsChainVerboseErrors: aws.Bool(true),
	}

	session, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	service := ec2.New(session)
	input := &ec2.DescribeInstancesInput{}

	return discoverInstances(service, input)
}

func discoverInstances(service ec2iface.EC2API, input *ec2.DescribeInstancesInput) ([]*AWSInstance, error) {
	result, err := service.DescribeInstances(input)

	if err != nil {
		return nil, err
	}

	var instances []*AWSInstance
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {

			// Organize tags
			tags := convertTags(instance.Tags)

			// Discover ports
			var ports []string
			groupIDs := make([]string, len(instance.SecurityGroups))
			if len(groupIDs) > 0 {
				for _, sg := range instance.SecurityGroups {
					groupIDs = append(groupIDs, *sg.GroupId)
				}

				input := &ec2.DescribeSecurityGroupsInput{
					GroupIds: aws.StringSlice(groupIDs),
				}

				ports, err = discoverPorts(service, input)

				if err != nil {
					return nil, err
				}
			}

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
				ports,
			)

			instances = append(instances, i)
		}
	}
	return instances, nil
}

func discoverPorts(s ec2iface.EC2API, input *ec2.DescribeSecurityGroupsInput) ([]string, error) {

	result, err := s.DescribeSecurityGroups(input)

	if err != nil {
		return nil, err
	}

	var ports []string
	for _, sg := range result.SecurityGroups {
		for _, p := range sg.IpPermissions {
			if p.ToPort != nil {
				ports = append(ports, strconv.FormatInt(*p.ToPort, 10))
			}
		}
	}

	// Order ports for testing purposes
	sort.Strings(ports)

	return ports, nil
}

// convertTags transforms a list of ec2.Tag to a map[string]string always alphabetically ordered.
// Order makes it easier to test.
func convertTags(tags []*ec2.Tag) map[string]string {
	cache := make(map[string]string, len(tags))
	ts := make(map[string]string, len(tags))
	keys := make([]string, len(tags))

	for _, kv := range tags {
		keys = append(keys, *kv.Key)
		cache[*kv.Key] = *kv.Value
	}

	sort.Strings(keys)

	for _, k := range keys {
		ts[k] = cache[k]
	}

	return ts
}
