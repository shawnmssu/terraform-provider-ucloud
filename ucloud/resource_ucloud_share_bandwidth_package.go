package ucloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func resourceUCloudShareBandwidthPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceUCloudShareBandwidthPackageCreate,
		Read:   resourceUCloudShareBandwidthPackageRead,
		Update: resourceUCloudShareBandwidthPackageUpdate,
		Delete: resourceUCloudShareBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateName,
			},

			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "month",
				ValidateFunc: validation.StringInSlice([]string{"year", "month", "dynamic"}, false),
			},

			"share_bandwidth": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(20, 5000),
			},

			"duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateDuration,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUCloudShareBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.unetconn
	req := conn.NewAllocateShareBandwidthRequest()
	req.ChargeType = ucloud.String(upperCamelCvt.unconvert(d.Get("charge_type").(string)))
	req.ShareBandwidth = ucloud.Int(d.Get("share_bandwidth").(int))
	if v, ok := d.GetOk("name"); ok {
		req.Name = ucloud.String(v.(string))
	} else {
		req.Name = ucloud.String(resource.PrefixedUniqueId("tf-share_bandwidth_package"))
	}

	if v, ok := d.GetOk("duration"); ok {
		req.Quantity = ucloud.Int(v.(int))
	} else {
		req.Quantity = ucloud.Int(1)
	}

	resp, err := conn.AllocateShareBandwidth(req)
	if err != nil {
		return fmt.Errorf("error on creating share bandwidth package, %s", err)
	}

	d.SetId(resp.ShareBandwidthId)

	return resourceUCloudShareBandwidthPackageRead(d, meta)
}

func resourceUCloudShareBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.unetconn

	d.Partial(true)

	if d.HasChange("share_bandwidth") && !d.IsNewResource() {
		req := conn.NewResizeShareBandwidthRequest()
		req.ShareBandwidthId = ucloud.String(d.Id())
		req.ShareBandwidth = ucloud.Int(d.Get("share_bandwidth").(int))

		_, err := conn.ResizeShareBandwidth(req)
		if err != nil {
			return fmt.Errorf("error on %s to share bandwidth package %q, %s", "ResizeShareBandwidth", d.Id(), err)
		}

		d.SetPartial("share_bandwidth")
	}

	d.Partial(false)

	return resourceUCloudShareBandwidthPackageRead(d, meta)
}

func resourceUCloudShareBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)

	sbp, err := client.describeShareBandwidthById(d.Id())

	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error on reading share bandwidth package %q, %s", d.Id(), err)
	}

	d.Set("name", sbp.Name)
	d.Set("charge_type", upperCamelCvt.convert(sbp.ChargeType))
	d.Set("share_bandwidth", sbp.ShareBandwidth)
	d.Set("create_time", timestampToString(sbp.CreateTime))
	d.Set("expire_time", timestampToString(sbp.ExpireTime))
	return nil
}

func resourceUCloudShareBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.unetconn

	req := conn.NewReleaseShareBandwidthRequest()
	req.ShareBandwidthId = ucloud.String(d.Id())

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.ReleaseShareBandwidth(req); err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on deleting share bandwidth package %q, %s", d.Id(), err))
		}

		_, err := client.describeShareBandwidthById(d.Id())
		if err != nil {
			if isNotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("error on reading share bandwidth package when deleting %q, %s", d.Id(), err))
		}

		return resource.RetryableError(fmt.Errorf("the specified share bandwidth package %q has not been deleted due to unknown error", d.Id()))
	})
}
