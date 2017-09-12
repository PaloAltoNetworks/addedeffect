package certification

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/aporeto-inc/gaia/barretmodels/v1/golang"
	"github.com/aporeto-inc/manipulate"
)

// KeyUsage is the type of key usage.
type KeyUsage int

// Various possible KeyUsage values
const (
	KeyUsageClient KeyUsage = iota + 1
	KeyUsageServer
	KeyUsageClientServer
)

// GenerateBase64PKCS12 generates a full PKCS certificate based on the input keys.
func GenerateBase64PKCS12(cert []byte, key []byte, ca []byte, passphrase string) (string, error) {

	// cert
	tmpcert, err := ioutil.TempFile("", "tmpcert")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpcert.Name()) // nolint: errcheck
	defer tmpcert.Close()           // nolint: errcheck
	if _, err = tmpcert.Write(cert); err != nil {
		return "", err
	}

	// key
	tmpkey, err := ioutil.TempFile("", "tmpkey")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpkey.Name()) // nolint: errcheck
	defer tmpkey.Close()           // nolint: errcheck
	if _, err = tmpkey.Write(key); err != nil {
		return "", err
	}

	// ca
	tmpca, err := ioutil.TempFile("", "tmpca")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpca.Name()) // nolint: errcheck
	defer tmpca.Close()           // nolint: errcheck
	if _, err = tmpca.Write(ca); err != nil {
		return "", err
	}

	// p12
	tmpp12, err := ioutil.TempFile("", "tmpp12")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpp12.Name()) // nolint: errcheck
	defer tmpp12.Close()           // nolint: errcheck

	args := []string{
		"pkcs12",
		"-export",
		"-out", tmpp12.Name(),
		"-inkey", tmpkey.Name(),
		"-in", tmpcert.Name(),
		"-certfile", tmpca.Name(),
		"-passout", "pass:" + passphrase,
	}

	if err = exec.Command("openssl", args...).Run(); err != nil {
		return "", err
	}

	p12data, err := ioutil.ReadAll(tmpp12)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(p12data), nil
}

// GenerateECPrivateKey quickly generate a ECDSA private key.
func GenerateECPrivateKey() (crypto.PrivateKey, error) {

	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// PrivateKeyToPEM converts the given crypto.PrivateKey as PEM.
func PrivateKeyToPEM(key crypto.PrivateKey) ([]byte, error) {

	var err error
	var derBytes []byte
	var keyType string

	switch k := key.(type) {

	case *ecdsa.PrivateKey:
		if derBytes, err = x509.MarshalECPrivateKey(k); err != nil {
			return nil, err
		}
		keyType = "EC PRIVATE KEY"

	case *rsa.PrivateKey:
		derBytes = x509.MarshalPKCS1PrivateKey(k)
		keyType = "RSA PRIVATE KEY"

	default:
		return nil, fmt.Errorf("Given key is not ECDSA")
	}

	return pem.EncodeToMemory(
		&pem.Block{
			Type:  keyType,
			Bytes: derBytes,
		},
	), nil
}

// GenerateSimpleCSR generate a CSR using the given parameters.
func GenerateSimpleCSR(orgs []string, units []string, commonName string, emails []string, privateKey crypto.PrivateKey) ([]byte, error) {

	csr := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:         commonName,
			Organization:       orgs,
			OrganizationalUnit: units,
		},
		EmailAddresses:     emails,
		SignatureAlgorithm: x509.ECDSAWithSHA384,
	}

	csrDerBytes, err := x509.CreateCertificateRequest(rand.Reader, csr, privateKey)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDerBytes}), nil
}

// IssueCert asks and returns an new certificate using the given barret Manipulator and given CSR.
// You can generate easily a CSR using GenerateSimpleCSR.
func IssueCert(m manipulate.Manipulator, csrPEM []byte, expiration time.Time, usage KeyUsage) (cert []byte, serialNumber string, exp time.Time, err error) {

	request := barretmodels.NewCertificate()
	request.ExpirationDate = expiration
	request.CSR = string(csrPEM)
	request.Usage = convertKeyUsage(usage)

	if err = m.Create(nil, request); err != nil {
		return
	}

	cert = []byte(request.Certificate)
	serialNumber = request.ID
	exp = request.ExpirationDate

	return
}

// IssueEncryptionToken asks and return a token from the given certificate using the given barret manipulator.
func IssueEncryptionToken(m manipulate.Manipulator, cert []byte) (token string, err error) {

	request := barretmodels.NewToken()
	request.Certificate = string(cert)

	if err = m.Create(nil, request); err != nil {
		return
	}

	token = request.Token
	return
}

func convertKeyUsage(usage KeyUsage) barretmodels.CertificateUsageValue {
	switch usage {
	case KeyUsageServer:
		return barretmodels.CertificateUsageServer
	case KeyUsageClientServer:
		return barretmodels.CertificateUsageServerclient
	default:
		return barretmodels.CertificateUsageClient
	}
}
