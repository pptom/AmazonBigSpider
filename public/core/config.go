package core

import (
	"encoding/json"
	"github.com/hunterhug/AmazonBigSpider/public/log"
	"github.com/hunterhug/GoSpider/store/myredis"
	"github.com/hunterhug/GoSpider/store/mysql"
	"github.com/hunterhug/GoSpider/spider"
	"github.com/hunterhug/GoSpider/util"
	"strings"
	"fmt"
)

var (
	Dir                                       string          = util.CurDir() // now root dir core
	DataDir                                   string                          //global data dir, diff from Myconfig
	RedisClient                               myredis.MyRedis                 // redis
	BasicDb                                   mysql.Mysql                     // url db
	DataDb                                    mysql.Mysql                     // data db
	HashDb                                    mysql.Mysql
	MyConfig                                  Config // some config.json
	AmazonListLog, AmazonAsinLog, AmazonIpLog *log.Logger
	Today                                     string = util.TodayString(3) // Today
	SpiderType                                int
	URL                                       string
)

var (
	Urlmap        = map[string]string{}
	Urlnummap     = map[string]string{}
	Urlnumdudumap = map[string]string{}
)

const (
	USA   = 1
	JP    = 2
	DE    = 3
	UK    = 4
	OTHER = 5
)

type IPSecret struct {
	//    "d": {
	//  "Port": "808",
	//  "Secret": "smart:smart2016"
	//},
	Port   string
	Secret string
}
type Config struct {
	Type              string
	Datadir           string
	Rank              int
	Proxymaxtrytimes  int
	Listtasknum       int
	Asintasknum       int
	Localtasknum      int
	Proxypool         string
	Proxyhashpool     string // record times and another message
	Proxyloophours    int    // if  when ip pool empty wait some hour
	Spidersleeptime   int
	Spidertimeout     int
	Spiderloglevel    string
	Redisconfig       myredis.RedisConfig
	Redispoolsize     int
	Basicdb           mysql.MysqlConfig
	Hashdb            mysql.MysqlConfig
	Datadb            mysql.MysqlConfig
	Ipuse             map[string]IPSecret
	Ips               map[string][]string
	Urlpool           string
	Urldealpool       string
	Urlhashpool       string // record times and another message
	Asinpool          string
	Asindealpool      string
	Asinhashpool      string // record times and another message
	Otherhashpool     string
	Extrafromredis    bool
	Asinautopool      bool // url auto sent asin to redis
	Urlsql            string
	Asinsql           string
	Proxyinit         bool // every proxy ip init, del all exist ip?
	Asinlocalkeep     bool
	Categorylocalkeep bool
	Proxyasin         bool // use proxy ip?
	Proxycategory     bool
	Hashnum           int
}

func InitConfig(cfpath string, logpath string) {
	// log
	NewLog(logpath)
	// config load
	configbytes, err := util.ReadfromFile(cfpath)
	if err != nil {
		panic(err.Error())
	}
	configbytes = []byte(strings.Replace(strings.Replace(string(configbytes), "\r", "", -1), "\n", "", -1))

	err = json.Unmarshal(configbytes, &MyConfig)

	// some path adding today string
	MyConfig.Proxypool = MyConfig.Proxypool + "-" + Today
	MyConfig.Proxyhashpool = MyConfig.Proxyhashpool + "-" + Today
	// you know it
	DataDir = MyConfig.Datadir
	MyConfig.Urldealpool = MyConfig.Urldealpool + "-" + Today
	MyConfig.Urlhashpool = MyConfig.Urlhashpool + "-" + Today
	MyConfig.Urlpool = MyConfig.Urlpool + "-" + Today
	MyConfig.Asinpool = MyConfig.Asinpool + "-" + Today
	MyConfig.Asindealpool = MyConfig.Asindealpool + "-" + Today
	MyConfig.Asinhashpool = MyConfig.Asinhashpool + "-" + Today
	MyConfig.Otherhashpool = MyConfig.Otherhashpool + "-" + Today
	if err != nil {
		panic(err.Error())
	}
	spidertype := strings.ToLower(MyConfig.Type)
	switch spidertype {
	case "usa":
		SpiderType = USA
		URL = "https://www.amazon.com"
	case "jp":
		SpiderType = JP
		URL = "https://www.amazon.co.jp"
	case "de":
		SpiderType = DE
		URL = "https://www.amazon.de"
	case "uk":
		SpiderType = UK
		URL = "https://www.amazon.co.uk"
	default:
		SpiderType = OTHER
	}

	if SpiderType == OTHER {
		panic("spider type error")
	}

	MapUrl(SpiderType)
	// create dir so that no error
	util.MakeDir(MyConfig.Datadir + "/list/" + Today)
	util.MakeDir(MyConfig.Datadir + "/asin/" + Today)

	// spider log init and timeout
	spider.SetLogLevel(MyConfig.Spiderloglevel)
	spider.SetGlobalTimeout(MyConfig.Spidertimeout)

	// redis init
	redisconfig := MyConfig.Redisconfig
	redisclient, err := myredis.NewRedisPool(redisconfig, MyConfig.Redispoolsize)
	if err != nil {
		//panic("REDIS ERROR")
	}
	RedisClient = redisclient

	// db init
	BasicDb = mysql.New(MyConfig.Basicdb)
	DataDb = mysql.New(MyConfig.Datadb)
	HashDb = mysql.New(MyConfig.Hashdb)
}

func OpenMysql() {
	fmt.Println("open basicdb")
	BasicDb.Open(1000,0)
	fmt.Println("open db")
	DataDb.Open(1000,0)
	fmt.Println("open hashdb")
	HashDb.Open(1000,0)
}

func MapUrl(spidertype int) {
	urlconfig := "url.csv"
	switch spidertype {
	case USA:
		urlconfig = "usa_url.csv"
	case JP:
		urlconfig = "jp_url.csv"
	case UK:
		urlconfig = "uk_url.csv"
	case DE:
		urlconfig = "de_url.csv"
	default:
		panic("spider type error")
	}
	con, err := util.ReadfromFile(Dir + "/config/" + urlconfig)
	if err != nil {
		panic(err.Error())
	} else {
		temp := string(con)
		temp1 := strings.Split(temp, "\n")
		for _, i := range temp1 {
			j := strings.Split(strings.Replace(i, "\r", "", -1), ",")
			if len(j) != 3 {
				continue
			}
			name := j[2]
			namenum := j[0]
			Urlmap[name] = j[1]
			Urlnummap[name] = namenum
			Urlnumdudumap[namenum] = name
		}
	}
}

func NewLog(filename string) {
	logsconf, _ := util.ReadfromFile(filename)
	err := log.Init(string(logsconf))
	if err != nil {
		panic(err)
	}
	AmazonListLog = log.Get("daylist")
	AmazonAsinLog = log.Get("dayasin")
	AmazonIpLog = log.Get("dayip")
}
