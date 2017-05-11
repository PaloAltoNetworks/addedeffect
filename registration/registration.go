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
	enforcer.LastSyncTime = time.Now().Add(-1 * time.Hour)

	mctx := manipulate.NewContext()
	mctx.Parameters.KeyValues.Add("tag", "$name="+enforcer.Name)

	n := 0
	err := manipulate.RetryManipulation(func() error {
		var err error
		n, err = manipulator.Count(mctx, squallmodels.EnforcerIdentity)
		return err
	}, nil, 15)
	if err != nil {
		return nil, fmt.Errorf("Unable to access servers list. Does the namespace exist? Do you have the correct permissions?")
	}

	// Check if the server already exists and delete if deleteIFExist flag set
	if n > 0 {
		if !deleteIfExist {
			return nil, fmt.Errorf("A server with the name %s already exists", enforcer.Name)
		}

		if err := manipulator.Delete(mctx, enforcer); err != nil {
			return nil, fmt.Errorf("Unable to delete enforcer %s that already exists", enforcer.Name)
		}
	}

	if err := manipulator.Create(nil, enforcer); err != nil {
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
