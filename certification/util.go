// Author: Suresh Ramamurthy
// See LICENSE file for full LICENSE
// Copyright 2016 Aporeto.

package certification

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

func pemBlockForKey(priv interface{}) *pem.Block {

	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}

	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}

	default:
		return nil
	}
}

func loadCertificateBundle(filename string) ([]*x509.Certificate, error) {

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	certificates := []*x509.Certificate{}
	var block *pem.Block
	block, b = pem.Decode(b)

	for ; block != nil; block, b = pem.Decode(b) {
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			certificates = append(certificates, cert)
		} else {
			return nil, fmt.Errorf("invalid pem block type: %s", block.Type)
		}
	}

	return certificates, nil
}
