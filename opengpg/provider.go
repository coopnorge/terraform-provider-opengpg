package opengpg

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider exports terraform-provider-opengpg, which can be used in tests
// for other providers.
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"opengpg_encrypted_message": resourceGPGEncryptedMessage(),
		},
	}
}
