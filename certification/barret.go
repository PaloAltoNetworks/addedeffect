package certification

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/aporeto-inc/addedeffect/discovery"
	"github.com/aporeto-inc/gaia/barretmodels/v1/golang"
	"github.com/aporeto-inc/manipulate"
	"github.com/aporeto-inc/manipulate/maniphttp"
	"github.com/aporeto-inc/tg/tglib"
	"go.uber.org/zap"

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

// SignerType is the type of signer name.
type SignerType int

// Various possible SignerType values
const (
	SignerTypeSystem SignerType = iota + 1
	SignerTypePublic
)

// IssueCert asks and returns an new certificate using the given barret Manipulator and given CSR.
// You can generate easily a CSR using GenerateSimpleCSR.
func IssueCert(m manipulate.Manipulator, csrPEM []byte, expiration time.Time, usage KeyUsage, signerType SignerType, span opentracing.Span) (cert []byte, serialNumber string, exp time.Time, err error) {

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
	request.Signer = convertSignerType(signerType)

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

// IssueServiceClientCertificate is used internally to get a client certificates for a service.
func IssueServiceClientCertificate(m manipulate.Manipulator, serviceName string, validity time.Duration) (*tls.Certificate, error) {

	privateKey, err := tglib.ECPrivateKeyGenerator()
	if err != nil {
		return nil, fmt.Errorf("client: Unable to generate private key: %s", err)
	}

	privateKeyPEM, err := tglib.KeyToPEM(privateKey)
	if err != nil {
		return nil, fmt.Errorf("client: Unable to convert private key to PEM: %s", err)
	}

	csr, err := tglib.GenerateSimpleCSR([]string{"system"}, []string{"root"}, serviceName, nil, privateKey)
	if err != nil {
		return nil, fmt.Errorf("client: Unable to prepare certificate request: %s", err)
	}

	clientCert, _, _, err := IssueCert(m, csr, time.Now().Add(validity), KeyUsageClient, SignerTypeSystem, nil)
	if err != nil {
		return nil, fmt.Errorf("client: Unable to get new certificate: %s", err)
	}

	X509CertKeyPairs, err := tls.X509KeyPair(clientCert, pem.EncodeToMemory(privateKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("client: Cannot load newly generated client certificates: %s", err)
	}

	return &X509CertKeyPairs, nil
}

// IssueServiceServerCertificate is used internally to get a server certificates for a service.
func IssueServiceServerCertificate(m manipulate.Manipulator, serviceName string, dns []string, ips []string, validity time.Duration) (*tls.Certificate, error) {

	privateKey, err := tglib.ECPrivateKeyGenerator()
	if err != nil {
		return nil, fmt.Errorf("server: Unable to generate private key: %s", err)
	}

	privateKeyPEM, err := tglib.KeyToPEM(privateKey)
	if err != nil {
		return nil, fmt.Errorf("server: Unable to convert private key to PEM: %s", err)
	}

	var parsedIPs []net.IP
	for _, ip := range ips {
		parsedIPs = append(parsedIPs, net.ParseIP(ip))
	}

	csr := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:         serviceName,
			Organization:       []string{"aporeto"},
			OrganizationalUnit: []string{"service"},
		},
		DNSNames:           dns,
		IPAddresses:        parsedIPs,
		SignatureAlgorithm: x509.ECDSAWithSHA384,
	}

	csrData, err := tglib.GenerateCSR(csr, privateKey)
	if err != nil {
		return nil, fmt.Errorf("server: Unable to prepare certificate request: %s", err)
	}

	serverCert, _, _, err := IssueCert(m, csrData, time.Now().Add(validity), KeyUsageServer, SignerTypeSystem, nil)
	if err != nil {
		return nil, fmt.Errorf("server: Unable to get new certificate: %s", err)
	}

	X509ServerCert, err := tls.X509KeyPair(serverCert, pem.EncodeToMemory(privateKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("server: Cannot load newly generated certificates: %s", err)
	}

	leaf, _ := pem.Decode(serverCert)
	X509ServerCert.Leaf, err = x509.ParseCertificate(leaf.Bytes)
	if err != nil {
		return nil, fmt.Errorf("server: Unable to reparse x509 certificate leaf: %s", err)
	}

	return &X509ServerCert, nil
}

// BuildCertificatesMaps returns to maps to get what certificate to use for which DNS or IPs.
// This can be used in a custom tls.Config.GetCertificate function.
func BuildCertificatesMaps(certs []tls.Certificate) (map[string]*tls.Certificate, map[string]*tls.Certificate, error) {

	certsNamesMap := map[string]*tls.Certificate{}
	certsIPsMap := map[string]*tls.Certificate{}

	for _, item := range certs {
		for _, subItem := range item.Certificate {
			x509Cert, err := x509.ParseCertificate(subItem)
			if err != nil {
				return nil, nil, err
			}
			for _, dns := range x509Cert.DNSNames {
				certsNamesMap[dns] = &item
			}
			for _, ip := range x509Cert.IPAddresses {
				certsIPsMap[ip.String()] = &item
			}
		}
	}

	return certsNamesMap, certsIPsMap, nil
}

// MakeRenewServiceServerCertificateFunc returns a function that will renew the certificate if needed. This can be used as TLSConfig.GetCertificate func.
// Internally, it uses IssueServiceServerCertificate.
func MakeRenewServiceServerCertificateFunc(
	m manipulate.Manipulator,
	serviceName string,
	dns []string,
	ips []string,
	validity time.Duration,
	additionalCertificates []tls.Certificate,
) (func(*tls.ClientHelloInfo) (*tls.Certificate, error), error) {

	cert, err := IssueServiceServerCertificate(m, serviceName, dns, ips, validity)
	if err != nil {
		return nil, err
	}

	lock := &sync.Mutex{}

	certsNameMap, certsIPsMap, err := BuildCertificatesMaps(additionalCertificates)
	if err != nil {
		return nil, err
	}

	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {

		if ac := certsNameMap[hello.ServerName]; ac != nil {
			return ac, nil
		}

		host, _, _ := net.SplitHostPort(hello.Conn.LocalAddr().String())
		if ac := certsIPsMap[host]; ac != nil {
			return ac, nil
		}

		lock.Lock()
		defer lock.Unlock()

		if time.Now().Add(time.Hour).After(cert.Leaf.NotAfter) {

			renewedCerts, err := IssueServiceServerCertificate(m, serviceName, dns, ips, validity)
			if err != nil {
				return nil, err
			}

			if err := RevokeCert(m, renewedCerts.Leaf.SerialNumber.String(), true, nil); err != nil {
				return nil, err
			}

			cert = renewedCerts
		}

		return cert, nil
	}, nil
}

// CreateServiceCertificates is a helper func that can be used during service initialization.
func CreateServiceCertificates(
	serviceName string,
	rootCAPool *x509.CertPool,
	pf *discovery.PlatformInfo,
	password string,
	getClientCert bool,
	getServerCertFunc bool,
	dns []string,
	ips []string,
	additionalCertKeyPass string,
) (clientCert *tls.Certificate, serverCertFunc func(*tls.ClientHelloInfo) (*tls.Certificate, error)) {

	issuingCertKeyPair, err := pf.IssuingServiceClientCertPair(password)
	if err != nil {
		zap.L().Fatal("Unable to decode issuing certificate key pair", zap.Error(err))
	}

	issuingManipulator := maniphttp.NewHTTPManipulatorWithMTLS(pf.BarretURL, "", rootCAPool, []tls.Certificate{issuingCertKeyPair}, false)

	if getClientCert {
		clientCert, err = IssueServiceClientCertificate(issuingManipulator, serviceName, 8760*time.Hour)
		if err != nil {
			zap.L().Fatal("Unable to retrieve client certificate key pair", zap.Error(err))
		}
		zap.L().Info("Client certificate issued", zap.String("name", serviceName))

	}

	if getServerCertFunc {
		var additionalServerCertificates []tls.Certificate
		if pf.PublicServicesCert != "" && additionalCertKeyPass != "" {
			pcert, e := pf.PublicServicesCertPair(additionalCertKeyPass)
			if e != nil {
				zap.L().Fatal("Unable to decrypt public certs key pair", zap.Error(e))
			}

			additionalServerCertificates = append(additionalServerCertificates, pcert)
		}

		serverCertFunc, err = MakeRenewServiceServerCertificateFunc(issuingManipulator, serviceName, dns, ips, 4380*time.Hour, additionalServerCertificates)
		if err != nil {
			zap.L().Fatal("Unable to retrieve server certificate key pair",
				zap.Strings("dns", dns),
				zap.Strings("ips", ips),
				zap.Error(err))
		}
		zap.L().Info("Server certificate issued", zap.Strings("dns", dns), zap.Strings("ips", ips))
	}

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

func convertSignerType(signerType SignerType) barretmodels.CertificateSignerValue {
	switch signerType {
	case SignerTypeSystem:
		return barretmodels.CertificateSignerSystem
	default:
		return barretmodels.CertificateSignerPublic
	}
}
