package core

import (
	"fmt"
	"github.com/hunterhug/GoSpider/spider"
	"github.com/hunterhug/GoSpider/query"
	"github.com/hunterhug/GoSpider/util"
	"regexp"
	"testing"
)

func TestListDownload1(t *testing.T) {
	util.MakeDir(Dir + "/test/list/")
	ip := "*104.128.124.122:808"

	// debug info will no appear |nothing
	spider.SetLogLevel("info")

	url := "https://www.amazon.com/dp/B01C8RN7VG"
	content, err := Download(ip, url)
	if err != nil {
		fmt.Printf("%#v", err.Error())
	}
	util.SaveToFile(Dir+"/test/list/xxx2.html", content)
}

func TestParseRank1(t *testing.T) {
	bytecontent, _ := util.ReadfromFile(Dir + "/test/list/xxx2.html")
	fmt.Printf("%#v\n", Urlmap)
	doc, _ := query.QueryBytes(bytecontent)
	test := doc.Find("body").Text()
	fmt.Printf("%#v\n", test)
	r, _ := regexp.Compile(`#([,\d]{1,10})[\s]{0,1}[A-Za-z0-9]{0,6} in ([^#;)(\n]{2,30})[\s\n]{0,1}[(]{0,1}`)
	god := r.FindAllStringSubmatch(test, -1)
	fmt.Printf("%#v\n", god)
}

func TestParsedd1(t *testing.T) {
	bytecontent, _ := util.ReadfromFile(Dir + "/test/list/xxx2.html")
	t.Logf("%#v", ParseDetail("/dp/dd", bytecontent))
}
