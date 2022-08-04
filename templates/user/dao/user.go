package dao

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/zhufuyi/goctl/templates/user/model"

	"github.com/zhufuyi/pkg/mysql"
)

// User dao 对象
type User struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Phone      uint64 `json:"phone"`
	Age        int    `json:"age"`
	Gender     int    `json:"gender"`
	Status     int    `json:"status"`     // 用户状态，1:未激活，2:已激活，3:封禁
	LoginState int    `json:"loginState"` // 登录状态，1:未登录，2:已登录
}

// CreateUser 创建一条记录
func (d *Dao) CreateUser(param *User) error {
	data := &model.User{
		Model:      mysql.Model{ID: param.ID},
		Name:       param.Name,
		Password:   param.Password,
		Email:      param.Email,
		Phone:      param.Phone,
		Age:        param.Age,
		Gender:     param.Gender,
		LoginState: param.LoginState,
	}
	return data.Create(d.db)
}

// CreateUsers 创建多条记录
func (d *Dao) CreateUsers(params []*User) (int, error) {
	errStr := []string{}
	count := 0
	for _, o := range params {
		err := d.CreateUser(o)
		if err != nil {
			errStr = append(errStr, err.Error())
			continue
		}
		count++
	}

	if len(errStr) == 0 {
		return count, nil
	}

	return count, errors.New(strings.Join(errStr, " || "))
}

// DeleteUser 根据id删除一条记录
func (d *Dao) DeleteUser(id uint64) error {
	obj := &model.User{}
	obj.ID = id
	return obj.DeleteByID(d.db)
}

// DeleteUsers 根据id删除多条记录
func (d *Dao) DeleteUsers(ids []uint64) error {
	obj := &model.User{}
	return obj.Delete(d.db, "id IN (?)", ids)
}

// UpdateUser 更新记录
func (d *Dao) UpdateUser(param *User) error {
	obj := &model.User{}
	update := mysql.KV{}

	if param.Name != "" {
		update["name"] = param.Name
	}
	if param.Password != "" {
		update["password"] = param.Password
	}
	//if param.Email != "" {
	//	update["email"] = param.Email
	//}
	if param.Phone > 0 {
		update["phone"] = param.Phone
	}
	if param.Age > 0 {
		update["age"] = param.Age
	}
	if param.Gender >= 1 && param.Gender <= 3 {
		update["gender"] = param.Gender
	}
	if param.Status > 0 && param.Status < 4 {
		update["status"] = param.Status
	}
	if param.LoginState > 0 && param.LoginState < 4 {
		update["login_state"] = param.LoginState
	}

	return obj.Updates(d.db, update, "id = ?", param.ID)
}

// UserSimpleInfo dao 对象
type UserSimpleInfo struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      uint64    `json:"phone"`
	Age        int       `json:"age"`
	Gender     int       `json:"gender"`
	LoginState int       `json:"loginState"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func convert2SimpleInfo(obj *model.User) *UserSimpleInfo {
	return &UserSimpleInfo{
		ID:         strconv.FormatUint(obj.ID, 10),
		Name:       obj.Name,
		Email:      obj.Email,
		Phone:      obj.Phone,
		Age:        obj.Age,
		Gender:     obj.Gender,
		LoginState: obj.LoginState,
		CreatedAt:  obj.CreatedAt,
		UpdatedAt:  obj.UpdatedAt,
	}
}

// GetUserByID 根据id获取一条记录
func (d *Dao) GetUserByID(id uint64) (*UserSimpleInfo, error) {
	obj := &model.User{}
	// todo 从缓存读取

	// 从mysql读取
	err := obj.GetByID(d.db, id)
	if err != nil {
		return nil, err
	}

	return convert2SimpleInfo(obj), err
}

// GetUser 根据任意列获取一条记录
func (d *Dao) GetUser(columns []*mysql.Column) (*model.User, error) {
	obj := &model.User{}

	query, args, err := mysql.GetQueryConditions(columns)
	if err != nil {
		return nil, err
	}
	err = obj.Get(d.db, query, args...)
	if err != nil {
		return nil, err
	}

	return obj, err
}

// GetUserStatus 获取用户状态
func (d *Dao) GetUserStatus(id uint64) (int, error) {
	obj := &model.User{}
	// todo 从缓存读取

	// 从mysql读取
	err := obj.GetByID(d.db, id)
	if err != nil {
		return 0, err
	}

	return obj.Status, nil
}

// GetUsersByColumns 根据列信息筛选多条记录
// columns 列信息，列名称、列值、表达式，列之间逻辑关系
// page表示页码，从0开始, size表示每页行数, sort排序字段，默认是id倒叙，可以在字段前添加-号表示倒序，无-号表示升序
// 查询年龄大于20的男性示例：
//	columns=[]*mysql.Column{
//		{
//			Name:  "gender",
//			Value: 1,
//		},
//		{
//			Name:    "age",
//			Value:   20,
//			ExpType: mysql.Gt,
//		},
//	}
func (d *Dao) GetUsersByColumns(columns []*mysql.Column, page int, pageSize int, sort string) ([]*model.User, int, error) {
	query, args, err := mysql.GetQueryConditions(columns)
	if err != nil {
		return nil, 0, err
	}

	obj := &model.User{}
	total, err := obj.Count(d.db, query, args...)
	if err != nil {
		return nil, total, err
	}

	pageSet := mysql.NewPage(page, pageSize, sort)
	rows, err := obj.Gets(d.db, pageSet, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return rows, total, err
}
