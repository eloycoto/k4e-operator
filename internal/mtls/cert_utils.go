package mtls

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

// https://shaneutt.com/blog/golang-ca-and-signed-cert-go/
type CertificateGroup struct {
	cert       *x509.Certificate
	privKey    *rsa.PrivateKey
	certBytes  []byte
	certPEM    *bytes.Buffer
	PrivKeyPEM *bytes.Buffer
}

func NewCertificateGroupFromCACM(configMap map[string][]byte) (*CertificateGroup, error) {
	certGroup := &CertificateGroup{
		certPEM:    bytes.NewBuffer(configMap["ca.crt"]),
		PrivKeyPEM: bytes.NewBuffer(configMap["ca.key"]),
	}

	block, _ := pem.Decode(certGroup.certPEM.Bytes())
	if block == nil {
		return nil, fmt.Errorf("Cannot get CA certificate")
	}
	ca, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Failing parsing cert: %v", err)
	}
	block, _ = pem.Decode(certGroup.PrivKeyPEM.Bytes())
	if block == nil {
		return nil, fmt.Errorf("Cannot get CA certificate key")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failing parsing key: %v", err)
	}

	certGroup.cert = ca
	certGroup.privKey = key
	return certGroup, nil
}

// CreatePem from the load certificates create the PEM file and stores in local
func (c *CertificateGroup) CreatePem() {

	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: c.certBytes,
	})

	c.certPEM = caPEM

	privKeyPEM := new(bytes.Buffer)
	pem.Encode(privKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(c.privKey),
	})
	c.PrivKeyPEM = privKeyPEM
}

// GetCertificate returns the certificate Group in tls.Certficicate format.
func (c *CertificateGroup) GetCertificate() (tls.Certificate, error) {
	return tls.X509KeyPair(c.certPEM.Bytes(), c.PrivKeyPEM.Bytes())
}

func (c *CertificateGroup) GetCert() *x509.Certificate {
	return c.cert
}

func getCACertificate() (*CertificateGroup, error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
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

	// create our private and public key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("Cannot generate CA Key")
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, err
	}

	certificateBundle := CertificateGroup{
		cert:      ca,
		privKey:   caPrivKey,
		certBytes: caBytes,
	}
	certificateBundle.CreatePem()
	return &certificateBundle, nil
}

func getKeyAndCSR(cert *x509.Certificate, caCert *CertificateGroup) (*CertificateGroup, error) {

	certKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("Cannot generate cert Key")
	}

	// sign the cert by the CA
	certBytes, err := x509.CreateCertificate(
		rand.Reader, cert, caCert.cert, &certKey.PublicKey, caCert.privKey)
	if err != nil {
		return nil, err
	}

	certificateBundle := CertificateGroup{
		cert:      cert,
		privKey:   certKey,
		certBytes: certBytes,
	}
	certificateBundle.CreatePem()
	return &certificateBundle, nil
}

func getServerCertificate(dnsNames []string, localhostEnabled bool, CACert *CertificateGroup) (*CertificateGroup, error) {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			CommonName:    "*",
			Organization:  []string{"K4e-agent"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{""},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		DNSNames:     dnsNames,
		IPAddresses:  []net.IP{},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	return getKeyAndCSR(cert, CACert)

	// certKey, err := rsa.GenerateKey(rand.Reader, 4096)
	// if err != nil {
	// 	return nil, fmt.Errorf("Cannot generate cert Key")
	// }

	// // sign the cert by the CA
	// certBytes, err := x509.CreateCertificate(
	// 	rand.Reader, cert, caCert.cert, &certKey.PublicKey, caCert.privKey)
	// if err != nil {
	// 	return nil, err
	// }

	// certificateBundle := CertificateGroup{
	// 	cert:      cert,
	// 	privKey:   certKey,
	// 	certBytes: certBytes,
	// }
	// certificateBundle.CreatePem()
	// return &certificateBundle, nil
}
