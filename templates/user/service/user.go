package service

import (
	"crypto"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zhufuyi/goctl/templates/user/dao"
	"github.com/zhufuyi/goctl/templates/user/global"
	"github.com/zhufuyi/goctl/templates/user/model"

	"github.com/zhufuyi/pkg/gocrypto"
	"github.com/zhufuyi/pkg/jwt"
	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/mysql"
	"github.com/zhufuyi/pkg/snowflake"
	"golang.org/x/crypto/bcrypt"
)

const (
	ErrTypeUserOrPassword   = 1 // 用户名或密码错误
	ErrTypeUserNotActivated = 2 // 用户未激活
	ErrTypeInternal         = 3 // 服务内部错误

	LoginStateNo  = 1 // 用户未登录
	LoginStateYes = 2 // 用户已登录

	StatusNotActivated = 1 // 用户未激活
	StatusActivated    = 2 // 用户已激活
	StatusForbidden    = 3 // 用户已封禁
)

// CreateUserRequest 请求参数
type CreateUserRequest struct {
	// binding使用说明 https://github.com/go-playground/validator

	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"md5"`
}

// CreateUserResponse 返回参数
type CreateUserResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// CreateUser 创建
func (s *Service) CreateUser(req *CreateUserRequest) (*CreateUserResponse, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id := uint64(snowflake.NewID())
	data := &dao.User{
		ID:       id,
		Email:    req.Email,
		Password: string(password),
	}
	err = s.dao.CreateUser(data)
	if err != nil {
		return nil, err
	}

	idStr := strconv.FormatUint(id, 10)
	token, err := jwt.GenerateToken(idStr)
	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{
		ID:    idStr,
		Token: token,
	}, nil
}

// ActivateUserRequest 激活用户请求参数
type ActivateUserRequest struct {
	ID        uint64 `form:"id" binding:"gt=0"`
	URLParams string `form:"-" binding:"-"`
}

// ActivateUser 激活用户
func (s *Service) ActivateUser(req *ActivateUserRequest) (int, error) {
	isActivated := 0

	err := verifyURLParams(req.URLParams)
	if err != nil {
		return isActivated, err
	}

	status, err := s.dao.GetUserStatus(req.ID)
	if err != nil {
		return isActivated, err
	}

	// 已经激活过了
	if status != 0 {
		isActivated = 1
		return isActivated, errors.New("already activated")
	}

	data := &dao.User{
		ID:     req.ID,
		Status: StatusActivated, // 已激活
	}
	return isActivated, s.dao.UpdateUser(data)
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"md5"`
}

// LoginResponse 返回参数
type LoginResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// Login 登录
func (s *Service) Login(req *LoginRequest) (*LoginResponse, int, error) {
	// 判断用户名和密码是否匹配
	columns := []*mysql.Column{
		{
			Name:  "email",
			Value: req.Email,
		},
	}
	obj, err := s.dao.GetUser(columns)
	if err != nil {
		if err.Error() == mysql.ErrNotFound.Error() {
			return nil, ErrTypeUserOrPassword, err
		}
		return nil, ErrTypeInternal, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(obj.Password), []byte(req.Password))
	if err != nil {
		return nil, ErrTypeUserOrPassword, err
	}

	if obj.Status != StatusActivated { // 判断用户是否已激活状态
		return nil, ErrTypeUserNotActivated, errors.New("user not activated")
	}

	// 更新登录状态
	err = s.dao.UpdateUser(&dao.User{
		ID:         obj.ID,
		LoginState: LoginStateYes, // 已登录
	})
	if err != nil {
		return nil, ErrTypeInternal, err
	}

	// 生成token
	idStr := strconv.FormatUint(obj.ID, 10)
	token, err := jwt.GenerateToken(idStr)
	if err != nil {
		return nil, ErrTypeInternal, err
	}

	return &LoginResponse{
		ID:    idStr,
		Token: token,
	}, 0, nil
}

