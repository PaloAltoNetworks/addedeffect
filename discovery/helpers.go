package discovery

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func loadCertificates(certData []byte, keyData []byte, password string) (tls.Certificate, error) {

	keyBlock, rest := pem.Decode(keyData)

	if len(rest) != 0 {
		return tls.Certificate{}, fmt.Errorf("Multiple private key is not supported.")
	}

	if !x509.IsEncryptedPEMBlock(keyBlock) {
		return tls.X509KeyPair(certData, keyData)
	}

	data, err := x509.DecryptPEMBlock(keyBlock, []byte(password))
	if err != nil {
		return tls.Certificate{}, err
	}

	buffer := &bytes.Buffer{}
	if err := pem.Encode(buffer, &pem.Block{Type: keyBlock.Type, Bytes: data}); err != nil {
		return tls.Certificate{}, err
	}

	return tls.X509KeyPair(certData, buffer.Bytes())
}
