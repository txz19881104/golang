package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/text/encoding/simplifiedchinese"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"
)

/*下载每天数据的数据库结构*/
type MatchInfo struct {
	ID            int    `orm:"column(pk_id)"`
	Name          string `orm:"column(name)"`
	HTName        string `orm:"column(ht_name)"`
	VTName        string `orm:"column(vt_name)"`
	HTTotalScore  string `orm:"column(ht_total_score)"`
	VTTotalScore  string `orm:"column(vt_total_score)"`
	HTHalfScore   string `orm:"column(ht_half_score)"`
	VTHalfScore   string `orm:"column(vt_half_score)"`
	HTTotalCorner string `orm:"column(ht_total_corner)"`
	VTTotalCorner string `orm:"column(vt_total_corner)"`
	DetailHref    string `orm:"column(detail_href)"`
	GoalTime      string `orm:"column(goal_time)"`
	MatchTime     string `orm:"column(match_time)"`
}

/*下载每天数据的数据库结构*/
type MatchInfoDB struct {
	ID            int    `orm:"column(pk_id)"`
	Name          string `orm:"column(name)"`
	HTName        string `orm:"column(ht_name)"`
	VTName        string `orm:"column(vt_name)"`
	HTTotalScore  int    `orm:"column(ht_total_score)"`
	VTTotalScore  int    `orm:"column(vt_total_score)"`
	HTHalfScore   int    `orm:"column(ht_half_score)"`
	VTHalfScore   int    `orm:"column(vt_half_score)"`
	HTTotalCorner int    `orm:"column(ht_total_corner)"`
	VTTotalCorner int    `orm:"column(vt_total_corner)"`
	DetailHref    string `orm:"column(detail_href)"`
	GoalTime      string `orm:"column(goal_time)"`
	MatchTime     string `orm:"column(match_time)"`
}

var wg sync.WaitGroup
var goroutine_cnt = make(chan int, 10) /*最大协程数量*/
var db *sql.DB

func DecodeToGBK(text string) (string, error) {

	dst := make([]byte, len(text)*2)
	tr := simplifiedchinese.GB18030.NewDecoder()
	nDst, _, err := tr.Transform(dst, []byte(text), true)
	if err != nil {
		return text, err
	}

	return string(dst[:nDst]), nil
}

func GetDetailInfo(strHref string, strData string) {
	var stMatchInfo MatchInfo
	response, err := getRep(strHref)
	if response == nil || err != nil || (response != nil && response.StatusCode != 200) {

		i := 0
		for ; i < 50; i++ {
			time.Sleep(1 * time.Millisecond)
			response, err = getRep(strHref)
			if response == nil || err != nil || (response != nil && response.StatusCode != 200) {
				continue
			} else {
				break
			}
		}
		if i == 50 {
			fmt.Println(strHref, "download failed")
			<-goroutine_cnt
			wg.Done()
			return
		}
	}

	docHref, err := goquery.NewDocumentFromResponse(response)
	if docHref == nil {
		fmt.Println(strHref, "download failed")
		<-goroutine_cnt
		wg.Done()
		return
	}

	/*获取房源的所有介绍，根据不同类型进行存储*/
	lstContentHref := docHref.Find("div#list").Find("ul.layout-grid-list").Find("li.list-item")

	lstContentHref.Each(func(i int, contentSelectionHref *goquery.Selection) {
		strHead := contentSelectionHref.Find("a.event-name").Find("span.display-i-b").Text()
		stMatchInfo.Name = strings.Replace(strHead, " ", "", -1)

		stMatchInfo.HTName = strings.Replace(contentSelectionHref.Find("a.name").Eq(0).Text(), " ", "", -1)
		stMatchInfo.VTName = strings.Replace(contentSelectionHref.Find("a.name").Eq(1).Text(), " ", "", -1)
		lstTotalScore := strings.Split(contentSelectionHref.Find("span.score").Text(), "-")
		if len(lstTotalScore) == 2 {
			stMatchInfo.HTTotalScore = lstTotalScore[0]
			stMatchInfo.VTTotalScore = lstTotalScore[1]
		} else {
			stMatchInfo.HTTotalScore = "-1"
			stMatchInfo.VTTotalScore = "-1"
		}

		lstHalfScore := strings.Split(contentSelectionHref.Find("span.lab-half").Text(), "-")
		if len(lstHalfScore) == 2 {
			stMatchInfo.HTHalfScore = lstHalfScore[0]
			stMatchInfo.VTHalfScore = lstHalfScore[1]
		} else {
			stMatchInfo.HTHalfScore = "-1"
			stMatchInfo.VTHalfScore = "-1"
		}

		lstTotalCorner := strings.Split(contentSelectionHref.Find("span.lab-corner").Text(), "-")
		if len(lstTotalCorner) == 2 {
			stMatchInfo.HTTotalCorner = lstTotalCorner[0]
			stMatchInfo.VTTotalCorner = lstTotalCorner[1]
		} else {
			stMatchInfo.HTTotalCorner = "-1"
			stMatchInfo.VTTotalCorner = "-1"
		}

		strDetailHref, _ := contentSelectionHref.Find("a.link").Eq(1).Attr("href")
		stMatchInfo.DetailHref = "http:" + strDetailHref
		//fmt.Println(stMatchInfo.DetailHref)

		bErr := false
		fmt.Println(stMatchInfo.DetailHref)
		responseDetail, err := getRep(stMatchInfo.DetailHref)
		if responseDetail == nil || err != nil || (responseDetail != nil && responseDetail.StatusCode != 200) {

			i := 0
			for ; i < 3; i++ {
				time.Sleep(1 * time.Millisecond)
				responseDetail, err = getRep(stMatchInfo.DetailHref)
				if responseDetail == nil || err != nil || (responseDetail != nil && responseDetail.StatusCode != 200) {
					continue
				} else {
					break
				}
			}
			if i == 3 {
				bErr = true
				fmt.Println("Error is ", err)
			}
		}

		if !bErr {
			docDetailHref, _ := goquery.NewDocumentFromResponse(responseDetail)
			if docDetailHref != nil {
				lstDetailHref := docDetailHref.Find("li.goal")
				strGoalTime := ""
				lstDetailHref.Each(func(i int, contentDetailHref *goquery.Selection) {
					strTime := strings.Replace(contentDetailHref.Find("div.time").Text(), "'", "", -1)
					if strings.Contains(contentDetailHref.Find("p.clearfix-row").Text(), stMatchInfo.HTName) {
						strGoalTime = strGoalTime + strTime + "-1;"
					} else {
						strGoalTime = strGoalTime + strTime + "-2;"
					}
				})

				stMatchInfo.GoalTime = strGoalTime
				stMatchInfo.MatchTime = strData

				/*插入数据库*/
				strSql := fmt.Sprintf("INSERT into Bet365.leisu_football_data_test (name, ht_name, vt_name, ht_total_score, vt_total_score, ht_half_score, vt_half_score, ht_total_corner, vt_total_corner, detail_href, goal_time, match_time) VALUES (\"%s\", \"%s\", \"%s\", %s, %s, %s, %s, %s, %s, \"%s\",\"%s\", \"%s\");",
					stMatchInfo.Name, stMatchInfo.HTName, stMatchInfo.VTName, stMatchInfo.HTTotalScore, stMatchInfo.VTTotalScore, stMatchInfo.HTHalfScore, stMatchInfo.VTHalfScore,
					stMatchInfo.HTTotalCorner, stMatchInfo.VTTotalCorner, stMatchInfo.DetailHref, stMatchInfo.GoalTime, stMatchInfo.MatchTime)
				_, err = db.Exec(strSql)
				if err != nil {

					i := 0
					for ; i < 5; i++ {
						_, err = db.Exec(strSql)
						if err != nil {
							time.Sleep(1 * time.Millisecond)
							continue
						} else {
							break
						}
					}
					if i == 5 {
						fmt.Println("Error is ", err)
						fmt.Println(strSql)
					}
				}
			}

		}

	})
	<-goroutine_cnt

	wg.Done()
}

