package thales

import (
	"github.com/odysseia-greek/agora/thales/models"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1cert "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"log"
	"strings"
)

type KubeClient interface {
	Access() Access
	Configuration() Configuration
	Cluster() Cluster
	Util() Util
	Workload() Workload
	Nodes() Nodes
	Namespaces() Namespace
	Storage() Storage
	Certificate() Certificate
	V1Alpha1() V1Alpha1
}

type V1Alpha1 interface {
	ServiceMapping() ServiceMapping
}

type Access interface {
	GetServiceAccounts() (*corev1.ServiceAccountList, error)
}

type Configuration interface {
	GetSecret(namespace, secretName string) (*corev1.Secret, error)
	ListSecrets(namespace string) (*corev1.SecretList, error)
	CreateSecret(namespace, secretName string, data map[string][]byte) error
	UpdateSecret(namespace, secretName string, data map[string][]byte) error
	DeleteSecret(namespace, secretName string) error
	CreateDockerSecret(namespace, secretName string, data map[string]string) error
	CreateTlSSecret(namespace, secretName string, data map[string][]byte, immutable bool) error
	UpdateTLSSecret(namespace, secretName string, data map[string][]byte, annotation map[string]string) error
}

type Cluster interface {
	GetHostServer() (string, error)
	GetHostCaCert() ([]byte, error)
	GetCurrentContext() (string, error)
}

type Namespace interface {
	Create(namespace string) error
	Delete(namespace string) error
	List() (*corev1.NamespaceList, error)
}

type Util interface {
	CopyFileToPod(podName, destPath, srcPath string) (string, error)
	CopyFileFromPod(srcPath, destPath, namespace, podName string) error
}

type Workload interface {
	List(namespace string) (*corev1.PodList, error)
	CreatePodSpec(namespace, name, podImage string, command []string) *corev1.Pod
	DeletePod(namespace, podName string) error
	CreatePod(namespace string, pod *corev1.Pod) (*corev1.Pod, error)
	ExecNamedPod(namespace, podName string, command []string) (string, error)
	GetStatefulSets(namespace string) (*appsv1.StatefulSetList, error)
	GetPodsBySelector(namespace, selector string) (*corev1.PodList, error)
	GetPodByName(namespace, name string) (*corev1.Pod, error)
	CreateDeployment(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error)
	ListDeployments(namespace string) (*appsv1.DeploymentList, error)
	UpdateDeploymentViaAnnotation(namespace, name string, annotation map[string]string) (*appsv1.Deployment, error)
	GetDeployment(namespace, name string) (*appsv1.Deployment, error)
	GetDeploymentStatus(namespace string) (bool, error)
	CreateJob(namespace string, spec *batchv1.Job) (*batchv1.Job, error)
	GetJob(namespace, name string) (*batchv1.Job, error)
	ListJobs(namespace string) (*batchv1.JobList, error)
	GetNewLock(lockName, podName, namespace string) *resourcelock.LeaseLock
}

type Storage interface {
	ListPvc(namespace string) (*corev1.PersistentVolumeClaimList, error)
	DeletePvc(namespace, name string) error
	ListPv() (*corev1.PersistentVolumeList, error)
	DeletePv(name string) error
}

type Certificate interface {
	ListCsr() (*v1cert.CertificateSigningRequestList, error)
	DeleteCsr(name string) error
}

type Nodes interface {
	List() (*corev1.NodeList, error)
}

type Kube struct {
	set           *kubernetes.Clientset
	v1alpha       *V1Alpha1Impl
	access        *AccessImpl
	util          *UtilImpl
	cluster       *ClusterImpl
	configuration *ConfigurationImpl
	workload      *WorkloadImpl
	nodes         *NodesImpl
	namespace     *NamespaceImpl
	storage       *StorageImpl
	certificate   *CertificateImpl
	config        []byte
}

func NewKubeClient(cfg []byte, ns string) (KubeClient, error) {
	var kube *Kube
	var err error

	if cfg != nil {
		kube, err = NewConfigBasedKube(cfg, ns)
		if err != nil {
			return nil, err
		}
	} else {
		kube, err = NewInClusterKube(ns)
		if err != nil {
			return nil, err
		}
	}

	return kube, err
}

func New(clientSet *kubernetes.Clientset, restConfig *rest.Config, config []byte, ns string) (*Kube, error) {
	access, err := NewAccessClient(clientSet, ns)
	if err != nil {
		return nil, err
	}

	cluster := &ClusterImpl{}

	if config != nil {
		cluster, err = NewClusterClient(config)
		if err != nil {
			return nil, err
		}
	}

	configuration, err := NewConfigurationClient(clientSet)
	if err != nil {
		return nil, err
	}

	namespaceClient, err := NewNamespaceClient(clientSet)
	if err != nil {
		return nil, err
	}

	workload, err := NewWorkloadClient(clientSet)
	if err != nil {
		return nil, err
	}

	util, err := NewUtilClient(clientSet, ns)
	if err != nil {
		return nil, err
	}

	nodes, err := NewNodesClient(clientSet)
	if err != nil {
		return nil, err
	}

	v1alphaClient, err := NewV1AlphaClient(restConfig)
	if err != nil {
		return nil, err
	}

	storage, err := NewStorageClient(clientSet)
	if err != nil {
		return nil, err
	}

	certificates, err := NewCertificateClient(clientSet)
	if err != nil {
		return nil, err
	}

	return &Kube{
		set:           clientSet,
		v1alpha:       v1alphaClient,
		access:        access,
		util:          util,
		cluster:       cluster,
		configuration: configuration,
		workload:      workload,
		nodes:         nodes,
		namespace:     namespaceClient,
		config:        config,
		storage:       storage,
		certificate:   certificates,
	}, nil
}

