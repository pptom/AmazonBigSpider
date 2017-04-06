package core

import (
	"fmt"
	"github.com/hunterhug/AmazonBigSpider"
	"github.com/hunterhug/GoSpider/spider"
	"github.com/hunterhug/GoSpider/util"
	"testing"
)

// https://www.amazon.com/Best-Sellers-Automotive-Performance-ABS-Brake-Parts/zgbs/automotive/15710931/ref=zg_bs_pg_1?_encoding=UTF8&pg=1&ajax=1
// https://www.amazon.com/Best-Sellers-Automotive-Performance-ABS-Brake-Parts/zgbs/automotive/15710931/ref=zg_bs_pg_2?_encoding=UTF8&pg=2&ajax=1
// https://www.amazon.com/Best-Sellers-Automotive-Performance-ABS-Brake-Parts/zgbs/automotive/15710931/ref=zg_bs_pg_3?_encoding=UTF8&pg=3&ajax=1
// https://www.amazon.com/Best-Sellers-Automotive-Performance-ABS-Brake-Parts/zgbs/automotive/15710931/ref=zg_bs_pg_4?_encoding=UTF8&pg=4&ajax=1
// https://www.amazon.com/Best-Sellers-Automotive-Performance-ABS-Brake-Parts/zgbs/automotive/15710931/ref=zg_bs_pg_5?_encoding=UTF8&pg=5&ajax=1
// https://www.amazon.com/dp/B001IHBLPC
func TestAsinDownload(t *testing.T) {
	util.MakeDir(dudu.Dir + "/test/asin/")
	ip := "104.128.124.122:808"
	// debug info will no appear |nothing
	spider.SetLogLevel("info")
	url := "https://www.amazon.com/dp/B016L36UZI"
	prefix := "asin"
	testtimes := 1000
	for {
		testtimes--
		if testtimes == 0 {
			break
		}
		robotime := 0
		maxtime := 1000
		times := 0
		for {
			if times > maxtime {
				break
			}
			temp := url
			content, err := Download(ip, temp)
			if err != nil {
				fmt.Printf("%#v", err.Error())
			} else {
				err = spider.TooSortSizes(content, 10)
				// robot continue
				if err != nil {
					robotime++
					times++
					break
				} else {
					// and then out
					fmt.Printf("The %d try Asin page :%d times | robbot max times:%d\n", testtimes, times, robotime)
				}
				util.SaveToFile(dudu.Dir+"/test/asin/"+prefix+util.IS(testtimes)+".html", content)
				break
			}
		}
	}
}

func TestListDownload(t *testing.T) {
	util.MakeDir(dudu.Dir + "/test/list/")
	ip := "104.128.124.122:808"
	// debug info will no appear |nothing
	spider.SetLogLevel("info")
	url := "https://www.amazon.co.jp/gp/bestsellers/dvd/ref=zg_bs_nav_0"
	content, err := Download(ip, url)
	if err != nil {
		fmt.Printf("%#v", err.Error())
	}
	util.SaveToFile(dudu.Dir+"/test/list/xxx2.html", content)
}
