package integration

import (
	"flag"
	"testing"

	. "github.com/onsi/gomega"
	"golang.org/x/oauth2/google"
)

var (
	// digitaloceanCredsFile - retrieve digitalocean credentials from the file
	digitaloceanCredsFile = flag.String("digitalocean-creds", "", "run integration tests that require crossplane-digitalocean-provider-key.json")
)

func init() {
	flag.Parse()
}

// CredsOrSkip - returns digitalocean configuration if environment is set, otherwise - skips this test
func CredsOrSkip(t *testing.T, scopes ...string) (*GomegaWithT, *google.Credentials) {
	if *digitaloceanCredsFile == "" {
		t.Skip()
	}

	g := NewGomegaWithT(t)

	creds, err := digitalocean.CredentialsFromFile(*digitaloceanCredsFile, scopes...)
	g.Expect(err).NotTo(HaveOccurred())

	return g, creds
}
