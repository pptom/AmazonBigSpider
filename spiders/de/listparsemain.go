package main

import (
	"github.com/hunterhug/AmazonBigSpider"
	"github.com/hunterhug/AmazonBigSpider/public/core"
)

func main() {
	if dudu.Local {
		core.InitConfig(dudu.Dir+"/config/de_local_config.json", dudu.Dir+"/config/de_log.json")
	} else {
		core.InitConfig(dudu.Dir+"/config/de_config.json", dudu.Dir+"/config/de_log.json")
	}
	core.LocalListParseTask()

}
