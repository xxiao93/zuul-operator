package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ZuulSpec defines the desired state of Zuul
// +k8s:openapi-gen=true
type ZuulSpec struct {
	ZuulVersion    string        `json:"zuulversion"`
	ZuulScheduler  ZuulScheduler `json:"zuulscheduler-spec"`
	Gerrit         Gerrit        `json:"gerrit-spec"`
}

// ZuulStatus defines the observed state of Zuul
// +k8s:openapi-gen=true
type ZuulStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Zuul is the Schema for the zuuls API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Zuul struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ZuulSpec   `json:"spec,omitempty"`
	Status ZuulStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ZuulList contains a list of Zuul
type ZuulList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Zuul `json:"items"`
}

type ZuulScheduler struct {
	Size int32 `json:"size"`
}

type Gerrit struct {
	Server string `json:"server"`
	Port   string `json:"port"`
	User   string `json:"user"`
}

func init() {
	SchemeBuilder.Register(&Zuul{}, &ZuulList{})
}
