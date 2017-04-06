package core

import (
	"fmt"
	"github.com/hunterhug/GoSpider/spider"
	"github.com/hunterhug/GoSpider/util"
	"strings"
)

// ip and url download,if error from zero again!
// ip := "104.128.124.122:808"
// url := "https://www.amazon.com/Best-Sellers-Appliances-Home-Appliance-Warranties/zgbs/appliances/2242350011"
// filename:=1-1-1-1 ---->1,1-1-1-1 2,1-1-1-1
func GetUrl(ip string, url string, filename string, name string, bigname string, page int) (string, error) {
	firstproxy := false
	proxy := true
	if strings.Contains(ip, "*") {
		proxy = false
	}
	if strings.Contains(ip, "127.0.0.1") {
		firstproxy = true
	}
	ipbegintimes := ""
	if firstproxy {
		ip = GetIP()
		// before use, send to hash pool
		ipbegintimes = util.GetSecord2DateTimes(util.GetSecordTimes())
		RedisClient.Hset(MyConfig.Proxyhashpool, ip, ipbegintimes)
	}
	for i := 1; i <= page; i++ {
		keepdirtemp := MyConfig.Datadir + "/list/" + Today + "/" + util.IS(i) + "," + filename + ".html"
		if MyConfig.Categorylocalkeep {
			if util.FileExist(keepdirtemp) || util.FileExist(keepdirtemp+"x") || util.FileExist(keepdirtemp+"sql") || util.FileExist(keepdirtemp+"404") {
				AmazonListLog.Debugf("FileExist:" + keepdirtemp)
				continue
			}
		}
		temp := url + fmt.Sprintf("?_encoding=UTF8&pg=%d&ajax=1", i)
		for {
			content, err := Download(ip, temp)
			dspider, ok := Spiders.Get(ip)
			if err != nil {
				AmazonListLog.Errorf("%s:%s", temp, err.Error())
				// if proxy ip err more than config, change ip
				if proxy && ok && dspider.Errortimes > MyConfig.Proxymaxtrytimes {
					// die sent
					ipendtimes := util.GetSecord2DateTimes(util.GetSecordTimes())
					insidetemp := ipbegintimes + "|" + ipendtimes + "|" + util.IS(dspider.Fetchtimes-dspider.Errortimes) + "|" + util.IS(dspider.Errortimes)
					RedisClient.Hset(MyConfig.Proxyhashpool, ip, insidetemp)
					// you know it
					Spiders.Delete(ip)
					// get new proxy again
					ip = GetIP()
					AmazonListLog.Errorf("更换IP：%s", ip)
					ipbegintimes = util.GetSecord2DateTimes(util.GetSecordTimes())
					RedisClient.Hset(MyConfig.Proxyhashpool, ip, ipbegintimes)
				}
				continue
			} else {

				// bigger than 3KB, may be 404
				if Is404(content) {
					if MyConfig.Categorylocalkeep {
						util.SaveToFile(keepdirtemp+"404", []byte(""))
					}
					return ip, nil
				}
				if IsRobot(content) {
					dspider.Errortimes = dspider.Errortimes + 1
					AmazonListLog.Errorf("get %s fail(%d),total(%d),ip:%s", temp, dspider.Errortimes, dspider.Fetchtimes, ip)
					if proxy && ok && dspider.Errortimes > MyConfig.Proxymaxtrytimes {
						// die sent
						ipendtimes := util.GetSecord2DateTimes(util.GetSecordTimes())
						insidetemp := ipbegintimes + "|" + ipendtimes + "|" + util.IS(dspider.Fetchtimes-dspider.Errortimes) + "|" + util.IS(dspider.Errortimes)
						RedisClient.Hset(MyConfig.Proxyhashpool, ip, insidetemp)
						// you know it
						Spiders.Delete(ip)
						// get new proxy again
						ip = GetIP()
						AmazonListLog.Errorf("更换IP：%s", ip)
						ipbegintimes = util.GetSecord2DateTimes(util.GetSecordTimes())
						RedisClient.Hset(MyConfig.Proxyhashpool, ip, ipbegintimes)
					}
					continue
				}
				// 3KB smaller mean page not enough
				err = spider.TooSortSizes(content, 3)
				if err != nil {
					AmazonListLog.Errorf("%s,%s", temp, err.Error())
					if MyConfig.Categorylocalkeep {
						util.SaveToFile(keepdirtemp+"x", []byte(""))
					}
					return ip, nil
				}

				if MyConfig.Categorylocalkeep {
					util.SaveToFile(keepdirtemp, []byte(string(content)+"<div id='hunterhug'>"+name+"|"+bigname+"|"+url+"|"+filename+"</div>"))
				}
				insertlist, err := ParseList([]byte(string(content) + "<div id='hunterhug'>" + name + "|" + bigname + "|" + url + "|" + filename + "</div>"))
				if err != nil {
					AmazonListLog.Errorf("Parse %s-error:%s", temp, err.Error())
					break
				}
				createtime := util.TodayString(6)
				err = InsertAsinMysql(insertlist, createtime, filename)
				if err != nil {
					AmazonListLog.Errorf("InsertMysql %s-error:%s", temp, err.Error())

				}
				break
			}
		}
	}
	return ip, nil
}

