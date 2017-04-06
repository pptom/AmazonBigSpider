package core

import (
	"github.com/hunterhug/GoSpider/util"
	"strings"
)

func CollectAsin(files []string) {
	// you can open in many place
	OpenMysql()
	for _, file := range files {
		AmazonIpLog.Debugf("deal file:%s", file)
		text, err := util.ReadfromFile(file)
		if err != nil {
			AmazonIpLog.Errorf("Read %s-error:%s", file, err.Error())
			continue
		}
		fileinfo, _ := util.GetFilenameInfo(file)
		createtime := util.GetSecord2DateTimes(fileinfo.ModTime().Unix())
		insertlist, err := ParseList(text)
		if err != nil {
			AmazonIpLog.Errorf("Parse %s-error:%s", file, err.Error())
			continue
		}
		category := strings.Split(fileinfo.Name(), ",")
		if len(category) != 2 {
			AmazonIpLog.Errorf("Filename %s:error", file)
			continue
		}
		category1 := strings.Split(category[1], ".")[0]
		err = InsertAsinMysql(insertlist, createtime, category1)
		if err != nil {
			AmazonIpLog.Errorf("InsertMysql %s-error:%s", file, err.Error())
			continue
		}
		err = util.Rename(file, file+"sql")

		//err = os.Remove(file)
		if err != nil {
			AmazonIpLog.Errorf("Rename %s-error:%s", file, err.Error())
		} else {
			AmazonIpLog.Debugf("Rename %s", file+"sql")
		}
	}
}
