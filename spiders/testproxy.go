package main

import (
	"fmt"
	"github.com/hunterhug/AmazonBigSpider"
	"github.com/hunterhug/GoSpider/spider"
	"github.com/hunterhug/GoSpider/util"
	"strings"
)

func GetIPfromglobal(ipstring string) []string {
	ipsmart2016 := []string{}
	tempips := strings.Split(ipstring, "\n")
	for _, tempip := range tempips {
		ip := strings.TrimSpace(strings.Replace(tempip, "\r", "", -1))
		dudu := strings.Split(ip, ".")
		if len(dudu) != 4 {
			continue
		} else {
			IPdudu := true
			for _, d := range dudu {
				tempd := d
				d1 := strings.Split(d, "@")
				if len(d1) == 2 {
					tempd = d1[1]
				}
				if len(d1) > 2 {
					IPdudu = false
					break
				}
				d2 := strings.Split(tempd, ":")
				if len(d2) > 2 {
					IPdudu = false
					break
				}
				tempd = d2[0]
				dnum, de := util.SI(tempd)
				if de != nil {
					IPdudu = false
					break
				}
				if dnum > 255 || dnum <= 0 {
					IPdudu = false
					break
				}
			}
			if IPdudu {
				ipsmart2016 = append(ipsmart2016, ip)
			}
		}
	}
	return ipsmart2016
}

func main() {
	d, e := util.ReadfromFile(dudu.Dir + "/ip.txt")
	if e != nil {
		fmt.Println("ip.txt problem" + e.Error())
	} else {
		dududu := GetIPfromglobal(string(d))
		if len(dududu) == 0 {
			fmt.Println("ip.txt为空")
		} else {
			fmt.Println(len(dududu))
			for _, ip := range dududu {
				client, _ := spider.NewSpider("http://" + ip)
				client.Url = "http://ip.42.pl/short"
				b, e := client.Get()
				if e != nil {
					fmt.Printf("%v:%v\n", ip, e.Error())
				} else {
					fmt.Printf("%v:%v\n", ip, string(b))
				}
			}
		}
	}
}