func GetJapanUrl(ip string, url string, filename string, name string, bigname string, page int) (string, error) {
	firstproxy := false
	proxy := true
	if strings.Contains(ip, "*") {
		proxy = false
	}
	if strings.Contains(ip, "127.0.0.1") {
		firstproxy = true
	}
	ipbegintimes := ""
	if firstproxy {
		ip = GetIP()
		// before use, send to hash pool
		ipbegintimes = util.GetSecord2DateTimes(util.GetSecordTimes())
		RedisClient.Hset(MyConfig.Proxyhashpool, ip, ipbegintimes)
	}
	for i := 1; i <= page; i++ {
		keepdirtemp := MyConfig.Datadir + "/list/" + Today + "/" + util.IS(i) + "," + filename + ".html"
		if MyConfig.Categorylocalkeep {
			if util.FileExist(keepdirtemp) || util.FileExist(keepdirtemp+"x") || util.FileExist(keepdirtemp+"sql") || util.FileExist(keepdirtemp+"404") {
				AmazonListLog.Debugf("FileExist:" + keepdirtemp)
				continue
			}
		}
		temp := url + fmt.Sprintf("?_encoding=UTF8&pg=%d&ajax=1", i)
		for {
			content, err := Download(ip, temp)
			dspider, ok := Spiders.Get(ip)
			if err != nil {
				AmazonListLog.Errorf("%s:%s", temp, err.Error())
				// if proxy ip err more than config, change ip
				if proxy && ok && dspider.Errortimes > MyConfig.Proxymaxtrytimes {
					// die sent
					ipendtimes := util.GetSecord2DateTimes(util.GetSecordTimes())
					insidetemp := ipbegintimes + "|" + ipendtimes + "|" + util.IS(dspider.Fetchtimes-dspider.Errortimes) + "|" + util.IS(dspider.Errortimes)
					RedisClient.Hset(MyConfig.Proxyhashpool, ip, insidetemp)
					// you know it
					Spiders.Delete(ip)
					// get new proxy again
					ip = GetIP()
					AmazonListLog.Errorf("更换IP：%s", ip)
					ipbegintimes = util.GetSecord2DateTimes(util.GetSecordTimes())
					RedisClient.Hset(MyConfig.Proxyhashpool, ip, ipbegintimes)
				}
				continue
			} else {

				// bigger than 3KB, may be 404
				if Is404(content) {
					if MyConfig.Categorylocalkeep {
						util.SaveToFile(keepdirtemp+"404", []byte(""))
					}
					return ip, nil
				}
				if IsRobot(content) {
					dspider.Errortimes = dspider.Errortimes + 1
					AmazonListLog.Errorf("get %s fail(%d),total(%d),ip:%s", temp, dspider.Errortimes, dspider.Fetchtimes, ip)
					if proxy && ok && dspider.Errortimes > MyConfig.Proxymaxtrytimes {
						// die sent
						ipendtimes := util.GetSecord2DateTimes(util.GetSecordTimes())
						insidetemp := ipbegintimes + "|" + ipendtimes + "|" + util.IS(dspider.Fetchtimes-dspider.Errortimes) + "|" + util.IS(dspider.Errortimes)
						RedisClient.Hset(MyConfig.Proxyhashpool, ip, insidetemp)
						// you know it
						Spiders.Delete(ip)
						// get new proxy again
						ip = GetIP()
						AmazonListLog.Errorf("更换IP：%s", ip)
						ipbegintimes = util.GetSecord2DateTimes(util.GetSecordTimes())
						RedisClient.Hset(MyConfig.Proxyhashpool, ip, ipbegintimes)
					}
					continue
				}
				// 3KB smaller mean page not enough
				err = spider.TooSortSizes(content, 3)
				if err != nil {
					AmazonListLog.Errorf("%s,%s", temp, err.Error())
					if MyConfig.Categorylocalkeep {
						util.SaveToFile(keepdirtemp+"x", []byte(""))
					}
					return ip, nil
				}
				dudu := []byte("")
				var e error = nil
				for {
					temtemp := url + fmt.Sprintf("?_encoding=UTF8&pg=%d&ajax=1&isAboveTheFold=0", i)
					dudu, e = Download(ip, temtemp)
					dspider, ok := Spiders.Get(ip)
					if e != nil {
						AmazonListLog.Errorf("%s:%s", temtemp, e.Error())
						if proxy && ok && dspider.Errortimes > MyConfig.Proxymaxtrytimes {
							// die sent
							ipendtimes := util.GetSecord2DateTimes(util.GetSecordTimes())
							insidetemp := ipbegintimes + "|" + ipendtimes + "|" + util.IS(dspider.Fetchtimes-dspider.Errortimes) + "|" + util.IS(dspider.Errortimes)
							RedisClient.Hset(MyConfig.Proxyhashpool, ip, insidetemp)
							// you know it
							Spiders.Delete(ip)
							// get new proxy again
							ip = GetIP()
							AmazonListLog.Errorf("更换IP：%s", ip)
							ipbegintimes = util.GetSecord2DateTimes(util.GetSecordTimes())
							RedisClient.Hset(MyConfig.Proxyhashpool, ip, ipbegintimes)
						}
						continue
					}
					if IsRobot(dudu) {
						dspider.Errortimes = dspider.Errortimes + 1
						AmazonListLog.Errorf("get %s fail(%d),total(%d),ip:%s", temtemp, dspider.Errortimes, dspider.Fetchtimes, ip)
						if proxy && ok && dspider.Errortimes > MyConfig.Proxymaxtrytimes {
							// die sent
							ipendtimes := util.GetSecord2DateTimes(util.GetSecordTimes())
							insidetemp := ipbegintimes + "|" + ipendtimes + "|" + util.IS(dspider.Fetchtimes-dspider.Errortimes) + "|" + util.IS(dspider.Errortimes)
							RedisClient.Hset(MyConfig.Proxyhashpool, ip, insidetemp)
							// you know it
							Spiders.Delete(ip)
							// get new proxy again
							ip = GetIP()
							AmazonListLog.Errorf("更换IP：%s", ip)
							ipbegintimes = util.GetSecord2DateTimes(util.GetSecordTimes())
							RedisClient.Hset(MyConfig.Proxyhashpool, ip, ipbegintimes)
						}
						continue
					}
					break
				}

				if MyConfig.Categorylocalkeep {
					util.SaveToFile(keepdirtemp, []byte(string(content)+string(dudu)+"<div id='hunterhug'>"+name+"|"+bigname+"|"+url+"|"+filename+"</div>"))
				}
				insertlist, err := ParseList([]byte(string(content) + string(dudu) + "<div id='hunterhug'>" + name + "|" + bigname + "|" + url + "|" + filename + "</div>"))
				if err != nil {
					AmazonListLog.Errorf("Parse %s-error:%s", temp, err.Error())
					break
				}
				createtime := util.TodayString(6)
				err = InsertAsinMysql(insertlist, createtime, filename)
				if err != nil {
					AmazonListLog.Errorf("InsertMysql %s-error:%s", temp, err.Error())

				}
				break
			}
		}
	}
	return ip, nil
}

