package main

import (
	"github.com/hunterhug/AmazonBigSpider"
	"github.com/hunterhug/AmazonBigSpider/public/core"
)

func main() {
	if dudu.Local {
		core.InitConfig(dudu.Dir+"/config/uk_local_config.json", dudu.Dir+"/config/uk_log.json")
	} else {
		core.InitConfig(dudu.Dir+"/config/uk_config.json", dudu.Dir+"/config/uk_log.json")
	}

	//Todo
	go func() {
		host := ":12346"
		ac := &core.AmazonController{Message: "uk spider running", SpiderType: "IP process is running"}
		err := core.ServePort(host, ac)
		if err != nil {
			panic(err.Error())
		}
	}()
	core.IPPool()
}
