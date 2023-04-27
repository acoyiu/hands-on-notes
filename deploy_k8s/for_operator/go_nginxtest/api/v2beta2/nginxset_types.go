/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v2beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NginxsetSpec defines the desired state of Nginxset
type NginxsetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// target reture string
	ReturnText string `json:"returnText,omitempty"`
}

// NginxsetStatus defines the observed state of Nginxset
type NginxsetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Ready represent the nginx behind this res is created or not
	Ready            string     `json:"ready,omitempty"`
	LinkedDeployment *types.UID `json:"linkedPod,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Ready",type="string",JSONPath=`.status.ready`
//+kubebuilder:printcolumn:name="LastModify",type="date",JSONPath=`.metadata.creationTimestamp`

// Nginxset is the Schema for the nginxsets API
type Nginxset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NginxsetSpec   `json:"spec,omitempty"`
	Status NginxsetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NginxsetList contains a list of Nginxset
type NginxsetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Nginxset `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Nginxset{}, &NginxsetList{})
}
