package logger

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"os"
	"strings"
)

// logger/log.go
// 日志选项结构体
type LogSetting struct {
	Filename   string
	Level      string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

var LogLevel = new(slog.LevelVar)

func InitLogger(logConfig *LogSetting) {
	log := lumberjack.Logger{
		Filename:   logConfig.Filename,   //日志文件的位置
		MaxSize:    logConfig.MaxSize,    //文件最大尺寸(以mb为单位)
		MaxBackups: logConfig.MaxBackups, //保留的最大文件个数
		MaxAge:     logConfig.MaxAge,     //保留旧文件的最大天数
		LocalTime:  true,                 //使用本地时间创建时间戳
	}
	log = log

	LogLevel.Set(ParseLogLevel(logConfig.Level)) //这样就可以在运行时更新日志等级

	//使用json格式, slog.NewJSONHandler(输出位置 io.Writer, 日志输出格式选项) 这里输出位置本应为&log，这里为了方便改为输出终端
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     LogLevel,
	}))

	slog.SetDefault(logger)
}
func SetLevel(level string) {
	LogLevel.Set(ParseLogLevel(level))
}

// Go 1.23版本后才允许NewMultiHandler，本次项目使用Go 1.21，为了简便，Log直接输出终端
// InitLogger 初始化日志：同时输出到文件（JSON）和终端（文本）
//func InitLogger(level string) error {
//	// 1. 打开日志文件（添加写权限）
//	//file, err := os.OpenFile("dianping.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
//	//if err != nil {
//	//	return err
//	//}
//
//	// 2. 定义日志级别
//	logLevel := parseLogLevel(level)
//
//	// 3. 创建文件处理器（JSON格式）
//	//fileHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{
//	//	AddSource: true, // 显示代码文件和行号
//	//	Level:     logLevel,
//	//})
//
//	// 4. 创建终端处理器（文本格式，更易读）
//	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
//		AddSource: true,
//		Level:     logLevel,
//	})
//
//	// 5. 组合处理器：日志同时输出到文件和终端
//	// multiHandler := slog.NewMultiHandler(fileHandler, consoleHandler)
//	logger := slog.New(consoleHandler)
//
//	// 6. 设置为默认日志器
//	slog.SetDefault(logger)
//	return nil
//}

// 获得日志等级
func ParseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}
	return slog.LevelInfo // 默认Info
}
