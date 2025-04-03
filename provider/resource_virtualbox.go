package provider

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Define the struct that represents the VM resource data
type VirtualBoxVMResourceData struct {
	Name           types.String                `tfsdk:"name"`
	ISOPath        types.String                `tfsdk:"iso_path"`
	Memory         types.Int64                 `tfsdk:"memory"`
	CPUs           types.Int64                 `tfsdk:"cpus"`
	NetworkAdapter *VirtualBoxVMNetworkAdapter `tfsdk:"network_adapter"`
}
type VirtualBoxVMNetworkAdapter struct {
	Type          types.String `tfsdk:"type"`
	Device        types.String `tfsdk:"device"`
	HostInterface types.String `tfsdk:"host_interface"`
	IPv4Address   types.String `tfsdk:"ipv4_address"`
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
			"network_adapter": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:    true,
						Description: "Type of network adapter.",
					},
					"device": schema.StringAttribute{
						Optional:    true,
						Description: "Device of the network adapter.",
					},
					"host_interface": schema.StringAttribute{
						Optional:    true,
						Description: "Host interface for the network adapter.",
					},
					"ipv4_address": schema.StringAttribute{
						Computed:    true,
						Description: "IPv4 address assigned to the VM.",
					},
				},
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

	name := data.Name.ValueString()
	isoPath := data.ISOPath.ValueString()
	memory := data.Memory.ValueInt64()
	cpus := data.CPUs.ValueInt64()

	var netType, hostInterface, ipv4Address string

	if data.NetworkAdapter != nil {
		netType = data.NetworkAdapter.Type.ValueString()
		hostInterface = data.NetworkAdapter.HostInterface.ValueString()
		ipv4Address = data.NetworkAdapter.IPv4Address.ValueString()
	} else {
		netType = "nat"
	}

	if _, err := os.Stat(isoPath); os.IsNotExist(err) {
		resp.Diagnostics.AddError("Invalid ISO Path", fmt.Sprintf("The specified ISO file does not exist: %s", isoPath))
		return
	}

	// Execute VBoxManage commands to create and configure the VM
	cmds := [][]string{
		{"VBoxManage", "createvm", "--name", name, "--register"},
		{"VBoxManage", "modifyvm", name, "--memory", fmt.Sprintf("%d", memory), "--cpus", fmt.Sprintf("%d", cpus)},
		{"VBoxManage", "storagectl", name, "--name", "SATA Controller", "--add", "sata", "--controller", "IntelAhci"},
		{"VBoxManage", "storageattach", name, "--storagectl", "SATA Controller", "--port", "0", "--device", "0", "--type", "dvddrive", "--medium", isoPath},
		{"VBoxManage", "modifyvm", name, "--nic1", netType},
	}

	// Apply network settings
	if netType == "bridged" && hostInterface != "" {
		cmds = append(cmds, []string{"VBoxManage", "modifyvm", name, "--bridgeadapter1", hostInterface})
	}

	for _, cmdArgs := range cmds {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			resp.Diagnostics.AddError("Error executing VBoxManage", fmt.Sprintf("Command failed: %v, Error: %s", cmdArgs, stderr.String()))
			return
		}
	}

	// Save the new state
	data.Name = types.StringValue(name)
	data.ISOPath = types.StringValue(isoPath)
	data.Memory = types.Int64Value(memory)
	data.CPUs = types.Int64Value(cpus)

	if data.NetworkAdapter != nil {
		data.NetworkAdapter.Type = types.StringValue(netType)
		if hostInterface != "" {
			data.NetworkAdapter.HostInterface = types.StringValue(hostInterface)
		}
		data.NetworkAdapter.IPv4Address = types.StringValue(ipv4Address)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddWarning("VM Created", fmt.Sprintf("VM '%s' created successfully.", name))
}

func (r *VirtualBoxVMResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VirtualBoxVMResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()

	// Get the IP address from VBoxManage
	cmd := exec.Command("VBoxManage", "guestproperty", "get", name, "/VirtualBox/GuestInfo/Net/0/V4/IP")
	output, err := cmd.Output()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VM IP", fmt.Sprintf("Failed to get VM IP address: %s", err))
		return
	}

	outputStr := strings.TrimSpace(string(output))
	if strings.HasPrefix(outputStr, "Value: ") {
		data.NetworkAdapter.IPv4Address = types.StringValue(strings.TrimPrefix(outputStr, "Value: "))
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *VirtualBoxVMResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Implementation for updating the resource
}

func (r *VirtualBoxVMResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VirtualBoxVMResourceData

	// Decode the current state into the struct
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()

	// Power off the VM if running
	cmd := exec.Command("VBoxManage", "controlvm", name, "poweroff")
	if err := cmd.Run(); err != nil {
		resp.Diagnostics.AddWarning("VM Power Off Failed", fmt.Sprintf("Could not power off VM '%s'. It may not be running.", name))
	}

	// Unregister and delete the VM
	cmd = exec.Command("VBoxManage", "unregistervm", name, "--delete")
	if err := cmd.Run(); err != nil {
		resp.Diagnostics.AddError("Error deleting VM", fmt.Sprintf("Failed to delete VM '%s': %s", name, err))
		return
	}

	// Terraform removes the resource from state automatically after this function completes
	resp.Diagnostics.AddWarning("VM Deleted", fmt.Sprintf("VM '%s' has been successfully deleted.", name))
}
