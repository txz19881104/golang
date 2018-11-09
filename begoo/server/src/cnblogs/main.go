package main

import (
	"cnblogs/cors"
	"cnblogs/models"
	_ "cnblogs/routers"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func main() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:   []string{"Content-Length", "Access-Control-Allow-Origin"},
	}))

	go GetBasketBallLive()

	beego.Run()
}

func GetBasketBallLive() {
	timer(GetBasketBallLiveInfo)
}

func timer(timer func()) {
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ticker.C:
			timer()
		}
	}
}

func GetBasketBallLiveInfo() {
	strHref := "https://live.leisu.com/lanqiu"
	var stMatchInfo models.MatchInfo

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
	} else {
		models.ArrMatchInfo = []models.MatchInfo{}
	}

	defer httpSearchRespHref.Body.Close()

	docHref, err := goquery.NewDocumentFromReader(httpSearchRespHref.Body)

	if err != nil {
		return
	}

	/*获取房源的所有介绍，根据不同类型进行存储*/
	lstContentHref := docHref.Find("div#live").Find("div.list-right")
	lstContentHref.Each(func(i int, contentSelectionHref *goquery.Selection) {
		strHead := contentSelectionHref.Find("span.h2-title").Find("font.alt-top").Text()

		stMatchInfo.Name = strHead

		stMatchInfo.VTName = contentSelectionHref.Find("div.tbody-one").Find("span.name").Find("span.o-hidden").Text()
		strVTFirst := strings.Replace(contentSelectionHref.Find("div.tbody-one").Find("div.race-node").Find("div.float-left").Eq(0).Text(), " ", "", -1)
		stMatchInfo.VTFirst, _ = strconv.ParseInt(strings.Replace(strVTFirst, "-", "0", -1), 10, 64)
		strVTSecond := strings.Replace(contentSelectionHref.Find("div.tbody-one").Find("div.race-node").Find("div.float-left").Eq(1).Text(), " ", "", -1)
		stMatchInfo.VTSecond, _ = strconv.ParseInt(strings.Replace(strVTSecond, "-", "0", -1), 10, 64)
		strVTThird := strings.Replace(contentSelectionHref.Find("div.tbody-one").Find("div.race-node").Find("div.float-left").Eq(2).Text(), " ", "", -1)
		stMatchInfo.VTThird, _ = strconv.ParseInt(strings.Replace(strVTThird, "-", "0", -1), 10, 64)
		strVTFourth := strings.Replace(contentSelectionHref.Find("div.tbody-one").Find("div.race-node").Find("div.float-left").Eq(3).Text(), " ", "", -1)
		stMatchInfo.VTFourth, _ = strconv.ParseInt(strings.Replace(strVTFourth, "-", "0", -1), 10, 64)

		strExpectTotal := contentSelectionHref.Find("div.tbody-one").Find("div.size").Children().Eq(0).Text()
		strExpectTotal = strings.Replace(strExpectTotal, "大", "", -1)
		strExpectTotal = strings.Replace(strExpectTotal, "小", "", -1)
		strExpectTotal = strings.Replace(strExpectTotal, " ", "", -1)
		stMatchInfo.ExpectTotal, _ = strconv.ParseFloat(strExpectTotal, 64)

		stMatchInfo.HTName = contentSelectionHref.Find("div.tbody-tow").Find("span.name").Find("span.o-hidden").Text()
		strHTFirst := strings.Replace(contentSelectionHref.Find("div.tbody-tow").Find("div.race-node").Find("div.float-left").Eq(0).Text(), " ", "", -1)
		stMatchInfo.HTFirst, _ = strconv.ParseInt(strings.Replace(strHTFirst, "-", "0", -1), 10, 64)
		strHTSecond := strings.Replace(contentSelectionHref.Find("div.tbody-tow").Find("div.race-node").Find("div.float-left").Eq(1).Text(), " ", "", -1)
		stMatchInfo.HTSecond, _ = strconv.ParseInt(strings.Replace(strHTSecond, "-", "0", -1), 10, 64)
		strHTThird := strings.Replace(contentSelectionHref.Find("div.tbody-tow").Find("div.race-node").Find("div.float-left").Eq(2).Text(), " ", "", -1)
		stMatchInfo.HTThird, _ = strconv.ParseInt(strings.Replace(strHTThird, "-", "0", -1), 10, 64)
		strHTFourth := strings.Replace(contentSelectionHref.Find("div.tbody-tow").Find("div.race-node").Find("div.float-left").Eq(3).Text(), " ", "", -1)
		stMatchInfo.HTFourth, _ = strconv.ParseInt(strings.Replace(strHTFourth, "-", "0", -1), 10, 64)

		stMatchInfo.MatchTime = contentSelectionHref.Find("span.race-trim").Text() + contentSelectionHref.Find("span.code-time").Text()

		models.ArrMatchInfo = append(models.ArrMatchInfo, stMatchInfo)

	})

}
