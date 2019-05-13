// +build ignore

package scaleway

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("scaleway_volume", &resource.Sweeper{
		Name: "scaleway_volume",
		F:    testSweepVolume,
	})
}

func testSweepVolume(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	scaleway := client.(*Client).scaleway
	log.Printf("[DEBUG] Destroying the volumes in (%s)", region)

	volumes, err := scaleway.GetVolumes()
	if err != nil {
		return fmt.Errorf("Error describing volumes in Sweeper: %s", err)
	}

	for _, volume := range *volumes {
		if err := scaleway.DeleteVolume(volume.Identifier); err != nil {
			return fmt.Errorf("Error deleting volume in Sweeper: %s", err)
		}
	}

	return nil
}

func TestAccScalewayVolume_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScalewayVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckScalewayVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayVolumeExists("scaleway_volume.test"),
					testAccCheckScalewayVolumeAttributes("scaleway_volume.test"),
				),
			},
		},
	})
}

func testAccCheckScalewayVolumeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).scaleway

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "scaleway" {
			continue
		}

		_, err := client.GetVolume(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Volume still exists")
		}
	}

	return nil
}

func testAccCheckScalewayVolumeAttributes(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Unknown resource: %s", n)
		}

		client := testAccProvider.Meta().(*Client).scaleway
		volume, err := client.GetVolume(rs.Primary.ID)

		if err != nil {
			return err
		}

		if volume.Name != "test" {
			return fmt.Errorf("volume has wrong name: %q", volume.Name)
		}
		if volume.Size != 2e+09 {
			return fmt.Errorf("volume has wrong size: %d", volume.Size)
		}
		if volume.VolumeType != "l_ssd" {
			return fmt.Errorf("volume has volume type: %q", volume.VolumeType)
		}

		return nil
	}
}

func testAccCheckScalewayVolumeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Volume ID is set")
		}

		client := testAccProvider.Meta().(*Client).scaleway
		volume, err := client.GetVolume(rs.Primary.ID)

		if err != nil {
			return err
		}

		if volume.Identifier != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		return nil
	}
}

var testAccCheckScalewayVolumeConfig = `
resource "scaleway_volume" "test" {
  name = "test"
  size_in_gb = 2
  type = "l_ssd"
}
`
