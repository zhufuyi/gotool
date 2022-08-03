package errcode

import "github.com/zhufuyi/pkg/gin/errcode"

const (
	userExampleName = "用户" // userExample对应的名称
	userExampleNO   = 1    // 每个资源名称对应唯一编号，值范围建议1~1000
)

// 自定义业务错误码
var (
	CreateUserExampleErr = errcode.NewError(genCode(userExampleNO)+1, "创建"+userExampleName+"失败")
	DeleteUserExampleErr = errcode.NewError(genCode(userExampleNO)+2, "删除"+userExampleName+"失败")
	UpdateUserExampleErr = errcode.NewError(genCode(userExampleNO)+3, "更新"+userExampleName+"失败")
	GetUserExampleErr    = errcode.NewError(genCode(userExampleNO)+4, "获取"+userExampleName+"失败")
	// 每添加一个错误码，在上一个错误码基础上+1
)
