package provider

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

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
				// Default:     2048,
			},
			"cpus": schema.Int64Attribute{
				Required:    true,
				Description: "Number of CPUs.",
				// Default:     2,
			},
		},
	}
}

func (r *VirtualBoxVMResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data map[string]interface{}
	req.Plan.Get(ctx, &data)

	name := data["name"].(string)
	isoPath := data["iso_path"].(string)
	memory := data["memory"].(int64)
	cpus := data["cpus"].(int64)

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
