package certification

import (
	"time"

	"github.com/aporeto-inc/gaia/barretmodels/v1/golang"
	"github.com/aporeto-inc/manipulate"

	opentracing "github.com/opentracing/opentracing-go"
)

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
