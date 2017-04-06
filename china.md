# 中文介绍

?Golang全自动亚马逊全网分布式爬虫（美国，日本，德国和英国）

## 2017.04.06

爬虫包升级， 亚马逊爬虫仍然能跑！！！依赖[https://www.github.com/hunterhug/GoSpider](https://www.github.com/hunterhug/GoSpider)

网站端见：Web is GoAmazonWeb See [https://github.com/hunterhug/AmazonBigSpiderWeb](https://github.com/hunterhug/AmazonBigSpiderWeb)

文件目录

```
--config
    -- *_local_* 本地跑爬虫配置
    -- ** 远程跑爬虫配置（在根目录放一个远程.txt文件即读取此配置）
--crontab  定时器与定时器脚本
--doc
    --sql  类目URL数据，需要几个月跑一次tool文件夹python代码
--public
    --core 爬虫核心包
--sh    运行脚本
--spiders 爬虫入口
--tool
    --python 代理测试
    --url  抓类目URL工具，需要手动几个月定时跑，防止类目失效，先抓取后解析塞入数据库，需运行后取消注释代码，再运行，再取消注释。。。
```

## 2016.12.10

亚马逊爬虫支持

1. 列表页和详情页可选择代理方式
2. 多浏览器保存cookie机制
3. 机器人检测达到阈值自动换代理
4. 检测日期过期自动停止程序
5. IP池扫描周期填充代理IP
6. 支持分布式跨平台抓取
7. 高并发进程设置抓取
8. 默认网页爬取去重
9. 日志记录功能
10. 配套可视化网站，支持多角度查看数据，小类数据，大类数据，Asin数据和类目数据，支持查看每件Asin商品的历史记录，如排名，价格，打分，reviews变化。部分数据支持导出，且网站支持RBAC权限，可分配每部分数据的查看和使用权限。
11. 网络端监控爬虫，可查看爬虫当前时段数据抓取状态，爬取的进度，IP的消耗程度。   **可支持网络端启动和停止爬虫，彻底成为Saas**（待做）
12. 可自定义填入IP，如塞入其他代理IP网站API获取的IP
13. 可选择HTML文件保存本地

分布式，高并发，跨平台，多站点，多种自定义配置，极强的容错能力是这个爬虫的特点。机器数量和IP代理足够情况下，每天每个站点可满足抓取几百万的商品数据。

## 2016.12.15

1. 数据库初始化脚本优化
2. IP配置文件优化
3. 网络WEB端自定义IP

# 中文使用

安装Go环境，MYSQL和Redis

一.Go安装

Go包安装 https://yun.baidu.com/s/1jHKUGZG 选择1.6后缀msi安装

环境变量设置，请根据定时器进行相应更改

```
Path G:\smartdogo\bin
GOBIN G:\smartdogo\bin
GOPATH G:\smartdogo
GOROOT C:\Go\
```

然后

```
go get -u -v github.com/hunterhug/AmazonBigSpider
```

二. Mysql安装

https://yun.baidu.com/s/1hrF0QC8 找到mysql文件夹，下面的5.6.17.0.msi根据说明安装。

三.Redis安装

https://yun.baidu.com/s/1jHKUGZG 选择redis64bit或32bit，解压 ，然后Shift+右键 在此次打开命令窗口 输入“redis-server.exe redis.conf ”，在redis.conf设置密码smart2016，可改，已经设置好

四.修改配置文件

可以修改GOPATH，即后缀为?_config.json的配置（默认不需要改，下面注释说明需要删除）

```
{
  "Type": "USA",     //美国站类型，有四种usa,jp,uk,de
  "Datadir": "/data/db/usa",   // 文件保存位置，可选择保存，/代表在本盘下
  "Proxymaxtrytimes": 40,     // 机器人错误最大次数，超过换IP
  "Rank": 80000,               // 只保存排名在这个数字之前的商品
  "Listtasknum": 30,        // 抓列表页进程数，建议不要太大，并发最多设置50
  "Asintasknum": 30,      // 抓详情页进程数，建议不要太大，并发最多设置50
  "Localtasknum": 150,  // 本地文件处理进程数，建议不要太大，并发最多设置50，可不管
  "Proxypool": "USAIPPOOL",   // Redis IP代理池名字
  "Proxyhashpool": "USAIPPOLLHASH",  // Redis IP已用池名字
  "Proxyloophours": 24,        // 重导IP时间（小时,Redis IP池用完）
  "Proxyasin": true,         // 详情页使用代理？
  "Proxycategory": false,    //列表页使用代理？
  "Proxyinit": false,   // IP池程序每次启动是否追加，可不管
  "Urlpool": "USAURLPOOL",  //列表页待抓池名字
  "Urldealpool": "USAURLDEALPOOL", //列表页处理中池
  "Urlhashpool": "USAURLHASHPOOL",  //列表页已抓池
  "Asinpool": "USAAsinPOOL",       // 同理
  "Asindealpool": "USAAsinDEALPOOL",
  "Asinhashpool": "USAAsinHASHPOOL",
  "Otherhashpool": "USAOtherHashPOOL",  // 小类数据额外redis池，方便填充大类数据，开关在ExtraFromRedis,如果关，大类数据填充查找小类数据库，大数据下会导致慢
  "Asinautopool": true,   //列表页抓取数据后自动把Asin数据塞在Asinpool,如果设置为false，需要手动运行asinpool.exe
  "ExtraFromRedis": true,  //搭配Otherhashpool
  "Asinlocalkeep": false,   //保存详情页在Datadir
  "Categorylocalkeep": false, //保存列表页在Datadir
  "Urlsql": "SELECT distinct url,id,bigpid ,name,bigpname,page FROM smart_category where isvalid=1 order by bigpid limit 100000",  //抓取那些列表页，可改
  "Asinsql": "SELECT distinct asin as id FROM `{?}` order by bigname limit 1000000", //抓取哪些Asin，{?}是程序预带占位符，被今天日期替代，可去掉
  "Spidersleeptime": 3, // 无用
  "Spidertimeout": 35,  //链接抓取超时时间
  "Spiderloglevel": "DEBUG",  //爬虫日志记录，可不管,建议设置为ERROR，注意！！！
  "Redisconfig": {  // redis配置
    "Host": "14.215.177.40:6379",  //主机
    "Password": "smart2016",   //密码
    "DB": 0
  },
  "Redispoolsize": 100,  // redis程序库连接池最大数量，应该比Listtasknum和Asintasknum都大
  "Basicdb": {   // 基础数据库配置
    "Ip": "14.215.177.38",
    "Port": "3306",
    "Username": "root",
    "Password": "smart2016",
    "Dbname": "smart_base"
  },
  "Hashdb": {   //历史数据库配置
    "Ip": "14.215.177.38",
    "Port": "3306",
    "Username": "root",
    "Password": "smart2016",
    "Dbname": "smart_hash"
  },
  "Hashnum": 500,   //历史数据库按hashcode分表，分表数量
  "Datadb": {     // 日期数据库，按天分表
    "Ip": "14.215.177.38",
    "Port": "3306",
    "Username": "root",
    "Password": "smart2016",
    "Dbname": "smartdb"
  },
  "Ipuse": {   //要用的IP组
    "d": {    //端口和密码，密码可留空，组名所在的IP在下面
      "Port": "808",
      "Secret": "smart:smart2016"
    },
    "e": {
      "Port": "808",
      "Secret": "smart:smart2016"
    },
    "f": {
      "Port": "808",
      "Secret": "smart:smart2016"
    },
    "h": {
      "Port": "808",
      "Secret": "smart:smart2016"
    }
  },
  "Ips": {
    "d": [   //组名为d的IP们
      "146.148.149.203-254",   // 连续Ip,也可以不连续，如146.148.149.203
    ]
  }
}
```

运行程序

Windows方式：

在spiders文件夹下,进去各个站点，运行go build *.go会得到exe文件

1. 点击initsql.exe初始化数据库
2. SQL文件导入数据库，列表URL
3. 点击ippool.exe填充代理IP
4. 点击urlpool.exe填充类目URL
5. 点击listmain.exe抓取列表页
6. 点击asinmain.exe抓取详情页
7. 如果配置中Asinautopool设置为false，那么需要自己导Asin进入Redis,运行asinpool.exe

Linux定时器（helpspider.sh这些文件里面的路径要改，如果你的GOPATH不是/data/www/web/go）：

```
0 */2 * * * /sbin/ntpdate time.windows.com
10 0 * * * killall go
10 0 * * * ps -ef|grep /tmp/go-build |awk '{print $2}'|xargs -i kill {}
20 0 * * * /data/app/redis-3.2.1/src/redis-cli -a smart2016 flushall
0 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/usa/ip.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
5 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/usa/urlpool.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
10 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/usa/helpspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
20 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/usa/asinspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
0 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/jp/ip.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
5 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/jp/urlpool.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
10 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/jp/helpspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
20 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/jp/asinspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
0 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/de/ip.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
5 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/de/urlpool.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
10 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/de/helpspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
20 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/de/asinspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
0 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/uk/ip.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
5 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/uk/urlpool.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
10 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/uk/helpspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
20 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/uk/asinspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
```

跑ip.sh 后可以打开127.0.0.1:12345浏览器查看美国爬虫监控，12346，12347，12348依次为日本，英国，德国（Maybe)

![](doc/img/3.png)

重蹈代理IP，账号和密码是smart，smart2016

跑urlpool.sh可以导入需要抓取的类目URL

跑helpspider.sh可以抓列表页

跑asinspider.sh可以抓详情页

## MYSQL主从

http://blog.csdn.net/faye0412/article/details/6280761


Master
```
1.vim /etc/my.cnf

[mysqld]
log-bin=mysql-bin
server-id=1
binlog-ignore-db=information_schema
binlog-ignore-db=beauty
binlog-ignore-db=mysql

2.service mysqld restart

3.grant all privileges on *.* to 'smart'@'%' identified by '123456';

4.flush tables with read lock;

5.show master status;

 mysql-bin.000001 |     6503 |              | information_schema,beauty,mysql |

6.unlock tables;
```

Slave

```
1.vim /etc/my.cnf
[mysqld]
server-id=2

2.stop slave ;
3.change master to master_host='192.168.2.119',master_user='smart',master_password='123456',master_log_file='mysql-bin.000001', master_log_pos=6503;
4.start slave ;
5.show slave status
```

----