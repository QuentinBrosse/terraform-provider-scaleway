package scaleway

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"scaleway": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
	testAccPreCheck(t)
}

func testAccPreCheck(t *testing.T) {
	if ak, exist := scwConfig.GetSecretKey(); !exist {
		t.Logf(">>>> %s\n", ak)
		t.Fatal("the Scaleway token must be set for acceptance tests.")
	}

	if oi, exist := scwConfig.GetDefaultOrganizationID(); !exist {
		t.Logf(">>>> %s\n", oi)
		t.Fatal("the Scaleway organization ID must be set for acceptance tests.")
	}
}
