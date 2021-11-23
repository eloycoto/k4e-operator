package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CAProvider interface {
	GetCACertificate() (map[string][]byte, error)
}

type TLSConfig struct {
	config           *tls.Config
	client           client.Client
	caProvider       []CAProvider
	Domains          []string
	LocalhostEnabled bool
}

func NewMTLSconfig(client client.Client, namespace string, domains []string, localhostEnabled bool) *TLSConfig {
	config := &TLSConfig{
		config:           nil,
		client:           client,
		Domains:          domains,
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

func (conf *TLSConfig) InitCertificates() (*tls.Config, error) {
	if len(conf.caProvider) == 0 {
		return nil, fmt.Errorf("No provider set")
	}

	result := []map[string][]byte{}
	for _, caProvider := range conf.caProvider {
		caCerts, err := caProvider.GetCACertificate()
		if err != nil {
			// TODO do something here like multipleerror
			return nil, err
		}
		result = append(result, caCerts)
	}
	var certGroup *CertificateGroup
	var err error

	caCertPool := x509.NewCertPool()

	for _, provider := range result {
		certGroup, err = NewCertificateGroupFromCACM(provider)
		if err != nil {
			continue
		}
		caCertPool.AppendCertsFromPEM(certGroup.certPEM.Bytes())
	}

	if certGroup == nil {
		return nil, fmt.Errorf("Cannot get CA certificate")
	}

	serverCert, err := getServerCertificate(conf.Domains, conf.LocalhostEnabled, certGroup)
	if err != nil {
		return nil, err
	}

	certificate, err := serverCert.GetCertificate()
	if err != nil {
		return nil, fmt.Errorf("Cannot create server certfificate: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    caCertPool,
		// ClientAuth:   tls.RequireAnyClientCert,
		ClientAuth: tls.NoClientCert,
	}
	return tlsConfig, nil
}
