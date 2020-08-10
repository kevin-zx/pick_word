package pick_word

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/kevin-zx/keysitematch"
	"github.com/kevin-zx/pick_word/fengchao_rpc"
	siteinfocrawler "github.com/kevin-zx/site-info-crawler"
	"github.com/kevin-zx/site-info-crawler/sitethrougher"
	"github.com/kevin-zx/wordproperty"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type ExpandLimit struct {
	MaxDeep               int
	Si                    *sitethrougher.SiteInfo
	ContentMatchRateLimit float64
}

func ExpandWords(seedWords []string, client fengchao_rpc.FengchaoServiceClient, area []string, el ExpandLimit) ([]*fengchao_rpc.WordInfo, error) {
	var isExpand = make(map[string]bool)
	seedWords = removeDuplicate(seedWords)
	maxTaskLen := len(seedWords) + 4
	tasks := make(chan task, maxTaskLen)
	results := make(chan result, 10000)
	wiMap := map[string]*fengchao_rpc.WordInfo{}
	for _, keyword := range seedWords {
		isExpand[keyword] = true
		tasks <- task{
			Keyword: keyword,
			deep:    1,
		}
	}
	go expand(tasks, results, client)

	var completeStatus = false
	timeOut := time.NewTimer(20 * time.Second)
	for !completeStatus {
		select {
		case result := <-results:
			qr := result.w
			if !keywordInfoIsNormal(*qr) {
				continue
			}
			all := int(qr.PcPv + qr.MobilePv)
			wiMap[result.w.Word] = result.w
			if len(wiMap)%97 == 0 {
				fmt.Printf("word len is: %d\n", len(wiMap))
			}
			if all > 5000 {
				if !isExpand[qr.Word] {
					isExpand[qr.Word] = true
					go addToTask(el, tasks, maxTaskLen, qr.Word, result.deep, all)
					// expand by area
					if len(area) != 0 && all > 10000 {
						ok, _, _ := wordproperty.IsChinaArea(qr.Word)
						if !ok {
							for i := range area {
								areaWord := area[i] + qr.Word
								if !isExpand[areaWord] {
									isExpand[areaWord] = true
									go addToTask(el, tasks, maxTaskLen, areaWord, result.deep, all)
								}

							}
						}
					}
				}

			}
			timeOut.Reset(20 * time.Second)
		case <-timeOut.C:
			completeStatus = true
		default:
			//fmt.Printf("-")
			time.Sleep(10 * time.Millisecond)
		}
	}
	close(tasks)
	close(results)
	timeOut.Stop()
	var rs []*fengchao_rpc.WordInfo
	for _, m := range wiMap {
		rs = append(rs, m)
	}
	return rs, nil
}

func addToTask(el ExpandLimit, tasks chan task, maxTaskLen int, word string, currentDeep int, pv int) {
	//go func(taskKeyword string, currentDeep int) {

	if currentDeep >= el.MaxDeep {
		return
	}
	if el.Si != nil {
		km := keysitematch.Match(el.Si, []string{word})
		if km[word].MaxContentMatchRate <= el.ContentMatchRateLimit {
			return
		}
	}

	for len(tasks) >= maxTaskLen-4 {
		time.Sleep(1*time.Second)
	}
	log.Printf("%d/%d %s %d\n", len(tasks), maxTaskLen, word, pv)
	tasks <- task{
		Keyword: word,
		deep:    currentDeep + 1,
	}
	//}(qr.Word, deep)
}

func PickWord(siteUrl string, seedWords []string, careWords []string, client fengchao_rpc.FengchaoServiceClient, area []string) {
	//siteUrl := "http://www.szchenxing.cn/"
	var keywordsPV = make(map[string]int)
	si, err := goThoughtSite(siteUrl, 500, sitethrougher.PortPC, "./")
	if err != nil {
		panic(err)
	}

	sitethrougher.FillSiteLinksDetailHrefText(si)
	r, err := os.Create("r1.csv")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	csvW := csv.NewWriter(r)
	words, err := ExpandWords(seedWords, client, area, ExpandLimit{
		MaxDeep:               2,
		Si:                    si,
		ContentMatchRateLimit: 0.97,
	})

	var keywords []string
	for _, word := range words {
		keywords = append(keywords, word.Word)
		keywordsPV[word.Word] = int(word.MobilePv + word.PcPv)
	}
	keywords = removeDuplicate(keywords)

	keywordMatchGet(si, keywords, csvW, careWords, client, keywordsPV)

}

