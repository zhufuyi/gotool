package errcode

import "github.com/zhufuyi/pkg/gin/errcode"

// 自定义业务错误码
var (
	CreateUserExampleErr = errcode.NewError(100101, "创建xxx失败")
	DeleteUserExampleErr = errcode.NewError(100102, "删除xxx失败")
	UpdateUserExampleErr = errcode.NewError(100103, "更新xxx失败")
	GetUserExampleErr    = errcode.NewError(100104, "获取xxx失败")
)
