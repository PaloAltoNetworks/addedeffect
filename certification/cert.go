package certification

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/aporeto-inc/elemental"
)

// A Signer is certificate signing object. It uses a CA certificate
// and a private key to create and/or sign a client certificate.
type Signer struct {
	cacert []*x509.Certificate
	key    crypto.PrivateKey
}

// NewSigner returns a pointer to a new Signer given a certificate path, a private key path and a password.
//
// If something went wrong the returned signer will be nil, and an error will be returned.
func NewSigner(CACertData, CACertKeyData []byte, keyPass string) (*Signer, error) {

	var c Signer

	// Load CA.pem.
	cacert, err := LoadCertificateBundle(CACertData)
	if err != nil {
		zap.L().Error("Failed to load ca certificate", zap.Error(err))
		return nil, elemental.NewError("Invalid CA certificate", "Failed to load the ca certificate", "certification", http.StatusUnprocessableEntity)
	}

	c.cacert = cacert

	// Load CA-Key.pem file.

	block, _ := pem.Decode(CACertKeyData)
	if block == nil {
		return nil, fmt.Errorf("Can not decode CA private key")
	}

	decryptedPEMBlock := block.Bytes

	if procType, ok := block.Headers["Proc-Type"]; ok && procType == "4,ENCRYPTED" {
		decryptedPEMBlock, err = x509.DecryptPEMBlock(block, []byte(keyPass))
		if err != nil {
			return nil, fmt.Errorf("Cannot decrypt CA private key")
		}
		block.Bytes = decryptedPEMBlock
	}

	// Parse the key
	if block.Type == "EC PRIVATE KEY" {
		c.key, err = x509.ParseECPrivateKey(decryptedPEMBlock)
	} else if block.Type == "RSA PRIVATE KEY" {
		c.key, err = x509.ParsePKCS1PrivateKey(decryptedPEMBlock)
	}

	if err != nil || c.key == nil {
		return nil, elemental.NewError("Unmarshal failed", "Failed to unmarshal the ca key file", "certification", http.StatusUnprocessableEntity)
	}
	return &c, nil
}

// IssueClientCertificate creates a new client certificate and signs it with the CA certificate in memory.
//
// It will return the private key in a PEM formatted string, the certificate or an error.
func (s *Signer) IssueClientCertificate(expiration time.Time, cn string, email string, org []string, units []string, dnsNames []string) (string, string, string, error) {

	for _, o := range org {
		if o == "system" {
			return "", "", "", fmt.Errorf("System organization is reserved")
		}
	}

	var key crypto.PrivateKey
	var err error

	// Generate the key.
	switch s.key.(type) {
	case *ecdsa.PrivateKey:
		key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	default:
		key, err = rsa.GenerateKey(rand.Reader, 2048)
	}

	if err != nil {
		zap.L().Error("Failed to generate private key", zap.Error(err))
		return "", "", "", elemental.NewError("Certificate generation failed", "Failed to generate private key", "certification", http.StatusInternalServerError)
	}

	// Generate random serial number.
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		zap.L().Error("Failed to generate serial number for the certificate", zap.Error(err))
		return "", "", "", elemental.NewError("Certificate generation failed", "Failed to generate serial number for the certificate", "certification", http.StatusInternalServerError)
	}

	// Create certfificate template.
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:       org,
			CommonName:         cn,
			OrganizationalUnit: units,
		},
		NotBefore:      time.Now(),
		NotAfter:       expiration,
		EmailAddresses: []string{email},
		DNSNames:       dnsNames,

		// KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},

		BasicConstraintsValid: true,
	}

	var derBytes []byte
	switch s.key.(type) {
	case *ecdsa.PrivateKey:
		derBytes, err = x509.CreateCertificate(rand.Reader, &template, s.cacert[0], &key.(*ecdsa.PrivateKey).PublicKey, s.key)
	default:
		derBytes, err = x509.CreateCertificate(rand.Reader, &template, s.cacert[0], &key.(*rsa.PrivateKey).PublicKey, s.key)
	}

	if err != nil {
		zap.L().Error("Failed to create certificate", zap.Error(err))
		return "", "", "", elemental.NewError("Failed to Create Certificate", err.Error(), "certification", http.StatusInternalServerError)
	}

	clientCertificate := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certificatePem := string(bytes.TrimSpace(clientCertificate))

	clientKeyCertificate := pem.EncodeToMemory(pemBlockForKey(key))
	keyPem := string(bytes.TrimSpace(clientKeyCertificate))

	return keyPem, certificatePem, serialNumber.String(), nil
}

// Secrets returns the current secrets used by the signer
func (s *Signer) Secrets() (crypto.PrivateKey, []*x509.Certificate) {
	return s.key, s.cacert
}
