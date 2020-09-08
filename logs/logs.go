// Package logs
package logs

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

//Logger .
var Logger *PPLog

// PPLog .
type PPLog struct {
	*zap.SugaredLogger
	atom   zap.AtomicLevel
	caller bool
	level  zapcore.Level
}

func init() {
	Logger = NewLog("console", 0, 0, 0, true)
	Logger.SetLevel(zap.DebugLevel)
}

func NewLog(pathFile string, maxSize, maxBackups, maxAge int, caller bool) (l *PPLog) {
	l = new(PPLog)
	l.atom = zap.NewAtomicLevel()
	l.caller = caller
	l.setLogFile(pathFile, maxSize, maxBackups, maxAge)
	l.level = zap.DebugLevel
	l.SetLevel(zap.DebugLevel)
	return
}

func (l *PPLog) SetLevel(lev zapcore.Level) {
	l.atom.SetLevel(lev)
}

func (l *PPLog) SetLogFile(pathFile string, maxSize, maxBackups, maxAge int, caller bool) {
	l.caller = caller
	l.setLogFile(pathFile, maxSize, maxBackups, maxAge)
	l.SetLevel(l.level)
}

func (l *PPLog) Flush() error {
	return l.Sync()
}

func (l *PPLog) setLogFile(pathFile string, maxSize, maxBackups, maxAge int) {
	cfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "lev",
		NameKey:        "log",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, //zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,                   // zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	var core zapcore.Core
	if pathFile == "console" {
		wfd := zapcore.Lock(os.Stdout)
		core = zapcore.NewCore(
			//zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			//zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
			//zapcore.NewJSONEncoder(cfg),
			zapcore.NewConsoleEncoder(cfg),
			wfd,
			l.atom,
		)
	} else {
		wfd := zapcore.AddSync(&lumberjack.Logger{
			Filename:   pathFile,
			MaxSize:    maxSize, // megabytes
			MaxBackups: maxBackups,
			MaxAge:     maxAge, // days
		})
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg),
			wfd,
			l.atom,
		)

	}

	if l.caller {
		l.SugaredLogger = zap.New(core, zap.AddCaller()).Sugar()
	} else {
		l.SugaredLogger = zap.New(core).Sugar()
	}
}

func (l *PPLog) consoleColor() error {
	_t, err := zap.NewDevelopment(zap.AddCaller())
	if err != nil {
		return err
	}
	l.SugaredLogger = _t.Sugar()
	return nil
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//enc.AppendString("[" + t.UTC().Format("2006-01-02 15:04:05") + "]")
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}

func consoleColor() {
	// logger, err := zap.NewDevelopment()
	// if err != nil {
	// 	return
	// }
	// logger.Sugar()
}
