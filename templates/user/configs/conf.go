package configs

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var conf *Conf

// Conf 服务配置信息
type Conf struct {
	// 服务名称
	ServerName string `json:"serverName" yaml:"serverName"`
	// 服务端口
	ServerPort int `json:"serverPort" yaml:"serverPort"`
	// 运行模式，dev,prod
	RunMode string `json:"runMode" yaml:"runMode"`

	Address string `json:"address" yaml:"address"`

	ReadTimeout  int `json:"readTimeout" yaml:"readTimeout"`
	WriteTimeout int `json:"writeTimeout" yaml:"writeTimeout"`

	// 是否开启go profile
	IsEnableProfile bool `json:"isEnableProfile" yaml:"isEnableProfile"`

	// 是否开启鉴权
	IsUseAuth bool `json:"isUseAuth" yaml:"isUseAuth"`

	// mysql配置
	MysqlURL string `json:"mysqlURL" yaml:"mysqlURL"`

	// jwt配置
	Jwt *Jwt `json:"jwt" yaml:"jwt"`

	// Email配置
	Email *Email `json:"email" yaml:"email"`

	// 输出日志级别
	Log *Log `json:"log" yaml:"log"`
}

// Jwt token配置
type Jwt struct {
	SigningKey string `json:"signingKey" yaml:"signingKey"`
	Expire     int    `json:"expire" yaml:"expire"` // 单位秒
}

// Email 配置
type Email struct {
	Sender   string `json:"sender" yaml:"sender"`
	Password string `json:"password" yaml:"password"`
}

// Log 日志配置
type Log struct {
	Level  string `json:"level" yaml:"level"`
	Format string `json:"format" yaml:"format"`
	IsSave bool   `json:"isSave" yaml:"isSave"`

	// 保存日志文件相关设置
	LogFileConfig *LogFileConfig `json:"logFileConfig" yaml:"logFileConfig"`
}

// LogFileConfig 日志文件配置
type LogFileConfig struct {
	Filename      string `json:"filename" yaml:"filename"`
	MaxSize       int    `json:"maxSize" yaml:"maxSize"`
	MaxBackups    int    `json:"maxBackups" yaml:"maxBackups"`
	MaxAge        int    `json:"maxAge" yaml:"maxAge"`
	IsCompression bool   `json:"isCompression" yaml:"isCompression"`
}

// Get 获取配置对象
func Get() *Conf {
	if conf == nil && (conf.ServerPort == 0 || conf.ServerName == "") {
		panic(`uninitialised profile, eg:configs.Init("conf.yml")`)
	}

	return conf
}

// Init 解析配置文件到struct，包括yaml、toml、json等文件
func Init(configFile string) error {
	confFileAbs, err := filepath.Abs(configFile)
	if err != nil {
		return err
	}

	filePathStr, filename := filepath.Split(confFileAbs)
	if filePathStr == "" {
		filePathStr = "."
	}
	ext := strings.TrimLeft(path.Ext(filename), ".")
	filename = strings.ReplaceAll(filename, "."+ext, "") // 不包括后缀名

	viper.AddConfigPath(filePathStr) // 路径
	viper.SetConfigName(filename)    // 名称
	viper.SetConfigType(ext)         // 从文件名中获取配置类型
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	conf = new(Conf)
	err = viper.Unmarshal(conf)
	if err != nil {
		return err
	}

	return nil
}

// WatchConfig 监听配置文件更新
func WatchConfig(fs ...func()) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		viper.Unmarshal(conf)
		// 更新初始化
	})
}

// IsProd 判断是否正式环境
func IsProd() bool {
	return strings.ToLower(conf.RunMode) == "prod"
}

// ShowConfig 打印配置信息(去掉敏感信息)
func ShowConfig() Conf {
	config := *conf

	// 去掉敏感信息
	config.MysqlURL = hideMysqlPassword(config.MysqlURL)
	config.Email.Password = "******"

	return config
}

// 隐藏mysql配置密码
func hideMysqlPassword(str string) string {
	mysqlPWD := []byte(str)
	start, end := 0, 0
	for k, v := range mysqlPWD {
		if v == ':' {
			start = k
		}
		if v == '@' {
			end = k
			break
		}
	}

	if start >= end {
		return str
	}

	return fmt.Sprintf("%s******%s", mysqlPWD[:start+1], mysqlPWD[end:])
}
