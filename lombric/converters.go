package lombric

import (
	"fmt"
	"net"
	"strconv"
)

func convertDefaultBool(defaultValue []string) (bools []bool, err error) {

	for _, item := range defaultValue {

		switch item {
		case "true", "True", "TRUE":
			bools = append(bools, true)
		case "false", "False", "FALSE":
			bools = append(bools, false)
		default:
			return nil, fmt.Errorf("default value must a bool got: '%s'", item)
		}
	}

	return
}

func convertDefaultInts(defaultValue []string) (ints []int, err error) {

	for _, item := range defaultValue {

		n, err := strconv.Atoi(item)

		if err != nil {
			return nil, fmt.Errorf("default value must be an int. got '%s'", item)
		}

		ints = append(ints, n)
	}

	return
}

func convertDefaultIPs(defaultValue []string) (ips []net.IP, err error) {

	for _, item := range defaultValue {

		ip := net.ParseIP(item)
		if ip == nil {
			return nil, fmt.Errorf("default value must be an IP. got '%s'", item)
		}

		ips = append(ips, ip)
	}

	return
}
