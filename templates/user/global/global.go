package global

import (
	"github.com/zhufuyi/goctl/templates/user/configs"
	"github.com/zhufuyi/pkg/email"

	"github.com/zhufuyi/pkg/mysql"
)

var (
	Conf    *configs.Conf
	MysqlDB *mysql.DB
	//Logger  *zap.Logger // 使用封装后的全局logger "github.com/zhufuyi/pkg/logger"
	EmailCli *email.Client
)
