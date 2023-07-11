package iam

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/providers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func getProviderResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMNoVersionClient(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SberCloud IAM without version: %s", err)
	}
	return providers.Get(client, state.Primary.ID)
}

func TestAccIdentityProvider_basic(t *testing.T) {
	var provider providers.Provider
	var name = acceptance.RandomAccResourceName()
	resourceName := "sbercloud_identity_provider.provider_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&provider,
		getProviderResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProvider_saml(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "protocol", "saml"),
				),
			},
			{
				Config: testAccIdentityProvider_saml_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "protocol", "saml"),
					resource.TestCheckResourceAttr(resourceName, "status", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIdentityProvider_oidc(t *testing.T) {
	var provider providers.Provider
	var name = acceptance.RandomAccResourceName()
	resourceName := "sbercloud_identity_provider.provider_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&provider,
		getProviderResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProvider_oidc(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "protocol", "oidc"),
					resource.TestCheckResourceAttr(resourceName, "access_config.0.access_type", "program_console"),
					resource.TestCheckResourceAttr(resourceName, "access_config.0.client_id", "client_id_example"),
				),
			},
			{
				Config: testAccIdentityProvider_oidc_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "protocol", "oidc"),
					resource.TestCheckResourceAttr(resourceName, "status", "false"),
					resource.TestCheckResourceAttr(resourceName, "access_config.0.access_type", "program"),
					resource.TestCheckResourceAttr(resourceName, "access_config.0.client_id", "client_id_demo"),
				),
			},
		},
	})
}

func testAccIdentityProvider_saml(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_identity_provider" "provider_1" {
  name     = "%s"
  protocol = "saml"
}
`, name)
}

func testAccIdentityProvider_saml_update(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_identity_provider" "provider_1" {
  name     = "%s"
  protocol = "saml"
  status   = false
}
`, name)
}

func testAccIdentityProvider_oidc(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_identity_provider" "provider_1" {
  name        = "%s"
  protocol    = "oidc"
  description = "unit test"

  access_config {
    access_type            = "program_console"
    provider_url           = "https://accounts.example.com"
    client_id              = "client_id_example"
    authorization_endpoint = "https://accounts.example.com/o/oauth2/v2/auth"
    scopes                 = ["openid"]
    signing_key            = jsonencode(
    {
      keys = [
        {
          alg = "RS256"
          e   = "AQAB"
          kid = "d05ef20c4512645vv1..."
          kty = "RSA"
          n   = "cws_cnjiwsbvweolwn_-vnl..."
          use = "sig"
        },
      ]
    }
    )
  }
}
`, name)
}

func testAccIdentityProvider_oidc_update(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_identity_provider" "provider_1" {
  name        = "%s"
  protocol    = "oidc"
  status      = false
  description = "unit test"

  access_config {
    access_type            = "program"
    provider_url           = "https://accounts.example.com"
    client_id              = "client_id_demo"
    signing_key            = jsonencode(
    {
      keys = [
        {
          alg = "RS256"
          e   = "AQAB"
          kid = "d05ef20c4512645vv1..."
          kty = "RSA"
          n   = "cws_cnjiwsbvweolwn_-vnl..."
          use = "sig"
        },
      ]
    }
    )
  }
}
`, name)
}
