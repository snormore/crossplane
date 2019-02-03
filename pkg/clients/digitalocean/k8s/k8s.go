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

package k8s

import (
	"context"

	computev1alpha1 "github.com/crossplaneio/crossplane/pkg/apis/digitalocean/compute/v1alpha1"
	"github.com/crossplaneio/crossplane/pkg/clients/digitalocean"
	"github.com/digitalocean/godo"
)

// Client interface to perform cluster operations
type Client interface {
	CreateCluster(string, computev1alpha1.KubernetesClusterSpec) (*godo.KubernetesCluster, error)
	GetCluster(name string) (*godo.KubernetesCluster, error)
	DeleteCluster(name string) error
}

// ClusterClient implementation
type ClusterClient struct {
	creds  *digitalocean.Credentials
	client *godo.Client
}

// NewClusterClient return new instance of the Client based on credentials
func NewClusterClient(creds *digitalocean.Credentials) (Client, error) {
	client, err := digitalocean.GetClient(context.Background(), creds)
	if err != nil {
		return nil, err
	}
	return &ClusterClient{
		creds:  creds,
		client: client,
	}, nil
}

// CreateCluster provisions a new Kubernetes cluster.
func (c *ClusterClient) CreateCluster(name string, spec computev1alpha1.KubernetesClusterSpec) (*godo.KubernetesCluster, error) {
	request := &godo.KubernetesClusterCreateRequest{
		Name:        spec.Name,
		RegionSlug:  spec.RegionSlug,
		VersionSlug: spec.VersionSlug,
		Tags:        spec.Tags,
	}
	if _, _, err := c.client.Kubernetes.Create(context.TODO(), request); err != nil {
		return nil, err
	}
	return c.GetCluster(name)
}

// GetCluster retrieves a Kubernetes Cluster based on provided name.
func (c *ClusterClient) GetCluster(name string) (*godo.KubernetesCluster, error) {
	cluster, _, err := c.client.Kubernetes.Get(context.TODO(), name)
	return cluster, err
}

// DeleteCluster in the given zone with the given name.
func (c *ClusterClient) DeleteCluster(name string) error {
	_, err := c.client.Kubernetes.Delete(context.TODO(), name)
	return err
}

// DefaultKubernetesVersion is the default Kubernetes Cluster version supported by DO for given project/zone
func (c *ClusterClient) DefaultKubernetesVersion(zone string) (string, error) {
	options, _, err := c.client.Kubernetes.GetOptions(context.TODO())
	if err != nil {
		return "", nil
	}
	return options.Versions[0].KubernetesVersion, nil
}
