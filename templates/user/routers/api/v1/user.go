package v1

import (
	"strconv"

	"github.com/zhufuyi/goctl/templates/user/errcode"
	"github.com/zhufuyi/goctl/templates/user/service"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/pkg/gin/render"
	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/mysql"
)

// Register 注册
func Register(c *gin.Context) {
	form := &service.CreateUserRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	resp, err := svc.CreateUser(form)
	if err != nil {
		logger.Error("CreateUser error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.RegisterErr)
		return
	}

	// 发送邮件
	err = service.SendEmail(form.Email, resp.ID)
	if err != nil {
		logger.Error("SendEmail error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.SendEmailErr)
		return
	}

	render.Success(c)
}

// ActivateUser 激活用户，通过收件邮箱点击链接
func ActivateUser(c *gin.Context) {
	form := &service.ActivateUserRequest{}
	err := c.ShouldBindQuery(form)
	if err != nil {
		logger.Error("ShouldBindQuery error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}
	form.URLParams = c.Request.URL.RawQuery

	svc := service.New(c.Request.Context())
	isActivated, err := svc.ActivateUser(form)
	if err != nil {
		logger.Error("ActivateUser error", logger.Err(err), logger.String("params", c.Request.URL.RawQuery))
		if isActivated > 0 {
			render.Error(c, errcode.AlreadyActivateUserErr)
			return
		}
		render.Error(c, errcode.ActivateUserErr)
		return
	}

	render.Success(c)
}

// Login 登录
func Login(c *gin.Context) {
	form := &service.LoginRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	resp, errType, err := svc.Login(form)
	if err != nil {
		logger.Error("Login error", logger.Err(err), logger.String("email", form.Email))
		ec := errcode.InternalServerError
		switch errType {
		case service.ErrTypeUserOrPassword: // 用户或密码错误
			ec = errcode.LoginErr
		case service.ErrTypeUserNotActivated: // 用户存在，但未激活
			ec = errcode.NotActivateUserErr
		}
		render.Error(c, ec)
		return
	}

	render.Success(c, gin.H{"authInfo": resp})
}

// Logout 退出登录
func Logout(c *gin.Context) {
	id, isAbout := getIDFromPath(c)
	if isAbout {
		return
	}

	svc := service.New(c.Request.Context())
	isAbout = checkIsLogin(c, svc, id)
	if isAbout {
		return
	}
	err := svc.Logout(id)
	if err != nil {
		logger.Error("Logout error", logger.Err(err), logger.Uint64("id", id))
		render.Error(c, errcode.LogoutErr)
		return
	}

	render.Success(c)
}

// CreateUser 创建
func CreateUser(c *gin.Context) {
	form := &service.CreateUserRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	resp, err := svc.CreateUser(form)
	if err != nil {
		logger.Error("CreateUser error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.CreateUserErr)
		return
	}

	render.Success(c, resp)
}

// DeleteUser 删除一条记录
func DeleteUser(c *gin.Context) {
	id, isAbout := getIDFromPath(c)
	if isAbout {
		return
	}

	form := &service.DeleteUserRequest{ID: id}
	svc := service.New(c.Request.Context())
	isAbout = checkIsLogin(c, svc, id)
	if isAbout {
		return
	}
	err := svc.DeleteUser(form.ID)
	if err != nil {
		logger.Error("DeleteUser error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.DeleteUserErr)
		return
	}

	render.Success(c)
}

// DeleteUsers 删除多条记录
func DeleteUsers(c *gin.Context) {
	form := &service.DeleteUsersRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	err = svc.DeleteUser(form.IDs...)
	if err != nil {
		logger.Error("DeleteUser error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.DeleteUserErr)
		return
	}

	render.Success(c)
}

// UpdateUser 更新
func UpdateUser(c *gin.Context) {
	id, isAbout := getIDFromPath(c)
	if isAbout {
		return
	}

	form := &service.UpdateUserRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}
	form.ID = id

	svc := service.New(c.Request.Context())
	isAbout = checkIsLogin(c, svc, id)
	if isAbout {
		return
	}
	err = svc.UpdateUser(form)
	if err != nil {
		logger.Error("CreateUser error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.UpdateUserErr)
		return
	}

	render.Success(c)
}

// GetUser 获取一条记录
func GetUser(c *gin.Context) {
	id, isAbout := getIDFromPath(c)
	if isAbout {
		return
	}

	form := &service.GetUserRequest{ID: id}
	svc := service.New(c.Request.Context())
	user, err := svc.GetUser(form)
	if err != nil {
		logger.Error("GetUserByID error", logger.Err(err), logger.Any("form", form))
		if err.Error() == mysql.ErrNotFound.Error() {
			render.Error(c, errcode.NotFound)
		} else {
			render.Error(c, errcode.GetUserErr)
		}
		return
	}
	if user.LoginState == service.LoginStateNo {
		render.Error(c, errcode.LoginStateNoErr)
		return
	}

	render.Success(c, gin.H{"user": user})
}

// GetUsers 获取多条记录
// 通过url参数作为查询条件，支持任意多个字段，下面以user表为例子get请求参数，不同条件查询第0页20条记录，默认是id降序
// 没有条件查询 ?page=0&size=20
// 名称查询 ?page=0&size=20&k=name&v=张三
// 名称模糊查询 ?page=0&size=20&k=name&v=张&exp=like
// 年龄为18岁的男性 ?page=0&size=20&k=age&v=22&gender=1
// 年龄小于18或者大于60的人 ?page=0&size=20&k=age&v=18&exp=gt&logic=or&k=age&v=60&exp=lt
// 排序可以在参数添加sort字段，例如sort=id表示id升序，sort=-id表示id降序
func GetUsers(c *gin.Context) {
	form := &service.GetUsersRequest{}
	err := c.ShouldBindQuery(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}
	form.URLParams = c.Request.URL.RawQuery

	svc := service.New(c.Request.Context())
	isAbout := checkIsLogin(c, svc, form.ID)
	if isAbout {
		return
	}
	users, total, err := svc.GetUsers(form)
	if err != nil {
		logger.Error("GetUserByID error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.GetUserErr)
		return
	}

	render.Success(c, gin.H{
		"users": users,
		"total": total,
	})
}

// GetUsers2 通过post获取多条记录
func GetUsers2(c *gin.Context) {
	form := &service.GetUsersRequest2{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Error("ShouldBindJSON error: ", logger.Err(err))
		render.Error(c, errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	isAbout := checkIsLogin(c, svc, form.ID)
	if isAbout {
		return
	}
	users, total, err := svc.GetUsers2(form)
	if err != nil {
		logger.Error("GetUserByID error", logger.Err(err), logger.Any("form", form))
		render.Error(c, errcode.GetUserErr)
		return
	}

	render.Success(c, gin.H{
		"users": users,
		"total": total,
	})
}

func getIDFromPath(c *gin.Context) (uint64, bool) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		logger.Error("ParseUint error: ", logger.String("idStr", idStr))
		render.Error(c, errcode.InvalidParams)
		return 0, true
	}

	return id, false
}

func checkIsLogin(c *gin.Context, svc *service.Service, id uint64) bool {
	isLogined, err := svc.IsLogined(id)
	if err != nil {
		logger.Error("IsLogined error", logger.Err(err), logger.Any("id", id))
		render.Error(c, errcode.LoginStateNoErr)
		return true
	}
	if !isLogined {
		render.Error(c, errcode.LoginStateNoErr)
		return true
	}

	return false
}
