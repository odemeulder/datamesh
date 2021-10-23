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

// InputPortS3Spec defines the desired state of InputPortS3
type InputPortS3Spec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	BucketArn  string `json:"bucketArn,omitempty"`
	ObjectName string `json:"objectName,omitempty"`
	FileTypem  string `json:"fileType,omitempty"`
}

// InputPortS3Status defines the observed state of InputPortS3
type InputPortS3Status struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// InputPortS3 is the Schema for the inputports3s API
type InputPortS3 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InputPortS3Spec   `json:"spec,omitempty"`
	Status InputPortS3Status `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// InputPortS3List contains a list of InputPortS3
type InputPortS3List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InputPortS3 `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InputPortS3{}, &InputPortS3List{})
}
