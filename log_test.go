package glog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestDebugToFile(t *testing.T) {
	config := &Config{
		Level:               "debug",
		Format:              "json",
		DisableTimestamp:    false,
		File:                FileLogConfig{
			LogDir:     "./logs",
			Filename:   "test",
			MaxSize:    1,
			MaxDays:    5,
			MaxBackups: 10,
		},
		Development:         false,
		DisableCaller:       false,
		DisableStacktrace:   false,
		DisableErrorVerbose: false,
		Sampling:            nil,
	}
	err := InitLogger(config)
	if err != nil {
		t.Fatal(err)
	}
	defer Sync()
	type args struct {
		message string
		field   []zap.Field
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name:"debug1",
			args:args{
				message: "test1",
				field: []zapcore.Field{zap.String("name", "dream1"), zap.Int("age", 20)},
			},
		},
		{
			name:"debug2",
			args:args{
				message: "test2",
				field: []zapcore.Field{zap.String("name", "dream2"), zap.Int("age", 21)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Debug(tt.args.message, tt.args.field...)
			Error(tt.args.message, tt.args.field...)
		})
	}
}

func TestDebugToConsole(t *testing.T) {
	config := &Config{
		Level:               "debug",
		Format:              "json",
		DisableTimestamp:    false,
		File:                FileLogConfig{
			LogDir:     "./logs",
			Filename:   "",
			MaxSize:    1,
			MaxDays:    5,
			MaxBackups: 10,
		},
		Development:         false,
		DisableCaller:       false,
		DisableStacktrace:   false,
		DisableErrorVerbose: false,
		Sampling:            nil,
	}
	err := InitLogger(config)
	if err != nil {
		t.Fatal(err)
	}
	defer Sync()
	type args struct {
		message string
		field   []zap.Field
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name:"debug1",
			args:args{
				message: "test1",
				field: []zapcore.Field{zap.String("name", "dream1"), zap.Int("age", 20)},
			},
		},
		{
			name:"debug2",
			args:args{
				message: "test2",
				field: []zapcore.Field{zap.String("name", "dream2"), zap.Int("age", 21)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Debug(tt.args.message, tt.args.field...)
			Error(tt.args.message, tt.args.field...)
		})
	}
}

func BenchmarkDebug(b *testing.B) {
	config := &Config{
		Level:               "debug",
		Format:              "json",
		DisableTimestamp:    false,
		File:                FileLogConfig{
			Filename:   "./logs/test.log",
			MaxSize:    10,
			MaxDays:    2,
			MaxBackups: 5,
		},
		Development:         false,
		DisableCaller:       false,
		DisableStacktrace:   false,
		DisableErrorVerbose: false,
		Sampling:            nil,
	}
	err := InitLogger(config)
	if err != nil {
		b.Fatal(err)
	}
	defer Sync()
	for i:=0; i<b.N; i++ {
		Debug("test bench mark", zap.String("name", "dream"))
		//Error("test bench mark", zap.String("name", "dream"))
	}
	//b.RunParallel(func(pb *testing.PB) {
	//	for pb.Next(){
	//		Debug("test bench mark", zap.String("name", "dream"))
	//	}
	//})
}

func TestDefaultLogger(t *testing.T){
	L().Info("default logger")
}