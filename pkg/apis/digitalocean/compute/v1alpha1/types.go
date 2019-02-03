/*
Copyright 2018 The Crossplane Authors.

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
	corev1alpha1 "github.com/crossplaneio/crossplane/pkg/apis/core/v1alpha1"
	"github.com/crossplaneio/crossplane/pkg/util"
	"github.com/digitalocean/godo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ClusterStateProvisioning = "PROVISIONING"
	ClusterStateRunning      = "RUNNING"

	DefaultReclaimPolicy = corev1alpha1.ReclaimRetain
	DefaultNumberOfNodes = int64(1)
)

// KubernetesClusterSpec
type KubernetesClusterSpec struct {
	godo.KubernetesClusterCreateRequest

	// Kubernetes object references
	ClaimRef            *corev1.ObjectReference      `json:"claimRef,omitempty"`
	ClassRef            *corev1.ObjectReference      `json:"classRef,omitempty"`
	ConnectionSecretRef *corev1.LocalObjectReference `json:"connectionSecretRef,omitempty"`
	ProviderRef         corev1.LocalObjectReference  `json:"providerRef,omitempty"`

	// ReclaimPolicy identifies how to handle the cloud resource after the deletion of this type
	ReclaimPolicy corev1alpha1.ReclaimPolicy `json:"reclaimPolicy,omitempty"`
}

// KubernetesClusterStatus
type KubernetesClusterStatus struct {
	corev1alpha1.ConditionedStatus
	corev1alpha1.BindingStatusPhase
	ClusterName string `json:"clusterName"`
	Endpoint    string `json:"endpoint"`
	State       string `json:"state,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubernetesCluster is the Schema for the instances API
// +k8s:openapi-gen=true
// +groupName=compute.digitalocean
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.state"
// +kubebuilder:printcolumn:name="CLUSTER-NAME",type="string",JSONPath=".status.clusterName"
// +kubebuilder:printcolumn:name="ENDPOINT",type="string",JSONPath=".status.endpoint"
// +kubebuilder:printcolumn:name="CLUSTER-CLASS",type="string",JSONPath=".spec.classRef.name"
// +kubebuilder:printcolumn:name="LOCATION",type="string",JSONPath=".spec.location"
// +kubebuilder:printcolumn:name="RECLAIM-POLICY",type="string",JSONPath=".spec.reclaimPolicy"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type KubernetesCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubernetesClusterSpec   `json:"spec,omitempty"`
	Status KubernetesClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubernetesClusterList contains a list of KubernetesCluster items
type KubernetesClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubernetesCluster `json:"items"`
}

// NewKubernetesClusterSpec from properties map
func NewKubernetesClusterSpec(properties map[string]string) *KubernetesClusterSpec {
	spec := &KubernetesClusterSpec{
		KubernetesClusterCreateRequest: godo.KubernetesClusterCreateRequest{
			Name:        properties["name"],
			RegionSlug:  properties["region"],
			VersionSlug: properties["version"],
		},
	}
	return spec
}

// ObjectReference to this RDSInstance
func (g *KubernetesCluster) ObjectReference() *corev1.ObjectReference {
	return util.ObjectReference(g.ObjectMeta, util.IfEmptyString(g.APIVersion, APIVersion), util.IfEmptyString(g.Kind, KubernetesClusterKind))
}

// OwnerReference to use this instance as an owner
func (g *KubernetesCluster) OwnerReference() metav1.OwnerReference {
	return *util.ObjectToOwnerReference(g.ObjectReference())
}

func (g *KubernetesCluster) ConnectionSecret() *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:       g.Namespace,
			Name:            g.ConnectionSecretName(),
			OwnerReferences: []metav1.OwnerReference{g.OwnerReference()},
		},
	}
}

// ConnectionSecretName returns a secret name from the reference
func (g *KubernetesCluster) ConnectionSecretName() string {
	if g.Spec.ConnectionSecretRef == nil {
		g.Spec.ConnectionSecretRef = &corev1.LocalObjectReference{
			Name: g.Name,
		}
	} else if g.Spec.ConnectionSecretRef.Name == "" {
		g.Spec.ConnectionSecretRef.Name = g.Name
	}

	return g.Spec.ConnectionSecretRef.Name
}

// State returns rds instance state value saved in the status (could be empty)
func (g *KubernetesCluster) State() string {
	return string(g.Status.State)
}

// IsAvailable for usage/binding
func (g *KubernetesCluster) IsAvailable() bool {
	return g.State() == ClusterStateRunning
}

// IsBound
func (g *KubernetesCluster) IsBound() bool {
	return g.Status.Phase == corev1alpha1.BindingStateBound
}

// SetBound
func (g *KubernetesCluster) SetBound(state bool) {
	if state {
		g.Status.Phase = corev1alpha1.BindingStateBound
	} else {
		g.Status.Phase = corev1alpha1.BindingStateUnbound
	}
}
