package log

import (
	"encoding/json"
	"errors"
	"strings"
)

type ConfigLogger struct {
	Appenders []string
	Level     interface{}
}

//日志的配置，可以通过配置json来配置日志
type Config struct {
	UseShortFile bool
	Appenders    map[string]struct {
		Type   string
		Target string
	}
	Loggers    map[string]ConfigLogger
	Root       ConfigLogger
	Roots      map[string]ConfigLogger
	_appenders map[string]Appender
}

//加载json配置
func LoadConf(jsConf string) (cfg *Config, err error) {
	cfg = &Config{}
	err = json.Unmarshal([]byte(jsConf), cfg)
	if err == nil {
		err = cfg.Verify()
	}
	UseShortFile = cfg.UseShortFile
	return
}

//实例化配置中的appender
func (self *Config) BuildAppenders() map[string]Appender {
	if self._appenders == nil {
		ap := make(map[string]Appender)
		for name, cfg := range self.Appenders {
			switch strings.ToLower(cfg.Type) {
			case "console":
				ap[name] = NewConsoleAppender(name)
			case "file":
				ap[name] = NewFileAppender(name, cfg.Target)
			case "level":
				ap[name] = NewLevelSeparationDailyAppender(name, cfg.Target)
			case "dailyfile":
				ap[name] = NewDailyAppender(name, cfg.Target)
			default:
				panic("配置中含有未知的AppenderType [" + cfg.Type + "]")
			}
		}
		self._appenders = ap
	}
	return self._appenders
}

func (self *Config) BuildLoggers() []*LoggerConf {
	var loggers []*LoggerConf
	for name, cfg := range self.Loggers {
		apppenders := self.BuildAppenders()
		var ap []Appender
		for _, name := range cfg.Appenders {
			ap = append(ap, apppenders[name])
		}
		loggerConf := &LoggerConf{
			Name:      name,
			Appenders: ap,
		}
		loggerConf.Levels = make(map[int]bool)
		switch levels := cfg.Level.(type) {
		case []interface{}:
			tmpArr := make([]int, 0)
			for _, v := range levels {
				tmpArr = append(tmpArr, LogLevelMap[strings.ToUpper(v.(string))])
			}
			loggerConf.SetOnlyLevels(tmpArr...)
		case string:
			intLevel := LogLevelMap[strings.ToUpper(levels)]
			loggerConf.SetLevel(intLevel)
		}
		loggers = append(loggers, loggerConf)
	}
	return loggers
}
func (self Config) RootLogger() *LoggerConf {
	apppenders := self.BuildAppenders()
	var ap []Appender
	for _, name := range self.Root.Appenders {
		ap = append(ap, apppenders[name])
	}

	root := &LoggerConf{
		Name:   "",
		Levels: parseLevels(self.Root.Level),
		// Level:     LogLevelMap[strings.ToUpper(self.Root.Level)],
		Appenders: ap,
	}
	return root
}

func (self Config) RootsLogger() map[string]*LoggerConf {
	apppenders := self.BuildAppenders()
	ret := map[string]*LoggerConf{}
	for key, clogger := range self.Roots {
		var ap []Appender
		for _, name := range clogger.Appenders {
			ap = append(ap, apppenders[name])
		}

		ret[key] = &LoggerConf{
			Name:   "",
			Levels: parseLevels(clogger.Level),
			// Level:     LogLevelMap[strings.ToUpper(self.Root.Level)],
			Appenders: ap,
		}
	}
	return ret
}

func parseLevels(levels interface{}) (levelMap map[int]bool) {
	levelMap = map[int]bool{}
	switch levels := levels.(type) {
	case []interface{}:
		for _, v := range levels {
			levelMap[LogLevelMap[strings.ToUpper(v.(string))]] = true
		}
	case string:
		for _, l := range LogLevelMap {
			if l <= LogLevelMap[strings.ToUpper(levels)] {
				levelMap[l] = true
			}
		}
	}
	return levelMap
}
func (self Config) Verify() error {
	var e = func(msg string) error {
		return errors.New("logger.conf错误：" + msg)
	}
	if len(self.Appenders) == 0 {
		return e("配置中没有[Appenders]配置")
	}

	for _, ac := range self.Appenders {
		switch strings.ToLower(ac.Type) {
		case "file", "level", "dailyfile":
			if ac.Target == "" {
				return e("fileAppender的 Target[文件名] 不能为空！")
			}
		case "console":
		default:
			return e("不支持的[Appender]类型[" + ac.Type + "]")
		}
	}

	if len(self.Root.Appenders) == 0 {
		return e("Root.Appenders 必须配置")
	}
	for _, name := range self.Root.Appenders {
		if _, ok := self.Appenders[name]; !ok {
			return e("Root.Appender." + name + " 在Appender中找不到")
		}
	}
	if self.Root.Level == "" {
		return e("Root.Level必须配置")
	}
	// if _, ok := LogLevelMap[strings.ToUpper(self.Root.Level)]; !ok {
	//     return e("Root.Level [" + self.Root.Level + "] 不是正确的level等级")
	// }
	return nil
}
