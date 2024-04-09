package thales

import (
	appsv1Api "k8s.io/api/apps/v1"
	corev1Api "k8s.io/api/core/v1"
	apiextensionsv1Api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	fakeapiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	admissionregistrationv1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1"
	fakeadmissionregistrationv1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1/fake"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	fakeappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1/fake"
	authorizationv1 "k8s.io/client-go/kubernetes/typed/authorization/v1"
	fakeauthorizationv1 "k8s.io/client-go/kubernetes/typed/authorization/v1/fake"
	batchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	fakebatchv1 "k8s.io/client-go/kubernetes/typed/batch/v1/fake"
	certificatesv1 "k8s.io/client-go/kubernetes/typed/certificates/v1"
	fakecertificatesv1 "k8s.io/client-go/kubernetes/typed/certificates/v1/fake"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	fakecorev1 "k8s.io/client-go/kubernetes/typed/core/v1/fake"
	rbacv1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	fakerbacv1 "k8s.io/client-go/kubernetes/typed/rbac/v1/fake"
	storagev1 "k8s.io/client-go/kubernetes/typed/storage/v1"
	fakestoragev1 "k8s.io/client-go/kubernetes/typed/storage/v1/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/testing"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	fakemetrics "k8s.io/metrics/pkg/client/clientset/versioned/fake"
)

type KubeClient struct {
	fake                    bool
	config                  *rest.Config
	coreV1                  corev1.CoreV1Interface
	appsV1                  appsv1.AppsV1Interface
	rbacV1                  rbacv1.RbacV1Interface
	batchV1                 batchv1.BatchV1Interface
	storageV1               storagev1.StorageV1Interface
	discovery               discovery.DiscoveryInterface
	authorizationV1         authorizationv1.AuthorizationV1Interface
	apiextensionsV1         apiextensionsv1.ApiextensionsV1Interface
	certificatesV1          certificatesv1.CertificatesV1Interface
	admissionregistrationv1 admissionregistrationv1.AdmissionregistrationV1Interface
	metricsClient           versioned.Interface

	dynamic dynamic.Interface
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func NewFakeKubeClient() *KubeClient {
	var fake testing.Fake

	scheme := runtime.NewScheme()

	must(metav1.AddMetaToScheme(scheme))
	must(corev1Api.AddToScheme(scheme))
	must(apiextensionsv1Api.AddToScheme(scheme))
	must(appsv1Api.AddToScheme(scheme))

	codecs := serializer.NewCodecFactory(scheme)
	tracker := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())

	fake.AddReactor("*", "*", testing.ObjectReaction(tracker))

	fakeMetricsClientset := &fakemetrics.Clientset{}

	return &KubeClient{
		fake:                    true,
		coreV1:                  &fakecorev1.FakeCoreV1{Fake: &fake},
		appsV1:                  &fakeappsv1.FakeAppsV1{Fake: &fake},
		rbacV1:                  &fakerbacv1.FakeRbacV1{Fake: &fake},
		batchV1:                 &fakebatchv1.FakeBatchV1{Fake: &fake},
		storageV1:               &fakestoragev1.FakeStorageV1{Fake: &fake},
		authorizationV1:         &fakeauthorizationv1.FakeAuthorizationV1{Fake: &fake},
		apiextensionsV1:         &fakeapiextensionsv1.FakeApiextensionsV1{Fake: &fake},
		certificatesV1:          &fakecertificatesv1.FakeCertificatesV1{Fake: &fake},
		admissionregistrationv1: &fakeadmissionregistrationv1.FakeAdmissionregistrationV1{Fake: &fake},
		metricsClient:           fakeMetricsClientset,
	}
}

