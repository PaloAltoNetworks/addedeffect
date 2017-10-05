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
)

// KeyUsage is the type of key usage.
type KeyUsage int

// Various possible KeyUsage values
const (
	KeyUsageClient KeyUsage = iota + 1
	KeyUsageServer
	KeyUsageClientServer
)

// Verify verifies the given certificate is signed by the given other certificate, and that
// the other certificate has the correct required key usage.
func Verify(signingCertPEMData []byte, certPEMData []byte, keyUsages []x509.ExtKeyUsage) error {

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(signingCertPEMData))
	if !ok {
		return fmt.Errorf("Unable to parse signing certificate")
	}

	block, rest := pem.Decode(certPEMData)
	if block == nil || len(rest) != 0 {
		return fmt.Errorf("Invalid child certificate")
	}

	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("Unable to parse child certificate: %s", err)
	}

	if keyUsages == nil {
		keyUsages = []x509.ExtKeyUsage{x509.ExtKeyUsageAny}
	}

	if _, err := x509Cert.Verify(
		x509.VerifyOptions{
			Roots:     roots,
			KeyUsages: keyUsages,
		},
	); err != nil {
		return fmt.Errorf("Unable to verify child certificate: %s", err)
	}

	return nil
}

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

// SplitChain splits the given certificate data into the actual *x509.Certificate and a list of
// CA chain in a []*x509.Certificate
func SplitChain(certData []byte) (cert *x509.Certificate, caChain []*x509.Certificate, err error) {

	block, rest := pem.Decode(certData)

	for ; block != nil; block, rest = pem.Decode(rest) {

		if block.Type != "CERTIFICATE" {
			continue
		}

		crt, err := x509.ParseCertificate(block.Bytes)

		if err != nil {
			return nil, nil, err
		}

		if !crt.IsCA {
			cert = crt
			continue
		}

		if len(rest) != 0 {
			caChain = append(caChain, crt)
		}
	}

	return
}

// SplitChainPEM splits the given cert PEM []byte as the actual certificate
// and []byte as the rest of the chain.
func SplitChainPEM(certData []byte) ([]byte, []byte) {

	block, rest := pem.Decode(certData)

	return pem.EncodeToMemory(block), rest
}
