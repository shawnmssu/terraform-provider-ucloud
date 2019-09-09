package ucloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"strings"
	"time"
	//"github.com/hashicorp/terraform/helper/validation"
)

func resourceUCloudShareBandwidthPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceUCloudShareBandwidthPackageAttachmentCreate,
		Read:   resourceUCloudShareBandwidthPackageAttachmentRead,
		Delete: resourceUCloudShareBandwidthPackageAttachmentDelete,

		Schema: map[string]*schema.Schema{

			"share_bandwidth_package_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"eip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			//"release_bandwidth": {
			//	Type:     schema.TypeInt,
			//	Default:  1,
			//},
		},
	}
}

func resourceUCloudShareBandwidthPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.unetconn

	eipId := d.Get("eip_id").(string)
	bandwidthPackageId := d.Get("share_bandwidth_package_id").(string)
	req := conn.NewAssociateEIPWithShareBandwidthRequest()
	req.EIPIds = []string{eipId}
	req.ShareBandwidthId = ucloud.String(bandwidthPackageId)

	_, err := conn.AssociateEIPWithShareBandwidth(req)
	if err != nil {
		return fmt.Errorf("error on creating share bandwidth package attachment, %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", bandwidthPackageId, eipId))
	return resourceUCloudShareBandwidthPackageAttachmentRead(d, meta)
}

func resourceUCloudShareBandwidthPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)

	resp, err := client.describeShareBandwidthAttachmentById(d.Id())
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error on reading share bandwidth package attachment %q, %s", d.Id(), err)
	}

	d.Set("bandwidth_package_id", resp.bandwidthPackageId)
	d.Set("eip_id", resp.eipId)

	return nil
}

func resourceUCloudShareBandwidthPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.unetconn

	p := strings.Split(d.Id(), ":")
	req := conn.NewDisassociateEIPWithShareBandwidthRequest()
	req.ShareBandwidthId = ucloud.String(p[0])
	req.Bandwidth = ucloud.Int(1)
	req.EIPIds = []string{p[1]}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.DisassociateEIPWithShareBandwidth(req); err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on deleting share bandwidth package attachment %q, %s", d.Id(), err))
		}

		_, err := client.describeShareBandwidthAttachmentById(d.Id())
		if err != nil {
			if isNotFoundError(err) {
				return nil
			}

			return resource.NonRetryableError(fmt.Errorf("error on reading share bandwidth package attachment when deleting %q, %s", d.Id(), err))
		}

		return resource.RetryableError(fmt.Errorf("the specified share bandwidth package attachment %q has not been deleted due to unknown error", d.Id()))
	})
}
