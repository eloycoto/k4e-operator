package mtls_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jakub-dzon/k4e-operator/internal/mtls"
)

const (
	certRegisterCN = "register" // Important, make a copy here to prevent breaking changes
)

var _ = Describe("MTLS CA test", func() {

	Context("VerifyRequest", func() {
		var (
			ca         []*certificate
			CACertPool *x509.CertPool
			CAChain    []*x509.Certificate
			opts       x509.VerifyOptions
		)

		BeforeEach(func() {
			ca = []*certificate{createCACert(), createCACert()}

			CACertPool = x509.NewCertPool()
			CAChain = []*x509.Certificate{}

			for _, cert := range ca {
				CACertPool.AddCert(cert.signedCert)
				CAChain = append(CAChain, cert.signedCert)
			}

			opts = x509.VerifyOptions{
				Roots:         CACertPool,
				Intermediates: x509.NewCertPool(),
			}
		})

		It("No peer certificates are present", func() {
			// given
			r := &http.Request{
				TLS: &tls.ConnectionState{
					PeerCertificates: []*x509.Certificate{},
				},
			}

			// when
			res := mtls.VerifyRequest(r, 0, opts, CAChain)

			// then
			Expect(res).To(BeFalse())
		})

		Context("Registration Auth", func() {
			const (
				AuthType = 1 // Equals to YggdrasilRegisterAuth, but it's important, so keep a copy here.
			)

			It("Peer certificate is valid", func() {
				// given
				cert := createRegistrationClientCert(ca[0])
				r := &http.Request{
					TLS: &tls.ConnectionState{
						PeerCertificates: []*x509.Certificate{cert.signedCert},
					},
				}

				// when
				res := mtls.VerifyRequest(r, AuthType, opts, CAChain)

				// then
				Expect(res).To(BeTrue())
			})

			It("Peer certificate is invalid", func() {
				// given
				cert := createRegistrationClientCert(createCACert())
				r := &http.Request{
					TLS: &tls.ConnectionState{
						PeerCertificates: []*x509.Certificate{cert.signedCert},
					},
				}

				// when
				res := mtls.VerifyRequest(r, AuthType, opts, CAChain)

				// then
				Expect(res).To(BeFalse())
			})

			It("Lastet CA certificate is valid", func() {
				// given
				cert := createRegistrationClientCert(ca[1])
				r := &http.Request{
					TLS: &tls.ConnectionState{
						PeerCertificates: []*x509.Certificate{cert.signedCert},
					},
				}

				// when
				res := mtls.VerifyRequest(r, AuthType, opts, CAChain)

				// then
				Expect(res).To(BeTrue())
			})

			It("Expired certificate is valid", func() {
				// given
				c := &x509.Certificate{
					SerialNumber: big.NewInt(time.Now().Unix()),
					Subject: pkix.Name{
						Organization: []string{"K4e-operator"},
					},
					NotBefore:             time.Now(),
					NotAfter:              time.Now().AddDate(0, 0, 0),
					ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
					KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
					BasicConstraintsValid: true,
				}

				cert := createGivenClientCert(c, ca[1])
				r := &http.Request{
					TLS: &tls.ConnectionState{
						PeerCertificates: []*x509.Certificate{cert.signedCert},
					},
				}

				// when
				res := mtls.VerifyRequest(r, AuthType, opts, CAChain)

				// then
				Expect(res).To(BeTrue())
			})

		})

		Context("Normal device Auth", func() {
			const (
				AuthType = 0 // Equals to YggdrasilCompleteAuth, but it's important, so keep a copy here.
			)

			It("Register certificate is invalid", func() {
				// given
				cert := createRegistrationClientCert(ca[0])
				r := &http.Request{
					TLS: &tls.ConnectionState{
						PeerCertificates: []*x509.Certificate{cert.signedCert},
					},
				}

				// when
				res := mtls.VerifyRequest(r, AuthType, opts, CAChain)

				// then
				Expect(res).To(BeFalse())
			})

			It("Certificate is correct", func() {
				// given
				cert := createClientCert(ca[0])
				r := &http.Request{
					TLS: &tls.ConnectionState{
						PeerCertificates: []*x509.Certificate{cert.signedCert},
					},
				}

				// when
				res := mtls.VerifyRequest(r, AuthType, opts, CAChain)

				// then
				Expect(res).To(BeTrue())
			})

			It("Invalid certificate is correct", func() {

				// given
				cert := createClientCert(createCACert())
				r := &http.Request{
					TLS: &tls.ConnectionState{
						PeerCertificates: []*x509.Certificate{cert.signedCert},
					},
				}

				// when
				res := mtls.VerifyRequest(r, AuthType, opts, CAChain)

				// then
				Expect(res).To(BeFalse())
			})

			It("Certificate valid with any CA position on the store.", func() {
				// given
				cert := createClientCert(ca[1])
				r := &http.Request{
					TLS: &tls.ConnectionState{
						PeerCertificates: []*x509.Certificate{cert.signedCert},
					},
				}

				// when
				res := mtls.VerifyRequest(r, AuthType, opts, CAChain)

				// then
				Expect(res).To(BeTrue())
			})

			It("Expired certificate is not working", func() {
				// given
				c := &x509.Certificate{
					SerialNumber: big.NewInt(time.Now().Unix()),
					Subject: pkix.Name{
						Organization: []string{"K4e-operator"},
						CommonName:   "test-device",
					},
					NotBefore:             time.Now(),
					NotAfter:              time.Now().AddDate(0, 0, 0),
					ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
					KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
					BasicConstraintsValid: true,
				}
				cert := createGivenClientCert(c, ca[0])

				r := &http.Request{
					TLS: &tls.ConnectionState{
						PeerCertificates: []*x509.Certificate{cert.signedCert},
					},
				}

				// when
				res := mtls.VerifyRequest(r, AuthType, opts, CAChain)

				// then
				Expect(res).To(BeFalse())
			})

		})

	})
})

