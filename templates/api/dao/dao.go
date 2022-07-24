package dao

import (
	"github.com/zhufuyi/pkg/mysql"
)

// Dao 数据
type Dao struct {
	db *mysql.DB
}

// New Dao示例化
func New(db *mysql.DB) *Dao {
	if db == nil {
		panic("db is nil")
	}
	return &Dao{db: db}
}
