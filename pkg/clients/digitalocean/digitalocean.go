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

package digitalocean

import (
	"context"

	digitaloceanv1alpha1 "github.com/crossplaneio/crossplane/pkg/apis/digitalocean/v1alpha1"
	"github.com/crossplaneio/crossplane/pkg/util"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"k8s.io/client-go/kubernetes"
)

// TokenSource encapsulates an access token for the DO client.
type TokenSource struct {
	AccessToken string
}

// Token returns an oauth2 token for this source.
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

// Credentials encapsulates the DO client credentials strategy.
type Credentials struct {
	TokenSource *TokenSource
}

// ValidateClient verifies if the given client is valid by testing if it can make an DO service API call
// TODO: is there a better way to validate the DO client?
func ValidateClient(creds *Credentials) error {
	client, err := GetClient(context.TODO(), creds)
	if err != nil {
		return err
	}
	_, _, err = client.Kubernetes.GetOptions(context.TODO())
	return err
}

// GetClient builds and returns a DO API client.
func GetClient(ctx context.Context, creds *Credentials) (*godo.Client, error) {
	hc := oauth2.NewClient(ctx, creds.TokenSource)
	client := godo.NewClient(hc)
	return client, nil
}

// ProviderCredentials return DO credentials based on the provider's credentials secret data
func ProviderCredentials(client kubernetes.Interface, p *digitaloceanv1alpha1.Provider) (*Credentials, error) {
	// retrieve provider secret data
	data, err := util.SecretData(client, p.Namespace, p.Spec.Secret)
	if err != nil {
		return nil, err
	}
	creds := &Credentials{
		TokenSource: &TokenSource{
			AccessToken: string(data),
		},
	}
	return creds, nil
}