func NewConfigBasedKube(config []byte, ns string) (*Kube, error) {
	c, err := clientcmd.NewClientConfigFromBytes(config)
	if err != nil {
		return nil, err
	}

	raw, err := c.RawConfig()
	if err != nil {
		return nil, err
	}

	if raw.CurrentContext == "" {
		cfg, err := models.UnmarshalKubeConfig(config)
		if err != nil {
			return nil, err
		}

		if len(raw.Contexts) == 1 {
			firstContext := raw.Contexts[cfg.Contexts[0].Name]
			raw.CurrentContext = firstContext.Cluster
			c = clientcmd.NewNonInteractiveClientConfig(raw, firstContext.Cluster, nil, nil)
			cfg.CurrentContext = firstContext.Cluster
		} else {
			for key, _ := range raw.Contexts {
				if strings.Contains(key, "desktop") {
					raw.CurrentContext = key
					c = clientcmd.NewNonInteractiveClientConfig(raw, key, nil, nil)
					cfg.CurrentContext = key
					break
				}
			}
		}

		newConfig, err := cfg.Marshal()
		if err != nil {
			return nil, err
		}

		config = newConfig
	}

	restConfig, err := c.ClientConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return New(clientSet, restConfig, config, ns)

}

func NewInClusterKube(ns string) (*Kube, error) {
	restConfig, err := rest.InClusterConfig()
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	log.Print("created in cluster kube client")

	return New(clientSet, restConfig, nil, ns)
}

func FakeKubeClient(ns string) (KubeClient, error) {
	clientSet := fake.NewSimpleClientset()

	log.Print("created fake kube client")

	access, err := NewAccessClient(clientSet, ns)
	if err != nil {
		return nil, err
	}

	cluster, err := NewClusterClient(nil)
	if err != nil {
		return nil, err
	}

	configuration, err := NewConfigurationClient(clientSet)
	if err != nil {
		return nil, err
	}

	namespaceClient, err := NewNamespaceClient(clientSet)
	if err != nil {
		return nil, err
	}

	workload, err := NewWorkloadClient(clientSet)
	if err != nil {
		return nil, err
	}

	util, err := NewUtilClient(clientSet, ns)
	if err != nil {
		return nil, err
	}

	nodes, err := NewNodesClient(clientSet)
	if err != nil {
		return nil, err
	}

	storage, err := NewStorageClient(clientSet)
	if err != nil {
		return nil, err
	}

	certificates, err := NewCertificateClient(clientSet)
	if err != nil {
		return nil, err
	}

	alpha, err := NewV1FakeAlphaClient()
	if err != nil {
		return nil, err
	}

	return &Kube{
		set:           nil,
		config:        nil,
		access:        access,
		cluster:       cluster,
		configuration: configuration,
		workload:      workload,
		util:          util,
		nodes:         nodes,
		namespace:     namespaceClient,
		v1alpha:       alpha,
		storage:       storage,
		certificate:   certificates,
	}, nil
}

func (k *Kube) V1Alpha1() V1Alpha1 {
	if k == nil {
		return nil
	}
	return k.v1alpha
}

func (k *Kube) Access() Access {
	if k == nil {
		return nil
	}
	return k.access
}

func (k *Kube) Util() Util {
	if k == nil {
		return nil
	}

	return k.util
}

func (k *Kube) Namespaces() Namespace {
	if k == nil {
		return nil
	}

	return k.namespace
}

func (k *Kube) Cluster() Cluster {
	if k == nil {
		return nil
	}

	return k.cluster
}

func (k *Kube) Configuration() Configuration {
	if k == nil {
		return nil
	}

	return k.configuration
}

func (k *Kube) Workload() Workload {
	if k == nil {
		return nil
	}

	return k.workload
}

func (k *Kube) Nodes() Nodes {
	if k == nil {
		return nil
	}

	return k.nodes
}

func (k *Kube) Storage() Storage {
	if k == nil {
		return nil
	}

	return k.storage
}

func (k *Kube) Certificate() Certificate {
	if k == nil {
		return nil
	}

	return k.certificate
}

//func (k *Kube) GetK8sClientSet() *kubernetes.Clientset {
//	return k.set
//}
//
//func (k *Kube) GetConfig() []byte {
//	return k.config
//}
