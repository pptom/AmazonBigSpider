package core

import (
	"errors"
	"github.com/hunterhug/GoSpider/util"
	"strings"
)

// ip and url download,if error from zero again!
// ip := "104.128.124.122:808"
// url := "https://www.amazon.com/dp/B01MTJ1E4Q"
// filename:B01MTJ1E4Q
func GetAsinUrl(ip string, url string) ([]byte, error) {
	filename := strings.Split(url, "/dp/")
	if len(filename) != 2 {
	}
	keepdirtemp := MyConfig.Datadir + "/asin/" + Today + "/" + filename[1] + ".html"
	if MyConfig.Asinlocalkeep {
		if util.FileExist(keepdirtemp) {
			AmazonAsinLog.Debugf("FileExist:%s", keepdirtemp)
			return util.ReadfromFile(keepdirtemp)
		}
		if util.FileExist(keepdirtemp + "sql") {
			AmazonAsinLog.Debugf("FileExist: % sql", keepdirtemp)
			return util.ReadfromFile(keepdirtemp + "sql")
		}
	}
	content, err := Download(ip, url)
	if err != nil {
		return nil, err
	}
	if IsRobot(content) {
		return nil, errors.New("robot")
	}
	if Is404(content) {
		return nil, errors.New("404")
	}
	if MyConfig.Asinlocalkeep {
		util.SaveToFile(keepdirtemp, content)
	}
	return content, nil

}

// most import
func GetAsinUrls() error {
	AmazonAsinLog.Log("Start Get Asin url")
	ip := GetIP()

	// before use, send to hash pool
	ipbegintimes := util.GetSecord2DateTimes(util.GetSecordTimes())
	RedisClient.Hset(MyConfig.Proxyhashpool, ip, ipbegintimes)

	// do a lot url still can't pop url
	for {
		// take url!block!!!
		// url such like https://www.amazon.com/dp/B01MTJ1E4Q
		// take a url and throw it into deal pool
		url, err := RedisClient.Brpoplpush(MyConfig.Asinpool, MyConfig.Asindealpool, 0)
		if err != nil {
			return err
		}
		exist, _ := RedisClient.Hexists(MyConfig.Asinhashpool, url)
		if exist {
			AmazonAsinLog.Errorf("exist %s", url)
			continue
		}

		urlbegintime := util.GetSecord2DateTimes(util.GetSecordTimes())

		content := []byte("")
		err = nil
		// when error loop
		for {
			content, err = GetAsinUrl(ip, url)
			spider, ok := Spiders.Get(ip)
			if err == nil {
				break
			} else {
				if strings.Contains(err.Error(), "404") {
					break
				}
				if strings.Contains(err.Error(), "robot") {
					if ok {
						spider.Errortimes = spider.Errortimes + 1
					}
				}
				if ok {
					AmazonAsinLog.Errorf("get %s fail(%d),total(%d) error:%s,ip:%s", url, spider.Errortimes, spider.Fetchtimes, err.Error(), ip)
				}
			}
			// if proxy ip err more than config, change ip
			if ok && spider.Errortimes > MyConfig.Proxymaxtrytimes {
				// die sent
				ipendtimes := util.GetSecord2DateTimes(util.GetSecordTimes())
				insidetemp := ipbegintimes + "|" + ipendtimes + "|" + util.IS(spider.Fetchtimes-spider.Errortimes) + "|" + util.IS(spider.Errortimes)
				RedisClient.Hset(MyConfig.Proxyhashpool, ip, insidetemp)
				// you know it
				Spiders.Delete(ip)
				// get new proxy again
				ip = GetIP()
				ipbegintimes = util.GetSecord2DateTimes(util.GetSecordTimes())
				RedisClient.Hset(MyConfig.Proxyhashpool, ip, ipbegintimes)
			}
		}
		if err != nil && strings.Contains(err.Error(), "404") {
			// 404 set asin invaild
			err = SetAsinInvalid(url)
			if err != nil {
				AmazonAsinLog.Errorf("%s set invalid error:%s", url, err.Error())
			}
		} else {
			// parse detail
			info := ParseDetail(url, content)
			// insert, error you still ignore
			err := InsertDetailMysql(info)
			if err != nil {
				AmazonAsinLog.Errorf("%s mysql error:%s", url, err.Error())
			}
		}
		// done! rem redis deal pool
		RedisClient.Lrem(MyConfig.Asindealpool, 0, url)
		// throw it to a hash pool
		urlendtimes := util.GetSecord2DateTimes(util.GetSecordTimes())
		RedisClient.Hset(MyConfig.Asinhashpool, url, urlbegintime+"|"+urlendtimes)
	}
	return nil
}

func GetNoneProxyAsinUrls(taskname string) error {
	AmazonAsinLog.Log("Start Get Asin url")
	ip := "*" + taskname

	// do a lot url still can't pop url
	for {
		// take url!block!!!
		// url such like https://www.amazon.com/dp/B01MTJ1E4Q
		// take a url and throw it into deal pool
		url, err := RedisClient.Brpoplpush(MyConfig.Asinpool, MyConfig.Asindealpool, 0)
		if err != nil {
			return err
		}

		exist, _ := RedisClient.Hexists(MyConfig.Asinhashpool, url)
		if exist {
			AmazonAsinLog.Errorf("exist %s", url)
			continue
		}
		urlbegintime := util.GetSecord2DateTimes(util.GetSecordTimes())

		content := []byte("")
		err = nil
		// when error loop
		for {
			content, err = GetAsinUrl(ip, url)
			spider, ok := Spiders.Get(ip)
			if err == nil {
				break
			} else {
				if strings.Contains(err.Error(), "404") {
					break
				}
				if strings.Contains(err.Error(), "robot") {
					if ok {
						spider.Errortimes = spider.Errortimes + 1
					}
				}
				if ok {
					AmazonAsinLog.Errorf("get %s fail(%d),total(%d) error:%s,ip:%s", url, spider.Errortimes, spider.Fetchtimes, err.Error(), ip)
				}
			}
			// if proxy ip err more than config, change ip
			if ok && spider.Errortimes > MyConfig.Proxymaxtrytimes {
				// you know it
				Spiders.Delete(ip)
			}
		}
		if err != nil && strings.Contains(err.Error(), "404") {
			// 404 set asin invaild
			err = SetAsinInvalid(url)
			if err != nil {
				AmazonAsinLog.Errorf("%s set invalid error:%s", url, err.Error())
			}
		} else {
			// parse detail
			info := ParseDetail(url, content)
			// insert, error you still ignore
			err := InsertDetailMysql(info)
			if err != nil {
				AmazonAsinLog.Errorf("%s mysql error:%s", url, err.Error())
			} else {
				AmazonAsinLog.Debug("Insert!!")
			}
		}
		// done! rem redis deal pool
		RedisClient.Lrem(MyConfig.Asindealpool, 0, url)
		// throw it to a hash pool
		urlendtimes := util.GetSecord2DateTimes(util.GetSecordTimes())
		RedisClient.Hset(MyConfig.Asinhashpool, url, urlbegintime+"|"+urlendtimes)
	}
	return nil
}
