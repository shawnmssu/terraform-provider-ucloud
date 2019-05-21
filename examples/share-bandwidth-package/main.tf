provider "ucloud" {
  region = "cn-bj2"
}

resource "ucloud_eip" "foo" {
  name          = "tf-acc-share-bandwidth-package-Attachment"
  tag           = "tf-acc"
  internet_type = "bgp"

  # charge_mode   = "bandwidth"
}

resource "ucloud_share_bandwidth_package" "foo" {
  name            = "tf-acc-share-bandwidth-package-Attachment"
  share_bandwidth = 20
}

resource "ucloud_share_bandwidth_package_attachment" "foo" {
  eip_id                     = "${ucloud_eip.foo.id}"
  share_bandwidth_package_id = "${ucloud_share_bandwidth_package.foo.id}"
}
