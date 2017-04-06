package main

// insert url into mysql
import (
	"fmt"
	"github.com/hunterhug/AmazonBigSpider"
	"github.com/hunterhug/AmazonBigSpider/public/core"
	"github.com/hunterhug/GoSpider/util"
	"strings"
)

func main() {
	if dudu.Local {
		core.InitConfig(dudu.Dir+"/config/"+"usa_local_config.json", dudu.Dir+"/config/"+"usa_log.json")
	} else {
		core.InitConfig(dudu.Dir+"/config/"+"usa_config.json", dudu.Dir+"/config/"+"usa_log.json")
	}
	dir := core.MyConfig.Datadir + "/url"
	files, e := util.WalkDir(dir, "md")
	filesxx, exx := util.WalkDir(dir, "mdxx")
	if exx != nil {
		fmt.Println(exx.Error())
		panic("dudu")
	}
	if e != nil {
		fmt.Println(e.Error())
		panic("dudu")
	} else {
		ismallbool := map[string]bool{}
		for _, v := range filesxx {
			xxtemp := strings.Split(v, "\\")
			xxlen := len(xxtemp)
			ismallbool[xxtemp[xxlen-1]] = true

		}
		for _, file := range files {
			fmt.Printf("处理%s\n", file)
			dudu, dudue := util.ReadfromFile(file)
			if dudue != nil {
				fmt.Printf("打开%s失败\n")
			} else {
				filecont := string(dudu)
				filelist := strings.Split(filecont, "\n")
				for _, onefile := range filelist {
					mysqllist := strings.Split(onefile, ",")
					if len(mysqllist) != 3 {
						continue
					}
					temp := strings.Split(mysqllist[0], "-")
					pid := "0"
					ismall := 0
					level := len(temp)
					if level > 1 {
						pid = strings.Join(temp[0:len(temp)-1], "-")
					}
					if level == 6 {
						ismall = 1
					} else {
						xx := mysqllist[0] + ".mdxx"
						if _, ok := ismallbool[xx]; ok {
							ismall = 1
						}
					}
					bigpid := temp[0]
					bigname, ok := core.Urlnumdudumap[bigpid]
					if ok == false {
						continue
					}
					url := mysqllist[1]
					name := mysqllist[2]
					//fmt.Printf("%d,%s,%s,%s,%s\n",level,bigpid,bigname,url,name)
					// url must set as unique
					//Todo robot!!!!!and url repeat
					sql := "INSERT IGNORE INTO `smart_category`(`id`,`url`,`name`,`level`,`pid`,`createtime`,`bigpname`,`bigpid`,`ismall`) VALUES(?,?,?,?,?,?,?,?,?);"
					_, mysqle := core.BasicDb.Insert(sql, mysqllist[0], url, name, level, pid, util.TodayString(6), bigname, bigpid, ismall)
					if mysqle != nil {
						fmt.Printf("插入错误:%s\n", mysqle.Error())
					} else {
						fmt.Printf("插入成功:%s\n", onefile)
					}
				}
			}
		}
	}
	//fmt.Printf("%#v",core.Urlnumdudumap)
}
