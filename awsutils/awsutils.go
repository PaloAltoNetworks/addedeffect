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
	// AWSPendingTime is the key for the pending time information
	AWSPendingTime = "pendingTime"
	// AWSInstanceID is the key for the instance id
	AWSInstanceID = "instanceId"
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

// getMetaData creates a Metadata map based on the available keys
func getMetaData() (map[string]string, error) {

	metadata := map[string]string{}

	keys, err := getMetadataKeys()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		value, err := getValue(metaDataPath + key)
		if err != nil {
			continue
		}
		metadata[key] = value
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

// getIdentity gets the PKCS7 identity document of an instance
func getIdentity() (string, error) {

	return getValue(dynamicDataPath + "/" + pkcs7Name)
}

// InstanceMetadata retrieves all the instance metadata in a map
// It returns the PKCS7 ID, a map of the meta-data and error
func InstanceMetadata() (string, map[string]string, error) {

	metadata, err := getMetaData()
	if err != nil {
		return "", nil, err
	}

	id, err := getIdentity()
	if err != nil {
		return "", nil, err
	}

	document, err := getDocument()
	if err != nil {
		return "", nil, err
	}

	for k, v := range document {
		metadata[k] = v
	}

	return id, metadata, nil
}
