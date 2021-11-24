package mtls

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	CASecretName = "k4e-ca"
)

type CASecretProvider struct {
	client    client.Client
	namespace string
}

func NewCASecretProvider(client client.Client, namespace string) *CASecretProvider {
	return &CASecretProvider{
		client:    client,
		namespace: namespace,
	}
}

func (config *CASecretProvider) GetCACertificate() (map[string][]byte, error) {
	var secret corev1.Secret

	err := config.client.Get(context.TODO(), client.ObjectKey{
		Namespace: config.namespace,
		Name:      CASecretName,
	}, &secret)

	if err == nil {
		return secret.Data, nil
	}

	if !errors.IsNotFound(err) {
		return nil, err
	}

	certificateGroup, err := getCACertificate()
	if err != nil {
		return nil, fmt.Errorf("cannot create sample certificate")
	}

	secret = corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Namespace: config.namespace,
			Name:      CASecretName,
		},
		Data: map[string][]byte{
			"ca.crt": certificateGroup.certPEM.Bytes(),
			"ca.key": certificateGroup.PrivKeyPEM.Bytes(),
		},
	}

	err = config.client.Create(context.TODO(), &secret)
	return secret.Data, err
}

func (config *CASecretProvider) CreateRegistrationCertificate(name string) (map[string][]byte, error) {
	CAData, err := config.GetCACertificate()
	if err != nil {
		return nil, fmt.Errorf("Cannot retrieve caCert")
	}

	CACert, err := NewCertificateGroupFromCACM(CAData)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse CA certificate")
	}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			CommonName:   fmt.Sprintf("registered-%s", name),
			Organization: []string{"K4e-agent"},
			Country:      []string{"US"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certGroup, err := getKeyAndCSR(cert, CACert)
	if err != nil {
		return nil, fmt.Errorf("Cannot sign certificate request")
	}
	certGroup.CreatePem()

	res := map[string][]byte{
		"client.crt": certGroup.certPEM.Bytes(),
		"client.key": certGroup.PrivKeyPEM.Bytes(),
	}
	return res, nil
}