func getdateArray(strStartTime string, strEndTime string) []string {
	//获取本地location
	strTimeLayout := "2006-01-02"                                        //转化所需模板
	loc, _ := time.LoadLocation("Local")                                 //重要：获取时区
	theTime, _ := time.ParseInLocation(strTimeLayout, strStartTime, loc) //使用模板在对应时区转化为time.time类型
	nStartTime := theTime.Unix()
	theTime, _ = time.ParseInLocation(strTimeLayout, strEndTime, loc) //使用模板在对应时区转化为time.time类型
	nEndTime := theTime.Unix()
	var dateArray []string
	i := nStartTime
	for ; i <= nEndTime; i += 24 * 60 * 60 {
		strDataTime := time.Unix(i, 0).Format("2006-01-02") //设置时间戳 使用模板格式化为日期字符串
		strDataTime = strings.Replace(strDataTime, "-", "", -1)
		dateArray = append(dateArray, strDataTime)
	}

	return dateArray
}

/*打开数据库，使用多核处理*/
func InitNumCpu() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	db, _ = sql.Open("mysql", "txz:passwd@tcp(198.13.54.7:3306)/LianjiaDB?charset=utf8")
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)
	db.Ping()
}

func main() {
	fmt.Println("Start =", time.Now().Format("2006-01-02 15:04:05"))

	InitNumCpu()

	runtime.GOMAXPROCS(runtime.NumCPU())
	if true {
		strUrlKeyword := "http://live.leisu.com/wanchang?date="
		//20150401-20151231 20161020-20161231
		dayarray := getdateArray("2016-10-20", "2016-12-31")
		//dayarray = []string{"20150319"}
		wg.Add(len(dayarray))
		for _, value := range dayarray {

			strUrl := strUrlKeyword + value
			fmt.Println(strUrl)
			goroutine_cnt <- 1
			go GetDetailInfo(strUrl, value)
		}
		wg.Wait()
	}
	fmt.Println("End =", time.Now().Format("2006-01-02 15:04:05"))
}

/**
* 返回response
 */
func getRep(urls string) (*http.Response, error) {

	ip_port := "http://transfer.mogumiao.com:9001"
	request, _ := http.NewRequest("GET", urls, nil)

	//随机返回User-Agent 信息
	request.Header.Set("User-Agent", getAgent())
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Proxy-Authorization", "Basic UzB5YnhjQzBOQWcwbmNmYTpEcnVQbnNmTXVnRHRNamx6")

	proxy, err := url.Parse(ip_port)
	//设置超时时间

	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 50,
	}

	response, err := client.Do(request)

	return response, err
}

/**
* 随机返回一个User-Agent
 */
func getAgent() string {
	agent := [...]string{
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0",
		"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
		"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; The World)",
		"User-Agent,Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"User-Agent, Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)",
		"User-Agent,Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	len := len(agent)
	return agent[r.Intn(len)]
}
