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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ServerSpec defines the desired state of Server
type ServerSpec struct {
	// The type of Minecraft Server we want to run. Includes VANILLA, PAPER, etc.
	ServerType string `json:"serverType,omitempty"`

	// The name of the server when a user connects
	ServerName string `json:"serverName,omitempty"`

	// The list of users that are ops/admin on the server
	Ops []string `json:"ops,omitempty"`

	// The list of users that can join. If this is empty then the server is open
	// to the public.
	Allowlist []string `json:"allowList,omitempty"`

	// Do you agree to the EULA?
	EULA bool `json:"EULA,omitempty"`
}

// ServerStatus defines the observed state of Server
type ServerStatus struct {
	PodName string `json:"podName,omitempty"`
	Running bool   `json:"running,omitempty"`
}

// +kubebuilder:object:root=true

// Server is the Schema for the servers API
// +kubebuilder:subresource:status
type Server struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServerSpec   `json:"spec,omitempty"`
	Status ServerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ServerList contains a list of Server
type ServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Server `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Server{}, &ServerList{})
}