// Logout 登录
func (s *Service) Logout(id uint64) error {
	isLogined, err := s.IsLogined(id)
	if err != nil {
		return err
	}
	if !isLogined { // 未登录状态
		return nil
	}

	return s.dao.UpdateUser(&dao.User{
		ID:         id,
		LoginState: LoginStateNo, // 未登录
	})
}

// IsLogined 是否已经登录
func (s *Service) IsLogined(id uint64) (bool, error) {
	obj, err := s.dao.GetUserByID(id)
	if err != nil {
		return false, err
	}
	if obj.LoginState == LoginStateNo { // 未登录状态
		return false, nil
	}

	return true, nil
}

// DeleteUserRequest 删除一个id时，从url参数
type DeleteUserRequest struct {
	ID uint64 `form:"id" binding:"gt=0"`
}

// DeleteUsersRequest 删除多个id时，从body获取
type DeleteUsersRequest struct {
	IDs []uint64 `form:"ids" binding:"min=1"`
}

// DeleteUser 删除记录
func (s *Service) DeleteUser(ids ...uint64) error {
	if len(ids) == 1 {
		return s.dao.DeleteUser(ids[0])
	}

	// 批量删除
	var (
		total = len(ids)
		start = 0
		end   = 0
		size  = 100
	)
	for start < total {
		end += size
		if end > total {
			end = total
		}
		err := s.dao.DeleteUsers(ids[start:end])
		if err != nil {
			return err
		}
		start = end
	}

	return nil
}

// UpdateUserRequest 请求参数
type UpdateUserRequest struct {
	ID   uint64 `form:"id" binding:"-"`
	Name string `form:"name" binding:"required"`
	//Phone  uint64 `form:"phone" binding:"e164"`
	Age    int `form:"age" binding:"gt=0,lt=120"`
	Gender int `json:"gender" binding:"gte=1,lte=3"`
}

// UpdateUser 更新
func (s *Service) UpdateUser(req *UpdateUserRequest) error {
	return s.dao.UpdateUser(&dao.User{
		ID:   req.ID,
		Name: req.Name,
		//Phone:  req.Phone,
		Age:    req.Age,
		Gender: req.Gender,
	})
}

// GetUserRequest 请求参数
type GetUserRequest struct {
	ID uint64 `form:"id" binding:"gt=0"`
}

// GetUser 根据id获取一条记录
func (s *Service) GetUser(req *GetUserRequest) (*dao.UserSimpleInfo, error) {
	return s.dao.GetUserByID(req.ID)
}

// GetUsersRequest 请求参数
type GetUsersRequest struct {
	ID   uint64 `form:"id" binding:"gt=0"`
	Page int    `form:"page" binding:"gte=0"`
	Size int    `form:"size" binding:"gt=0"`
	Sort string `form:"sort" binding:"-"`

	// 参数填写方式一：从request请求url中获取参数(form.URLParams = c.Request.URL.RawQuery)，
	// 用来自动填充exp、logic的默认值，为了在url参数减少填写exp和logic的默认值，例如url参数?page=0&size=20&k=age&exp=gt&v=22&k=gender&v=1，表示查询年龄大于22岁的男性
	// 参数填写方式二：没有从请求url中获取参数，也就是ParamSrc为空时，请求url参数必须满足len(k)=len(v)=len(exp)=len(logic)，
	// 可以同时存在多个，也可以同时不存在，例如url参数?page=0&size=20&k=age&v=22&exp=gt&logic=and&k=gender&v=1&exp=eq&logic=and，也是表示查询年龄大于22岁的男性
	// 两种url参数都是合法，建议使用第一种
	URLParams string   `form:"-" binding:"-"`
	Keys      []string `form:"k" binding:"-"`
	Values    []string `form:"v" binding:"-"`
	Exps      []string `form:"exp" binding:"-"`
	Logics    []string `form:"logic" binding:"-"`
}

