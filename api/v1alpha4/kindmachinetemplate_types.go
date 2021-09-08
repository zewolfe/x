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

package v1alpha4

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KINDMachineTemplateSpec defines the desired state of KINDMachineTemplate
type KINDMachineTemplateSpec struct {
	Template KINDMachineTemplateResource `json:"template"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KINDMachineTemplate is the Schema for the kindmachinetemplates API
type KINDMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec KINDMachineTemplateSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// KINDMachineTemplateList contains a list of KINDMachineTemplate
type KINDMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KINDMachineTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KINDMachineTemplate{}, &KINDMachineTemplateList{})
}

type KINDMachineTemplateResource struct {
	// Spec is the specification of the desired behavior of the machine.
	Spec KINDMachineSpec `json:"spec"`
}
