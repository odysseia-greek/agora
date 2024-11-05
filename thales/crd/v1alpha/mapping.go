package v1alpha

import (
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//go:generate controller-gen object paths=$GOFILE

func (m *Mapping) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MappingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Mapping `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Mapping struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	APIVersion        string `json:"apiVersion,omitempty"`
	Kind              string `json:"kind,omitempty"`
	Spec              Spec   `json:"spec"`
}

// +k8s:deepcopy-gen=true
type Spec struct {
	Services []Service `json:"services"`
}

// +k8s:deepcopy-gen=true
type Service struct {
	Name       string   `json:"name"`
	KubeType   string   `json:"kubeType"`
	SecretName string   `json:"secretName"`
	Namespace  string   `json:"namespace"`
	Active     bool     `json:"active"`
	Created    string   `json:"created"`
	Validity   int      `json:"validity"`
	Clients    []Client `json:"clients"`
}

// +k8s:deepcopy-gen=true
type Client struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	KubeType  string `json:"kubeType"`
}
