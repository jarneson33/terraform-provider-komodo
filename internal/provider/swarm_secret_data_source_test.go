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

func TestAccSwarmSecretDataSource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSwarmSecretDataSourceConfig_byName("tf-acc-swarm-secret-ds", "tf-acc-secret-ds"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.komodo_swarm_secret.test", "name", "tf-acc-secret-ds"),
					resource.TestCheckResourceAttrSet("data.komodo_swarm_secret.test", "id"),
				),
			},
		},
	})
}

func TestAccSwarmSecretDataSource_bothSet_isError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSwarmSecretDataSourceConfig_bothSet(),
				ExpectError: regexp.MustCompile(`Only one of`),
			},
		},
	})
}

func TestAccSwarmSecretDataSource_neitherSet_isError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSwarmSecretDataSourceConfig_neitherSet(),
				ExpectError: regexp.MustCompile(`One of`),
			},
		},
	})
}

func testAccSwarmSecretDataSourceConfig_byName(swarmName, secretName string) string {
	return fmt.Sprintf(`
resource "komodo_swarm" "src" {
  name = %q
}

resource "komodo_swarm_secret" "src" {
  swarm = komodo_swarm.src.name
  name  = %q
  data  = "test-secret-data"
}

data "komodo_swarm_secret" "test" {
  swarm      = komodo_swarm.src.name
  name       = komodo_swarm_secret.src.name
  depends_on = [komodo_swarm_secret.src]
}
`, swarmName, secretName)
}

func testAccSwarmSecretDataSourceConfig_bothSet() string {
	return `
data "komodo_swarm_secret" "test" {
  swarm = "my-swarm"
  id    = "some-id"
  name  = "some-name"
}
`
}

func testAccSwarmSecretDataSourceConfig_neitherSet() string {
	return `
data "komodo_swarm_secret" "test" {
  swarm = "my-swarm"
}
`
}

func TestUnitSwarmSecretDataSource_configure(t *testing.T) {
	d := &SwarmSecretDataSource{}
	resp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: "wrong"}, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected diagnostic error for wrong provider data type")
	}
}
