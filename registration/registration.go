package registration

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/aporeto-inc/manipulate"
	"github.com/aporeto-inc/trireme/crypto"

	gaia "github.com/aporeto-inc/gaia/golang"
	uuid "github.com/satori/go.uuid"
)

const serverAPIAuthorizationName = "[auto] aporeto agent api authorizations"

// RegisterAgent registers a new agent server with given name, description and tags in Squall using the given Manipulator.
func RegisterAgent(
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
) (*gaia.Server, error) {

	server := gaia.NewServer()
	server.Name = serverName
	server.FQDN = serverFQDN
	server.Description = serverDescription
	server.AssociatedTags = serverTags
	server.LastSyncTime = time.Now()

	// Check if the server already exists
	mctx := manipulate.NewContext()
	mctx.Parameters.KeyValues["tag"] = "$name=" + server.Name

	if n, err := manipulator.Count(mctx, gaia.ServerIdentity); err != nil || n > 0 {
		if err != nil {
			return nil, fmt.Errorf("Unable to access servers list. Does the namespace exist? Do you have the correct permissions?")
		}

		return nil, fmt.Errorf("A server with the name %s already exists.", server.Name)
	}

	if err := manipulator.Create(nil, server); err != nil {
		return nil, err
	}

	policies := gaia.APIAuthorizationPoliciesList{}
	mctx = manipulate.NewContextWithFilter(manipulate.NewFilterComposer().WithKey("name").Equals(serverAPIAuthorizationName).Done())

	if err := manipulator.RetrieveMany(mctx, gaia.APIAuthorizationPolicyIdentity, &policies); err != nil {
		return nil, err
	}

	if len(policies) == 0 {
		APIAuthorizationPolicy := gaia.NewAPIAuthorizationPolicy()
		APIAuthorizationPolicy.Description = "Auto generated policy to allow agents to communicate with Squall."
		APIAuthorizationPolicy.Subject = [][]string{[]string{"auth:commonname=" + server.ID}}
		APIAuthorizationPolicy.AuthorizedNamespace = namespace
		APIAuthorizationPolicy.AuthorizedIdentities = []string{"*"}
		APIAuthorizationPolicy.Name = serverAPIAuthorizationName
		APIAuthorizationPolicy.Propagate = true
		APIAuthorizationPolicy.AllowsDelete = true
		APIAuthorizationPolicy.AllowsGet = true
		APIAuthorizationPolicy.AllowsHead = true
		APIAuthorizationPolicy.AllowsPatch = true
		APIAuthorizationPolicy.AllowsPost = true
		APIAuthorizationPolicy.AllowsPut = true

		if err := manipulator.Create(nil, APIAuthorizationPolicy); err != nil {
			return nil, err
		}
	} else {
		APIAuthorizationPolicy := policies[0]
		APIAuthorizationPolicy.Subject = append(APIAuthorizationPolicy.Subject, []string{"auth:commonname=" + server.ID})

		if err := manipulator.Update(nil, APIAuthorizationPolicy); err != nil {
			return nil, err
		}
	}

	certData := []byte(fmt.Sprintf("%s\n", server.Certificate))
	keyData := []byte(fmt.Sprintf("%s\n", server.CertificateKey))

	if err := writeCertificate(folderPath, certificateName, keyName, 0700, 0600, certData, keyData); err != nil {
		return nil, err
	}

	return server, nil
}

// ServerInfoFromCertificate retrieves and verifies the serverID and namespace stored in the
// certificate at the given path using the given x509.CertPool.
func ServerInfoFromCertificate(certPath string, CAPool *x509.CertPool) (uuid.UUID, string, error) {

	certificate, err := ioutil.ReadFile(certPath)
	if err != nil {
		return uuid.UUID{}, "", err
	}

	cert, err := crypto.LoadAndVerifyCertificate(certificate, CAPool)
	if err != nil {
		return uuid.UUID{}, "", err
	}

	if len(cert.Subject.OrganizationalUnit) == 0 {
		return uuid.UUID{}, "", fmt.Errorf("Missing Organizational Unit field.")
	}

	serverID := uuid.FromStringOrNil(cert.Subject.CommonName)
	namespace := cert.Subject.OrganizationalUnit[0]

	if err != nil {
		return uuid.UUID{}, "", err
	}

	return serverID, namespace, nil
}

// SendServerHeartBeat sends a heartbeat message for the given server.
func SendServerHeartBeat(manipulator manipulate.Manipulator, server *gaia.Server, t time.Time) error {

	if err := manipulate.RetryManipulation(func() error { return manipulator.Retrieve(nil, server) }, nil, 10); err != nil {
		return err
	}

	server.LastSyncTime = t
	return manipulate.RetryManipulation(func() error { return manipulator.Update(nil, server) }, nil, 10)
}

// RetrieveServerProfile retrieves the profile to use according to the given serverID.
func RetrieveServerProfile(manipulator manipulate.Manipulator, serverID string) (*gaia.ServerProfile, error) {

	server := gaia.NewServer()
	server.ID = serverID

	profile := gaia.ServerProfilesList{}

	ctx := manipulate.NewContext()
	ctx.Parent = server

	if err := manipulator.RetrieveMany(ctx, gaia.ServerProfileIdentity, &profile); err != nil {
		return nil, err
	}

	if len(profile) == 0 {
		return nil, fmt.Errorf("Could not find any profile")
	}

	return profile[0], nil
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
