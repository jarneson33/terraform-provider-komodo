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

var _ datasource.DataSource = &SwarmSecretDataSource{}
var _ datasource.DataSourceWithValidateConfig = &SwarmSecretDataSource{}

func NewSwarmSecretDataSource() datasource.DataSource {
	return &SwarmSecretDataSource{}
}

type SwarmSecretDataSource struct {
	client *client.Client
}

type SwarmSecretDataSourceModel struct {
	Swarm types.String `tfsdk:"swarm"`
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
}

func (d *SwarmSecretDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_swarm_secret"
}

func (d *SwarmSecretDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads a Docker Swarm secret on a Komodo swarm by name or ID.",
		Attributes: map[string]schema.Attribute{
			"swarm": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The swarm name or ID to query.",
			},
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The swarm secret ID. One of `name` or `id` must be set.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The swarm secret name. One of `name` or `id` must be set.",
			},
		},
	}
}

func (d *SwarmSecretDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SwarmSecretDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data SwarmSecretDataSourceModel
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

func (d *SwarmSecretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SwarmSecretDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lookup := data.Name.ValueString()
	if lookup == "" {
		lookup = data.ID.ValueString()
	}

	tflog.Debug(ctx, "Reading swarm secret", map[string]interface{}{"swarm": data.Swarm.ValueString(), "lookup": lookup})

	items, err := d.client.ListSwarmSecrets(ctx, client.ListSwarmSecretsRequest{Swarm: data.Swarm.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list swarm secrets, got error: %s", err))
		return
	}

	var found *client.SwarmSecretListItem
	for i := range items {
		if (items[i].Name != nil && *items[i].Name == lookup) || (items[i].ID != nil && *items[i].ID == lookup) {
			found = &items[i]
			break
		}
	}

	if found == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Swarm secret %q not found on swarm %q", lookup, data.Swarm.ValueString()))
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
