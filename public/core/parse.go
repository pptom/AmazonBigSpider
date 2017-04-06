package core

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/hunterhug/GoSpider/util"
	"regexp"
	"strings"
	"github.com/hunterhug/GoSpider/query"
)

func QueryBytes(content []byte) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(content)))
	return doc, err
}

// list page parse
func ParseList(content []byte) ([]map[string]string, error) {
	returnmap := []map[string]string{}
	doc, _ := QueryBytes(content)
	nametemp := strings.Split(doc.Find("#hunterhug").Text(), "|")
	namelen := len(nametemp)
	name := nametemp[0]
	bigname := ""
	url := ""
	id := ""
	if namelen >= 2 {
		bigname = nametemp[1]
	}
	if namelen > 2 {
		url = nametemp[2]
	}
	if namelen > 3 {
		id = nametemp[3]
	}
	goodclass := ".zg_itemImmersion"
	if SpiderType == JP {
		goodclass = ".zg_itemRow"
	}
	doc.Find(goodclass).Each(func(i int, node *goquery.Selection) {
		dudu := node.Find("a")
		text, exist := dudu.Attr("href")
		if exist {
			temp := map[string]string{}
			if strings.Contains(text, "/dp/") {
				t1 := strings.Split(text, "/dp/")
				if len(t1) != 2 {
					return
				}
				t2 := strings.Split(t1[1], "/")
				temp["asin"] = t2[0]
				if temp["asin"] == "" {
					return
				}

			} else if strings.Contains(text, "/product/") {
				temp := map[string]string{}
				t1 := strings.Split(text, "/product/")
				if len(t1) != 2 {
					return
				}
				t2 := strings.Split(t1[1], "/")
				temp["asin"] = t2[0]
				if temp["asin"] == "" {
					return
				}

			} else {
				return
			}
			imag := dudu.Find("img")
			temp["title"], _ = imag.Attr("alt")
			temp["img"], _ = imag.Attr("src")
			returnmap = append(returnmap, temp)

			score := strings.TrimSpace(node.Find(".a-icon-row").Text())
			temp["reviews"] = "0"
			temp["score"] = "0"
			if SpiderType == USA || SpiderType == UK {
				scoretemp := strings.Split(score, "star")
				switch len(scoretemp) {
				case 1:
					temp["score"] = strings.Replace(scoretemp[0], "out of", "", -1)
				case 2:
					temp["score"] = strings.Replace(scoretemp[0], "out of", "", -1)
					temp["reviews"] = strings.TrimSpace(strings.Replace(scoretemp[1], "s", "", -1))
				}
				temp["score"] = strings.TrimSpace(strings.Replace(temp["score"], "5", "", -1))
				temp["reviews"] = strings.Replace(temp["reviews"], ",", "", -1)
			} else if SpiderType == JP {
				scoretemp := strings.Split(strings.Replace(score, "5つ星のうち ", "", -1), "\n")
				temp["score"] = scoretemp[0]
				if len(scoretemp) > 2 {
					temp["reviews"] = scoretemp[len(scoretemp)-1]
				}
				temp["score"] = strings.TrimSpace(strings.Replace(temp["score"], " ", "", -1))
				temp["reviews"] = strings.Replace(strings.Replace(temp["reviews"], " ", "", -1), ",", "", -1)
			} else if SpiderType == DE {
				scoretemp := strings.Split(strings.Replace(score, "von 5 Sternen", "", -1), "\n")
				temp["reviews"] = "0"
				temp["score"] = scoretemp[0]
				if len(scoretemp) > 2 {
					temp["reviews"] = scoretemp[len(scoretemp)-1]
				}
				temp["score"] = strings.TrimSpace(strings.Replace(temp["score"], " ", "", -1))
				temp["score"] = strings.Replace(temp["score"], ",", ".", -1)
				temp["reviews"] = strings.Replace(strings.Replace(temp["reviews"], " ", "", -1), ",", "", -1)
			}

			if temp["reviews"] == "" {
				temp["reviews"] = "0"
			}
			if temp["score"] == "" {
				temp["score"] = "0"
			}
			temp["name"] = name
			temp["bigname"] = bigname
			temp["url"] = url
			temp["purl"] = ""
			if v, ok := Urlmap[bigname]; ok {
				temp["purl"] = v
			}
			temp["id"] = temp["asin"] + "|" + id

			temp["smallrank"] = strings.TrimSpace(strings.Replace(node.Find(".zg_rankNumber").Text(), ".", "", -1))
			temp["price"] = strings.TrimSpace(node.Find(".a-color-price").Text())
			//fmt.Printf("%v:%v\n", temp["score"], temp["reviews"])
			return
		}
	})
	if len(returnmap) == 0 {
		return nil, errors.New("parse get null")
	}
	return returnmap, nil
}