func NewKubeClient(c *rest.Config) (*KubeClient, error) {
	httpClient, err := rest.HTTPClientFor(c)
	if err != nil {
		return nil, err
	}

	var kc KubeClient

	kc.config = c

	kc.coreV1, err = corev1.NewForConfigAndClient(c, httpClient)
	if err != nil {
		return nil, err
	}

	kc.batchV1, err = batchv1.NewForConfigAndClient(c, httpClient)
	if err != nil {
		return nil, err
	}

	kc.storageV1, err = storagev1.NewForConfigAndClient(c, httpClient)
	if err != nil {
		return nil, err
	}

	kc.appsV1, err = appsv1.NewForConfigAndClient(c, httpClient)
	if err != nil {
		return nil, err
	}

	kc.rbacV1, err = rbacv1.NewForConfigAndClient(c, httpClient)
	if err != nil {
		return nil, err
	}

	kc.authorizationV1, err = authorizationv1.NewForConfigAndClient(c, httpClient)
	if err != nil {
		return nil, err
	}

	kc.discovery, err = discovery.NewDiscoveryClientForConfigAndClient(c, httpClient)
	if err != nil {
		return nil, err
	}

	kc.apiextensionsV1, err = apiextensionsv1.NewForConfigAndClient(c, httpClient)
	if err != nil {
		return nil, err
	}

	kc.certificatesV1, err = certificatesv1.NewForConfigAndClient(c, httpClient)
	if err != nil {
		return nil, err
	}

	kc.admissionregistrationv1, err = admissionregistrationv1.NewForConfigAndClient(c, httpClient)
	if err != nil {
		return nil, err
	}

	kc.metricsClient, err = versioned.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	kc.dynamic, err = dynamic.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	return &kc, nil
}

func (c *KubeClient) CoreV1() corev1.CoreV1Interface {
	return c.coreV1
}

func (c *KubeClient) BatchV1() batchv1.BatchV1Interface {
	return c.batchV1
}

func (c *KubeClient) StorageV1() storagev1.StorageV1Interface {
	return c.storageV1
}

func (c *KubeClient) AppsV1() appsv1.AppsV1Interface {
	return c.appsV1
}

func (c *KubeClient) AuthorizationV1() authorizationv1.AuthorizationV1Interface {
	return c.authorizationV1
}

func (c *KubeClient) RbacV1() rbacv1.RbacV1Interface {
	return c.rbacV1
}

func (c *KubeClient) Discovery() discovery.DiscoveryInterface {
	return c.discovery
}

func (c *KubeClient) ApiextensionsV1() apiextensionsv1.ApiextensionsV1Interface {
	return c.apiextensionsV1
}

func (c *KubeClient) CertificatesV1() certificatesv1.CertificatesV1Interface {
	return c.certificatesV1
}

func (c *KubeClient) AdmissionRegistrationV1() admissionregistrationv1.AdmissionregistrationV1Interface {
	return c.admissionregistrationv1
}

func (c *KubeClient) Dynamic() dynamic.Interface {
	return c.dynamic
}

func (c *KubeClient) Host() string {
	return c.config.Host
}

func (c *KubeClient) RestConfig() *rest.Config {
	return c.config
}

func (c *KubeClient) Cert() corev1.CoreV1Interface {
	return c.coreV1
}

func (c *KubeClient) MetricsClient() versioned.Interface {
	return c.metricsClient
}

func (c *KubeClient) KubeCliConfig(namespace string) (*genericclioptions.ConfigFlags, error) {
	kubeConfig := genericclioptions.NewConfigFlags(false)
	kubeConfig.APIServer = &c.config.Host
	kubeConfig.BearerToken = &c.config.BearerToken
	kubeConfig.CAFile = &c.config.CAFile
	kubeConfig.Namespace = &namespace

	return kubeConfig, nil
}

func NewFromConfig(config []byte) (*KubeClient, error) {
	c, err := clientcmd.NewClientConfigFromBytes(config)
	if err != nil {
		return nil, err
	}

	restConfig, err := c.ClientConfig()
	if err != nil {
		return nil, err
	}

	return NewKubeClient(restConfig)
}

func NewInClusterKube() (*KubeClient, error) {
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return NewKubeClient(restConfig)
}
