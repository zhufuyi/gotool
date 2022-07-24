package errcode

import "github.com/zhufuyi/pkg/gin/errcode"

// 公共错误码
var (
	Success             = errcode.Success
	InvalidParams       = errcode.InvalidParams
	Unauthorized        = errcode.Unauthorized
	InternalServerError = errcode.InternalServerError
	NotFound            = errcode.NotFound
	AlreadyExists       = errcode.AlreadyExists
	Timeout             = errcode.Timeout
	TooManyRequests     = errcode.TooManyRequests
	Forbidden           = errcode.Forbidden
	LimitExceed         = errcode.LimitExceed
)
