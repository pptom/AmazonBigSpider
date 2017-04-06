package core

import (
	"errors"
	"fmt"
	"github.com/hunterhug/GoSpider/util"
	"strconv"
	"strings"
)

func InsertAsinMysql(items []map[string]string, createtime string, category string) error {
	var returnerr error
	category = strings.Split(category, "-")[0]
	sql := `INSERT INTO smart_asin(updatetime,id,createtime,isvalid,times,category) VALUES(?,?,?,1,1,?) on duplicate key update isvalid=1,updatetime=?,times=times+1,category=?`
	for _, item := range items {
		if item["asin"] == "" {
			continue
		}
		_, err := BasicDb.Insert(sql, createtime, item["asin"], createtime, category, createtime, category)
		if err != nil {
			AmazonListLog.Errorf("ASIN update Mysql error:%s,%s,%s", item["asin"], createtime, err.Error())
		} else {
			AmazonListLog.Debugf("ASIN update Mysql:%s,%s", item["asin"], createtime)
		}

	}
	for _, item := range items {
		if item["asin"] == "" {
			continue
		}
		sql = "INSERT IGNORE INTO `" + Today + "` (id,purl,img,smallrank,name,bigname,rbigname,title,asin,url,price,score,reviews,createtime,iscatch) "
		sql = sql + "VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,0)"
		if len(item["title"]) > 200 {
			item["title"] = util.Substr(item["title"], 0, 200)
		}
		_, err := DataDb.Insert(sql, item["id"], item["purl"], item["img"], item["smallrank"], item["name"], item["bigname"], item["bigname"], item["title"], item["asin"], item["url"], item["price"], item["score"], item["reviews"], createtime)
		if err != nil {
			AmazonListLog.Errorf("20161111 Mysql err:%s,%s,%s", item["asin"], createtime, err.Error())
		} else {
			AmazonListLog.Debugf("20161111 Mysql:%s,%s", item["asin"], createtime)
		}
		goodinfo := fmt.Sprintf("%s|%s|%s|%s", item["img"], item["price"], item["score"], item["reviews"])
		if MyConfig.Extrafromredis {
			RedisClient.Hset(MyConfig.Otherhashpool, item["asin"], goodinfo)
		}
		// sent to redis
		if MyConfig.Asinautopool {
			_, errp := RedisClient.Lpush(MyConfig.Asinpool, URL+"/dp/"+item["asin"])
			if errp != nil {
				AmazonListLog.Errorf("Sent to Asinpool Redis error:%v,%v", item["asin"], errp)
			} else {
				AmazonListLog.Debugf("Sent to Asinpool Redis :%s", item["asin"])
			}
		}
	}
	return returnerr
}

func SetAsinInvalid(url string) error {
	return nil
	// Todo not need
	temp := strings.Split(url, "/dp/")
	if len(temp) != 2 {
		return errors.New(url + " is error")
	}
	sql := "UPDATE smart_asin SET isvalid=0 where id=? limit 1"
	_, err := BasicDb.Insert(sql, temp[1])
	return err
}

func SetAsinToRightCategory(asin, num string) error {
	return nil
	// Todo not need
	sql := "UPDATE smart_asin SET category=? where id=? limit 1"
	_, err := BasicDb.Insert(sql, num, asin)
	return err
}

func CreateAsinTables() error {
	sql := `CREATE TABLE IF NOT EXISTS ` + "`%s`" + `(
  id VARCHAR(150),
  purl varchar(255) DEFAULT NULL COMMENT '父类类目链接',
  col1 varchar(255) DEFAULT NULL COMMENT '预留字段',
  col2 varchar(255) DEFAULT NULL,
  img varchar(255) DEFAULT NULL,
  iscatch tinyint(4) DEFAULT '0' COMMENT '已抓取是1',
  smallrank INT NULL COMMENT '小类排名',
  name VARCHAR(255) NULL COMMENT '小类名',
  bigname VARCHAR(255) NULL COMMENT '大类名',
  rbigname VARCHAR(255) NULL COMMENT '真实大类名',
  title TEXT NULL COMMENT '商品标题',
  asin VARCHAR(255) NULL,
  url VARCHAR(255) NULL,
  rank INT NULL COMMENT '大类排名',
  soldby VARCHAR(255) NULL COMMENT '卖家',
  shipby VARCHAR(255) NULL COMMENT '物流',
  price VARCHAR(255) NULL COMMENT '价格',
  score FLOAT NULL COMMENT '打分',
  reviews INT NULL COMMENT '评论数',
  commenttime VARCHAR(255) NULL COMMENT '第一条评论时间',
  createtime VARCHAR(255) NULL,
  updatetime VARCHAR(255) NULL,
  PRIMARY KEY (id))ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err := DataDb.Create(fmt.Sprintf(sql, Today))
	return err
}

func CreateAsinRankTables() error {
	sql := `CREATE TABLE IF NOT EXISTS ` + "`Asin%s`" + `(
  id VARCHAR(150),
  bigname VARCHAR(255) NULL COMMENT '大类名',
  title TEXT NULL COMMENT '商品标题',
  rank INT NULL COMMENT '大类排名',
  price FLOAT NULL,
  sold VARCHAR(255) NULL COMMENT '自营',
  ship VARCHAR(255) NULL COMMENT 'FBA',
  score FLOAT NULL COMMENT '打分',
  reviews INT NULL COMMENT '评论数',
  createtime VARCHAR(255) NULL,
  img VARCHAR(255) NULL,
  PRIMARY KEY (id))ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err := DataDb.Create(fmt.Sprintf(sql, Today))
	return err
}

