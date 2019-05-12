// +build ignore

package scaleway

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccScalewayUserData_importBasic(t *testing.T) {
	resourceName := "scaleway_user_data.base"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScalewayUserDataDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckScalewayUserDataConfig,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
