package errcode

import "github.com/zhufuyi/pkg/gin/errcode"

const (
	userName = "用户" // user对应的名称
	userNO   = 1    // 每个资源名称对应唯一编号，值范围建议1~1000
)

// 自定义业务错误码
var (
	CreateUserErr          = errcode.NewError(genCode(userNO)+1, "创建"+userName+"失败")
	DeleteUserErr          = errcode.NewError(genCode(userNO)+2, "删除"+userName+"失败")
	UpdateUserErr          = errcode.NewError(genCode(userNO)+3, "更新"+userName+"失败")
	GetUserErr             = errcode.NewError(genCode(userNO)+4, "获取"+userName+"失败")
	RegisterErr            = errcode.NewError(genCode(userNO)+5, "注册失败")
	ActivateUserErr        = errcode.NewError(genCode(userNO)+6, "激活用户失败")
	NotActivateUserErr     = errcode.NewError(genCode(userNO)+7, "用户未激活")
	AlreadyActivateUserErr = errcode.NewError(genCode(userNO)+8, "已经激活过了")
	LoginErr               = errcode.NewError(genCode(userNO)+9, "用户或密码错误")
	SendEmailErr           = errcode.NewError(genCode(userNO)+10, "发送邮件失败")
	LogoutErr              = errcode.NewError(genCode(userNO)+11, "登出失败")
	LoginStateNoErr        = errcode.NewError(genCode(userNO)+12, "请先登录")
	// 每添加一个错误码，在上一个错误码基础上+1
)
