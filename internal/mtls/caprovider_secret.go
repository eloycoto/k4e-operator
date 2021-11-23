package mtls

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
		Name:      "test",
	}, &secret)

	if err == nil {
		// TODO return something here
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
			Name:      "test",
		},
		Data: map[string][]byte{
			"ca.crt": certificateGroup.certPEM.Bytes(),
			"ca.key": certificateGroup.PrivKeyPEM.Bytes(),
		},
	}

	err = config.client.Create(context.TODO(), &secret)
	return secret.Data, err
}

// func (config *CASecretProvider) GetServerCertificate(domains []string{}) (*CertificateGroup, error){
// 	return nil, nil
// }
