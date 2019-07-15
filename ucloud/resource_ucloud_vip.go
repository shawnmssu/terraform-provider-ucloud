package ucloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/hashicorp/terraform/helper/validation"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func resourceUCloudVIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceUCloudVIPCreate,
		Read:   resourceUCloudVIPRead,
		Update: resourceUCloudVIPUpdate,
		Delete: resourceUCloudVIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateName,
			},

			"tag": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      defaultTag,
				ValidateFunc: validateTag,
				StateFunc:    stateFuncTag,
			},

			"remark": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceUCloudVIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.unetconn

	req := conn.NewAllocateVIPRequest()
	req.VPCId = ucloud.String(d.Get("vpc_id").(string))
	req.SubnetId = ucloud.String(d.Get("subnet_id").(string))
	req.Count = ucloud.Int(1)

	if v, ok := d.GetOk("remark"); ok {
		req.Remark = ucloud.String(v.(string))
	}

	// if tag is empty string, use default tag
	if v, ok := d.GetOk("tag"); ok {
		req.Tag = ucloud.String(v.(string))
	} else {
		req.Tag = ucloud.String(defaultTag)
	}

	if v, ok := d.GetOk("name"); ok {
		req.Name = ucloud.String(v.(string))
	} else {
		req.Name = ucloud.String(resource.PrefixedUniqueId("tf-vip"))
	}

	resp, err := conn.AllocateVIP(req)
	if err != nil {
		return fmt.Errorf("error on creating vip, %s", err)
	}

	d.SetId(resp.VIPSet[0].VIPId)

	return resourceUCloudVIPRead(d, meta)
}

func resourceUCloudVIPUpdate(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*UCloudClient)
	//conn := client.unetconn
	//
	//d.Partial(true)
	//
	////if d.HasChange("share_bandwidth") && !d.IsNewResource() {
	////	req := conn.NewResizeVIPRequest()
	////	req.VIPId = ucloud.String(d.Id())
	////
	////	_, err := conn.ResizeVIP(req)
	////	if err != nil {
	////		return fmt.Errorf("error on %s to vip %q, %s", "ResizeVIP", d.Id(), err)
	////	}
	////
	////	d.SetPartial("share_bandwidth")
	////}
	////
	////d.Partial(false)

	return resourceUCloudVIPRead(d, meta)
}

func resourceUCloudVIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	vpcId := d.Get("vpc_id").(string)
	subnetId := d.Get("subnet_id").(string)

	vip, err := client.describeVIPByIdAndVPC(d.Id(), vpcId, subnetId)

	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error on reading vip %q, %s", d.Id(), err)
	}

	d.Set("name", vip.Name)
	d.Set("vpc_id", vip.VPCId)
	d.Set("subnet_id", vip.SubnetId)
	d.Set("ip_address", vip.VIP)
	return nil
}

func resourceUCloudVIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.unetconn
	vpcId := d.Get("vpc_id").(string)
	subnetId := d.Get("subnet_id").(string)

	req := conn.NewReleaseVIPRequest()
	req.VIPId = ucloud.String(d.Id())

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.ReleaseVIP(req); err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on deleting vip %q, %s", d.Id(), err))
		}
		//TODO: [4412 Error: subnet is not in this account] vpcId and subnetId already be deleted
		_, err := client.describeVIPByIdAndVPC(d.Id(), vpcId, subnetId)
		if err != nil {
			if isNotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("error on reading vip when deleting %q, %s", d.Id(), err))
		}

		return resource.RetryableError(fmt.Errorf("the specified vip %q has not been deleted due to unknown error", d.Id()))
	})
}
