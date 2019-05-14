//Code is generated by ucloud code generator, don't modify it by hand, it will cause undefined behaviors.
//go:generate ucloud-gen-go-api UDB DescribeUDBInstanceState

package udb

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

// DescribeUDBInstanceStateRequest is request schema for DescribeUDBInstanceState action
type DescribeUDBInstanceStateRequest struct {
	request.CommonBase

	// [公共参数] 地域。 参见 [地域和可用区列表](../summary/regionlist.html)
	// Region *string `required:"true"`

	// [公共参数] 可用区。参见 [可用区列表](../summary/regionlist.html)
	// Zone *string `required:"false"`

	// [公共参数] 项目ID。不填写为默认项目，子帐号必须填写。 请参考[GetProjectList接口](../summary/get_project_list.html)
	// ProjectId *string `required:"false"`

	// 实例的Id,该值可以通过DescribeUDBInstance获取
	DBId *string `required:"true"`
}

// DescribeUDBInstanceStateResponse is response schema for DescribeUDBInstanceState action
type DescribeUDBInstanceStateResponse struct {
	response.CommonBase

	// DB状态标记 Init：初始化中；Fail：安装失败； Starting：启动中； Running ： 运行 ；Shutdown：关闭中； Shutoff ：已关闭； Delete：已删除； Upgrading：升级中； Promoting： 提升为独库进行中； Recovering： 恢复中； Recover fail：恢复失败。
	State string
}

// NewDescribeUDBInstanceStateRequest will create request of DescribeUDBInstanceState action.
func (c *UDBClient) NewDescribeUDBInstanceStateRequest() *DescribeUDBInstanceStateRequest {
	req := &DescribeUDBInstanceStateRequest{}

	// setup request with client config
	c.Client.SetupRequest(req)

	// setup retryable with default retry policy (retry for non-create action and common error)
	req.SetRetryable(true)
	return req
}

// DescribeUDBInstanceState - 获取UDB实例状态
func (c *UDBClient) DescribeUDBInstanceState(req *DescribeUDBInstanceStateRequest) (*DescribeUDBInstanceStateResponse, error) {
	var err error
	var res DescribeUDBInstanceStateResponse

	err = c.Client.InvokeAction("DescribeUDBInstanceState", req, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}
