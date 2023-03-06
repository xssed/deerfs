package loger

import (
	"fmt"

	"github.com/xssed/deerfs/deerfs_service/core/system/config"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Lg *zap.Logger

// InitLogger 初始化Logger
func Init() (err error) {

	//执行步骤信息
	fmt.Println("deerfs system logger initialization...")

	writeSyncer := getLogWriter(config.LogFilename(), config.LogMaxsize(), config.LogMaxbackups(), config.LogMaxage())
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(config.LogLevel()))
	if err != nil {
		return
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)

	Lg = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(Lg) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  //日志文件的位置
		MaxSize:    maxSize,   //在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: maxBackup, //保留旧文件的最大个数
		MaxAge:     maxAge,    //保留旧文件的最大天数
	}
	return zapcore.AddSync(lumberJackLogger)
}
