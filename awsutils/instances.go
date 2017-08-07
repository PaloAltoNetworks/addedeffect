package awsutils

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
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
	Tags         []string
	Ports        []string
	Namespace    string
	AccountID    string
}

// NewAWSInstance creates a new instance from an ec2 instance
func NewAWSInstance(
	id string,
	accountID string,
	instanceType string,
	name string,
	privateDNS string,
	privateIP string,
	publicDNS string,
	publicIP string,
	state string,
	tags []string,
	ports []string,
) *AWSInstance {

	// Order the tags...
	sort.Strings(tags)

	namespace := extractNamespaceFromTags(tags)

	return &AWSInstance{
		ID:           id,
		AccountID:    accountID,
		InstanceType: instanceType,
		Name:         name,
		PrivateDNS:   privateDNS,
		PrivateIP:    privateIP,
		PublicDNS:    publicDNS,
		PublicIP:     publicIP,
		State:        state,
		Tags:         tags,
		Ports:        ports,
		Namespace:    namespace,
	}
}

// String returns a string representation of the instance
func (i *AWSInstance) String() string {

	return fmt.Sprintf("<instance id=%s accountID=%s instancetype=%s name=%s privateDNS=%s privateIP=%s publicDNS=%s publicIP=%s state=%s tags=%s ports=%s",
		i.ID,
		i.AccountID,
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
func DiscoverInstances(accessKey string, secretKey string, region string) ([]*AWSInstance, error) {
	config := &aws.Config{
		Region: aws.String(region),
		CredentialsChainVerboseErrors: aws.Bool(true),
		Credentials:                   credentials.NewStaticCredentials(accessKey, secretKey, ""),
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

			if !shouldManageInstance(instance) {
				continue
			}

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
			accountID := ""
			if instance.IamInstanceProfile != nil {
				accountID = *instance.IamInstanceProfile.Id
			}
			publicIP := ""
			if instance.PublicIpAddress != nil {
				publicIP = *instance.PublicIpAddress
			}
			publicDNS := ""
			if instance.PublicDnsName != nil {
				publicDNS = *instance.PublicDnsName
			}
			privateIP := ""
			if instance.PrivateIpAddress != nil {
				privateIP = *instance.PrivateIpAddress
			}
			privateDNS := ""
			if instance.PrivateDnsName != nil {
				privateDNS = *instance.PrivateDnsName
			}

			// Convert instance
			i := NewAWSInstance(
				*instance.InstanceId,
				accountID,
				*instance.InstanceType,
				fmt.Sprintf("aws-%s", *instance.InstanceId),
				privateDNS,
				privateIP,
				publicDNS,
				publicIP,
				*instance.State.Name,
				tags,
				ports,
			)

			instances = append(instances, i)
		}
	}
	return instances, nil
}

func shouldManageInstance(instance *ec2.Instance) bool {
	if instance.State == nil {
		return false
	}

	state := *instance.State.Name
	if state == ec2.InstanceStateNameTerminated || state == ec2.InstanceStateNameStopped {
		return false
	}

	if instance.PublicIpAddress == nil {
		return false
	}

	return true
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

// convertTags transforms a list of ec2.Tag to a []string
func convertTags(tags []*ec2.Tag) []string {
	ts := []string{}

	for _, kv := range tags {
		ts = append(ts, fmt.Sprintf("%s=%s", *kv.Key, *kv.Value))
	}

	sort.Strings(ts)
	return ts
}

func extractNamespaceFromTags(tags []string) string {

	for _, s := range tags {
		if strings.HasPrefix(s, "namespace") {
			infos := strings.SplitN(s, "=", 2)

			if len(infos) == 2 {
				return infos[1]
			}
		}
	}

	return ""
}
