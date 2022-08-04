package service

import (
	"context"

	"github.com/zhufuyi/goctl/templates/user/dao"
	"github.com/zhufuyi/goctl/templates/user/global"
)

// Service 方法
type Service struct {
	ctx context.Context
	dao *dao.Dao
}

// New 实例化
func New(ctx context.Context) *Service {
	svc := &Service{ctx: ctx}
	svc.dao = dao.New(global.MysqlDB)
	return svc
}
