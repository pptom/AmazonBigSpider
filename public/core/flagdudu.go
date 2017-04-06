package core

import (
	"flag"
)

var (
	day      = flag.String("d", Today, "Date(such as 20161110)")
	//sitetype = flag.String("t", "usa", "Site Type(USA/JP/UK/DE)")
)

func init() {
	flag.Parse()
	//prefix := "usa"
	//switch strings.ToLower(*sitetype) {
	//case "uk":
	//	prefix = "uk"
	//case "jp":
	//	prefix = "jp"
	//case "de":
	//	prefix = "de"
	//case "usa":
	//	prefix = "usa"
	//default:
	//
	//}
	//if dudu.Local {
	//	InitConfig(dudu.Dir+"/config/"+prefix+"_local_config.json", dudu.Dir+"/config/"+prefix+"_log.json")
	//} else {
	//	InitConfig(dudu.Dir+"/config/"+prefix+"_config.json", dudu.Dir+"/config/"+prefix+"_log.json")
	//}
}
