package log

import "github.com/hunterhug/GoSpider/util"

var AmazonListLog, AmazonAsinLog, AmazonIpLog *Logger

func New(filename string) {
	logsconf, _ := util.ReadfromFile(filename)
	err := Init(string(logsconf))
	if err != nil {
		panic(err)
	}
	AmazonListLog = Get("daylist")
	AmazonAsinLog = Get("dayasin")
	AmazonIpLog = Get("dayip")
}
