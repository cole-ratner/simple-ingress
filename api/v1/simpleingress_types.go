/*


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

// SimpleIngressSpec defines the desired state of SimpleIngress
type SimpleIngressSpec struct {
	// +kubebuilder:validation:Required
	Host string `json:"host"`
	// +kubebuilder:validation:Required
	ServiceName string `json:"serviceName"`
}

// SimpleIngressStatus defines the observed state of SimpleIngress
type SimpleIngressStatus struct {
	Host        string `json:"host"`
	ServiceName string `json:"serviceName"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SimpleIngress is the Schema for the simpleingresses API
type SimpleIngress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SimpleIngressSpec   `json:"spec,omitempty"`
	Status SimpleIngressStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SimpleIngressList contains a list of SimpleIngress
type SimpleIngressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SimpleIngress `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SimpleIngress{}, &SimpleIngressList{})
}
