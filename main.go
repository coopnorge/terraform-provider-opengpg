package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/coopnorge/terraform-provider-opengpg/gpg"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gpg.Provider,
	})
}
