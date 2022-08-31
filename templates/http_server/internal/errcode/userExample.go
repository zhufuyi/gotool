package errcode

import "github.com/zhufuyi/pkg/gin/errcode"

const (
	// todo must be modified manually
	// 每个资源名称对应唯一编号，值范围建议1~1000，如果存在编号相同，启动服务会报错
	userExampleNO = 1
	// userExample对应的中文名称
	userExampleName = "userExample_cn_name"
)

// 服务级别错误码，有Err前缀
var (
	ErrCreateUserExample = errcode.NewError(genCode(userExampleNO)+1, "创建"+userExampleName+"失败") // todo 补充错误码注释，例如 20101
	ErrDeleteUserExample = errcode.NewError(genCode(userExampleNO)+2, "删除"+userExampleName+"失败")
	ErrUpdateUserExample = errcode.NewError(genCode(userExampleNO)+3, "更新"+userExampleName+"失败")
	ErrGetUserExample    = errcode.NewError(genCode(userExampleNO)+4, "获取"+userExampleName+"失败")
	// 每添加一个错误码，在上一个错误码基础上+1
)
