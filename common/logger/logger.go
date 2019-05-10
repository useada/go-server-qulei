package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"a.com/go-server/common/configor"
)

var defaultEncoderConfig = zapcore.EncoderConfig{
	CallerKey:      "caller",
	StacktraceKey:  "stack",
	TimeKey:        "time",
	MessageKey:     "msg",
	LevelKey:       "level",
	NameKey:        "logger",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeCaller:   zapcore.ShortCallerEncoder,
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeName:     zapcore.FullNameEncoder,
	EncodeTime:     MilliSecondTimeEncoder,
}

func InitLogger(conf configor.LoggerConfigor) *zap.SugaredLogger {
	return zap.New(newZapCore(conf), zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}

func MilliSecondTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func newZapCore(conf configor.LoggerConfigor) zapcore.Core {
	//日志文件路径配置
	hook := lumberjack.Logger{
		Filename:   conf.FilePath,   // 日志文件路径
		MaxSize:    conf.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: conf.MaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     conf.MaxAge,     // 文件最多保存多少天
		Compress:   conf.Compress,   // 是否压缩
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapcore.Level(conf.Level))

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(defaultEncoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)
}
