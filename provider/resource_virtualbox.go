package provider

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Define the struct that represents the VM resource data
type VirtualBoxVMResourceData struct {
	Name    types.String `tfsdk:"name"`
	ISOPath types.String `tfsdk:"iso_path"`
	Memory  types.Int64  `tfsdk:"memory"`
	CPUs    types.Int64  `tfsdk:"cpus"`
}

type VirtualBoxVMResource struct{}

func NewVirtualBoxVMResource() resource.Resource {
	return &VirtualBoxVMResource{}
}

func (r *VirtualBoxVMResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "virtualbox_vm"
}

func (r *VirtualBoxVMResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a VirtualBox VM.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VM.",
			},
			"iso_path": schema.StringAttribute{
				Required:    true,
				Description: "Path to the Ubuntu ISO file.",
			},
			"memory": schema.Int64Attribute{
				Required:    true,
				Description: "Memory size in MB.",
			},
			"cpus": schema.Int64Attribute{
				Required:    true,
				Description: "Number of CPUs.",
			},
		},
	}
}

func (r *VirtualBoxVMResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VirtualBoxVMResourceData

	// Decode the request into the struct
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract values from the struct
	name := data.Name.ValueString()
	isoPath := data.ISOPath.ValueString()
	memory := data.Memory.ValueInt64()
	cpus := data.CPUs.ValueInt64()

	// Execute VBoxManage commands
	cmd := exec.Command("VBoxManage", "createvm", "--name", name, "--register")
	if err := cmd.Run(); err != nil {
		resp.Diagnostics.AddError("Error creating VM", fmt.Sprintf("Failed to create VM: %s", err))
		return
	}

	cmd = exec.Command("VBoxManage", "modifyvm", name, "--memory", fmt.Sprintf("%d", memory), "--cpus", fmt.Sprintf("%d", cpus))
	if err := cmd.Run(); err != nil {
		resp.Diagnostics.AddError("Error configuring VM", fmt.Sprintf("Failed to configure VM: %s", err))
		return
	}

	cmd = exec.Command("VBoxManage", "storagectl", name, "--name", "SATA Controller", "--add", "sata", "--controller", "IntelAhci")
	if err := cmd.Run(); err != nil {
		resp.Diagnostics.AddError("Error adding storage controller", fmt.Sprintf("Failed to add storage controller: %s", err))
		return
	}

	cmd = exec.Command("VBoxManage", "storageattach", name, "--storagectl", "SATA Controller", "--port", "0", "--device", "0", "--type", "dvddrive", "--medium", isoPath)
	if err := cmd.Run(); err != nil {
		resp.Diagnostics.AddError("Error attaching ISO", fmt.Sprintf("Failed to attach ISO: %s", err))
		return
	}

	// Set the state so Terraform can track the resource
	data.Name = types.StringValue(name)
	data.ISOPath = types.StringValue(isoPath)
	data.Memory = types.Int64Value(memory)
	data.CPUs = types.Int64Value(cpus)

	// Save the new state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddWarning("VM Created", fmt.Sprintf("VM '%s' created successfully.", name))
}

func (r *VirtualBoxVMResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Implementation for reading the resource state
}

func (r *VirtualBoxVMResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Implementation for updating the resource
}

func (r *VirtualBoxVMResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Implementation for deleting the resource
}
