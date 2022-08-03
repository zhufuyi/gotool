package global

import (
	"github.com/zhufuyi/goctl/templates/api/configs"
	"github.com/zhufuyi/pkg/mysql"
)

var (
	Conf    *configs.Conf
	MysqlDB *mysql.DB
	//Logger  *zap.Logger // 使用封装后的全局logger "github.com/zhufuyi/pkg/logger"
)
