// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/jarneson33/terraform-provider-komodo/internal/client"
)

func TestAccSwarmConfigResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSwarmConfigResourceConfig("tf-acc-swarm-config-basic", "tf-acc-config-basic", "initial-value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("komodo_swarm_config.test", "name", "tf-acc-config-basic"),
					resource.TestCheckResourceAttr("komodo_swarm_config.test", "swarm", "tf-acc-swarm-config-basic"),
					resource.TestCheckResourceAttrSet("komodo_swarm_config.test", "id"),
				),
			},
			{
				Config:   testAccSwarmConfigResourceConfig("tf-acc-swarm-config-basic", "tf-acc-config-basic", "initial-value"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccSwarmConfigResource_rotate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSwarmConfigResourceConfig("tf-acc-swarm-config-rotate", "tf-acc-config-rotate", "value-one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("komodo_swarm_config.test", "name", "tf-acc-config-rotate"),
				),
			},
			{
				Config: testAccSwarmConfigResourceConfig("tf-acc-swarm-config-rotate", "tf-acc-config-rotate", "value-two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("komodo_swarm_config.test", "name", "tf-acc-config-rotate"),
				),
			},
		},
	})
}

func TestAccSwarmConfigResource_importState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSwarmConfigResourceConfig("tf-acc-swarm-config-import", "tf-acc-config-import", "import-value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("komodo_swarm_config.test", "id"),
				),
			},
			{
				Config:                  testAccSwarmConfigResourceConfig("tf-acc-swarm-config-import", "tf-acc-config-import", "import-value"),
				ResourceName:            "komodo_swarm_config.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data"},
			},
		},
	})
}

func TestAccSwarmConfigResource_disappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSwarmConfigResourceConfig("tf-acc-swarm-config-disappears", "tf-acc-config-disappears", "disappear-value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("komodo_swarm_config.test", "id"),
					testAccSwarmConfigDisappears("komodo_swarm_config.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccSwarmConfigDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}
		swarm := rs.Primary.Attributes["swarm"]
		name := rs.Primary.Attributes["name"]

		c := client.NewClient(
			os.Getenv("KOMODO_ENDPOINT"),
			os.Getenv("KOMODO_USERNAME"),
			os.Getenv("KOMODO_PASSWORD"),
		)
		return c.RemoveSwarmConfigs(context.Background(), client.RemoveSwarmConfigsRequest{
			Swarm:   swarm,
			Configs: []string{name},
		})
	}
}

func testAccSwarmConfigResourceConfig(swarmName, configName, data string) string {
	return fmt.Sprintf(`
resource "komodo_swarm" "test" {
  name = %q
}

resource "komodo_swarm_config" "test" {
  swarm = komodo_swarm.test.name
  name  = %q
  data  = %q
}
`, swarmName, configName, data)
}

func wrongRawSwarmConfigPlan(t *testing.T, r *SwarmConfigResource) tfsdk.Plan {
	t.Helper()
	ctx := context.Background()
	schemaResp := &fwresource.SchemaResponse{}
	r.Schema(ctx, fwresource.SchemaRequest{}, schemaResp)
	return tfsdk.Plan{
		Raw:    tftypes.NewValue(tftypes.String, "invalid"),
		Schema: schemaResp.Schema,
	}
}

func wrongRawSwarmConfigState(t *testing.T, r *SwarmConfigResource) tfsdk.State {
	t.Helper()
	ctx := context.Background()
	schemaResp := &fwresource.SchemaResponse{}
	r.Schema(ctx, fwresource.SchemaRequest{}, schemaResp)
	return tfsdk.State{
		Raw:    tftypes.NewValue(tftypes.String, "invalid"),
		Schema: schemaResp.Schema,
	}
}

func TestUnitSwarmConfigResource_configure(t *testing.T) {
	t.Run("wrong_type", func(t *testing.T) {
		r := &SwarmConfigResource{}
		req := fwresource.ConfigureRequest{ProviderData: "not-a-client"}
		resp := &fwresource.ConfigureResponse{}
		r.Configure(context.Background(), req, resp)
		if !resp.Diagnostics.HasError() {
			t.Fatal("expected diagnostic error for wrong ProviderData type")
		}
	})
}

func TestUnitSwarmConfigResource_createPlanGetError(t *testing.T) {
	r := &SwarmConfigResource{client: &client.Client{}}
	req := fwresource.CreateRequest{Plan: wrongRawSwarmConfigPlan(t, r)}
	resp := &fwresource.CreateResponse{}
	r.Create(context.Background(), req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected diagnostic error for malformed plan")
	}
}

func TestUnitSwarmConfigResource_readStateGetError(t *testing.T) {
	r := &SwarmConfigResource{client: &client.Client{}}
	req := fwresource.ReadRequest{State: wrongRawSwarmConfigState(t, r)}
	resp := &fwresource.ReadResponse{}
	r.Read(context.Background(), req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected diagnostic error for malformed state")
	}
}

func TestUnitSwarmConfigResource_updatePlanGetError(t *testing.T) {
	r := &SwarmConfigResource{client: &client.Client{}}
	req := fwresource.UpdateRequest{Plan: wrongRawSwarmConfigPlan(t, r)}
	resp := &fwresource.UpdateResponse{}
	r.Update(context.Background(), req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected diagnostic error for malformed plan")
	}
}

func TestUnitSwarmConfigResource_deleteStateGetError(t *testing.T) {
	r := &SwarmConfigResource{client: &client.Client{}}
	req := fwresource.DeleteRequest{State: wrongRawSwarmConfigState(t, r)}
	resp := &fwresource.DeleteResponse{}
	r.Delete(context.Background(), req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected diagnostic error for malformed state")
	}
}

func TestUnitSwarmConfigResource_importState_invalidFormat(t *testing.T) {
	r := &SwarmConfigResource{}
	req := fwresource.ImportStateRequest{ID: "no-colon-separator"}
	resp := &fwresource.ImportStateResponse{}
	r.ImportState(context.Background(), req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected diagnostic error for invalid import ID format")
	}
}
