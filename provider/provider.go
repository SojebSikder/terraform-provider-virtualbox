package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func New() provider.Provider {
	return &VirtualBoxProvider{}
}

type VirtualBoxProvider struct {
	version string
}

// Configure implements provider.Provider.
func (p *VirtualBoxProvider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {

}

// DataSources implements provider.Provider.
func (p *VirtualBoxProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}

// Resources implements provider.Provider.
func (p *VirtualBoxProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVirtualBoxVMResource,
	}
}

func (p *VirtualBoxProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "virtualbox"
	resp.Version = p.version
}

func (p *VirtualBoxProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provider for managing VirtualBox VMs.",
		Attributes: map[string]schema.Attribute{
			"vboxmanage_path": schema.StringAttribute{
				Optional:    true,
				Description: "Path to the VBoxManage executable.",
			},
		},
	}
}
