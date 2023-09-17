package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	serverConfig *ServerConfig
	initialized  bool
)

type ServerConfig struct {
	Port          int           `mapstructure:"port"`
	Twitter       TwitterConfig `mapstructure:"twitter"`
	MySql         MySqlConfig   `mapstructure:"mysql"`
	AKMDTwitterId string        `mapstructure:"akmdTwitterId"`
	SpiderToken   string        `mapstructure:"spiderToken"`
}

type TwitterConfig struct {
	ApiKey       string `mapstructure:"apiKey"`
	ApiSecret    string `mapstructure:"apiSecret"`
	AccessToken  string `mapstructure:"accessToken"`
	AccessSecret string `mapstructure:"accessSecret"`
	ClientID     string `mapstructure:"clientID"`
	ClientSecret string `mapstructure:"clientSecret"`
	CallbackURL  string `mapstructure:"callbackURL"`
}

type MySqlConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Address  string `mapstructure:"address"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
}

func GetServerConfig() *ServerConfig {
	if !initialized {
		panic("should initialize first!")
	}
	return serverConfig
}

// 内网穿透地址：http://keven007.vipgz1.91tunnel.com/
func Initialize(isDebug bool) {
	// 初始化配置文件
	v := viper.New()
	if isDebug {
		v.SetConfigFile("config-dev.yaml")
	} else {
		v.SetConfigFile("config-prod.yaml")
	}
	if err := v.ReadInConfig(); err != nil {
		zap.L().Panic("读取配置失败", zap.Error(err))
	}
	serverConfig = new(ServerConfig)
	if err := v.Unmarshal(serverConfig); err != nil {
		zap.L().Panic("解析配置失败", zap.Error(err))
	}
	zap.L().Info("server config:", zap.Any("config", serverConfig))

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.L().Info("配置文件修改了...")
		if err := v.Unmarshal(serverConfig); err != nil {
			panic(err)
		}
		zap.L().Info("server config:", zap.Any("config", serverConfig))
	})
	initialized = true
}
