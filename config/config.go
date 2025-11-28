package config

import (
	"Go_hmdp/commons/logger"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log/slog"
	"time"
)

var (
	//config.go
	LogOption    *logger.LogSetting
	ServerOption *ServerSetting
	MySQLOption  *MySQLSetting
)

type ServerSetting struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type MySQLSetting struct {
	UserName     string
	Password     string
	Host         string
	DbName       string
	MaxIdleConns int
	MaxOpenConns int
}

// 打开配置文件进行读取
func ReadConfigFile(path string) error {
	//viper是可以开箱即用的，这样写法就类似单例模式
	//也可以创建viper 比如 vp:=viper.New()
	viper.SetConfigFile(path) // 指定配置文件名和位置

	viper.WatchConfig() //该函数内部是开启了一个新协程去监听配置文件是否更新
	//设置回调函数
	viper.OnConfigChange(func(in fsnotify.Event) {
		reloadAllSection()
		//查看是否有更新了日志等级
		if logger.ParseLogLevel(LogOption.Level) != logger.LogLevel.Level() {
			fmt.Println("热更新loggerLevel: " + LogOption.Level)
			logger.SetLevel(LogOption.Level)
		}
	})

	return viper.ReadInConfig()

}

// 分段读取
func ReadSection(key string, v any) error {
	return viper.UnmarshalKey(key, v)
}

// config/config.go
func InitConfig(path string) {
	if err := ReadConfigFile(path); err != nil {
		panic(err)
	}

	err := ReadSection("server", &ServerOption)
	if err != nil {
		panic(err)
	}
	err = ReadSection("mysql", &MySQLOption)
	if err != nil {
		panic(err)
	}
	err = ReadSection("log", &LogOption)
	if err != nil {
		slog.Error("read log section error", "err", err)
	}

}

// 就只是再次读取数据而已
func reloadAllSection() {
	err := ReadSection("server", &ServerOption)
	if err != nil {
		slog.Error("read server section error", "err", err)
	}
	err = ReadSection("mysql", &MySQLOption)
	if err != nil {
		slog.Error("read mysql section error", "err", err)
	}
	err = ReadSection("log", &LogOption)
	if err != nil {
		slog.Error("read log section error", "err", err)
	}
}