type result struct {
	w    *fengchao_rpc.WordInfo
	deep int
}

type task struct {
	Keyword string
	deep    int
}

func expand(tasks <-chan task, resultChan chan<- result, client fengchao_rpc.FengchaoServiceClient) {
	c := 0
	for task := range tasks {
		c++
		if c%7 == 0 {
			log.Printf("执行到了 %d 个关键词任务 当前深度: %d\n", c, task.deep)
		}
		res, err := client.Expand(context.Background(), &fengchao_rpc.ExpandRequest{
			Word: task.Keyword,
		})
		//qrs, err := es.ExpandWordsByQuery(keyword, 0)
		if err != nil {
			panic(err)
		}
		if res == nil {
			log.Printf("expand word:%s err\n", task.Keyword)
			continue
		}

		for _, qr := range res.WordExpand.WordInfos {
			if ok, _ := wordproperty.EnvWordProperty(qr.Word); ok {
				continue
			}
			resultChan <- result{
				w:    qr,
				deep: task.deep,
			}
		}

	}

}

func keywordInfoIsNormal(qr fengchao_rpc.WordInfo) bool {
	rate := float64(qr.PcPv+1) / float64(qr.MobilePv+1)
	de := math.Abs(float64(qr.PcPv - qr.PcPv))
	return !((rate > 10 || rate < 0.1) && de > 300)
	//if (rate > 10 || rate < 0.1) && de > 300 {
	//	continue
	//}
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

func keywordMatchGet(si *sitethrougher.SiteInfo, keywords []string, csvWrite *csv.Writer, careWords []string, client fengchao_rpc.FengchaoServiceClient, keywordsPV map[string]int) {
	km := keysitematch.Match(si, keywords)
	//qer := apiUtil.NewQueryService(ePAuthHeader)

	//fmt.Printf("关键词匹配信息：\n")
	//csvWrite.WriteString(fmt.Sprintf("关键词匹配信息：\n"))
	//fmt.Printf("关键词,主页匹配类型,主页匹配度,锚文本全匹配数,title全匹配数,内容全匹配数,匹配指数\n")
	_ = csvWrite.Write([]string{"关键词", "主页匹配类型", "最大内容匹配度", "SPV", "主页匹配度", "锚文本全匹配数", "title全匹配数", "内容全匹配数", "匹配指数", "是否关注", "关注词", "地域"})
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
			wifR, err := client.Query(context.Background(), &fengchao_rpc.QueryRequest{
				Words: []string{keyword},
			})
			if err != nil {
				log.Printf("%+v\n", err)
			}
			if wifR != nil {
				for _, info := range wifR.WordInfos {
					spv = int(info.PcPv + info.MobilePv)
				}
			}

		}
		careInfo := "0"
		cw := ""
		for _, word := range careWords {
			if keyword == word {
				careInfo = "2"
				cw = word
				break
			}

		}
		if cw == "" {
			for _, word := range careWords {
				if strings.Contains(keyword, word) {
					cw = word
					careInfo = "1"
					break
				}
			}
		}
		is, _, a := wordproperty.IsChinaArea(keyword)
		area := ""
		if is {
			area = a[0].ShortName
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
			careInfo,
			cw,
			area,
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
		pathRoads = pathRoads[0 : len(pathRoads)-1]
	} else {
		pathRoads = pathRoads[0:len(pathRoads)]
	}
	return siteinfocrawler.RunSiteWithCache(siteURLRaw, strings.Join(pathRoads, "/"), 24, myOption)
}
