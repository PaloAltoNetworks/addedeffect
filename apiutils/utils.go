package apiutils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aporeto-inc/addedeffect/retry"
	"go.uber.org/zap"
)

// ServiceVersion holds the version of a servie
type ServiceVersion struct {
	Libs    map[string]string
	Version string
	Sha     string
}

// GetServiceVersions returns the version of the services.
func GetServiceVersions(ctx context.Context, api string, tlsConfig *tls.Config) (map[string]ServiceVersion, error) {

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	url := fmt.Sprintf("%s/_meta/versions", api)
	out, err := retry.Retry(
		ctx,
		makeJobFunc(client, url),
		makeRetryFunc("Unable to retrieve versions. Retrying in 3s", url),
	)

	if err != nil {
		return nil, err
	}

	resp := out.(*http.Response)

	config := map[string]ServiceVersion{}

	defer resp.Body.Close() // nolint: errcheck
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// GetConfig returns the additional config exposed by the gateway.
func GetConfig(ctx context.Context, api string, tlsConfig *tls.Config) (map[string]string, error) {

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	url := fmt.Sprintf("%s/_meta/config", api)
	out, err := retry.Retry(
		ctx,
		makeJobFunc(client, url),
		makeRetryFunc("Unable to retrieve config. Retrying in 3s", url),
	)

	if err != nil {
		return nil, err
	}

	resp := out.(*http.Response)

	config := map[string]string{}

	defer resp.Body.Close() // nolint: errcheck
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// GetPublicCA returns the public CA used by the api.
func GetPublicCA(ctx context.Context, api string, tlsConfig *tls.Config) ([]byte, error) {

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	url := fmt.Sprintf("%s/_meta/ca", api)
	out, err := retry.Retry(
		ctx,
		makeJobFunc(client, url),
		makeRetryFunc("Unable to retrieve public ca. Retrying in 3s", url),
	)

	if err != nil {
		return nil, err
	}

	resp := out.(*http.Response)

	defer resp.Body.Close() // nolint: errcheck
	return ioutil.ReadAll(resp.Body)
}

// GetPublicCAPool returns the public CA used by the api as a *x509.CertPool.
func GetPublicCAPool(ctx context.Context, api string, tlsConfig *tls.Config) (*x509.CertPool, error) {

	cadata, err := GetPublicCA(ctx, api, tlsConfig)
	if err != nil {
		return nil, err
	}

	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	pool.AppendCertsFromPEM(cadata)

	return pool, nil
}

// GetJWTCert returns the public certificate used to sign jwt.
func GetJWTCert(ctx context.Context, api string, tlsConfig *tls.Config) ([]byte, error) {

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	url := fmt.Sprintf("%s/_meta/jwtcert", api)
	out, err := retry.Retry(
		ctx,
		makeJobFunc(client, url),
		makeRetryFunc("Unable to retrieve jwt certificate. Retrying in 3s", url),
	)

	if err != nil {
		return nil, err
	}

	resp := out.(*http.Response)

	defer resp.Body.Close() // nolint: errcheck
	return ioutil.ReadAll(resp.Body)
}

// GetJWTX509Cert returns the public certificate used to sign jwt as an *x509.Certificate.
func GetJWTX509Cert(ctx context.Context, api string, tlsConfig *tls.Config) (*x509.Certificate, error) {

	data, err := GetJWTCert(ctx, api, tlsConfig)
	if err != nil {
		return nil, err
	}

	block, rest := pem.Decode(data)
	if block == nil {
		return nil, errors.New("unable to parse certificate data")
	}
	if len(rest) != 0 {
		return nil, errors.New("multiple certificates found in the certificate")
	}

	return x509.ParseCertificate(block.Bytes)
}

// GetManifestURL returns the url of the manifest.
func GetManifestURL(ctx context.Context, api string, tlsConfig *tls.Config) ([]byte, error) {

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	url := fmt.Sprintf("%s/_meta/manifest", api)
	out, err := retry.Retry(
		ctx,
		makeJobFunc(client, url),
		makeRetryFunc("Unable to retrieve manifest url. Retrying in 3s", url),
	)

	if err != nil {
		return nil, err
	}

	resp := out.(*http.Response)

	defer resp.Body.Close() // nolint: errcheck
	return ioutil.ReadAll(resp.Body)
}

// GetGoogleOAuthClientID returns the Google oauth client ID used bby the platform.
func GetGoogleOAuthClientID(ctx context.Context, api string, tlsConfig *tls.Config) ([]byte, error) {

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	url := fmt.Sprintf("%s/_meta/googleclientid", api)
	out, err := retry.Retry(
		ctx,
		makeJobFunc(client, url),
		makeRetryFunc("Unable to retrieve google client id. Retrying in 3s", url),
	)

	if err != nil {
		return nil, err
	}

	resp := out.(*http.Response)

	defer resp.Body.Close() // nolint: errcheck
	return ioutil.ReadAll(resp.Body)
}

func makeJobFunc(client *http.Client, url string) func() (interface{}, error) {

	return func() (interface{}, error) {

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("Bad response status: %s", resp.Status)
		}

		return resp, nil
	}
}

func makeRetryFunc(message string, url string) func(error) error {

	return func(err error) error {
		zap.L().Warn(message, zap.String("url", url), zap.Error(err))
		return nil
	}
}
