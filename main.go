package main

import (
	"Go_hmdp/config"
	"Go_hmdp/pkg/logger"
	"fmt"
	"github.com/spf13/pflag"
	"log/slog"
	"time"
)

func init() {
	// 命令行操作
	// 参数config是对应的长选项，c是短选项。"configs/config.yaml"是默认值。
	configPath := pflag.StringP("config", "c", "./resource/config.yaml", "config file path")
	pflag.Parse()

	config.InitConfig(*configPath)      //初始化配置
	logger.InitLogger(config.LogOption) // 初始化日志配置
}

func main() {
	fmt.Println(config.LogOption)

	go func() {
		for {
			slog.Info("Error")
			slog.Error("Error")
			slog.Debug("debug")
			time.Sleep(10 * time.Second)
		}
	}()

	time.Sleep(100 * time.Second)

}
