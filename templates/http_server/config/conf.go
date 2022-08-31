package config

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
	// 服务基本信息
	ServiceName string `json:"serviceName" yaml:"serviceName"`
	ServicePort int    `json:"servicePort" yaml:"servicePort"`
	ServiceEnv  string `json:"serviceEnv" yaml:"serviceEnv"`
	Version     string `json:"version" yaml:"version"`

	// http超时设置，单位(s)
	ReadTimeout  int `json:"readTimeout" yaml:"readTimeout"`
	WriteTimeout int `json:"writeTimeout" yaml:"writeTimeout"`

	// 是否开启go profile
	EnableProfile bool `json:"enableProfile" yaml:"enableProfile"`

	// 是否开启指标采集
	EnableMetrics bool `json:"enableMetrics" yaml:"enableMetrics"`

	// 是否开启限流
	EnableLimit bool        `json:"enableLimit" yaml:"enableLimit"`
	Ratelimiter Ratelimiter `json:"ratelimiter" yaml:"ratelimiter"`

	// 是否开启链路跟踪
	EnableTracing bool `json:"enableTracing" yaml:"enableTracing"`
	// jaeger配置
	Jaeger Jaeger `json:"jaeger" yaml:"jaeger"`

	// 输出日志级别
	Log *Log `json:"log" yaml:"log"`

	// mysql配置
	MysqlURL string `json:"mysqlUrl" yaml:"mysqlUrl"`

	// redis配置
	RedisURL string `json:"redisUrl" yaml:"redisUrl"`
}

// Ratelimiter 配置
type Ratelimiter struct {
	Dimension string `json:"dimension" yaml:"dimension"`
	QPSLimit  int    `json:"qpsLimit" yaml:"qpsLimit"`
	MaxLimit  int    `json:"maxLimit" yaml:"maxLimit"`
}

// IsIP 判断是否使用ip维度限流
func (r *Ratelimiter) IsIP() bool {
	return strings.ToUpper(r.Dimension) == "IP"
}

// Jaeger 配置
type Jaeger struct {
	AgentHost    string  `json:"agentHost" yaml:"agentHost"`
	AgentPort    string  `json:"agentPort" yaml:"agentPort"`
	SamplingRate float64 `json:"samplingRate" yaml:"samplingRate"` // 采样率0~1之间
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
	if conf == nil && (conf.ServicePort == 0 || conf.ServiceName == "") {
		panic(`uninitialised profile, eg:config.Init("conf.yml")`)
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
		err := viper.Unmarshal(conf)
		if err != nil {
			fmt.Println("viper.Unmarshal error: ", err)
		} else {
			// 更新初始化
			for _, f := range fs {
				f()
			}
		}
	})
}

// IsProd 判断是否正式环境
func IsProd() bool {
	return strings.ToLower(conf.ServiceEnv) == "prod"
}

// Show 打印配置信息(去掉敏感信息)
func Show() Conf {
	config := *conf

	// 去掉敏感信息
	config.MysqlURL = ReplaceDbURL(config.MysqlURL)
	config.RedisURL = ReplaceDbURL(config.RedisURL)

	return config
}

// ReplaceDbURL 替换密码
func ReplaceDbURL(str string) string {
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
