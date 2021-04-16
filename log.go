package glog

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"sync/atomic"
)

var _globalL, _globalS atomic.Value
var LevelList = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
var (
	debugLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel
	})

	infoLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel
	})

	warnLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel
	})
	errorLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel
	})

	fatalLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.FatalLevel
	})

	levelList = []zap.LevelEnablerFunc{debugLevel, infoLevel, warnLevel, errorLevel, fatalLevel}
)

func L() *zap.Logger{
	return _globalL.Load().(*zap.Logger)
}

func S() *zap.SugaredLogger{
	return _globalS.Load().(*zap.SugaredLogger)
}

// InitLogger initializes a zap logger.
func InitLogger(cfg *Config, opts ...zap.Option) error {
	var output = make([]zapcore.WriteSyncer, 0)
	var err error
	if len(cfg.File.Filename) > 0 {
		output, err = GetFileLogs(cfg)
		if err != nil {
			return errors.Wrap(err, "init log file failed")
		}
	} else {
		stdOut, _, err := zap.Open([]string{"stdout"}...)
		if err != nil {
			return errors.Wrap(err, "init stdout log failed")
		}
		output = []zapcore.WriteSyncer{stdOut}
	}
	logger, err := InitLoggerWithWriteSyncer(cfg, output, opts...)
	if err != nil {
		return errors.Wrap(err, "InitLoggerWithWriteSyncer failed")
	}
	_globalL.Store(logger)
	sugar := logger.Sugar()
	_globalS.Store(sugar)
	return nil
}

func GetFileLogs(cfg *Config) ([]zapcore.WriteSyncer, error) {
	fileLogs := make([]zapcore.WriteSyncer, 0)
	fileName := cfg.File.Filename
	for _, level := range LevelList{
		cfg.File.Filename = fmt.Sprintf("%s.%s.log", fileName, level)
		lg, err := initFileLog(&cfg.File)
		if err != nil {
			return nil, errors.Wrap(err, "init file log failed")
		}
		output := zapcore.AddSync(lg)
		fileLogs = append(fileLogs, output)
	}
	return fileLogs, nil
}

// initFileLog initializes file based logging options.
func initFileLog(cfg *FileLogConfig) (*lumberjack.Logger, error) {
	filename := filepath.Join(cfg.LogDir, cfg.Filename)
	if st, err := os.Stat(filename); err == nil {
		if st.IsDir() {
			return nil, errors.New("can't use directory as log file name")
		}
	}
	if cfg.MaxSize == 0 {
		cfg.MaxSize = defaultLogMaxSize
	}

	// use lumberjack to logrotate
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxDays,
		LocalTime:  true,
	}, nil
}

// InitLoggerWithWriteSyncer initializes a zap logger with specified write syncer.
func InitLoggerWithWriteSyncer(cfg *Config, output []zapcore.WriteSyncer, opts ...zap.Option) (*zap.Logger, error) {
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	cores := make([]zapcore.Core, 0)
	if len(output) == 1 {
		level := zap.NewAtomicLevel()
		err := level.UnmarshalText([]byte(cfg.Level))
		if err != nil {
			return nil, err
		}
		core := zapcore.NewCore(zapcore.NewJSONEncoder(ec), output[0], level)
		cores = append(cores, core)
	} else {
		for i := range LevelList{
			core := zapcore.NewCore(zapcore.NewJSONEncoder(ec), output[i], levelList[i])
			cores = append(cores, core)
		}
	}
	core := zapcore.NewTee(cores...)
	//core := zapcore.NewCore(zapcore.NewJSONEncoder(ec), output, level)
	//core := zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()), output, level)
	opts = append(cfg.buildOptions(output[len(output)-1]), opts...)
	lg := zap.New(core, opts...)
	return lg, nil
}

func IsDebugEnabled() bool {
	return L().Core().Enabled(zapcore.DebugLevel)
}

func IsInfoEnabled() bool {
	return L().Core().Enabled(zapcore.InfoLevel)
}

func IsWarnEnabled() bool {
	return L().Core().Enabled(zapcore.WarnLevel)
}

func Debug(message string, field ...zap.Field) {
	L().Debug(message, field...)
}

func Info(message string, field ...zap.Field) {
	L().Info(message, field...)
}

func Warn(message string, field ...zap.Field) {
	L().Warn(message, field...)
}

func Error(message string, field ...zap.Field) {
	L().Error(message, field...)
}

func Fatal(message string, field ...zap.Field) {
	L().Fatal(message, field...)
}

func Debugf(format string, v ...interface{}){
	L().Debug(fmt.Sprintf(format, v...))
}

func Infof(format string, v ...interface{}){
	L().Info(fmt.Sprintf(format, v...))
}

func Warnf(format string, v ...interface{}){
	L().Warn(fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...interface{}){
	L().Error(fmt.Sprintf(format, v...))
}

func Fatalf(format string, v ...interface{}){
	L().Fatal(fmt.Sprintf(format, v...))
}

func Sync(){
	S().Sync()
}

func Painc(format string, v ...interface{}){
	fmt.Println(fmt.Sprintf(format, v...))
	os.Exit(1)
}


