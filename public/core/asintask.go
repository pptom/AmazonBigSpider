package core

import (
	"github.com/hunterhug/GoSpider/util"
	"math/rand"
)

var (
	asintasknum int
	asinendchan chan string
)

func asintask(taskname string) {
	second := rand.Intn(5)
	AmazonAsinLog.Debugf("%s:%d second", taskname, second)
	util.Sleep(second)
	if MyConfig.Proxyasin {
		err := GetAsinUrls()
		if err != nil {
			AmazonAsinLog.Errorf(taskname + "-error:" + err.Error())
		}
	} else {
		err := GetNoneProxyAsinUrls(taskname)
		if err != nil {
			AmazonAsinLog.Errorf(taskname + "-error:" + err.Error())
		}
	}
	asinendchan <- "done!"
}

func AsinGo() {
	OpenMysql()
	err := CreateAsinTables()
	if err != nil {
		AmazonAsinLog.Errorf("createtables:%s,error:%s", Today, err.Error())
	}
	err = CreateAsinRankTables()
	if err != nil {
		AmazonAsinLog.Errorf("createtables:%s,error:%s", "Asin"+Today, err.Error())
	}
	asintasknum = MyConfig.Asintasknum
	asinendchan = make(chan string, asintasknum)
	for i := 0; i < asintasknum; i++ {
		go asintask("atask" + util.IS(i))
	}
	go Clean()
	for i := 0; i < asintasknum; i++ {
		<-asinendchan
	}
	AmazonAsinLog.Log("List All done")
}