func InsertDetailMysql(item map[string]string) error {
	sql := "REPLACE INTO `Asin" + Today + "`(id,bigname,title,rank,createtime,price,reviews,ship,sold,score,img) VALUES(?,?,?,?,?,?,?,?,?,?,?)"
	// Todo
	if item["id"] == "" || item["id"] == "null" {
		return nil
	}
	if item["bigname"] == "null" || item["bigname"] == "" {
		item["bigname"] = "NULL"
	}
	rank, de := util.SI(item["rank"])
	if de != nil {
		item["rank"] = "987654321"
	} else {
		if rank != 987654321 && rank > MyConfig.Rank {
			return nil
		}
	}
	extrainfo := map[string]string{}
	var extraerr error = nil
	if MyConfig.Extrafromredis {
		tempinfo, e := RedisClient.Hget(MyConfig.Otherhashpool, item["id"])
		if e != nil {
			extraerr = e
		} else {
			tempdudu := strings.Split(tempinfo, "|")
			if len(tempdudu) != 4 {
				extraerr = errors.New("extra info redis error")
			} else {
				AmazonAsinLog.Debugf("Asin %s extrainfo:%s", item["id"], tempinfo)
				extrainfo["img"] = tempdudu[0]
				//goodinfo := fmt.Sprintf("%s|%s|%s|%s", item["img"], item["price"], item["score"], item["reviews"])
				extrainfo["price"] = tempdudu[1]
				extrainfo["score"] = tempdudu[2]
				extrainfo["reviews"] = tempdudu[3]
			}
		}
	} else {
		querysql := "SELECT price,score,reviews,img FROM `" + Today + "` where asin=? limit 1;"
		r, e := DataDb.Select(querysql, item["id"])
		if e != nil {
			extraerr = e
		} else {
			if len(r) == 0 {
				extraerr = errors.New("no small info")
			} else {
				rr := r[0]
				extrainfo["img"] = rr["img"].(string)
				extrainfo["reviews"] = rr["reviews"].(string)
				extrainfo["score"] = rr["score"].(string)
				extrainfo["price"] = rr["price"].(string)
			}
		}
	}
	pricetemp := 0.0
	reviews := 0
	score := 0.0
	img := ""
	if extraerr != nil {
		AmazonAsinLog.Error(extraerr.Error())
	} else {
		img = extrainfo["img"]
		reviewst := extrainfo["reviews"]
		reviews, _ = util.SI(reviewst)

		if SpiderType == DE {
			scoret := strings.Replace(extrainfo["score"], ",", ".", -1)
			score, _ = strconv.ParseFloat(scoret, 32)
			asinr := strings.Split(extrainfo["price"], "-")[0]
			asinr = strings.TrimSpace(asinr)
			//EUR 8,90
			price := strings.Replace(strings.Replace(asinr, "EUR", "", -1), " ", "", -1)
			price = strings.Replace(price, ".", "", -1)
			price = strings.Replace(price, ",", ".", -1)
			price = strings.Replace(price, " ", "", -1)
			pricetemp, _ = strconv.ParseFloat(price, 32)
		} else if SpiderType == JP {
			scoret := extrainfo["score"]
			score, _ = strconv.ParseFloat(scoret, 32)
			asinr := strings.Split(extrainfo["price"], "-")[0]
			asinr = strings.TrimSpace(asinr)
			price := strings.Replace(strings.Replace(asinr, "￥", "", -1), ",", "", -1)
			price = strings.Replace(price, " ", "", -1)
			pricetemp, _ = strconv.ParseFloat(price, 32)
		} else if SpiderType == UK {
			scoret := extrainfo["score"]
			score, _ = strconv.ParseFloat(scoret, 32)
			asinr := strings.Split(extrainfo["price"], "-")[0]
			asinr = strings.TrimSpace(asinr)
			price := strings.Replace(strings.Replace(asinr, "£", "", -1), ",", "", -1)
			price = strings.Replace(price, " ", "", -1)
			pricetemp, _ = strconv.ParseFloat(price, 32)
		} else {
			scoret := extrainfo["score"]
			score, _ = strconv.ParseFloat(scoret, 32)
			asinr := strings.Split(extrainfo["price"], "-")[0]
			price := strings.Replace(strings.Replace(asinr, "$", "", -1), ",", "", -1)
			price = strings.Replace(price, " ", "", -1)
			pricetemp, _ = strconv.ParseFloat(price, 32)
		}
	}
	if _, ok := item["ship"]; ok == false {
		item["ship"] = ""
	}
	if _, ok := item["sold"]; ok == false {
		item["sold"] = ""
	}
	if len(item["title"]) > 200 {
		item["title"] = util.Substr(item["title"], 0, 200)
	}
	_, err := DataDb.Insert(sql, item["id"], item["bigname"], item["title"], item["rank"], item["createtime"], pricetemp, reviews, item["ship"], item["sold"], score, img)
	// Todo
	tablename := hashcode(item["id"])
	AmazonAsinLog.Debugf("insert A %s: %s", tablename, item["id"])
	hahasql := "REPLACE INTO `A" + tablename + "`(id,day,bigname,title,rank,createtime,price,reviews,ship,sold,score,img) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)"
	_, errr := HashDb.Insert(hahasql, item["id"], Today, item["bigname"], item["title"], item["rank"], item["createtime"], pricetemp, reviews, item["ship"], item["sold"], score, img)
	if err == nil {
		return errr
	}
	return err
}

func hashcode(asin string) string {
	dd := []byte(util.Md5("iloveyou"+asin+"hunterhug") + util.Md5(asin))
	sum := 0
	for _, i := range dd {
		sum = sum + int(i)
	}
	hashcode := sum % 501
	return util.IS(hashcode)
}
