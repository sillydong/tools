package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/httplib"
	"github.com/toqueteos/webbrowser"
	"strings"
)

func ip(ipaddr string) {
	ipip(ipaddr)
	ipip_sub("tencent", ipaddr)
	ipip_sub("taobao", ipaddr)
	ipip_sub("sina", ipaddr)
	ipip_sub("baidu", ipaddr)
	//ipbaidu(ipaddr)
}

type iptimezone struct {
	CityCode    int64    `json:"city_code"`
	Continent   string   `json:"continent"`
	CountryCode string   `json:"country_code"`
	En          []string `json:"en"`
	Latitude    float64  `json:"latitude,string"`
	Longitude   float64  `json:"longitude,string"`
	Timezone    string   `json:"timezone"`
	Timezone2   string   `json:"timezone2"`
}

func ipip(ipaddr string) {
	fmt.Println("ipip")
	url := "https://www.ipip.net/ip.html"
	request := httplib.Post(url)
	request.Param("ip", ipaddr)
	request.Header("Host", "www.ipip.net")
	request.Header("Origin", "https://www.ipip.net")
	request.Header("Referer", "https://www.ipip.net/ip.html")
	request.Header("Upgrade-Insecure-Requests", "1")
	request.Header("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36")
	resp, err := request.Response()
	if err != nil {
		fmt.Println(err)
	} else {
		if resp.StatusCode != 200 {
			fmt.Println(resp.Status)
		} else {
			page, err := goquery.NewDocumentFromResponse(resp)
			if err != nil {
				fmt.Println(err)
			} else {
				myself := page.Find("#myself").Text()
				fmt.Println(strings.TrimSpace(myself))

				//reg := regexp.MustCompile(`var ip_data = ({(.*)})?`)
				//ipdata := reg.FindString(page.Text())
				//if len(ipdata) > 0 {
				//	jstring := ipdata[strings.Index(ipdata, "{") : strings.Index(ipdata, "}")+1]
				//
				//	var timezone iptimezone
				//	err := json.Unmarshal([]byte(jstring), &timezone)
				//	if err != nil {
				//		fmt.Println(err)
				//	} else {
				//		bytes, err := json.MarshalIndent(timezone, "", "    ")
				//		if err != nil {
				//			fmt.Println(err)
				//		} else {
				//			fmt.Println(string(bytes))
				//		}
				//	}
				//}

				//chinacode := page.Find(".china_code").Text()
				//fmt.Println(strings.TrimSpace(chinacode))
			}
		}
	}
}

type iplocation struct {
	Data struct {
		Area string      `json:"area"`
		Data interface{} `json:"data"`
		Isp  string      `json:"isp"`
		Type string      `json:"type"`
	} `json:"data"`
	Ip    string `json:"ip"`
	State int64  `json:"state"`
}

func ipip_sub(site, ipaddr string) {
	fmt.Println(site)
	url := "https://www.ipip.net/ip.php?a=ajax"
	request := httplib.Post(url)
	request.Param("type", site)
	request.Param("ip", ipaddr)
	request.Header("Host", "www.ipip.net")
	request.Header("Origin", "https://www.ipip.net")
	request.Header("Referer", "https://www.ipip.net/ip.html")
	request.Header("X-Requested-With", "XMLHttpRequest")
	request.Header("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36")
	resp, err := request.Response()
	if err != nil {
		fmt.Println(err)
	} else {
		if resp.StatusCode != 200 {
			fmt.Println(resp.Status)
		} else {
			var location iplocation
			err := request.ToJSON(&location)
			if err != nil {
				fmt.Println(err)
			} else {
				bytes, err := json.MarshalIndent(location.Data.Data, "", "    ")
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(string(bytes))
				}
			}
		}
	}
}

type iplocationbaidu struct {
	Address string `json:"address"`
	Content struct {
		Address       string `json:"address"`
		AddressDetail struct {
			City         string `json:"city"`
			CityCode     int64  `json:"city_code"`
			District     string `json:"district"`
			Province     string `json:"province"`
			Street       string `json:"street"`
			StreetNumber string `json:"street_number"`
		} `json:"address_detail"`
		Point struct {
			X float64 `json:"x,string"`
			Y float64 `json:"y,string"`
		} `json:"point"`
	} `json:"content"`
	Status int64 `json:"status"`
}

func ipbaidu(ipaddr string) {
	fmt.Println("baidu api")
	ak := "GeWeL6ghCAG7WoBjmxHim0Gp"
	url := fmt.Sprintf("http://api.map.baidu.com/location/ip?ip=%v&ak=%v&coor=bd09ll", ipaddr, ak)
	request := httplib.Get(url)
	resp, err := request.Response()
	if err != nil {
		fmt.Println(err)
	} else {
		if resp.StatusCode != 200 {
			fmt.Println(resp.Status)
		} else {
			var location iplocationbaidu
			err := request.ToJSON(&location)
			if err != nil {
				fmt.Println(err)
			} else {
				view(location.Content.Point.Y, location.Content.Point.X, ipaddr)
				bytes, err := json.MarshalIndent(location, "", "    ")
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(string(bytes))
				}
			}
		}
	}
}

func view(lat, log float64, ip string) {
	if lat > 0 && log > 0 {
		url := fmt.Sprintf("http://api.map.baidu.com/geocoder?location=%v,%v&coord_type=bd09ll&output=html", lat, log)
		webbrowser.Open(url)
	}
}
