package main

import (
	"github.com/carwow/terraform-provider-gtm/gtm"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gtm.Provider})
}
