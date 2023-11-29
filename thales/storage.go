package thales

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"time"
)

type StorageImpl struct {
	client v1.CoreV1Interface
}

func NewStorageClient(kube kubernetes.Interface) (*StorageImpl, error) {
	coreClient := kube.CoreV1()

	return &StorageImpl{client: coreClient}, nil
}

func (s *StorageImpl) ListPvc(namespace string) (*corev1.PersistentVolumeClaimList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	pvc, err := s.client.PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return pvc, err
}

func (s *StorageImpl) DeletePvc(namespace, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	return s.client.PersistentVolumeClaims(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (s *StorageImpl) ListPv() (*corev1.PersistentVolumeList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	pv, err := s.client.PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return pv, err
}

func (s *StorageImpl) DeletePv(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	return s.client.PersistentVolumes().Delete(ctx, name, metav1.DeleteOptions{})
}
