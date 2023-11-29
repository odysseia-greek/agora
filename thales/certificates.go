package thales

import (
	"context"
	v1cert "k8s.io/api/certificates/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/certificates/v1"
	"time"
)

type CertificateImpl struct {
	client v1.CertificatesV1Interface
}

func NewCertificateClient(kube kubernetes.Interface) (*CertificateImpl, error) {
	certificateClient := kube.CertificatesV1()

	return &CertificateImpl{client: certificateClient}, nil
}

func (c *CertificateImpl) ListCsr() (*v1cert.CertificateSigningRequestList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	return c.client.CertificateSigningRequests().List(ctx, metav1.ListOptions{})
}

func (c *CertificateImpl) DeleteCsr(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	return c.client.CertificateSigningRequests().Delete(ctx, name, metav1.DeleteOptions{})
}