type certificate struct {
	cert       *x509.Certificate
	key        *rsa.PrivateKey
	certBytes  []byte
	signedCert *x509.Certificate
}

func createRegistrationClientCert(ca *certificate) *certificate {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Organization: []string{"K4e-operator"},
			CommonName:   certRegisterCN,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	return createGivenClientCert(cert, ca)
}

func createClientCert(ca *certificate) *certificate {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Organization: []string{"K4e-operator"},
			CommonName:   "device-UUID",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	return createGivenClientCert(cert, ca)
}

func createGivenClientCert(cert *x509.Certificate, ca *certificate) *certificate {
	certKey, err := rsa.GenerateKey(rand.Reader, 1024)
	ExpectWithOffset(1, err).To(BeNil(), "Fail on key generation")

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca.cert, &certKey.PublicKey, ca.key)
	ExpectWithOffset(1, err).To(BeNil(), "Fail on sign generation")

	signedCert, err := x509.ParseCertificate(certBytes)
	ExpectWithOffset(1, err).To(BeNil(), "Fail on parsing certificate")

	err = signedCert.CheckSignatureFrom(ca.signedCert)
	ExpectWithOffset(1, err).To(BeNil(), "Fail on check signature")

	return &certificate{cert, certKey, certBytes, signedCert}
}

func createCACert() *certificate {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Organization: []string{"K4e-operator"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 1024)
	ExpectWithOffset(1, err).To(BeNil(), "Fail on key generation")

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	ExpectWithOffset(1, err).To(BeNil(), "Fail on sign generation")

	signedCert, err := x509.ParseCertificate(caBytes)
	ExpectWithOffset(1, err).To(BeNil(), "Fail on parsing certificate")

	return &certificate{ca, caPrivKey, caBytes, signedCert}
}
