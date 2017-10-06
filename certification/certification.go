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
	opentracing "github.com/opentracing/opentracing-go"
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

func GeneratePKCS12(cert []byte, key []byte, ca []byte, passphrase string) ([]byte, error) {
	// cert
	tmpcert, err := ioutil.TempFile("", "tmpcert")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpcert.Name()) // nolint: errcheck
	defer tmpcert.Close()           // nolint: errcheck
	if _, err = tmpcert.Write(cert); err != nil {
		return nil, err
	}

	// key
	tmpkey, err := ioutil.TempFile("", "tmpkey")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpkey.Name()) // nolint: errcheck
	defer tmpkey.Close()           // nolint: errcheck
	if _, err = tmpkey.Write(key); err != nil {
		return nil, err
	}

	// ca
	tmpca, err := ioutil.TempFile("", "tmpca")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpca.Name()) // nolint: errcheck
	defer tmpca.Close()           // nolint: errcheck
	if _, err = tmpca.Write(ca); err != nil {
		return nil, err
	}

	// p12
	tmpp12, err := ioutil.TempFile("", "tmpp12")
	if err != nil {
		return nil, err
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
		return nil, err
	}

	p12data, err := ioutil.ReadAll(tmpp12)
	if err != nil {
		return nil, err
	}
	return p12data, nil
}

// GenerateBase64PKCS12 generates a full PKCS certificate based on the input keys.
func GenerateBase64PKCS12(cert []byte, key []byte, ca []byte, passphrase string) (string, error) {

	p12data, err := GeneratePKCS12(cert, key, ca, passphrase)

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

// IssueCert asks and returns an new certificate using the given barret Manipulator and given CSR.
// You can generate easily a CSR using GenerateSimpleCSR.
func IssueCert(m manipulate.Manipulator, csrPEM []byte, expiration time.Time, usage KeyUsage, span opentracing.Span) (cert []byte, serialNumber string, exp time.Time, err error) {

	var sp opentracing.Span
	if span != nil {
		sp = opentracing.StartSpan("addedeffect.certification.issuecert", opentracing.ChildOf(span.Context()))
	} else {
		sp = opentracing.StartSpan("addedeffect.certification.issuecert")
	}
	defer sp.Finish()

	request := barretmodels.NewCertificate()
	request.ExpirationDate = expiration
	request.CSR = string(csrPEM)
	request.Usage = convertKeyUsage(usage)

	mctx := manipulate.NewContext()
	mctx.TrackingSpan = sp

	if err = manipulate.Retry(func() error { return m.Create(mctx, request) }, nil, 10); err != nil {
		return
	}

	cert = []byte(request.Certificate)
	serialNumber = request.ID
	exp = request.ExpirationDate

	return
}

// RevokeCert sets the revocation status of the given certificate identified by its serial number.
func RevokeCert(m manipulate.Manipulator, serialNumber string, revoked bool, span opentracing.Span) (err error) {

	var sp opentracing.Span
	if span != nil {
		sp = opentracing.StartSpan("addedeffect.certification.revokecert", opentracing.ChildOf(span.Context()))
	} else {
		sp = opentracing.StartSpan("addedeffect.certification.revokecert")
	}
	defer sp.Finish()

	request := barretmodels.NewRevocation()
	request.Revoked = revoked
	request.ID = serialNumber

	mctx := manipulate.NewContext()
	mctx.TrackingSpan = sp

	return manipulate.Retry(func() error { return m.Update(mctx, request) }, nil, 10)
}

// CheckRevocation checks if the given certificate serial number is revoked.
func CheckRevocation(m manipulate.Manipulator, serialNumber string, span opentracing.Span) (err error) {

	var sp opentracing.Span
	if span != nil {
		sp = opentracing.StartSpan("addedeffect.certification.checkrevocation", opentracing.ChildOf(span.Context()))
	} else {
		sp = opentracing.StartSpan("addedeffect.certification.checkrevocation")
	}
	defer sp.Finish()

	request := barretmodels.NewCheck()
	request.ID = serialNumber

	mctx := manipulate.NewContext()
	mctx.TrackingSpan = sp

	return manipulate.Retry(func() error { return m.Retrieve(mctx, request) }, nil, 10)
}

// IssueEncryptionToken asks and return a token from the given certificate using the given barret manipulator.
func IssueEncryptionToken(m manipulate.Manipulator, cert []byte, span opentracing.Span) (token string, err error) {

	var sp opentracing.Span
	if span != nil {
		sp = opentracing.StartSpan("addedeffect.certification.issueencryptiontoken", opentracing.ChildOf(span.Context()))
	} else {
		sp = opentracing.StartSpan("addedeffect.certification.issueencryptiontoken")
	}
	defer sp.Finish()

	request := barretmodels.NewToken()
	request.Certificate = string(cert)

	mctx := manipulate.NewContext()
	mctx.TrackingSpan = sp

	if err = manipulate.Retry(func() error { return m.Create(mctx, request) }, nil, 10); err != nil {
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
