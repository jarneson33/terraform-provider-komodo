// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"testing"

	datasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSwarmSecretsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSwarmSecretsDataSourceConfig("tf-acc-swarm-secrets-ds", "tf-acc-secrets-ds"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.komodo_swarm_secrets.test", "secrets.#"),
					resource.TestCheckResourceAttrSet("data.komodo_swarm_secrets.test", "secrets.0.name"),
				),
			},
		},
	})
}

func testAccSwarmSecretsDataSourceConfig(swarmName, secretName string) string {
	return fmt.Sprintf(`
resource "komodo_swarm" "src" {
  name = %q
}

resource "komodo_swarm_secret" "src" {
  swarm = komodo_swarm.src.name
  name  = %q
  data  = "test-secret-data"
}

data "komodo_swarm_secrets" "test" {
  swarm      = komodo_swarm.src.name
  depends_on = [komodo_swarm_secret.src]
}
`, swarmName, secretName)
}

func TestUnitSwarmSecretsDataSource_configure(t *testing.T) {
	d := &SwarmSecretsDataSource{}
	resp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: "wrong"}, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected diagnostic error for wrong provider data type")
	}
}
