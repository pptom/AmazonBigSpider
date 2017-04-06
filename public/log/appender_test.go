package log

import "testing"

import (
	. "github.com/smartystreets/goconvey/convey"
)

func TestFileAppender(t *testing.T) {
	Convey("测试FileAppender", t, func() {
		return
		fa := NewConsoleAppender("goodappender")
		fa.Log(1, "ERROR", "中国是个好国家")
		fa.Logf(1, "INFO", "中国是个好国家 %d", 111)
		fa.Logln(1, "ERROR", "中国是个坏国家 ", "真的")
		fa.Logln(1, "ERROR", "中国是个坏国家 ", "真的")
	})
	Convey("测试Separation", t, func() {
		return
		a := NewLevelSeparationDailyAppender("spp", "/tmp/spp.log")
		a.Log(1, "ERROR", "ERRORFILE")
		a.Log(1, "LOG", "INFOFILE")
	})
}
