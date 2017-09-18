package awsutils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	metaDataPath    = "http://169.254.169.254/latest/meta-data/"
	dynamicDataPath = "http://169.254.169.254/latest/dynamic/instance-identity"
	pkcs7Name       = "pkcs7"
	// AWSName is the key for the name of the instance
	AWSName = "ami-id"
	// AWSPublicHostName is the key for the public hostname of the instance
	AWSPublicHostName = "public-hostname"
	// AWSLocalHostName is the key for the local hostname of the instance
	AWSLocalHostName = "local-hostname"
	// AWSLocalIP is the key for the local ipv4
	AWSLocalIP = "local-ipv4"
	// AWSPublicIP is the key for the public ipv4
	AWSPublicIP = "public-ipv4"
	// AWSPrivateIP is the key for the local private ip
	AWSPrivateIP = "privateIp"
	// AWSPendingTime is the key for the pending time information
	AWSPendingTime = "pendingTime"
	// AWSInstanceID is the key for the instance id
	AWSInstanceID = "instanceId"
	// AWSInstanceType is the key for the instance type
	AWSInstanceType = "instanceType"
	// AWSAccountID is the key for the account id
	AWSAccountID = "accountId"
	// AWSSecurityGroups is the key for security groups
	AWSSecurityGroups = "security-groups"
)

// getValue retrieves the value from a metadata URI - just single value
func getValue(uri string) (string, error) {

	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close() //nolint

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// getMetadataKeys retrieves all the metadaata keys that are not path levels
// we ignore paths at this point
func getMetadataKeys() ([]string, error) {
	keyString, err := getValue(metaDataPath)
	if err != nil {
		return nil, err
	}

	keys := strings.Split(keyString, "\n")
	valueKeys := []string{}

	for _, key := range keys {
		if strings.HasSuffix(key, "/") {
			continue
		}
		valueKeys = append(valueKeys, key)
	}

	return valueKeys, nil
}

// InstanceMetadata creates a Metadata map based on the available keys
func InstanceMetadata() (map[string]string, error) {

	metadata := map[string]string{}

	keys, err := getMetadataKeys()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		value, err := getValue(metaDataPath + key) // nolint
		if err != nil {
			continue
		}
		metadata[key] = value
	}

	document, err := getDocument()
	if err != nil {
		return nil, err
	}

	for k, v := range document {
		metadata[k] = v
	}

	return metadata, nil
}

// getDocument gets the Identity Document
func getDocument() (map[string]string, error) {
	document := map[string]string{}

	doc, err := getValue(dynamicDataPath + "/document")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(doc), &document); err != nil {
		return nil, err
	}

	return document, nil
}

// InstanceIdentity gets the PKCS7 identity document of an instance
func InstanceIdentity() (string, error) {

	return getValue(dynamicDataPath + "/" + pkcs7Name)
}
