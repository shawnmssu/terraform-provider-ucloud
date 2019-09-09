package ucloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/ucloud/ucloud-sdk-go/services/unet"
)

func TestAccUCloudShareBandwidthPackage_basic(t *testing.T) {
	var sbpSet unet.UnetShareBandwidthSet

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ucloud_share_bandwidth_package.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckShareBandwidthPackageDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccShareBandwidthPackageConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckShareBandwidthPackageExists("ucloud_share_bandwidth_package.foo", &sbpSet),
					testAccCheckShareBandwidthPackageAttributes(&sbpSet),
					resource.TestCheckResourceAttr("ucloud_share_bandwidth_package.foo", "bandwidth", "20"),
				),
			},

			{
				Config: testAccShareBandwidthPackageConfigUpdate,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckShareBandwidthPackageExists("ucloud_share_bandwidth_package.foo", &sbpSet),
					testAccCheckShareBandwidthPackageAttributes(&sbpSet),
					resource.TestCheckResourceAttr("ucloud_share_bandwidth_package.foo", "bandwidth", "25"),
				),
			},
		},
	})
}

func testAccCheckShareBandwidthPackageExists(n string, sbpSet *unet.UnetShareBandwidthSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("share bandwidth package id is empty")
		}

		client := testAccProvider.Meta().(*UCloudClient)
		ptr, err := client.describeShareBandwidthById(rs.Primary.ID)

		log.Printf("[INFO] share bandwidth package id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		*sbpSet = *ptr
		return nil
	}
}

func testAccCheckShareBandwidthPackageAttributes(sbpSet *unet.UnetShareBandwidthSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if sbpSet.ShareBandwidthId == "" {
			return fmt.Errorf("share bandwidth package id is empty")
		}
		return nil
	}
}

func testAccCheckShareBandwidthPackageDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ucloud_share_bandwidth_package" {
			continue
		}

		client := testAccProvider.Meta().(*UCloudClient)
		d, err := client.describeShareBandwidthById(rs.Primary.ID)

		if err != nil {
			if isNotFoundError(err) {
				continue
			}
			return err
		}

		if d.ShareBandwidthId != "" {
			return fmt.Errorf("share bandwidth package still exist")
		}
	}

	return nil
}

const testAccShareBandwidthPackageConfig = `
resource "ucloud_share_bandwidth_package" "foo" {
	name           = "tf-acc-share_bandwidth_package-basic"
	bandwidth 	   = 20
}
`

const testAccShareBandwidthPackageConfigUpdate = `
resource "ucloud_share_bandwidth_package" "foo" {
	name              = "tf-acc-share_bandwidth_package-basic"
	bandwidth 		  = 25
}
`
