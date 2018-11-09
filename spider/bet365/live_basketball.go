package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/*下载每天数据的数据库结构*/
type MatchInfo struct {
	Name        string `orm:"column(name)"`
	HTName      string `orm:"column(ht_name)"`
	HTFirst     string `orm:"column(ht_first_score)"`
	HTSecond    string `orm:"column(ht_second_score)"`
	HTThird     string `orm:"column(ht_third_score)"`
	HTFourth    string `orm:"column(ht_fourth_score)"`
	VTName      string `orm:"column(vt_name)"`
	VTFirst     string `orm:"column(vt_first_score)"`
	VTSecond    string `orm:"column(vt_second_score)"`
	VTThird     string `orm:"column(vt_third_score)"`
	VTFourth    string `orm:"column(vt_fourth_score)"`
	ExpectTotal int    `orm:"column(expect_total)"`
	MatchTime   string `orm:"column(match_time)"`
}

func GetDetailInfo(strHref string) {
	var stMatchInfo MatchInfo

	parResultHref, _ := url.Parse(strHref)

	parResultHref.RawQuery = parResultHref.Query().Encode()
	httpSearchRespHref, err := http.Get(parResultHref.String())

	if err != nil {

		i := 0
		for ; i < 5; i++ {
			httpSearchRespHref, err = http.Get(parResultHref.String())
			if err != nil {
				continue
			} else {
				break
			}
		}
		if i == 5 {
			fmt.Println("Error is ", err)
			return
		}
	}

	defer httpSearchRespHref.Body.Close()

	docHref, _ := goquery.NewDocumentFromReader(httpSearchRespHref.Body)

	/*获取房源的所有介绍，根据不同类型进行存储*/
	lstContentHref := docHref.Find("div#live").Find("div.list-right")
	lstContentHref.Each(func(i int, contentSelectionHref *goquery.Selection) {
		strHead := contentSelectionHref.Find("span.h2-title").Find("font.alt-top").Text()

		stMatchInfo.Name = strHead

		stMatchInfo.HTName = contentSelectionHref.Find("div.tbody-one").Find("span.name").Find("span.o-hidden").Text()
		stMatchInfo.HTFirst = strings.Replace(contentSelectionHref.Find("div.tbody-one").Find("div.race-node").Find("div.float-left").Eq(0).Text(), " ", "", -1)
		stMatchInfo.HTSecond = strings.Replace(contentSelectionHref.Find("div.tbody-one").Find("div.race-node").Find("div.float-left").Eq(1).Text(), " ", "", -1)
		stMatchInfo.HTThird = strings.Replace(contentSelectionHref.Find("div.tbody-one").Find("div.race-node").Find("div.float-left").Eq(2).Text(), " ", "", -1)
		stMatchInfo.HTFourth = strings.Replace(contentSelectionHref.Find("div.tbody-one").Find("div.race-node").Find("div.float-left").Eq(3).Text(), " ", "", -1)

		strExpectTotal := contentSelectionHref.Find("div.tbody-one").Find("div.size").Children().Eq(0).Text()
		strExpectTotal = strings.Replace(strExpectTotal, "大", "", -1)
		strExpectTotal = strings.Replace(strExpectTotal, "小", "", -1)
		strExpectTotal = strings.Replace(strExpectTotal, " ", "", -1)
		fExpectTotal, _ := strconv.ParseFloat(strExpectTotal, 64)
		stMatchInfo.ExpectTotal = int(fExpectTotal + 0.5)

		stMatchInfo.VTName = contentSelectionHref.Find("div.tbody-tow").Find("span.name").Find("span.o-hidden").Text()
		stMatchInfo.VTFirst = strings.Replace(contentSelectionHref.Find("div.tbody-tow").Find("div.race-node").Find("div.float-left").Eq(0).Text(), " ", "", -1)
		stMatchInfo.VTSecond = strings.Replace(contentSelectionHref.Find("div.tbody-tow").Find("div.race-node").Find("div.float-left").Eq(1).Text(), " ", "", -1)
		stMatchInfo.VTThird = strings.Replace(contentSelectionHref.Find("div.tbody-tow").Find("div.race-node").Find("div.float-left").Eq(2).Text(), " ", "", -1)
		stMatchInfo.VTFourth = strings.Replace(contentSelectionHref.Find("div.tbody-tow").Find("div.race-node").Find("div.float-left").Eq(3).Text(), " ", "", -1)

		if stMatchInfo.HTSecond == "-" {
			stMatchInfo.HTSecond = "0"
			stMatchInfo.VTSecond = "0"
		}
		if stMatchInfo.HTThird == "-" {
			stMatchInfo.HTThird = "0"
			stMatchInfo.VTThird = "0"
		}
		if stMatchInfo.HTFourth == "-" {
			stMatchInfo.HTFourth = "0"
			stMatchInfo.VTFourth = "0"
		}

		stMatchInfo.MatchTime = time.Now().Format("2006-01-02")

		fmt.Println(stMatchInfo)

	})

}

func main() {
	fmt.Println("Start =", time.Now().Format("2006-01-02 15:04:05"))

	strUrl := "https://live.leisu.com/lanqiu"
	GetDetailInfo(strUrl)

	fmt.Println("End =", time.Now().Format("2006-01-02 15:04:05"))
}
