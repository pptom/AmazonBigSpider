package core

// not use!!!!!down!
import (
	"github.com/hunterhug/GoSpider/util"
	"math/rand"
)

var (
	num  int
	cnum chan string
)

func collectasintask(taskname string, files []string) {
	second := rand.Intn(5)
	AmazonListLog.Debugf("%s:%d second", taskname, second)
	util.Sleep(second)
	CollectAsin(files)
	cnum <- taskname
}

func LocalListParseTask() {
	OpenMysql()
	err := CreateAsinTables()
	if err != nil {
		AmazonListLog.Errorf("createtables:%s,error:%s", Today, err.Error())
	}
	err = CreateAsinRankTables()
	if err != nil {
		AmazonListLog.Errorf("createtables:%s,error:%s", "Asin"+Today, err.Error())
	}
	num = MyConfig.Localtasknum
	cnum = make(chan string, num)
	stoptimes := 10
	// loop loop loop loop
	for {

		if stoptimes == 0 {
			break
		}
		files, err := util.ListDir(DataDir+"/list/"+*day, "html")
		if err != nil {
			AmazonListLog.Errorf("open dir error:%s", err.Error())
			break
		}
		filess, err := util.DevideStringList(files, num)
		if err != nil {
			if len(files) > 0 {
				collectasintask("dudu", files)
				<-cnum
			}
			AmazonListLog.Error(err.Error())
			stoptimes--
			// sleep 10 minute 60/10*24=144
			AmazonListLog.Error("wait for 10 minute。。")
			util.Sleep(600)
			continue
		}
		for i := 0; i < num; i++ {
			go collectasintask("Collect Asin task"+util.IS(i), filess[i])
		}
		for i := 0; i < num; i++ {
			AmazonListLog.Log("Collect Asin %s done!\n", <-cnum)
		}

	}
	AmazonListLog.Log("Collect Asin All done")

}
