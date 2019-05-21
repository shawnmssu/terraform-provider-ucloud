package ucloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/ucloud/ucloud-sdk-go/services/unet"
)

func TestAccUCloudShareBandwidthPackageAttachment_basic(t *testing.T) {
	var eip unet.UnetEIPSet
	var sbp unet.UnetShareBandwidthSet

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ucloud_share_bandwidth_package_attachment.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckShareBandwidthPackageAttachmentDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccShareBandwidthPackageAttachmentConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists("ucloud_eip.foo", &eip),
					testAccCheckShareBandwidthPackageExists("ucloud_share_bandwidth_package.foo", &sbp),
					testAccCheckShareBandwidthPackageAttachmentExists("ucloud_share_bandwidth_package_attachment.foo", &eip, &sbp),
				),
			},
		},
	})
}

func testAccCheckShareBandwidthPackageAttachmentExists(n string, eip *unet.UnetEIPSet, sbp *unet.UnetShareBandwidthSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("share bandwidth package attachment id is empty")
		}

		client := testAccProvider.Meta().(*UCloudClient)
		return resource.Retry(3*time.Minute, func() *resource.RetryError {
			d, err := client.describeShareBandwidthAttachmentById(rs.Primary.ID)

			if err != nil {
				return resource.NonRetryableError(err)
			}

			if d.bandwidthPackageId == sbp.ShareBandwidthId && d.eipId == eip.EIPId {
				return nil
			}

			return resource.NonRetryableError(fmt.Errorf("share bandwidth package attachment not found"))
		})
	}
}

func testAccCheckShareBandwidthPackageAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ucloud_share_bandwidth_package_attachment" {
			continue
		}

		client := testAccProvider.Meta().(*UCloudClient)
		d, err := client.describeShareBandwidthAttachmentById(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			if isNotFoundError(err) {
				continue
			}
			return err
		}

		if d.eipId == rs.Primary.Attributes["eip_id"] {
			return fmt.Errorf("eip associatoin still exists")
		}
	}

	return nil
}

const testAccShareBandwidthPackageAttachmentConfig = `
resource "ucloud_eip" "foo" {
	name          = "tf-acc-share-bandwidth-package-Attachment"
	tag           = "tf-acc"
	internet_type = "bgp"
	bandwidth     = 1
	duration      = 1
}

resource "ucloud_share_bandwidth_package" "foo" {
	name              = "tf-acc-share-bandwidth-package-Attachment"
	share_bandwidth   = 20
}

resource "ucloud_share_bandwidth_package_attachment" "foo" {
	eip_id                       = "${ucloud_eip.foo.id}"
	share_bandwidth_package_id   = "${ucloud_share_bandwidth_package.foo.id}"
}
`
