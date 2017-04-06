package log

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestLogger(t *testing.T) {
	convey.Convey("测试Logger", t, func() {
		return
		logger := &Logger{
			LoggerConf: &LoggerConf{
				Name:   "wahhhh",
				Levels: map[int]bool{LOG: true, ERROR: true},
				Appenders: []Appender{
					NewConsoleAppender("test1"),
					NewFileAppender("ffff", "/tmp/testtt.log"),
				},
			},
		}
		logger.Debug("大创绩！")
		logger.Log("大创绩Log！")
		logger.Errorf("大创绩 [%s]！\n", "ERRRRRR")
	})
}

func TestLoggerMananager(t *testing.T) {
	convey.Convey("测试LoggerMananager", t, func() {
		root := &LoggerConf{
			Name:   "",
			Levels: map[int]bool{DEBUG: true},
			Appenders: []Appender{
				NewConsoleAppender("test1"),
				NewConsoleAppender("test2"),
			},
		}
		manager := NewLoggerManager(root)
		logger := manager.Logger("a/b/c")
		convey.So(logger.Name, convey.ShouldEqual, "a/b/c")
		convey.So(logger.Levels[DEBUG], convey.ShouldEqual, true)
		convey.So(len(logger.Appenders), convey.ShouldEqual, 2)

		manager.SetLogger(&Logger{
			LoggerConf: &LoggerConf{
				Name:   "a/b",
				Levels: map[int]bool{ERROR: true},
				Appenders: []Appender{
					NewConsoleAppender("test1"),
				},
			},
		})

		logger = manager.Logger("a/b/c")
		convey.So(logger.Name, convey.ShouldEqual, "a/b/c")
		convey.So(logger.Levels[ERROR], convey.ShouldEqual, true)
		convey.So(len(logger.Appenders), convey.ShouldEqual, 1)
	})

	convey.Convey("测试加载配置的manager", t, func() {
		var js = `
		{
			"Appenders":{
				"test_appender":{
					"Type":"file",
					"Target":"/tmp/test.log"
				},
				"a_appender":{
					"Type":"console"
				},
				"b_appender":{
					"Type":"file",
					"Target":"/tmp/test.b.log"
				}
			},
			"Loggers":{
				"sunteng/commons/log/a":{
					"Appenders":["test_appender","a_appender"],
					"Level":"debug"
				},
				"sunteng/commons/log/b":{
					"Level":["debug","error"]
				}
			},
			"Root":{
				"Level":"log",
				"Appenders":["test_appender"]
			}
		}
		`
		manager, err := NewLoggerManagerWithJsconf(js)
		convey.So(err, convey.ShouldEqual, nil)
		aLogger := manager.Logger("sunteng/commons/log/a")
		convey.So(aLogger.Levels[DEBUG], convey.ShouldEqual, true)
		convey.So(len(aLogger.Appenders), convey.ShouldEqual, 2)

		son := manager.Logger("sunteng/commons/log/b")
		convey.So(son.Levels[ERROR], convey.ShouldEqual, true)
		convey.So(len(son.Appenders), convey.ShouldEqual, 1)
		return
		aLogger.Debug("alogger debug")
		// fmt.Println(aLogger.Levels)
		// fmt.Println(son.Levels)
		son.Error("BBBBB logger debug")

	})
}
