package registration

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aporeto-inc/manipulate"
	"github.com/aporeto-inc/trireme/crypto"

	squallmodels "github.com/aporeto-inc/gaia/squallmodels/current/golang"
)

// RegisterEnforcer registers a new enforcer server with given name, description and tags in Squall using the given Manipulator.
func RegisterEnforcer(
	manipulator manipulate.Manipulator,
	namespace string,
	serverName string,
	serverFQDN string,
	serverDescription string,
	serverTags []string,
	folderPath string,
	certificateName string,
	keyName string,
	certificateExpirationDate time.Time,
) (*squallmodels.Enforcer, error) {

	server := squallmodels.NewEnforcer()
	server.Name = serverName
	server.FQDN = serverFQDN
	server.Description = serverDescription
	server.AssociatedTags = serverTags
	server.LastSyncTime = time.Now()

	// Check if the server already exists
	mctx := manipulate.NewContext()
	mctx.Parameters.KeyValues.Add("tag", "$name="+server.Name)

	if n, err := manipulator.Count(mctx, squallmodels.EnforcerIdentity); err != nil || n > 0 {
		if err != nil {
			return nil, fmt.Errorf("Unable to access servers list. Does the namespace exist? Do you have the correct permissions?")
		}

		return nil, fmt.Errorf("A server with the name %s already exists", server.Name)
	}

	if err := manipulator.Create(nil, server); err != nil {
		return nil, err
	}

	certData := []byte(fmt.Sprintf("%s\n", server.Certificate))
	keyData := []byte(fmt.Sprintf("%s\n", server.CertificateKey))

	if err := writeCertificate(folderPath, certificateName, keyName, 0700, 0600, certData, keyData); err != nil {
		return nil, err
	}

	return server, nil
}

// ServerInfoFromCertificate retrieves and verifies the enforcerID and namespace stored in the
// certificate at the given path using the given x509.CertPool.
func ServerInfoFromCertificate(certPath string, CAPool *x509.CertPool) (string, string, error) {

	certificate, err := ioutil.ReadFile(certPath)
	if err != nil {
		return "", "", err
	}

	cert, err := crypto.LoadAndVerifyCertificate(certificate, CAPool)
	if err != nil {
		return "", "", err
	}

	if len(cert.Subject.OrganizationalUnit) == 0 {
		return "", "", fmt.Errorf("Missing Organizational Unit field")
	}

	parts := strings.SplitN(cert.Subject.CommonName, "@", 2)
	enforcerID := parts[0]
	namespace := parts[1]

	if err != nil {
		return "", "", err
	}

	return enforcerID, namespace, nil
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
