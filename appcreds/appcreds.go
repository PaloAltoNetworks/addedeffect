package appcreds

import (
	"context"
	"encoding/base64"
	"encoding/pem"

	"go.aporeto.io/gaia"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/tg/tglib"
)

// New creates a new *gaia.AppCredential.
func New(ctx context.Context, m manipulate.Manipulator, namespace string, name string, roles []string) (*gaia.AppCredential, error) {

	creds := gaia.NewAppCredential()
	creds.Name = name
	creds.Roles = roles
	creds.Namespace = namespace

	return NewWithAppCredential(ctx, m, creds)
}

// NewWithAppCredential creates a new *gaia.AppCredential from an *AppCredential
func NewWithAppCredential(ctx context.Context, m manipulate.Manipulator, template *gaia.AppCredential) (*gaia.AppCredential, error) {

	creds := gaia.NewAppCredential()
	creds.Name = template.Name
	creds.Description = template.Description
	creds.Roles = template.Roles
	creds.Protected = template.Protected
	creds.Metadata = template.Metadata

	if err := m.Create(
		manipulate.NewContext(
			ctx,
			manipulate.ContextOptionNamespace(template.Namespace),
		),
		creds,
	); err != nil {
		return nil, err
	}

	return Renew(ctx, m, creds)
}

// Renew renews the given appcred.
func Renew(ctx context.Context, m manipulate.Manipulator, creds *gaia.AppCredential) (*gaia.AppCredential, error) {

	// Then we generate a private key and a CSR from the appcred info.
	csr, pk, err := makeCSR(creds.Name, creds.ID, creds.Namespace)
	if err != nil {
		return nil, err
	}

	// And we update the appcred with the csr
	creds.CSR = string(csr)

	if err = m.Update(
		manipulate.NewContext(
			ctx,
			manipulate.ContextOptionNamespace(creds.Namespace),
		),
		creds,
	); err != nil {
		return nil, err
	}

	// And we write the private key in the appcred.
	creds.Credentials.CertificateKey = base64.StdEncoding.EncodeToString(pk)

	return creds, nil
}

func makeCSR(name string, id string, namespace string) (csr []byte, key []byte, err error) {

	pk, err := tglib.ECPrivateKeyGenerator()
	if err != nil {
		return nil, nil, err
	}

	csr, err = tglib.GenerateSimpleCSR([]string{namespace}, nil, "app:credential:"+id+":"+name, nil, pk)
	if err != nil {
		return nil, nil, err
	}

	keyBlock, err := tglib.KeyToPEM(pk)
	if err != nil {
		return nil, nil, err
	}

	return csr, pem.EncodeToMemory(keyBlock), nil
}
