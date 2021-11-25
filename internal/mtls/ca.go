package mtls

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/jakub-dzon/k4e-operator/internal/yggdrasil"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	regClientSecretNameRandomLen = 10
	regClientSecretNamePrefix    = "reg-client-ca"
	regClientSecretLabelKey      = "reg-client-ca"
)

type CAProvider interface {
	GetCACertificate() (map[string][]byte, error)
	CreateRegistrationCertificate(name string) (map[string][]byte, error)
}

type TLSConfig struct {
	config           *tls.Config
	client           client.Client
	caProvider       []CAProvider
	Domains          []string
	LocalhostEnabled bool
	namespace        string
}

func NewMTLSconfig(client client.Client, namespace string, domains []string, localhostEnabled bool) *TLSConfig {
	config := &TLSConfig{
		config:           nil,
		client:           client,
		Domains:          domains,
		namespace:        namespace,
		LocalhostEnabled: localhostEnabled,
	}

	// Secret providers here
	secretProvider := NewCASecretProvider(client, namespace)
	config.caProvider = append(config.caProvider, secretProvider)

	return config
}

// @TODO mainly used for testing, maybe not needed at all
func (conf *TLSConfig) SetCAProvider(caProviders []CAProvider) {
	conf.caProvider = caProviders
}

func (conf *TLSConfig) InitCertificates() (*tls.Config, []*x509.Certificate, error) {
	if len(conf.caProvider) == 0 {
		return nil, nil, fmt.Errorf("No provider set")
	}

	result := []map[string][]byte{}
	for _, caProvider := range conf.caProvider {
		caCerts, err := caProvider.GetCACertificate()
		if err != nil {
			// TODO do something here like multipleerror
			return nil, nil, err
		}
		result = append(result, caCerts)
	}

	CACertChain := []*x509.Certificate{}
	var certGroup *CertificateGroup
	var err error

	caCertPool := x509.NewCertPool()

	for _, provider := range result {
		certGroup, err = NewCertificateGroupFromCACM(provider)
		if err != nil {
			continue
		}
		CACertChain = append(CACertChain, certGroup.GetCert())
		caCertPool.AppendCertsFromPEM(certGroup.certPEM.Bytes())
	}

	if certGroup == nil {
		return nil, nil, fmt.Errorf("Cannot get CA certificate")
	}

	serverCert, err := getServerCertificate(conf.Domains, conf.LocalhostEnabled, certGroup)
	if err != nil {
		return nil, nil, err
	}

	certificate, err := serverCert.GetCertificate()
	if err != nil {
		return nil, nil, fmt.Errorf("Cannot create server certfificate: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAnyClientCert,
		// ClientAuth: tls.NoClientCert,
		MinVersion: tls.VersionTLS13,
	}
	return tlsConfig, CACertChain, nil
}

func (conf *TLSConfig) CreateRegistrationClient() error {

	name := fmt.Sprintf("%s-%s",
		regClientSecretNamePrefix,
		utilrand.String(regClientSecretNameRandomLen))

	if len(conf.caProvider) == 0 {
		return fmt.Errorf("Cannot get ca provider")
	}

	certData, err := conf.caProvider[0].CreateRegistrationCertificate(name)
	if err != nil {
		return fmt.Errorf("Cannot create client certificate")
	}

	secret := corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Namespace: conf.namespace,
			Name:      name,
			Labels:    map[string]string{regClientSecretLabelKey: "true"},
		},
		Data: certData,
	}

	err = conf.client.Create(context.TODO(), &secret)
	return err
}

func IsClientCertificateSigned(PeerCertificates []*x509.Certificate, CAChain []*x509.Certificate) bool {
	for _, cert := range PeerCertificates {
		certValid := false
		for _, caCert := range CAChain {

			err := cert.CheckSignatureFrom(caCert)
			if err == nil {
				certValid = true
				break
			}
		}
		if !certValid {
			return false
		}
	}
	return true
}

func VerifyRequest(r *http.Request, verifyType int, verifyOpts x509.VerifyOptions, CACertChain []*x509.Certificate) bool {
	if verifyType == yggdrasil.YggdrasilRegisterAuth {
		res := IsClientCertificateSigned(r.TLS.PeerCertificates, CACertChain)
		return res
	}

	valid := true
	for _, cert := range r.TLS.PeerCertificates {
		if cert.Subject.CommonName == "registered" {
			valid = false
		}
		if _, err := cert.Verify(verifyOpts); err != nil {
			return false
		}
	}

	return valid
}
