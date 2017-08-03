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
	Response *ec2.DescribeInstancesOutput
}

func (m mockEC2) DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	if m.Response == nil {
		return nil, fmt.Errorf("No response provided")
	}
	return m.Response, nil
}

func Test_discoverInstances(t *testing.T) {

	Convey("Given an instance running", t, func() {

		reservation := &ec2.Reservation{
			Instances: []*ec2.Instance{
				&ec2.Instance{
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
				},
			},
		}

		Convey("when the EC2 service is responding", func() {
			mockService := &mockEC2{
				Response: &ec2.DescribeInstancesOutput{
					Reservations: []*ec2.Reservation{
						reservation,
					},
				},
			}
			instances, err := discoverInstances(mockService, &ec2.DescribeInstancesInput{})
			Convey("I should be able to discover the instance", func() {
				So(err, ShouldBeNil)
				So(len(instances), ShouldEqual, 1)
				discoveredInstance := instances[0]
				expectedInstance := reservation.Instances[0]
				So(discoveredInstance.ID, ShouldEqual, *expectedInstance.InstanceId)
				So(discoveredInstance.InstanceType, ShouldEqual, *expectedInstance.InstanceType)
				So(discoveredInstance.Name, ShouldEqual, *expectedInstance.KeyName)
				So(discoveredInstance.PrivateDNS, ShouldEqual, *expectedInstance.PrivateDnsName)
				So(discoveredInstance.PrivateIP, ShouldEqual, *expectedInstance.PrivateIpAddress)
				So(discoveredInstance.PublicDNS, ShouldEqual, *expectedInstance.PublicDnsName)
				So(discoveredInstance.PublicIP, ShouldEqual, *expectedInstance.PublicIpAddress)
				So(discoveredInstance.State, ShouldEqual, *expectedInstance.State.Name)
				So(discoveredInstance.Tags, ShouldResemble, convertTags(expectedInstance.Tags))
			})
		})

		Convey("when the EC2 service is not responding", func() {
			mockService := &mockEC2{
				Response: nil,
			}
			instances, err := discoverInstances(mockService, &ec2.DescribeInstancesInput{})
			Convey("I should be able to discover the instance", func() {
				So(err, ShouldNotBeNil)
				So(len(instances), ShouldEqual, 0)
			})
		})
	})
}

func Test_AWSInstance_String(t *testing.T) {

	Convey("Given an AWS Instance", t, func() {
		instance := NewAWSInstance("i-xxxxx", "VM", "VM Test", "private.dns.com", "192.168.0.1", "public.dns.com", "10.0.0.1", "running", map[string]string{
			"tag1": "value1",
			"tag2": "value2",
		})
		expectedString := "<instance id=i-xxxxx instancetype=VM name=VM Test privateDNS=private.dns.com privateIP=192.168.0.1 publicDNS=public.dns.com publicIP=10.0.0.1 state=running tags=map[tag1:value1 tag2:value2]"
		So(instance.String(), ShouldEqual, expectedString)
	})
}
