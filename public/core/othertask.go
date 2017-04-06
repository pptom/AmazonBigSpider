package core

import (
	"fmt"
	"github.com/hunterhug/GoSpider/util"
	"os"
)

func Montior() {
	for {
		urltotal, e1 := RedisClient.Llen(MyConfig.Urlpool)
		urldone, e2 := RedisClient.Hlen(MyConfig.Urlhashpool)
		asintotal, e3 := RedisClient.Llen(MyConfig.Asinpool)
		asindone, e4 := RedisClient.Hlen(MyConfig.Asinhashpool)
		ipremain, e5 := RedisClient.Llen(MyConfig.Proxypool)
		sql := "INSERT INTO `smart_monitor`(id,redistype,doing,done,createtime)VALUES(?,?,?,?,?)on duplicate key update doing=?,done=?,updatetime=?"
		fmt.Printf("e1:%v,e2:%v,e3:%v,e4:%v,e5:%v\n", e1, e2, e3, e4, e5)
		if e1 == nil && e2 == nil {
			_, err := BasicDb.Insert(sql, Today, "urlpool", urltotal, urldone, util.TodayString(6), urltotal, urldone, util.TodayString(6))
			if err == nil {
				fmt.Printf("dbur-%v,%v,%v,%v,%v\n", urltotal, urldone, asintotal, asindone, ipremain)
			} else {
				fmt.Println("urlpoolinsert:" + err.Error())
			}
		}
		if e3 == nil && e4 == nil {
			_, err := BasicDb.Insert(sql, Today, "asinpool", asintotal, asindone, util.TodayString(6), asintotal, asindone, util.TodayString(6))
			if err == nil {
				fmt.Printf("dbas-%v,%v,%v,%v,%v\n", urltotal, urldone, asintotal, asindone, ipremain)
			} else {
				fmt.Println("asinpoolinsert:" + err.Error())
			}
		}
		if e5 == nil {
			//_, err := BasicDb.Insert(sql, Today, "ippool", ipremain, 0, util.TodayString(6), ipremain, 0, util.TodayString(6))
			//if err == nil {
				fmt.Printf("dbip-%v,%v,%v,%v,%v\n", urltotal, urldone, asintotal, asindone, ipremain)
			//} else {
			//	fmt.Println("ippoolinsert:" + err.Error())
			//}
		}
		fmt.Println("----------------")
		util.Sleep(30)

	}

}

func Clean() {
	today, _ := util.SI(Today)
	for {

		Newday := util.TodayString(3)
		newday, _ := util.SI(Newday)
		if newday > today {
			fmt.Println("today out!,now is " + Newday)
			os.Exit(1)
		}
		util.Sleep(1800)
	}
}

func smart2016() string {
	urltotal, _ := RedisClient.Llen(MyConfig.Urlpool)
	urldone, _ := RedisClient.Hlen(MyConfig.Urlhashpool)
	asintotal, _ := RedisClient.Llen(MyConfig.Asinpool)
	asindone, _ := RedisClient.Hlen(MyConfig.Asinhashpool)
	ipremain, _ := RedisClient.Llen(MyConfig.Proxypool)
	return fmt.Sprintf(`
	<table border="1" style="text-align:center;font-size:2.0em">
	<tr>
	<th>URLPOOL</th><th>URLDONE</th><th>ASINPOOL</th><th>ASINDONE</th><th>IPPOOL</th>
	</tr>
	<tr>
	<td>%v</td>
	<td>%v</td>
	<td>%v</td>
	<td>%v</td>
	<td>%v</td>
	</tr>
	</table>
	`, urltotal, urldone, asintotal, asindone, ipremain)

}
