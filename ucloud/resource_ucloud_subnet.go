package ucloud

import (
	"fmt"
	uerr "github.com/ucloud/ucloud-sdk-go/ucloud/error"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func resourceUCloudSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceUCloudSubnetCreate,
		Update: resourceUCloudSubnetUpdate,
		Read:   resourceUCloudSubnetRead,
		Delete: resourceUCloudSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRBlock,
			},

			"vpc_id": {
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
				ForceNew: true,
				Computed: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUCloudSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.vpcconn

	req := conn.NewCreateSubnetRequest()
	req.VPCId = ucloud.String(d.Get("vpc_id").(string))
	cidrBlock := d.Get("cidr_block").(string)

	// skip parse error, because has been validated at schema validator
	cidr, _ := parseCidrBlock(cidrBlock)
	req.Subnet = ucloud.String(cidr.Network)
	req.Netmask = ucloud.Int(cidr.Mask)

	if v, ok := d.GetOk("name"); ok {
		req.SubnetName = ucloud.String(v.(string))
	} else {
		req.SubnetName = ucloud.String(resource.PrefixedUniqueId("tf-subnet-"))
	}

	// if tag is empty string, use default tag
	if v, ok := d.GetOk("tag"); ok {
		req.Tag = ucloud.String(v.(string))
	} else {
		req.Tag = ucloud.String(defaultTag)
	}

	if v, ok := d.GetOk("remark"); ok {
		req.Remark = ucloud.String(v.(string))
	}

	resp, err := conn.CreateSubnet(req)
	if err != nil {
		return fmt.Errorf("error on creating subnet, %s", err)
	}

	d.SetId(resp.SubnetId)

	// after create subnet, we need to wait it initialized
	stateConf := subnetWaitForState(client, d.Id())

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error on waiting for subnet %q complete creating, %s", d.Id(), err)
	}

	return resourceUCloudSubnetRead(d, meta)
}

func resourceUCloudSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.vpcconn

	d.Partial(true)

	isChanged := false
	req := conn.NewUpdateSubnetAttributeRequest()
	req.SubnetId = ucloud.String(d.Id())

	if d.HasChange("name") && !d.IsNewResource() {
		isChanged = true
		req.Name = ucloud.String(d.Get("name").(string))
	}

	if d.HasChange("tag") && !d.IsNewResource() {
		isChanged = true

		// if tag is empty string, use default tag
		if v, ok := d.GetOk("tag"); ok {
			req.Tag = ucloud.String(v.(string))
		} else {
			req.Tag = ucloud.String(defaultTag)
		}
	}

	if isChanged {
		_, err := conn.UpdateSubnetAttribute(req)
		if err != nil {
			return fmt.Errorf("error on %s to subnet %q, %s", "UpdateSubnetAttribute", d.Id(), err)
		}

		d.SetPartial("name")
		d.SetPartial("tag")

		// after update subnet attribute, we need to wait it completed
		stateConf := subnetWaitForState(client, d.Id())
		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("error on waiting for %s complete to subnet %q, %s", "UpdateSubnetAttribute", d.Id(), err)
		}
	}

	d.Partial(false)

	return resourceUCloudSubnetRead(d, meta)
}

func resourceUCloudSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)

	subnetSet, err := client.describeSubnetById(d.Id())
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error on reading subnet %q, %s", d.Id(), err)
	}

	d.Set("name", subnetSet.SubnetName)
	d.Set("cidr_block", fmt.Sprintf("%s/%s", subnetSet.Subnet, subnetSet.Netmask))
	d.Set("vpc_id", subnetSet.VPCId)
	d.Set("tag", subnetSet.Tag)
	d.Set("remark", subnetSet.Remark)
	d.Set("create_time", timestampToString(subnetSet.CreateTime))

	return nil
}

func resourceUCloudSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.vpcconn

	reqDelete := conn.NewDeleteSubnetRequest()
	reqDelete.SubnetId = ucloud.String(d.Id())
	reqDes := conn.NewDescribeSubnetResourceRequest()
	reqDes.SubnetId = ucloud.String(d.Id())
	return resource.Retry(10*time.Minute, func() *resource.RetryError {
		respDes, err := conn.DescribeSubnetResource(reqDes)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on %s before deleting subnet %q, %s", "DescribeSubnetResource", d.Id(), err))
		}

		if len(respDes.DataSet) > 0 {
			var resourceData []string
			for _, v := range respDes.DataSet {
				resourceData = append(resourceData, v.ResourceId)
			}
			return resource.NonRetryableError(fmt.Errorf("error on deleting subnet %q, we find the resource %v bind to it", d.Id(), resourceData))
		}

		if _, err := conn.DeleteSubnet(reqDelete); err != nil {
			if uErr, ok := err.(uerr.Error); ok && uErr.Code() == 4411 {
				return resource.RetryableError(fmt.Errorf("error on deleting subnet %q, %s", d.Id(), err))
			}
			return resource.NonRetryableError(fmt.Errorf("error on deleting subnet %q, %s", d.Id(), err))
		}

		_, err = client.describeSubnetById(d.Id())
		if err != nil {
			if isNotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("error on reading subnet when deleting %q, %s", d.Id(), err))
		}

		return resource.RetryableError(fmt.Errorf("the specified subnet %q has not been deleted due to unknown error", d.Id()))
	})
}

func subnetWaitForState(client *UCloudClient, subnetId string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{statusPending},
		Target:     []string{statusInitialized},
		Timeout:    3 * time.Minute,
		Delay:      2 * time.Second,
		MinTimeout: 1 * time.Second,
		Refresh: func() (interface{}, string, error) {
			subnetSet, err := client.describeSubnetById(subnetId)
			if err != nil {
				if isNotFoundError(err) {
					return nil, statusPending, nil
				}
				return nil, "", err
			}

			return subnetSet, statusInitialized, nil
		},
	}
}
