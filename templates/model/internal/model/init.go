package model

import (
	"sync"

	"github.com/zhufuyi/goctl/templates/model/config"

	"github.com/go-redis/redis/v8"
	"github.com/zhufuyi/pkg/goredis"
	"github.com/zhufuyi/pkg/mysql"
	"gorm.io/gorm"
)

var (
	// ErrRecordNotFound 没有找到记录
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

var (
	db    *gorm.DB
	once1 sync.Once

	redisCli *redis.Client
	once2    sync.Once
)

// InitMysql 连接mysql
func InitMysql() {
	opts := []mysql.Option{mysql.WithLog()}
	if config.Get().EnableTracing {
		opts = append(opts, mysql.WithEnableTrace())
	}

	var err error
	db, err = mysql.Init(config.Get().MysqlURL, opts...)
	if err != nil {
		panic("mysql.Init error: " + err.Error())
	}
}

// GetDB 返回db对象
func GetDB() *gorm.DB {
	if db == nil {
		once1.Do(func() {
			InitMysql()
		})
	}

	return db
}

// CloseDB 关闭连接
func CloseDB() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		if sqlDB != nil {
			return sqlDB.Close()
		}
	}

	return nil
}

// InitRedis 连接redis
func InitRedis() {
	opts := []goredis.Option{}
	if config.Get().EnableTracing {
		opts = append(opts, goredis.WithEnableTrace())
	}

	var err error
	redisCli, err = goredis.Init(config.Get().RedisURL, opts...)
	if err != nil {
		panic("goredis.Init error: " + err.Error())
	}
}

// GetRedis 返回redis client
func GetRedis() *redis.Client {
	if redisCli == nil {
		once2.Do(func() {
			InitRedis()
		})
	}

	return redisCli
}
