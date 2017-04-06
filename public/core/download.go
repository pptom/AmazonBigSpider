package core

import (
	"github.com/hunterhug/GoSpider/spider"
	"github.com/hunterhug/GoSpider/util"
	"math/rand"
	"strings"
	"sync"
)

var (
	Spiders = &_Spider{brower: make(map[string]*spider.Spider)}
	Ua      = map[int]string{}
)

type _Spider struct {
	mux    sync.RWMutex
	brower map[string]*spider.Spider
}

func (sb *_Spider) Get(name string) (b *spider.Spider, ok bool) {
	sb.mux.RLock()
	b, ok = sb.brower[name]
	sb.mux.RUnlock()
	return
}

func (sb *_Spider) Set(name string, b *spider.Spider) {
	sb.mux.Lock()
	sb.brower[name] = b
	sb.mux.Unlock()
	return
}

func (sb *_Spider) Delete(name string) {
	sb.mux.Lock()
	delete(sb.brower, name)
	sb.mux.Unlock()
	return
}
func init() {
	Ua[0] = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.71 Safari/537.36"
	temp, err := util.ReadfromFile(Dir + "/config/ua.txt")
	if err != nil {
		panic(err.Error())
	} else {
		uas := strings.Split(string(temp), "\n")

		for i, ua := range uas {
			Ua[i] = strings.TrimSpace(strings.Replace(ua, "\r", "", -1))
		}
	}

}
func Download(ip string, url string) ([]byte, error) {
	if strings.Contains(ip, "*") {
		return NonProxyDownload(ip, url)
	}
	browser, ok := Spiders.Get(ip)
	if ok {
		browser.Url = url
		content, err := browser.Get()
		spider.Logger.Debugf("url:%s,status:%d,ip:%s,ua:%s", url, browser.UrlStatuscode, ip, browser.Header.Get("User-Agent"))
		return content, err
	} else {
		proxy := "http://" + ip
		browser, _ := spider.NewSpider(proxy)
		browser.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		browser.Header.Set("Accept-Language", "en-US;q=0.8,en;q=0.5")
		browser.Header.Set("Connection", "keep-alive")
		if strings.Contains(url, "www.amazon.co.jp") {
			browser.Header.Set("Host", "www.amazon.co.jp")
		} else if strings.Contains(url, "www.amazon.de") {
			browser.Header.Set("Host", "www.amazon.de")
		} else if strings.Contains(url, "www.amazon.co.uk") {
			browser.Header.Set("Host", "www.amazon.co.uk")
		} else {
			browser.Header.Set("Host", "www.amazon.com")
		}
		browser.Header.Set("Upgrade-Insecure-Requests", "1")
		browser.Header.Set("User-Agent", Ua[rand.Intn(len(Ua)-1)])
		browser.Url = url
		Spiders.Set(ip, browser)
		content, err := browser.Get()
		spider.Logger.Debugf("url:%s,status:%d,ip:%s,ua:%s", url, browser.UrlStatuscode, ip, browser.Header.Get("User-Agent"))
		return content, err
	}
}

func NonProxyDownload(ip string, url string) ([]byte, error) {
	browser, ok := Spiders.Get(ip)
	if ok {
		browser.Url = url
		content, err := browser.Get()
		spider.Logger.Debugf("url:%s,status:%d,ip:%s,ua:%s", url, browser.UrlStatuscode, ip, browser.Header.Get("User-Agent"))
		return content, err
	} else {
		browser, _ := spider.NewSpider(nil)
		browser.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		browser.Header.Set("Accept-Language", "en-US;q=0.8,en;q=0.5")
		browser.Header.Set("Connection", "keep-alive")
		if strings.Contains(url, "www.amazon.co.jp") {
			browser.Header.Set("Host", "www.amazon.co.jp")
		} else if strings.Contains(url, "www.amazon.de") {
			browser.Header.Set("Host", "www.amazon.de")
		} else if strings.Contains(url, "www.amazon.co.uk") {
			browser.Header.Set("Host", "www.amazon.co.uk")
		} else {
			browser.Header.Set("Host", "www.amazon.com")
		}
		browser.Header.Set("Upgrade-Insecure-Requests", "1")
		browser.Header.Set("User-Agent", Ua[rand.Intn(len(Ua)-1)])
		browser.Url = url
		Spiders.Set(ip, browser)
		content, err := browser.Get()
		spider.Logger.Debugf("url:%s,status:%d,ip:%s,ua:%s", url, browser.UrlStatuscode, ip, browser.Header.Get("User-Agent"))
		return content, err
	}
}

func GetIP() string {
	spider.Logger.Debug("Get IP...")
	iptemp, ierr := RedisClient.Brpop(0, MyConfig.Proxypool)
	// ip null return,maybe forever not happen
	if ierr != nil {
		panic("ip:" + ierr.Error())
	}
	ip := iptemp[1]
	spider.Logger.Debug("Get IP done:" + ip)
	return ip
}
