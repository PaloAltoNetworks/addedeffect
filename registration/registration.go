package registration

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"

	"github.com/aporeto-inc/gaia/squallmodels/v1/golang"
	"github.com/aporeto-inc/manipulate"
	"go.uber.org/zap"
)

// RegisterEnforcer registers a new enforcer server with given name, description and tags in Squall using the given Manipulator.
func RegisterEnforcer(
	manipulator manipulate.Manipulator,
	namespace string,
	name string,
	fqdn string,
	description string,
	tags []string,
	folderPath string,
	certificateName string,
	keyName string,
	certificateExpirationDate time.Time,
	deleteIfExist bool,
) (*squallmodels.Enforcer, error) {

	enforcer := squallmodels.NewEnforcer()
	enforcer.Name = name
	enforcer.FQDN = fqdn
	enforcer.Description = description
	enforcer.AssociatedTags = tags
	enforcer.LastSyncTime = time.Now().Add(1 - time.Hour)

	mctx := manipulate.NewContext()
	mctx.Parameters.KeyValues.Add("tag", "$name="+enforcer.Name)

	var err error
	enforcers := squallmodels.EnforcersList{}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for i := 0; i < 12; i++ {

		err = manipulate.Retry(func() error { return manipulator.RetrieveMany(mctx, &enforcers) }, nil, 5)

		if err == nil {
			break
		}

		zap.L().Warn("Unable to register. Retrying in 5s", zap.Error(err))

		select {
		case <-time.After(5 * time.Second):
		case <-c:
			return nil, manipulate.NewErrDisconnected("Disconnected per signal")
		}
	}

	if err != nil {
		return nil, fmt.Errorf("Unable to access servers list: %s", err)
	}

	// Check if the server already exists and delete if deleteIFExist flag set
	if len(enforcers) > 0 {
		if !deleteIfExist {
			return nil, fmt.Errorf("A server with the name %s already exists", enforcer.Name)
		}

		for _, existingEnforcer := range enforcers {
			if err := manipulate.Retry(func() error { return manipulator.Delete(mctx, existingEnforcer) }, nil, 5); err != nil {
				return nil, fmt.Errorf("Unable to delete enforcer %s that already exists: %s", enforcer.Name, err)
			}
		}
	}

	if err := manipulate.Retry(func() error { return manipulator.Create(nil, enforcer) }, nil, 5); err != nil {
		return nil, err
	}

	certData := []byte(fmt.Sprintf("%s\n", enforcer.Certificate))
	keyData := []byte(fmt.Sprintf("%s\n", enforcer.CertificateKey))

	if err := writeCertificate(folderPath, certificateName, keyName, 0700, 0600, certData, keyData); err != nil {
		return nil, err
	}

	return enforcer, nil
}

// ServerInfoFromCertificate retrieves and verifies the enforcerID and namespace stored in the
// certificate at the given path using the given x509.CertPool.
func ServerInfoFromCertificate(certPath string, CAPool *x509.CertPool) (string, string, []byte, error) {

	b, err := ioutil.ReadFile(certPath)
	if err != nil {
		return "", "", nil, err
	}

	var cert *x509.Certificate
	intermediates := x509.NewCertPool()
	block, rest := pem.Decode(b)

	var caChain []byte

	for ; block != nil; block, rest = pem.Decode(rest) {
		if block.Type == "CERTIFICATE" {
			crt, e := x509.ParseCertificate(block.Bytes)
			if e != nil {
				return "", "", nil, e
			}
			if !crt.IsCA {
				cert = crt
				if caChain == nil {
					caChain = rest
				}
				continue
			}
			if len(rest) != 0 {
				intermediates.AddCert(crt)
			}
		} else {
			return "", "", nil, fmt.Errorf("Invalid pem block type: %s", block.Type)
		}
	}

	_, err = cert.Verify(x509.VerifyOptions{Roots: CAPool, Intermediates: intermediates})
	if err != nil {
		return "", "", nil, err
	}

	if len(cert.Subject.OrganizationalUnit) == 0 {
		return "", "", nil, fmt.Errorf("Missing Organizational Unit field")
	}

	parts := strings.SplitN(cert.Subject.CommonName, "@", 2)
	enforcerID := parts[0]
	namespace := parts[1]

	if err != nil {
		return "", "", nil, err
	}

	return enforcerID, namespace, caChain, nil
}

func writeCertificate(folder, certName, keyName string, folderPerm os.FileMode, certPerm os.FileMode, certData, keyData []byte) error {

	if err := os.MkdirAll(folder, folderPerm); err != nil {
		return fmt.Errorf("Unable to create %s: %s", folder, err)
	}

	certPath := path.Join(folder, certName)
	if err := ioutil.WriteFile(certPath, certData, certPerm); err != nil {
		return fmt.Errorf("Unable to write certificate to %s: %s", certPath, err)
	}

	keyPath := path.Join(folder, keyName)
	if err := ioutil.WriteFile(keyPath, keyData, certPerm); err != nil {
		return fmt.Errorf("Unable to write key to %s: %s", keyPath, err)
	}

	return nil
}
