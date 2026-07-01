// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	datasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSwarmConfigDataSource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSwarmConfigDataSourceConfig_byName("tf-acc-swarm-config-ds", "tf-acc-config-ds"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.komodo_swarm_config.test", "name", "tf-acc-config-ds"),
					resource.TestCheckResourceAttrSet("data.komodo_swarm_config.test", "id"),
				),
			},
		},
	})
}

func TestAccSwarmConfigDataSource_bothSet_isError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSwarmConfigDataSourceConfig_bothSet(),
				ExpectError: regexp.MustCompile(`Only one of`),
			},
		},
	})
}

func TestAccSwarmConfigDataSource_neitherSet_isError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSwarmConfigDataSourceConfig_neitherSet(),
				ExpectError: regexp.MustCompile(`One of`),
			},
		},
	})
}

func testAccSwarmConfigDataSourceConfig_byName(swarmName, configName string) string {
	return fmt.Sprintf(`
resource "komodo_swarm" "src" {
  name = %q
}

resource "komodo_swarm_config" "src" {
  swarm = komodo_swarm.src.name
  name  = %q
  data  = "test-config-data"
}

data "komodo_swarm_config" "test" {
  swarm      = komodo_swarm.src.name
  name       = komodo_swarm_config.src.name
  depends_on = [komodo_swarm_config.src]
}
`, swarmName, configName)
}

func testAccSwarmConfigDataSourceConfig_bothSet() string {
	return `
data "komodo_swarm_config" "test" {
  swarm = "my-swarm"
  id    = "some-id"
  name  = "some-name"
}
`
}

func testAccSwarmConfigDataSourceConfig_neitherSet() string {
	return `
data "komodo_swarm_config" "test" {
  swarm = "my-swarm"
}
`
}

func TestUnitSwarmConfigDataSource_configure(t *testing.T) {
	d := &SwarmConfigDataSource{}
	resp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: "wrong"}, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected diagnostic error for wrong provider data type")
	}
}
