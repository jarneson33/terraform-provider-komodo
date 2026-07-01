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

var _ datasource.DataSource = &SwarmSecretsDataSource{}

func NewSwarmSecretsDataSource() datasource.DataSource {
	return &SwarmSecretsDataSource{}
}

type SwarmSecretsDataSource struct {
	client *client.Client
}

type SwarmSecretsDataSourceModel struct {
	Swarm   types.String                  `tfsdk:"swarm"`
	Secrets []SwarmSecretDataSourceModel2 `tfsdk:"secrets"`
}

type SwarmSecretDataSourceModel2 struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *SwarmSecretsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_swarm_secrets"
}

func (d *SwarmSecretsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists Docker Swarm secrets on a Komodo swarm.",
		Attributes: map[string]schema.Attribute{
			"swarm": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The swarm name or ID to query.",
			},
			"secrets": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of swarm secrets.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The swarm secret ID.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The swarm secret name.",
						},
					},
				},
			},
		},
	}
}

func (d *SwarmSecretsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SwarmSecretsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SwarmSecretsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Listing swarm secrets", map[string]interface{}{"swarm": data.Swarm.ValueString()})

	items, err := d.client.ListSwarmSecrets(ctx, client.ListSwarmSecretsRequest{Swarm: data.Swarm.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list swarm secrets, got error: %s", err))
		return
	}

	data.Secrets = make([]SwarmSecretDataSourceModel2, len(items))
	for i := range items {
		if items[i].ID != nil {
			data.Secrets[i].ID = types.StringValue(*items[i].ID)
		} else {
			data.Secrets[i].ID = types.StringValue("")
		}
		if items[i].Name != nil {
			data.Secrets[i].Name = types.StringValue(*items[i].Name)
		} else {
			data.Secrets[i].Name = types.StringValue("")
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
