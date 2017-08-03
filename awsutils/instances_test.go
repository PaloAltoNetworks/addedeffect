package awsutils

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	. "github.com/smartystreets/goconvey/convey"
)

type mockEC2 struct {
	ec2iface.EC2API

	Instances      *ec2.DescribeInstancesOutput
	SecurityGroups *ec2.DescribeSecurityGroupsOutput
}

func (m mockEC2) DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	if m.Instances == nil {
		return nil, fmt.Errorf("No response provided")
	}
	return m.Instances, nil
}

func (m mockEC2) DescribeSecurityGroups(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	if m.SecurityGroups == nil {
		return nil, fmt.Errorf("No response provided")
	}
	return m.SecurityGroups, nil
}

// func Test_DiscoverInstances(t *testing.T) {
// 	Convey("Given a real aws instance", t, func() {
// 		instances, err := DiscoverInstances("us-east-1")
//
// 		So(err, ShouldBeNil)
// 		So(len(instances), ShouldEqual, 1)
// 	})
// }

func Test_discoverInstances(t *testing.T) {

	Convey("Given an instance running", t, func() {

		expectedGroup := &ec2.SecurityGroup{
			GroupId: aws.String("group-8392"),
			IpPermissions: []*ec2.IpPermission{
				&ec2.IpPermission{
					FromPort: aws.Int64(443),
					ToPort:   aws.Int64(80),
				},
			},
		}

		expectedInstance := &ec2.Instance{
			InstanceId:       aws.String("i-1234"),
			InstanceType:     aws.String("Test"),
			KeyName:          aws.String("Test Instance"),
			PrivateDnsName:   aws.String("test.demo.private"),
			PrivateIpAddress: aws.String("192.168.0.10"),
			PublicDnsName:    aws.String("test.demo.public"),
			PublicIpAddress:  aws.String("10.0.0.1"),
			State: &ec2.InstanceState{
				Code: aws.Int64(0),
				Name: aws.String(ec2.StatePending),
			},
			Tags: []*ec2.Tag{
				&ec2.Tag{
					Key:   aws.String("purpose"),
					Value: aws.String("demo"),
				},
				&ec2.Tag{
					Key:   aws.String("env"),
					Value: aws.String("test"),
				},
			},
			SecurityGroups: []*ec2.GroupIdentifier{
				&ec2.GroupIdentifier{
					GroupId: aws.String("group-8392"),
				},
			},
		}

		reservation := &ec2.Reservation{
			Instances: []*ec2.Instance{
				expectedInstance,
			},
		}

		Convey("when the EC2 service is responding", func() {
			mockService := &mockEC2{
				Instances: &ec2.DescribeInstancesOutput{
					Reservations: []*ec2.Reservation{
						reservation,
					},
				},
				SecurityGroups: &ec2.DescribeSecurityGroupsOutput{
					SecurityGroups: []*ec2.SecurityGroup{
						expectedGroup,
					},
				},
			}
			instances, err := discoverInstances(mockService, &ec2.DescribeInstancesInput{})
			Convey("I should be able to discover the instance", func() {
				So(err, ShouldBeNil)
				So(len(instances), ShouldEqual, 1)
				discoveredInstance := instances[0]

				So(discoveredInstance.ID, ShouldEqual, *expectedInstance.InstanceId)
				So(discoveredInstance.InstanceType, ShouldEqual, *expectedInstance.InstanceType)
				So(discoveredInstance.Name, ShouldEqual, *expectedInstance.KeyName)
				So(discoveredInstance.PrivateDNS, ShouldEqual, *expectedInstance.PrivateDnsName)
				So(discoveredInstance.PrivateIP, ShouldEqual, *expectedInstance.PrivateIpAddress)
				So(discoveredInstance.PublicDNS, ShouldEqual, *expectedInstance.PublicDnsName)
				So(discoveredInstance.PublicIP, ShouldEqual, *expectedInstance.PublicIpAddress)
				So(discoveredInstance.State, ShouldEqual, *expectedInstance.State.Name)
				So(discoveredInstance.Tags, ShouldResemble, convertTags(expectedInstance.Tags))
				So(discoveredInstance.Ports, ShouldResemble, []string{"80"})
			})
		})

		Convey("when the EC2 service is responding without security groups", func() {
			mockService := &mockEC2{
				Instances: &ec2.DescribeInstancesOutput{
					Reservations: []*ec2.Reservation{
						reservation,
					},
				},
				SecurityGroups: &ec2.DescribeSecurityGroupsOutput{},
			}
			instances, err := discoverInstances(mockService, &ec2.DescribeInstancesInput{})
			Convey("I should be able to discover the instance", func() {
				So(err, ShouldBeNil)
				So(len(instances), ShouldEqual, 1)
				discoveredInstance := instances[0]
				So(len(discoveredInstance.Ports), ShouldEqual, 0)
			})
		})

		Convey("when the EC2 service is not responding to instances request", func() {
			mockService := &mockEC2{
				Instances: nil,
			}
			instances, err := discoverInstances(mockService, &ec2.DescribeInstancesInput{})
			Convey("I should be able to discover the instance", func() {
				So(err, ShouldNotBeNil)
				So(len(instances), ShouldEqual, 0)
			})
		})

		Convey("when the EC2 service is not responding to security groups requests", func() {
			mockService := &mockEC2{
				Instances: &ec2.DescribeInstancesOutput{
					Reservations: []*ec2.Reservation{
						reservation,
					},
				},
				SecurityGroups: nil, // will return an error
			}
			instances, err := discoverInstances(mockService, &ec2.DescribeInstancesInput{})
			Convey("I should be able to discover the instance", func() {
				So(err, ShouldNotBeNil)
				So(len(instances), ShouldEqual, 0)
			})
		})
	})
}

func Test_discoverPorts(t *testing.T) {

	Convey("Given an instance running", t, func() {
	})
}

func Test_AWSInstance_String(t *testing.T) {

	Convey("Given an AWS Instance", t, func() {
		instance := NewAWSInstance("i-xxxxx", "VM", "VM Test", "private.dns.com", "192.168.0.1", "public.dns.com", "10.0.0.1", "running", map[string]string{
			"tag1": "value1",
			"tag2": "value2",
		}, []string{"80", "443"})
		expectedString := "<instance id=i-xxxxx instancetype=VM name=VM Test privateDNS=private.dns.com privateIP=192.168.0.1 publicDNS=public.dns.com publicIP=10.0.0.1 state=running tags=map[tag1:value1 tag2:value2] ports=[80 443]"
		So(instance.String(), ShouldEqual, expectedString)
	})
}
