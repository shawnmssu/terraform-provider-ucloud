---
layout: "ucloud"
page_title: "UCloud: ucloud_share_bandwidth_package"
sidebar_current: "docs-ucloud-resource-share-bandwidth-package"
description: |-
  Provides a Share Bandwidth Package resource.
---

# ucloud_share_bandwidth_package

Provides a Share Bandwidth Package resource. Share Bandwidth is a bandwidth model in which multiple instances share the total amount of network bandwidth.

## Example Usage

```hcl
resource "ucloud_share_bandwidth_package" "foo" {
	name              = "tf-example-share-bandwidth-package"
	bandwidth   = 20
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth` - (Required) Maximum bandwidth to the elastic public network, measured in Mbps (Mega bit per second). the ranges for bandwidth are: 20-5000.

- - -

* `name` - (Optional) The name of the Share Bandwidth Package which contains 1-63 characters and only support Chinese, English, numbers, '-', '_' and '.'. If not specified, terraform will auto-generate a name beginning with `tf-share-bandwidth-package`.
* `charge_type` - (Optional) Share Bandwidth Package charge type. Possible values are: `year` as pay by year, `month` as pay by month, `dynamic` as pay by hour (specific permission required). (Default: `month`).
* `duration` - (Optional) The duration that you will buy the resource. (Default: `1`). It is not required when `dynamic` (pay by hour), the value is `0` when `month`(pay by month) and the instance will be valid till the last day of that month.
* `remark` - (Optional) The remarks of the Share Bandwidth Package. (Default: `""`).
* `tag` - (Optional) A tag assigned to Share Bandwidth Package, which contains at most 63 characters and only support Chinese, English, numbers, '-', '_', and '.'. If it is not filled in or a empty string is filled in, then default tag will be assigned. (Default: `Default`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation of Share Bandwidth Package, formatted in RFC3339 time string.
* `expire_time` - The expiration time for Share Bandwidth Package, formatted in RFC3339 time string.

## Import

Share Bandwidth Package can be imported using the `id`, e.g.

```
$ terraform import ucloud_share_bandwidth_package.example bwshare-abc123456
```