// GetUsers 获取多条记录
func (s *Service) GetUsers(req *GetUsersRequest) ([]*model.User, int, error) {
	var values []interface{}
	for _, v := range req.Values {
		values = append(values, v)
	}
	columns, err := mysql.GetColumns(req.Keys, values, req.Exps, req.Logics, req.URLParams)
	if err != nil {
		return nil, 0, err
	}

	return s.dao.GetUsersByColumns(columns, req.Page, req.Size, req.Sort)
}

// 通过post方法提交表单进行查询
type column struct {
	Name  string      `json:"name"`  // 列名
	Value interface{} `json:"value"` // 值
	Exp   string      `json:"exp"`   // 表达式，值为空时默认为eq，有eq、neq、gt、gte、lt、lte、like七种类型
	Logic string      `json:"logic"` // 逻辑类型，值为空时默认为and，有and、or两种类型
}

// GetUsersRequest2 请求参数
type GetUsersRequest2 struct {
	ID   uint64 `form:"id" binding:"gt=0"`
	Page int    `form:"page" binding:"gte=0" json:"page"`
	Size int    `form:"size" binding:"gt=0" json:"size"`
	Sort string `form:"sort" binding:"" json:"sort"`

	Columns []column `json:"columns"`
}

// GetUsers2 获取多条记录
func (s *Service) GetUsers2(req *GetUsersRequest2) ([]*model.User, int, error) {
	var columns []*mysql.Column
	for _, v := range req.Columns {
		if v.Value == "" {
			continue
		}
		columns = append(columns, &mysql.Column{
			Name:      v.Name,
			Value:     v.Value,
			ExpType:   v.Exp,
			LogicType: v.Logic,
		})
	}

	return s.dao.GetUsersByColumns(columns, req.Page, req.Size, req.Sort)
}

// ------------------------------------------------------------------------------------------

// SendEmail 发送邮件
func SendEmail(toEmail string, id string) error {
	kvs := map[string]interface{}{
		"id": id,
		"ts": time.Now().Unix(),
	}

	params, err := signURLParams(kvs)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s:%d/api/v1/auth/activate?%s", global.Conf.Address, global.Conf.ServerPort, params)

	//content := fmt.Sprintf("点击url激活用户 %s", url)
	//err = global.EmailCli.SendMessage(&email.Message{
	//	To:          []string{toEmail},
	//	Cc:          nil,
	//	Subject:     "欢迎注册xxx",
	//	ContentType: "text/plain",
	//	Content:     content,
	//	Attach:      "",
	//})
	//if err != nil {
	//	return err
	//}

	logger.Info("mock send email successfully", logger.String("url", url))

	return nil
}

// 签名url参数
func signURLParams(kvs map[string]interface{}) (string, error) {
	if len(kvs) == 0 {
		return "", errors.New("params is empty")
	}

	// 保证参数有序
	var sortParams []string
	for k, v := range kvs {
		kv := fmt.Sprintf("%s=%v", k, v)
		sortParams = append(sortParams, kv)
	}
	sort.Strings(sortParams)

	rawStr := strings.Join(sortParams, "&")
	ciphertext, err := gocrypto.Hash(crypto.BLAKE2b_384, []byte(rawStr))
	if err != nil {
		return "", err
	}

	return rawStr + "&sign=" + ciphertext, nil
}

// 验证url参数签名
func verifyURLParams(urlParams string) error {
	if strings.Count(urlParams, "&") == 0 || !strings.Contains(urlParams, "sign=") {
		return fmt.Errorf("'%s' is not legal, must contain sign", urlParams)
	}

	// 保证参数有序
	var sortParams []string
	params := strings.Split(urlParams, "&")
	signField := ""
	for _, param := range params {
		if strings.Contains(param, "sign=") {
			signField = param
			continue
		}
		sortParams = append(sortParams, param)
	}
	sort.Strings(sortParams)

	rawData := []byte(strings.Join(sortParams, "&"))
	ciphertext, err := gocrypto.Hash(crypto.BLAKE2b_384, rawData)
	if err != nil {
		return err
	}
	if ciphertext != strings.ReplaceAll(signField, "sign=", "") {
		return fmt.Errorf("the parameters have been modified")
	}

	return nil
}
