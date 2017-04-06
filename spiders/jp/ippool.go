package main

import (
	"github.com/hunterhug/AmazonBigSpider"
	"github.com/hunterhug/AmazonBigSpider/public/core"
)

func main() {
	if dudu.Local {
		core.InitConfig(dudu.Dir+"/config/jp_local_config.json", dudu.Dir+"/config/jp_log.json")
	} else {
		core.InitConfig(dudu.Dir+"/config/jp_config.json", dudu.Dir+"/config/jp_log.json")
	}

	//Todo
	go func() {
		host := ":12347"
		ac := &core.AmazonController{Message: "jp spider running", SpiderType: "IP process is running"}
		err := core.ServePort(host, ac)
		if err != nil {
			panic(err.Error())
		}
	}()
	core.IPPool()
}
