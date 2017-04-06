package core

import (
	"fmt"
	"github.com/hunterhug/AmazonBigSpider"
	"github.com/hunterhug/GoSpider/query"
	"github.com/hunterhug/GoSpider/util"
	"testing"
)

func TestRobot(t *testing.T) {
	bytecontent, _ := util.ReadfromFile(dudu.Dir + "/test/list/404.html")
	t.Log(IsRobot(bytecontent))
	bytecontent, _ = util.ReadfromFile(dudu.Dir + "/test/list/categorynotexist.html")
	t.Log(IsRobot(bytecontent))
	bytecontent, _ = util.ReadfromFile(dudu.Dir + "/test/list/listnull.html")
	t.Log(IsRobot(bytecontent))
	bytecontent, _ = util.ReadfromFile(dudu.Dir + "/test/list/listnormal.html")
	t.Log(IsRobot(bytecontent))
	bytecontent, _ = util.ReadfromFile(dudu.Dir + "/test/robot/robot.html")
	t.Log(IsRobot(bytecontent))
}

func Test404(t *testing.T) {
	bytecontent, _ := util.ReadfromFile(dudu.Dir + "/test/list/404.html")
	t.Log(Is404(bytecontent))
	bytecontent, _ = util.ReadfromFile(dudu.Dir + "/test/list/categorynotexist.html")
	t.Log(Is404(bytecontent))
	bytecontent, _ = util.ReadfromFile(dudu.Dir + "/test/list/listnull.html")
	t.Log(Is404(bytecontent))
	bytecontent, _ = util.ReadfromFile(dudu.Dir + "/test/list/listnormal.html")
	t.Log(Is404(bytecontent))
	bytecontent, _ = util.ReadfromFile(dudu.Dir + "/test/robot/robot.html")
	t.Log(Is404(bytecontent))
}

func TestParselist(t *testing.T) {
	bytecontent, _ := util.ReadfromFile(dudu.Dir + "/test/list/1,18-2-5-1-10.html")
	results, err := ParseList(bytecontent)
	for _, result := range results {
		t.Logf("%v:%v", result, err)
	}
}

func TestParseRank(t *testing.T) {
	bytecontent, _ := util.ReadfromFile(dudu.Dir + "/test/list/xxx2.html")
	doc, _ := query.QueryBytes(bytecontent)
	test := doc.Find("body").Text()
	fmt.Printf("%#v\n", test)
	t.Logf("%#v", ParseRank(test))
}

func TestParsedd(t *testing.T) {
	bytecontent, _ := util.ReadfromFile(dudu.Dir + "/test/list/xxx2.html")
	t.Logf("%#v", ParseDetail("/dp/dd", bytecontent))
}

func TestManyRank(t *testing.T) {
	files, _ := util.ListDir(DataDir+"/asin/20161114", "html")
	for _, file := range files {
		fmt.Printf("%s\n", file)
		bytecontent, _ := util.ReadfromFile(file)
		doc, _ := query.QueryBytes(bytecontent)
		test := doc.Find("body").Text()
		//fmt.Printf("%#v\n", test)
		fmt.Printf("%#v\n", ParseRank(test))
	}
}

func TestParserankk(t *testing.T) {
	fmt.Printf("%#v", ParseRank("#1 in Computers & Accessories > Computer Accessories > Computer Cable Adapters > Serial Adapters "))
}

func TestBig(t *testing.T) {
	fmt.Println(BigReallyName("artscr_afts"))
}
