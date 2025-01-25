package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/sojebsikder/terraform-provider-virtualbox/provider"
)

func main() {
	providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/sojebsikder/virtualbox",
	})
}
