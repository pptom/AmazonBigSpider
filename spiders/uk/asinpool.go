package main

import (
	"github.com/hunterhug/AmazonBigSpider/public/core"
	"github.com/hunterhug/AmazonBigSpider"
)

func main() {
	if dudu.Local {
		core.InitConfig(dudu.Dir+"/config/uk_local_config.json", dudu.Dir+"/config/uk_log.json")
	} else {
		core.InitConfig(dudu.Dir+"/config/uk_config.json", dudu.Dir+"/config/uk_log.json")
	}
	core.AsinPool()
}
