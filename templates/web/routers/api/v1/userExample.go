package v1

import (
	"strconv"

	"github.com/zhufuyi/goctl/templates/web/errcode"
	"github.com/zhufuyi/goctl/templates/web/service"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/pkg/gin/render"
	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/mysql"
)

// CreateUserExample 创建
func CreateUserExample(c *gin.Context) {
	form := &service.CreateUserExampleRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	err = svc.CreateUserExample(form)
	if err != nil {
		logger.Error("CreateUserExample error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.CreateUserExampleErr)
		return
	}

	render.Success(c)
}

// DeleteUserExample 删除一条记录
func DeleteUserExample(c *gin.Context) {
	form := &service.DeleteUserExampleRequest{}
	err := c.ShouldBindQuery(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	err = svc.DeleteUserExample(form.ID)
	if err != nil {
		logger.Error("DeleteUserExample error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.DeleteUserExampleErr)
		return
	}

	render.Success(c)
}

// DeleteUserExamples 删除多条记录
func DeleteUserExamples(c *gin.Context) {
	form := &service.DeleteUserExamplesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	err = svc.DeleteUserExample(form.IDs...)
	if err != nil {
		logger.Error("DeleteUserExample error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.DeleteUserExampleErr)
		return
	}

	render.Success(c)
}

// UpdateUserExample 更新
func UpdateUserExample(c *gin.Context) {
	form := &service.UpdateUserExampleRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	err = svc.UpdateUserExample(form)
	if err != nil {
		logger.Error("CreateUserExample error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.UpdateUserExampleErr)
		return
	}

	render.Success(c)
}

// GetUserExample 获取一条记录
func GetUserExample(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	form := &service.GetUserExampleRequest{ID: id}
	err := c.ShouldBindQuery(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	userExample, err := svc.GetUserExample(form)
	if err != nil {
		if err.Error() == mysql.ErrNotFound.Error() {
			logger.Warn("GetUserExampleByID warn", logger.Err(err), logger.Any("form", form))
			render.Error(c, errcode.NotFound)
		} else {
			logger.Error("GetUserExampleByID error", logger.Err(err), logger.Any("form", form))
			render.Error(c, errcode.GetUserExampleErr)
		}
		return
	}

	render.Success(c, gin.H{"userExample": userExample})
}

// GetUserExamples 获取多条记录
// 通过url参数作为查询条件，支持任意多个字段，下面以user表为例子get请求参数，不同条件查询第0页20条记录，默认是id降序
// 没有条件查询 ?page=0&size=20
// 名称查询 ?page=0&size=20&k=name&v=张三
// 名称模糊查询 ?page=0&size=20&k=name&v=张&exp=like
// 年龄为18岁的男性 ?page=0&size=20&k=age&v=22&gender=1
// 年龄小于18或者大于60的人 ?page=0&size=20&k=age&v=18&exp=gt&logic=or&k=age&v=60&exp=lt
// 排序可以在参数添加sort字段，例如sort=id表示id升序，sort=-id表示id降序
func GetUserExamples(c *gin.Context) {
	form := &service.GetUserExamplesRequest{}
	err := c.ShouldBindQuery(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}
	form.URLParams = c.Request.URL.RawQuery

	svc := service.New(c.Request.Context())
	userExamples, total, err := svc.GetUserExamples(form)
	if err != nil {
		logger.Error("GetUserExampleByID error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.GetUserExampleErr)
		return
	}

	render.Success(c, gin.H{
		"userExamples": userExamples,
		"total":        total,
	})
}

// GetUserExamples2 通过post获取多条记录
func GetUserExamples2(c *gin.Context) {
	form := &service.GetUserExamplesRequest2{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	userExamples, total, err := svc.GetUserExamples2(form)
	if err != nil {
		logger.Error("GetUserExampleByID error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.GetUserExampleErr)
		return
	}

	render.Success(c, gin.H{
		"userExamples": userExamples,
		"total":        total,
	})
}
