package log

import (
	"ginana/library/log/hook"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

var (
	_level        logrus.Level
	_path         string
	_maxAge       uint
	_rotationTime uint
)

type logger struct {
	log     *logrus.Logger
	logType int
}

func (l *logger) GetOutFile() (out io.Writer) {
	return l.log.Out
}

func (l *logger) isStdOut() bool {
	return l.logType == 0
}

func (l *logger) NewStdOut() {
	cli := logrus.New()
	cli.Out = os.Stdout
	cli.Formatter = &GiNanaStdFormatter{}
	cli.AddHook(&hook.DefaultFieldHook{})
	cli.AddHook(&hook.LineHook{})
	l.log = cli
	return
}

func (l *logger) NewFile() (cf func()) {
	cli := logrus.New()
	out, _ := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	cf = func() {
		if out != nil {
			_ = out.Close()
		}
	}
	cli.Out = out
	cli.SetLevel(_level)
	logWriter, err := rotatelogs.New(
		_path+"-%Y-%m-%d-%H-%M.log",
		//rotatelogs.WithLinkName(d.path),
		rotatelogs.WithMaxAge(time.Duration(_maxAge)*time.Hour),             // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Duration(_rotationTime)*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		return
	}
	writeMap := lfshook.WriterMap{
		logrus.DebugLevel: logWriter,
		logrus.InfoLevel:  logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.FatalLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.TextFormatter{DisableColors: true})
	cli.AddHook(&hook.DefaultFieldHook{})
	cli.AddHook(&hook.LineHook{})
	cli.AddHook(lfHook)
	l.log = cli
	return
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *logger) PrintErrf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}
