package opengpg_test

import (
	"testing"

	"github.com/coopnorge/terraform-provider-opengpg/opengpg"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"opengpg": func() (*schema.Provider, error) {
		return opengpg.Provider(), nil
	},
}

func TestProvider(t *testing.T) {
	t.Parallel()

	provider := opengpg.Provider()

	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("validating provider internally: %v", err)
	}
}
