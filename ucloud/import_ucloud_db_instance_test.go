package ucloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccUCloudDBInstance_import(t *testing.T) {
	resourceName := "ucloud_db_instance.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBInstanceConfig,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{
					"password",
					"duration",
					"backup_id",
					"vpc_id",
					"subnet_id",
				},
			},
		},
	})
}
