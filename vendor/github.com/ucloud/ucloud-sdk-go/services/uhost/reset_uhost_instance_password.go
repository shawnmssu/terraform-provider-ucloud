//Code is generated by ucloud code generator, don't modify it by hand, it will cause undefined behaviors.
//go:generate ucloud-gen-go-api UHost ResetUHostInstancePassword

package uhost

import (
	"encoding/base64"

	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

// ResetUHostInstancePasswordRequest is request schema for ResetUHostInstancePassword action
type ResetUHostInstancePasswordRequest struct {
	request.CommonBase

	// [公共参数] 地域。 参见 [地域和可用区列表](../summary/regionlist.html)
	// Region *string `required:"true"`

	// [公共参数] 可用区。参见 [可用区列表](../summary/regionlist.html)
	// Zone *string `required:"false"`

	// [公共参数] 项目ID。不填写为默认项目，子帐号必须填写。 请参考[GetProjectList接口](../summary/get_project_list.html)
	// ProjectId *string `required:"false"`

	// UHost实例ID
	UHostId *string `required:"true"`

	// UHost新密码（密码格式使用BASE64编码）
	Password *string `required:"true"`
}

// ResetUHostInstancePasswordResponse is response schema for ResetUHostInstancePassword action
type ResetUHostInstancePasswordResponse struct {
	response.CommonBase

	// UHost实例ID
	UhostId string
}

// NewResetUHostInstancePasswordRequest will create request of ResetUHostInstancePassword action.
func (c *UHostClient) NewResetUHostInstancePasswordRequest() *ResetUHostInstancePasswordRequest {
	req := &ResetUHostInstancePasswordRequest{}

	// setup request with client config
	c.Client.SetupRequest(req)

	// setup retryable with default retry policy (retry for non-create action and common error)
	req.SetRetryable(true)
	return req
}

// ResetUHostInstancePassword - 重置UHost实例的管理员密码。
func (c *UHostClient) ResetUHostInstancePassword(req *ResetUHostInstancePasswordRequest) (*ResetUHostInstancePasswordResponse, error) {
	var err error
	var res ResetUHostInstancePasswordResponse
	var reqImmutable = *req
	reqImmutable.Password = ucloud.String(base64.StdEncoding.EncodeToString([]byte(ucloud.StringValue(req.Password))))

	err = c.Client.InvokeAction("ResetUHostInstancePassword", &reqImmutable, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}
