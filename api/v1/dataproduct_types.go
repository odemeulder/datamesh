/*
Copyright 2021.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DataProductSpec defines the desired state of DataProduct
type DataProductSpec struct {
	Image       string   `json:"image,omitempty"`
	Cronexpr    string   `json:"cronexpr,omitempty"`
	InputPorts  []string `json:"inputPorts,omitempty"`
	OutputPorts []string `json:"outputPorts,omitempty"`
	Size        int32    `json:"size,omitempty"`
}

// DataProductStatus defines the observed state of DataProduct
type DataProductStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes []string `json:"nodes,omitempty"`
}

//+kubebuilder:object:root=true

// DataProduct is the Schema for the dataproducts API
//+kubebuilder:subresource:status
type DataProduct struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataProductSpec   `json:"spec,omitempty"`
	Status DataProductStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DataProductList contains a list of DataProduct
type DataProductList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DataProduct `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DataProduct{}, &DataProductList{})
}
