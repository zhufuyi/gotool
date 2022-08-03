package service

import (
	"github.com/zhufuyi/goctl/templates/web/dao"
	"github.com/zhufuyi/goctl/templates/web/model"

	"github.com/zhufuyi/pkg/mysql"
)

// CreateUserExampleRequest 请求参数
type CreateUserExampleRequest struct {
	// binding使用说明 https://github.com/go-playground/validator

	/* todo
	Name   string `form:"name" binding:"min=1"`
	Age    int    `form:"age" binding:"gt=0,lt=120"`
	Gender string `form:"gender" binding:"min=1"`
	*/
}

// CreateUserExample 创建
func (s *Service) CreateUserExample(req *CreateUserExampleRequest) error {
	data := &dao.UserExample{
		/* todo
		Name:   req.Name,
		Age:    req.Age,
		Gender: req.Gender,
		*/
	}
	return s.dao.CreateUserExample(data)
}

// DeleteUserExampleRequest 删除一个id时，从url参数
type DeleteUserExampleRequest struct {
	ID uint64 `form:"id" binding:"gt=0"`
}

// DeleteUserExamplesRequest 删除多个id时，从body获取
type DeleteUserExamplesRequest struct {
	IDs []uint64 `form:"ids" binding:"min=1"`
}

// DeleteUserExample 删除记录
func (s *Service) DeleteUserExample(ids ...uint64) error {
	if len(ids) == 1 {
		return s.dao.DeleteUserExample(ids[0])
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
		err := s.dao.DeleteUserExamples(ids[start:end])
		if err != nil {
			return err
		}
		start = end
	}

	return nil
}

// UpdateUserExampleRequest 请求参数
type UpdateUserExampleRequest struct {
	/* todo
	ID     uint64 `form:"id" binding:"gt=0"`
	Name   string `form:"name" binding:""`
	Age    int    `form:"age" binding:""`
	Gender string `form:"gender" binding:""`
	*/
}

// UpdateUserExample 更新
func (s *Service) UpdateUserExample(req *UpdateUserExampleRequest) error {
	return s.dao.UpdateUserExample(&dao.UserExample{
		/* todo
		ID:     req.ID,
		Name:   req.Name,
		Age:    req.Age,
		Gender: req.Gender,
		*/
	})
}

// GetUserExampleRequest 请求参数
type GetUserExampleRequest struct {
	ID uint64 `form:"id" binding:"gt=0"`
}

// GetUserExample 根据id获取一条记录
func (s *Service) GetUserExample(req *GetUserExampleRequest) (*model.UserExample, error) {
	return s.dao.GetUserExample(req.ID)
}

// GetUserExamplesRequest 请求参数
type GetUserExamplesRequest struct {
	Page int    `form:"page" binding:"gte=0"`
	Size int    `form:"size" binding:"gt=0"`
	Sort string `form:"sort" binding:"-"`

	// 参数填写方式一：从request请求url中获取参数(form.URLParams = c.Request.URL.RawQuery)，
	// 用来自动填充exp、logic的默认值，为了在url参数减少填写exp和logic的默认值，例如url参数?page=0&size=20&exp=gt&k=age&v=22&k=gender&v=1，表示查询年龄大于22岁的男性
	// 参数填写方式二：没有从请求url中获取参数，也就是ParamSrc为空时，请求url参数必须满足len(k)=len(v)=len(exp)=len(logic)，
	// 可以同时存在多个，也可以同时不存在，例如url参数?page=0&size=20&k=age&v=22&exp=gt&logic=and&k=gender&v=1&exp=eq&logic=and，也是表示查询年龄大于22岁的男性
	// 两种url参数都是合法，建议使用第一种
	URLParams string   `form:"-" binding:"-"`
	Keys      []string `form:"k" binding:"-"`
	Values    []string `form:"v" binding:"-"`
	Exps      []string `form:"exp" binding:"-"`
	Logics    []string `form:"logic" binding:"-"`
}

// GetUserExamples 获取多条记录
func (s *Service) GetUserExamples(req *GetUserExamplesRequest) ([]*model.UserExample, int, error) {
	var values []interface{}
	for _, v := range req.Values {
		values = append(values, v)
	}
	columns, err := mysql.GetColumns(req.Keys, values, req.Exps, req.Logics, req.URLParams)
	if err != nil {
		return nil, 0, err
	}

	return s.dao.GetUserExamplesByColumns(columns, req.Page, req.Size, req.Sort)
}

// 通过post方法提交表单进行查询
type column struct {
	Name  string      `json:"name"`  // 列名
	Value interface{} `json:"value"` // 值
	Exp   string      `json:"exp"`   // 表达式，值为空时默认为eq，有eq、neq、gt、gte、lt、lte、like七种类型
	Logic string      `json:"logic"` // 逻辑类型，值为空时默认为and，有and、or两种类型
}

// GetUserExamplesRequest2 请求参数
type GetUserExamplesRequest2 struct {
	Columns []column `json:"columns"`

	Page int    `form:"page" binding:"gte=0" json:"page"`
	Size int    `form:"size" binding:"gt=0" json:"size"`
	Sort string `form:"sort" binding:"" json:"sort"`
}

// GetUserExamples2 获取多条记录
func (s *Service) GetUserExamples2(req *GetUserExamplesRequest2) ([]*model.UserExample, int, error) {
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

	return s.dao.GetUserExamplesByColumns(columns, req.Page, req.Size, req.Sort)
}
