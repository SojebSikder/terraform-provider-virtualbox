package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/sojebsikder/terraform-provider-virtualbox/provider"
)

func main() {
	providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		// comment this if you want do development
		Address: "registry.terraform.io/sojebsikder/virtualbox",
		// uncomment if you want to do development
		// Address: "hashicorp.com/edu/virtualbox",
	})
}