// most import
// hash pool is just to record times so that we can run a program to clear deal pool which url is timeout
func GetUrls() error {
	AmazonListLog.Log("Start Get List url")
	ip := "127.0.0.1"
	// do a lot url still can't pop url
	for {
		// take url!block!!!
		// urlmap such like 1-10-1|https://www.amazon.com/Best-Sellers-Appliances-Clothes-Dryers/zgbs/appliances/13397481
		// take a url and throw it into deal pool
		urlmap, err := RedisClient.Brpoplpush(MyConfig.Urlpool, MyConfig.Urldealpool, 0)
		if err != nil {
			return err
		}
		exist, _ := RedisClient.Hexists(MyConfig.Urlhashpool, urlmap)
		if exist {
			AmazonListLog.Errorf("exist %s", urlmap)
			continue
		}
		urlbegintime := util.GetSecord2DateTimes(util.GetSecordTimes())
		// start catch url
		temp := strings.Split(urlmap, "|")

		// if url not right,continue
		if len(temp) != 5 {
			continue
		}
		url := temp[1]
		filename := temp[0]
		name := temp[2]
		bigname := temp[3]
		pagetemp := temp[4]
		page, err := util.SI(pagetemp)
		if err != nil {
			continue
		}
		if page > 5 || page < 0 {
			page = 5
		}
		//Todo
		//if SpiderType == JP {
		ip, _ = GetJapanUrl(ip, url, filename, name, bigname, page)
		//} else {
		//ip, _ = GetUrl(ip, url, filename, name, bigname, page)
		//}
		// done! rem redis deal pool
		RedisClient.Lrem(MyConfig.Urldealpool, 0, urlmap)
		// throw it to a hash pool
		urlendtimes := util.GetSecord2DateTimes(util.GetSecordTimes())
		RedisClient.Hset(MyConfig.Urlhashpool, urlmap, urlbegintime+"|"+urlendtimes)
	}
	return nil
}

