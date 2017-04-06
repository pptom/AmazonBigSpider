package core

import (
	"github.com/hunterhug/GoSpider/util"
	"math/rand"
)

var (
	listnum int
	endchan chan string
)

func listtask(taskname string) {
	second := rand.Intn(5)
	AmazonListLog.Debugf("%s:%d second", taskname, second)
	util.Sleep(second)
	if MyConfig.Proxycategory {
		err := GetUrls()
		if err != nil {
			AmazonListLog.Error(taskname + "-error:" + err.Error())
		}
	} else {
		err := GetNoneProxyUrls(taskname)
		if err != nil {
			AmazonListLog.Error(taskname + "-error:" + err.Error())
		}
	}
	endchan <- "done!"
}

func ListTask() {
	OpenMysql()
	err := CreateAsinTables()
	if err != nil {
		AmazonListLog.Errorf("createtables:%s,error:%s", Today, err.Error())
	}
	err = CreateAsinRankTables()
	if err != nil {
		AmazonListLog.Errorf("createtables:%s,error:%s", "Asin"+Today, err.Error())
	}
	listnum = MyConfig.Listtasknum
	endchan = make(chan string, listnum)
	for i := 0; i < listnum; i++ {
		go listtask("ltask" + util.IS(i))
	}
	go Clean()
	for i := 0; i < listnum; i++ {
		<-endchan
	}
	AmazonListLog.Log("List All done")
}
