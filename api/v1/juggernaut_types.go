/*
Copyright 2024.

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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// JuggernautSpec defines the desired state of Juggernaut.
type JuggernautSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//容器的资源
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	//svc的类型
	Service JuggernautService `json:"service,omitempty"`
	//image
	Image string `json:"image,omitempty"`
}
type JuggernautService struct {
	//service修改
	Type corev1.ServiceType `json:"type,omitempty"`
}

// JuggernautStatus defines the observed state of Juggernaut.
type JuggernautStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Juggernaut is the Schema for the juggernauts API.
type Juggernaut struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JuggernautSpec   `json:"spec,omitempty"`
	Status JuggernautStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JuggernautList contains a list of Juggernaut.
type JuggernautList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Juggernaut `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Juggernaut{}, &JuggernautList{})
}
