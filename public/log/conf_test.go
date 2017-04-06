package log

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestLoadConf(t *testing.T) {
	convey.Convey("测试配置", t, func() {
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
					"Level":"error"
				}
			},
			"Root":{
				"Level":"log",
				"Appenders":["test_appender"]
			}
		}
		`
		_, err := LoadConf(js)
		convey.So(err, convey.ShouldEqual, nil)
	})
}
