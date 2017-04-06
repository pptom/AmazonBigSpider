package log

import (
	"errors"
	"fmt"
)

//日志类
type Logger struct {
	*LoggerConf
	Callpath int // 不放在LoggerConf中，因为同一个名称的Logger实例，Conf相同，可能被封装的层次不同，用Callpath的不同
}

func (l *Logger) SetCallpath(callpath int) {
	l.Callpath = callpath
}

func (l *Logger) IsAll() bool    { return l.IsLogLevel(ALL) }
func (l *Logger) IsInfo() bool   { return l.IsLogLevel(LOG) }
func (l *Logger) IsDebug() bool  { return l.IsLogLevel(DEBUG) }
func (l *Logger) IsNotice() bool { return l.IsLogLevel(NOTICE) }
func (l *Logger) IsWarn() bool   { return l.IsLogLevel(WARN) }
func (l *Logger) IsError() bool  { return l.IsLogLevel(ERROR) }

func (l *Logger) IsLogLevel(level int) bool {
	_, ok := l.Levels[level]
	return ok
}

func (l *Logger) checkAndLog(level int, args ...interface{}) {
	if shouldPrint, ok := l.Levels[level]; !ok || !shouldPrint {
		return
	}
	if l.Callpath == 0 {
		l.SetCallpath(DefaultLoggerCallpath)
	}
	levelStr := logLevelStringMap[level]
	for _, appender := range l.Appenders {
		appender.Log(l.Callpath, levelStr, args...)
	}
}

func (l *Logger) checkAndLogf(level int, format string, args ...interface{}) {
	if shouldPrint, ok := l.Levels[level]; !ok || !shouldPrint {
		return
	}
	if l.Callpath == 0 {
		l.SetCallpath(DefaultLoggerCallpath)
	}
	levelStr := logLevelStringMap[level]
	for _, appender := range l.Appenders {
		appender.Logf(l.Callpath, levelStr, format, args...)
	}
}

func (l *Logger) Debug(args ...interface{}) {
	l.checkAndLog(DEBUG, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.checkAndLogf(DEBUG, format, args...)
}

func (l *Logger) Log(args ...interface{}) {
	l.checkAndLog(LOG, args...)
}

func (l *Logger) Logf(format string, args ...interface{}) {
	l.checkAndLogf(LOG, format, args...)
}

func (l *Logger) Notice(args ...interface{}) {
	l.checkAndLog(NOTICE, args...)
}

func (l *Logger) Noticef(format string, args ...interface{}) {
	l.checkAndLogf(NOTICE, format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.checkAndLog(WARN, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.checkAndLogf(WARN, format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.checkAndLog(ERROR, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.checkAndLogf(ERROR, format, args...)
}

//日志的管理类，用于加载配置，返回正确的logger
type LoggerManager struct {
	*Tree

	cfg *Config
}

func NewLoggerManager(root *LoggerConf) *LoggerManager {
	manager := &LoggerManager{
		Tree: NewTree(root),
	}
	return manager
}

func NewLoggerManagerWithConf(cfg *Config) (*LoggerManager, error) {
	err := cfg.Verify()
	if err != nil {
		return nil, err
	}
	manager := NewLoggerManager(cfg.RootLogger())
	manager.cfg = cfg
	for _, loggerCfg := range cfg.BuildLoggers() {
		manager.insert(loggerCfg)
	}
	return manager, nil
}

func NewLoggerManagerWithJsconf(jsConf string) (*LoggerManager, error) {
	cfg, err := LoadConf(jsConf)
	if err != nil {
		return nil, err
	}
	return NewLoggerManagerWithConf(cfg)
}

func (self *LoggerManager) UpdateConf(cfg *Config) error {
	err := cfg.Verify()
	if err != nil {
		return err
	}
	newTree := self.Tree.clone()
	newTree.updateConf(cfg.RootLogger())
	for _, loggerCfg := range cfg.BuildLoggers() {
		newTree.updateConf(loggerCfg)
	}
	self.cfg = cfg
	self.Tree = newTree
	return nil
}
func (self *LoggerManager) Logger(name string) (logger *Logger) {
	cfg := self.inheritConf(name)
	cfg.Name = name
	return &Logger{
		LoggerConf: cfg,
	}
}

func (self *LoggerManager) SetLogger(logger *Logger) {
	self.Tree.updateConf(logger.LoggerConf)
}

func (self *LoggerManager) SetRootAppender(appenders ...Appender) {
	self.Root.current.SetAppender(appenders...)
	self.Root.resetFinalConf()
}

func (self *LoggerManager) UseRoot(name string) (err error) {
	if self.cfg == nil {
		return errors.New("LoggerManager 缺少配置")
	}
	if root, ok := self.cfg.RootsLogger()[name]; ok {
		self.Root.current.Appenders = root.Appenders
		self.Root.current.Levels = root.Levels
		self.Root.resetFinalConf()
	} else {
		return fmt.Errorf("LoggerManager 找不到 [%s] Root Logger", name)
	}
	return nil
}

func (self *LoggerManager) SetRootLevel(l int) {
	self.Root.current.SetLevel(l)
	self.Root.resetFinalConf()
}

func (self *LoggerManager) SetRootOnlyLevel(ls ...int) {
	self.Root.current.SetOnlyLevels(ls...)
	self.Root.resetFinalConf()
}