func IsRobot(content []byte) bool {
	doc, _ := query.QueryBytes(content)
	text := doc.Find("title").Text()
	// uk usa
	if strings.Contains(text, "Robot Check") {
		return true
	}
	//jp
	if strings.Contains(text, "CAPTCHA") {
		return true
	}
	//de
	if strings.Contains(text, "Bot Check") {
		return true
	}
	return false
}

func Is404(content []byte) bool {
	doc, _ := query.QueryBytes(content)
	text := doc.Find("title").Text()
	if strings.Contains(text, "Page Not Found") {
		return true
	}
	if strings.Contains(text, "404") {
		return true
	}
	//uk
	if strings.Contains(string(content), "The Web address you entered is not a functioning page on our site") {
		return true
	}
	//de
	if strings.Contains(string(content), "Suchen Sie bestimmte Informationen") {
		return true
	}
	if strings.Contains(string(content), "Suchen Sie etwas bestimmtes") {
		return true
	}
	return false
}

func ParseDetail(url string, content []byte) map[string]string {
	returnmap := map[string]string{}
	doc, _ := query.QueryBytes(content)

	// title bigname
	titletrip := "Amazon.com:"
	if SpiderType == UK {
		titletrip = "Amazon.co.uk:"
	} else if SpiderType == JP {
		titletrip = "Amazon.co.jp："
	} else if SpiderType == DE {
		titletrip = "Amazon.de:"
	}
	title := strings.Replace(doc.Find("title").Text(), titletrip, "", -1)
	bigname := "null"
	returnmap["bigname"] = bigname
	if SpiderType == JP {
		temp := strings.Split(title, "：")
		templength := len(temp)
		if templength >= 2 {
			bigname = strings.TrimSpace(temp[templength-1])
			bigname = strings.Replace(bigname, "&amp;", "&", -1)
		} else {
			temp = strings.Split(title, "|")
			templength = len(temp)
			if templength >= 2 {
				bigname = strings.TrimSpace(temp[templength-1])
				bigname = strings.Replace(bigname, "&amp;", "&", -1)
			}
		}
	} else {
		temp := strings.Split(title, ":")
		templength := len(temp)
		if templength >= 2 {
			bigname = strings.TrimSpace(temp[templength-1])
		}
	}
	if len(title) > 230 {
		title = title[0:230]
	}
	returnmap["title"] = title

	// asin
	id := "null"
	temp := strings.Split(url, "/")
	templength := len(temp)
	if templength >= 2 {
		id = strings.Split(temp[templength-1], ".")[0]
	}
	returnmap["id"] = id

	//rank
	rank := 987654321
	body := doc.Find("body").Text()
	rankall := ParseRank(body)
	ranktemp, err := util.SI(rankall[1])
	if err != nil {

	} else {
		rank = ranktemp
	}
	returnmap["rank"] = util.IS(rank)
	returnmap["createtime"] = util.GetSecord2DateTimes(util.GetSecordTimes())

	//purl
	purl := ""
	// if bottom bigname is right!
	dudu := BigReallyName(rankall[2])
	if v, ok := Urlmap[dudu]; ok {
		returnmap["bigname"] = dudu
		purl = v
		// just update category!
		num, duduok := Urlnummap[dudu]
		if duduok {
			SetAsinToRightCategory(id, num)
		}
	} else {
		//if title bigname is right
		if v, ok := Urlmap[bigname]; ok {
			// find it!
			returnmap["bigname"] = bigname
			purl = v
		}
		//-----------------
	}
	returnmap["purl"] = purl

	fba := strings.TrimSpace(doc.Find("#merchant-info").Text())
	returnmap["ship"] = "other"
	returnmap["sold"] = "非自营"
	fba = strings.Replace(fba, " ", "", -1)
	if SpiderType == USA {
		if strings.Contains(fba, "FulfilledbyAmazon") {
			returnmap["ship"] = "FBA"
		}
		if strings.Contains(fba, "ShipsfromandsoldbyAmazon.com") {
			returnmap["ship"] = "FBA"
			returnmap["sold"] = "自营"
		}
		if strings.Contains(fba, "sold by Amazon.com") {
			returnmap["sold"] = "自营"
		}
	} else if SpiderType == UK {
		if strings.Contains(fba, "FulfilledbyAmazon") {
			returnmap["ship"] = "FBA"
		}
		if strings.Contains(fba, "DispatchedfromandsoldbyAmazon") {
			returnmap["ship"] = "FBA"
			returnmap["sold"] = "自营"
		}
	} else if SpiderType == DE {
		if strings.Contains(fba, "VersanddurchAmazon") {
			returnmap["ship"] = "FBA"
		}
		if strings.Contains(fba, "VerkaufundVersanddurchAmazon") {
			returnmap["ship"] = "FBA"
			returnmap["sold"] = "自营"
		}
	} else if SpiderType == JP {
		fba = strings.Replace(fba, `"`, "", -1)
		fba = strings.Replace(fba, " ", "", -1)
		if strings.Contains(fba, "Amazon.co.jpが発送") {
			returnmap["ship"] = "FBA"
		}
		if strings.Contains(fba, "Amazon.co.jpが販売、発送") {
			returnmap["ship"] = "FBA"
			returnmap["sold"] = "自营"
		}
	} else {
		panic("spider type error")
	}
	//fmt.Printf("%#v", returnmap)
	return returnmap
}

