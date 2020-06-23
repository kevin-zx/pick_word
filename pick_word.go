package pick_word

import (
	"encoding/csv"
	"fmt"
	"github.com/kevin-zx/baiduApiSDK/apiUtil"
	"github.com/kevin-zx/baiduApiSDK/baiduSDK"
	"github.com/kevin-zx/keysitematch"
	siteinfocrawler "github.com/kevin-zx/site-info-crawler"
	"github.com/kevin-zx/site-info-crawler/sitethrougher"
	"github.com/kevin-zx/wordproperty"
	"os"
	"strconv"
	"strings"
)

var ePAuthHeader = &baiduSDK.AuthHeader{
	Username: "baidu-酷讯2732150-7",
	Password: "Hotel^Kuxun789",
	Token:    "d0a3c5f9ea56ab0e4e73db39f9c8bc36",
	Action:   "API-SDK",
}
var keywordsPV = make(map[string]int)

func PickWord(siteUrl string, seedWords []string) {
	//siteUrl := "http://www.szchenxing.cn/"
	si, err := goThoughtSite(siteUrl, 1000, sitethrougher.PortPC, "./data/selectWords/")
	if err != nil {
		panic(err)
	}
	es := apiUtil.NewQueryExpandService(ePAuthHeader)

	r, err := os.Create("./data/selectWords/r1.csv")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	csvW := csv.NewWriter(r)
	var keywords []string
	for _, keyword := range seedWords {
		qrs, err := es.ExpandWordsByQuery(keyword, 0)
		if err != nil {
			panic(err)
		}
		for _, qr := range *qrs {
			if ok, _ := wordproperty.EnvWordProperty(qr.Word); ok {
				continue
			}
			keywordsPV[qr.Word] = qr.Pv
			keywords = append(keywords, qr.Word)
		}
	}
	keywords = removeDuplicate(keywords)

	keywordMatchGet(si, keywords, csvW)
}

func removeDuplicate(keys []string) []string {
	km := map[string]int{}
	for _, key := range keys {
		km[key] = 1
	}
	var r []string
	for k, _ := range km {
		if k == "" {
			continue
		}
		r = append(r, k)
	}
	return r
}

func keywordMatchGet(si *sitethrougher.SiteInfo, keywords []string, csvWrite *csv.Writer) {
	km := keysitematch.Match(si, keywords)
	qer := apiUtil.NewQueryService(ePAuthHeader)
	//fmt.Printf("关键词匹配信息：\n")
	//csvWrite.WriteString(fmt.Sprintf("关键词匹配信息：\n"))
	fmt.Printf("关键词,主页匹配类型,主页匹配度,锚文本全匹配数,title全匹配数,内容全匹配数,匹配指数\n")
	_ = csvWrite.Write([]string{"关键词", "主页匹配类型", "最大内容匹配度", "SPV", "主页匹配度", "锚文本全匹配数", "title全匹配数", "内容全匹配数", "匹配指数"})
	for _, keyword := range keywords {
		kmInfo := km[keyword]
		if kmInfo == nil {
			fmt.Printf("关键词 %s,获取错误\n", keyword)
			//csvWrite.WriteString(fmt.Sprintf("关键词 %s,获取错误\n", keyword.Word))
			continue
		}
		if kmInfo.MaxContentMatchRate <= 0.4 {
			continue
		}
		spv := 0
		if pv, ok := keywordsPV[keyword]; ok {
			spv = pv
		} else {
			sis, _ := qer.Query([]string{keyword})
			for _, info := range *sis {
				spv = info.All.Pv
			}
		}

		err := csvWrite.Write([]string{
			keyword,
			kmInfo.HomePageMatchType,
			fmt.Sprintf("%.2f", kmInfo.MaxContentMatchRate),
			strconv.Itoa(spv),
			strconv.FormatFloat(kmInfo.HomePageMatchRate, 'f', 2, 32),
			strconv.Itoa(kmInfo.HrefTextMatchCount),
			strconv.Itoa(kmInfo.TitleMatchCount),
			strconv.Itoa(kmInfo.ContentMatchCount),
			strconv.Itoa(kmInfo.MatchIndex),
		})
		if err != nil {
			panic(err)
		}
		csvWrite.Flush()
	}
}
func goThoughtSite(siteURLRaw string, limitCount int, port sitethrougher.DevicePort, cachePath string) (*sitethrougher.SiteInfo, error) {
	myOption := sitethrougher.DefaultOption
	myOption.LimitCount = limitCount
	myOption.Port = port
	myOption.NeedDocument = true
	//remove path date
	pathRoads := strings.Split(cachePath, "/")
	if strings.HasSuffix(cachePath, "/") {
		pathRoads = pathRoads[0 : len(pathRoads)-2]
	} else {
		pathRoads = pathRoads[0 : len(pathRoads)-1]
	}
	return siteinfocrawler.RunSiteWithCache(siteURLRaw, strings.Join(pathRoads, "/"), 24, myOption)
}