// most import
// hash pool is just to record times so that we can run a program to clear deal pool which url is timeout
func GetNoneProxyUrls(taskname string) error {
	AmazonListLog.Log("Start Get List url")
	ip := "*" + taskname
	// do a lot url still can't pop url
	for {
		// take url!block!!!
		// urlmap such like 1-10-1|https://www.amazon.com/Best-Sellers-Appliances-Clothes-Dryers/zgbs/appliances/13397481
		// take a url and throw it into deal pool
		urlmap, err := RedisClient.Brpoplpush(MyConfig.Urlpool, MyConfig.Urldealpool, 0)
		if err != nil {
			return err
		}
		exist, _ := RedisClient.Hexists(MyConfig.Urlhashpool, urlmap)
		if exist {
			AmazonListLog.Errorf("exist %s", urlmap)
			continue
		}
		urlbegintime := util.GetSecord2DateTimes(util.GetSecordTimes())

		// start catch url
		temp := strings.Split(urlmap, "|")

		// if url not right,continue
		if len(temp) != 5 {
			continue
		}
		url := temp[1]
		filename := temp[0]
		name := temp[2]
		bigname := temp[3]
		pagetemp := temp[4]
		page, err := util.SI(pagetemp)
		if err != nil {
			continue
		}
		if page > 5 || page < 0 {
			page = 5
		}
		//japan you can!!
		GetJapanUrl(ip, url, filename, name, bigname, page)

		// done! rem redis deal pool
		RedisClient.Lrem(MyConfig.Urldealpool, 0, urlmap)
		// throw it to a hash pool
		urlendtimes := util.GetSecord2DateTimes(util.GetSecordTimes())
		RedisClient.Hset(MyConfig.Urlhashpool, urlmap, urlbegintime+"|"+urlendtimes)
	}
	return nil
}
