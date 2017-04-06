package core

import (
	"github.com/hunterhug/GoSpider/spider"
	"testing"
)

func TestDetalurl(t *testing.T) {
	spider.SetLogLevel("debug")
	ip := "104.128.124.122:808"
	url := "https://www.amazon.com/Best-Sellers-Appliances-Home-Appliance-Warranties/zgbs/appliances/2242350011"
	_, err := GetAsinUrl(ip, url)
	if err != nil {
		t.Error(err)
	}
}
