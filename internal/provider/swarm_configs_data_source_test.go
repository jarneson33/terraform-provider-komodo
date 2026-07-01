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

func TestAccSwarmConfigsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSwarmConfigsDataSourceConfig("tf-acc-swarm-configs-ds", "tf-acc-configs-ds"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.komodo_swarm_configs.test", "configs.#"),
					resource.TestCheckResourceAttrSet("data.komodo_swarm_configs.test", "configs.0.name"),
				),
			},
		},
	})
}

func testAccSwarmConfigsDataSourceConfig(swarmName, configName string) string {
	return fmt.Sprintf(`
resource "komodo_swarm" "src" {
  name = %q
}

resource "komodo_swarm_config" "src" {
  swarm = komodo_swarm.src.name
  name  = %q
  data  = "test-config-data"
}

data "komodo_swarm_configs" "test" {
  swarm      = komodo_swarm.src.name
  depends_on = [komodo_swarm_config.src]
}
`, swarmName, configName)
}

func TestUnitSwarmConfigsDataSource_configure(t *testing.T) {
	d := &SwarmConfigsDataSource{}
	resp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: "wrong"}, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected diagnostic error for wrong provider data type")
	}
}
