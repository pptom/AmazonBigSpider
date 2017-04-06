// 日志记录类
package log

//用法：
//
// import "sunteng/commons/log"
// var jsconf = `
// {
// 	"Appenders":{
// 		"test_appender":{
// 			"Type":"file",
// 			"Target":"/tmp/test.log"
// 		},
// 		"a_appender":{
// 			"Type":"console"
// 		}
// 	},
// 	"Loggers":{
// 		"sunteng/commons/log/a":{
//          "@Appenders":"日志输出到test_appender和a_appender",
// 			"Appenders":["test_appender","a_appender"],
//          "@Level":"记录debug和debug等级以上的数据",
// 			"Level":"debug"
// 		},
// 		"sunteng/commons/log/b":{
//          "@Appenders":"日志输出到最近上级的appender,即Root的Appenders",
//          "@Level":"只记录debug和error等级的数据",
// 			"Level":["debug","error"]
// 		}
// 	},
// 	"Root":{
// 		"Level":"log",
// 		"Appenders":["test_appender"]
// 	}
// }
// `
// log.Init(jsconf)
// logger := log.Get("sunteng/commons/log/a")
// logger.Log("hello logger")
// logger := log.Get("sunteng/commons/log/a/b")
// logger.Log("hello logger")
