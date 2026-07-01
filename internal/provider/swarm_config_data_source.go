// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/jarneson33/terraform-provider-komodo/internal/client"
)

var _ datasource.DataSource = &SwarmConfigDataSource{}
var _ datasource.DataSourceWithValidateConfig = &SwarmConfigDataSource{}

func NewSwarmConfigDataSource() datasource.DataSource {
	return &SwarmConfigDataSource{}
}

type SwarmConfigDataSource struct {
	client *client.Client
}

type SwarmConfigDataSourceModel struct {
	Swarm types.String `tfsdk:"swarm"`
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
}

func (d *SwarmConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_swarm_config"
}

func (d *SwarmConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads a Docker Swarm config on a Komodo swarm by name or ID.",
		Attributes: map[string]schema.Attribute{
			"swarm": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The swarm name or ID to query.",
			},
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The swarm config ID. One of `name` or `id` must be set.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The swarm config name. One of `name` or `id` must be set.",
			},
		},
	}
}

func (d *SwarmConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = c
}

func (d *SwarmConfigDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data SwarmConfigDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if data.Name.IsUnknown() || data.ID.IsUnknown() {
		return
	}
	nameSet := !data.Name.IsNull() && !data.Name.IsUnknown()
	idSet := !data.ID.IsNull() && !data.ID.IsUnknown()
	if nameSet && idSet {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Only one of `name` or `id` may be set, not both.",
		)
		return
	}
	if !nameSet && !idSet {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"One of `name` or `id` must be set.",
		)
	}
}

func (d *SwarmConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SwarmConfigDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lookup := data.Name.ValueString()
	if lookup == "" {
		lookup = data.ID.ValueString()
	}

	tflog.Debug(ctx, "Reading swarm config", map[string]interface{}{"swarm": data.Swarm.ValueString(), "lookup": lookup})

	items, err := d.client.ListSwarmConfigs(ctx, client.ListSwarmConfigsRequest{Swarm: data.Swarm.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list swarm configs, got error: %s", err))
		return
	}

	var found *client.SwarmConfigListItem
	for i := range items {
		if (items[i].Name != nil && *items[i].Name == lookup) || (items[i].ID != nil && *items[i].ID == lookup) {
			found = &items[i]
			break
		}
	}

	if found == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Swarm config %q not found on swarm %q", lookup, data.Swarm.ValueString()))
		return
	}

	if found.ID != nil {
		data.ID = types.StringValue(*found.ID)
	} else {
		data.ID = types.StringValue("")
	}
	if found.Name != nil {
		data.Name = types.StringValue(*found.Name)
	} else {
		data.Name = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