func ParseRank(content string) []string {
	content = strings.Replace(content, `"`, "", -1)
	if SpiderType == JP {
		r, _ := regexp.Compile(`Amazon 売れ筋ランキング:?[ \n]{1,30}(.*) - ([,\d]{1,10})位`)
		//売れ筋ランキング: \n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\nドラッグストア - 10,207位 (
		back := []string{"", "", ""}
		god := r.FindAllStringSubmatch(content, -1)
		if len(god) >= 1 {
			back = god[0]
		}
		back[2] = strings.Replace(back[2], ",", "", -1)
		back[1] = BigReallyName(back[1])
		temp := back[2]
		back[2] = back[1]
		back[1] = temp
		return back
	} else if SpiderType == DE {
		r, _ := regexp.Compile(`Nr. ([,.\d]{1,10}) in (.*)[(]`)
		//Nr. 41.398 in Fremdsprachige B&uuml;cher (
		back := []string{"", "", ""}
		god := r.FindAllStringSubmatch(content, -1)
		if len(god) >= 1 {
			back = god[0]
		}
		back[1] = strings.Replace(strings.Replace(back[1], ",", "", -1), ".", "", -1)
		back[2] = BigReallyName(back[2])
		return back
	} else if SpiderType == UK {
		r, _ := regexp.Compile(`([,.\d]{1,10}) in (.*)[(]See`)
		back := []string{"", "", ""}
		god := r.FindAllStringSubmatch(content, -1)
		if len(god) >= 1 {
			back = god[0]
		}
		back[1] = strings.Replace(strings.Replace(back[1], ",", "", -1), ".", "", -1)
		back[2] = BigReallyName(back[2])
		return back
	} else {
		content = strings.Replace(content, `"`, "", -1)
		r, _ := regexp.Compile(`#([,\d]{1,10})[\s]{0,1}[A-Za-z0-9]{0,6} in ([^#;)(\n]{2,30})[\s\n]{0,1}[(]{0,1}`)
		back := []string{"", "", ""}
		god := r.FindAllStringSubmatch(content, -1)
		if len(god) >= 1 {
			back = god[0]
		}
		for i, v := range god {
			bigname := BigReallyName(v[2])
			//if title bigname is right
			if _, ok := Urlmap[bigname]; ok {
				back = god[i]
				break
			} else {
				AmazonAsinLog.Error("bobobo:" + bigname)
			}
		}

		back[1] = strings.Replace(back[1], ",", "", -1)
		back[2] = BigReallyName(back[2])
		if strings.Contains(back[2], ">") {
			return []string{"", "", ""}
		}
		return back
	}
}

func BigReallyName(name string) string {
	bigname := strings.TrimSpace(name)
	bigname = strings.Replace(bigname, ",", "", -1)
	if _, ok := Urlmap[bigname]; ok {
		return bigname
	}
	patterns := []string{
		" ", "",
		"-", "",
		"&", "",
		"_", "",
		",", "",
		".", "",
	}
	r := strings.NewReplacer(patterns...)
	dudu := strings.ToLower(r.Replace(bigname))
	switch dudu {
	case "artscrafts":
		bigname = "Arts Crafts & Sewing"
	case strings.ToLower("ArtsCraftsSewing"):
		bigname = "Arts Crafts & Sewing"
	case strings.ToLower("HomeImprovements"):
		bigname = "Home Improvement"
	case strings.ToLower("HomeandKitchen"):
		bigname = "Home & Kitchen"
	case strings.ToLower("PatioLawnGarden"):
		bigname = "Patio Lawn & Garden"
	case strings.ToLower("ToysandGames"):
		bigname = "Toys & Games"
	case strings.ToLower("videogames"):
		bigname = "Video Games"
	case "homeandgarden":
		bigname = "Home & Kitchen"
	case "homegarden":
		bigname = "Home & Kitchen"
	case "furniture":
		bigname = "Home & Kitchen"
	case "kitchen":
		bigname = "Home & Kitchen"
	case "hi":
		bigname = "Home Improvement"
	case "lawngarden":
		bigname = "Patio Lawn & Garden"
	case "photo":
		bigname = "Camera & Photo"
	case "wireless":
		bigname = "Cell Phones & Accessories"
	case "hometheater":
		bigname = "Electronics"
	case "hpc":
		bigname = "Health & Personal Care"
	case "industrial":
		bigname = "Industrial & Scientific"
	}
	return bigname
}
