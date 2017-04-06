package core

import (
	"fmt"
	"github.com/hunterhug/AmazonBigSpider"
	"testing"
)

func TestGetIP(t *testing.T) {
	if dudu.Local {
		InitConfig(dudu.Dir+"/config/usa_local_config.json", dudu.Dir+"/config/usa_log.json")
	} else {
		InitConfig(dudu.Dir+"/config/usa_config.json", dudu.Dir+"/config/usa_log.json")
	}
	fmt.Printf("%#v", getips())
}
