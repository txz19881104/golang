package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/Tang-RoseChild/mahonia"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/text/encoding/simplifiedchinese"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*下载每天数据的数据库结构*/
type MatchInfoDB struct {
	ID            int    `orm:"column(pk_id)"`
	MatchID       string `orm:"column(match_id)"`
	Name          string `orm:"column(name)"`
	HTName        string `orm:"column(ht_name)"`
	VTName        string `orm:"column(vt_name)"`
	HTTotalScore  int64  `orm:"column(ht_total_score)"`
	VTTotalScore  int64  `orm:"column(vt_total_score)"`
	HTHalfScore   int64  `orm:"column(ht_half_score)"`
	VTHalfScore   int64  `orm:"column(vt_half_score)"`
	HTTotalCorner int64  `orm:"column(ht_total_corner)"`
	VTTotalCorner int64  `orm:"column(vt_total_corner)"`
	HTHalfCorner  int64  `orm:"column(ht_half_corner)"`
	VTHalfCorner  int64  `orm:"column(vt_half_corner)"`
	HTShoot       int64  `orm:"ht_shoot"`
	VTShoot       int64  `orm:"vt_shoot"`
	HTShooton     int64  `orm:"ht_shoot_on"`
	VTShooton     int64  `orm:"vt_shoot_on"`
	HTRed         int64  `orm:"ht_red"`
	VTRed         int64  `orm:"vt_red"`
	Asionodd      string `orm:"asion_odd"`
	Cornerodd     string `orm:"corner_odd"`
	Numodd        string `orm:"num_odd"`
	DetailHref    string `orm:"column(detail_href)"`
	GoalTime      string `orm:"column(goal_time)"`
	MatchTime     string `orm:"column(match_time)"`
}

var wg sync.WaitGroup
var goroutine_cnt = make(chan int, 100) /*最大协程数量*/
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

func GetBaseInfo(strHref string, strMatchDate string, arrMatchInfoDB *[]MatchInfoDB) {

	parResultHref, _ := url.Parse(strHref)

	parResultHref.RawQuery = parResultHref.Query().Encode()
	httpSearchRespHref, err := http.Get(parResultHref.String())

	if httpSearchRespHref == nil || err != nil || (httpSearchRespHref != nil && httpSearchRespHref.StatusCode != 200) {

		i := 0
		for ; i < 50; i++ {
			time.Sleep(1 * time.Millisecond)
			httpSearchRespHref, err := http.Get(parResultHref.String())
			if httpSearchRespHref == nil || err != nil || (httpSearchRespHref != nil && httpSearchRespHref.StatusCode != 200) {
				continue
			} else {
				break
			}
		}
		if i == 50 {
			fmt.Println(strHref, "download failed")

			return
		}
	}

	dec := mahonia.NewDecoder("gbk")
	rd := dec.NewReader(httpSearchRespHref.Body)
	docHref, err := goquery.NewDocumentFromReader(rd)

	if docHref == nil || err != nil {
		fmt.Println(strHref, "download failed")

		return
	}
	defer httpSearchRespHref.Body.Close()

	docHref.Find("tbody table tbody tbody tr").Each(func(i int, s *goquery.Selection) {

		var stMatchInfoDB MatchInfoDB
		var strData [10]string
		if i > 0 {
			s.Find("td").Each(func(j int, ss *goquery.Selection) {
				if j == 9 {
					strInput, _ := ss.Find("a").Attr("onclick")
					strData[j] = strings.Split(strings.Split(strInput, "(")[1], ")")[0]
				} else {
					strData[j] = ss.Text()
				}
			})

			/*非青U17
			23日16:00
			完
			[C4]马拉维 U17(中)
			5-0
			津巴布韦 U17[C3]
			2-0
			1581224*/
			if (strData[9] != " ") && (strData[2] == "完") {

				stMatchInfoDB.Name = strData[0]
				stMatchInfoDB.MatchTime = strMatchDate
				stMatchInfoDB.HTName = strData[3]

				temp := strings.Split(strData[4], "-")
				if len(temp) == 2 {
					stMatchInfoDB.HTTotalScore, _ = strconv.ParseInt(temp[0], 10, 64)
					stMatchInfoDB.VTTotalScore, _ = strconv.ParseInt(temp[1], 10, 64)
				} else {
					stMatchInfoDB.HTTotalScore = 0
					stMatchInfoDB.VTTotalScore = 0
				}

				temp = strings.Split(strData[6], "-")
				if len(temp) == 2 {
					stMatchInfoDB.HTHalfScore, _ = strconv.ParseInt(temp[0], 10, 64)
					stMatchInfoDB.VTHalfScore, _ = strconv.ParseInt(temp[1], 10, 64)
				} else {
					stMatchInfoDB.HTHalfScore = 0
					stMatchInfoDB.VTHalfScore = 0
				}

				stMatchInfoDB.VTName = strData[5]

				stMatchInfoDB.Asionodd = strData[7]
				stMatchInfoDB.Numodd = strData[8]
				stMatchInfoDB.MatchID = strData[9]

				*arrMatchInfoDB = append(*arrMatchInfoDB, stMatchInfoDB)

			}

		}

	})
	fmt.Println(len(*arrMatchInfoDB), "场")
}

func GetDetailInfo(strHref string, stMatchInfoDB MatchInfoDB) {

	parResultHref, _ := url.Parse(strHref)

	parResultHref.RawQuery = parResultHref.Query().Encode()
	httpSearchRespHref, err := http.Get(parResultHref.String())

	if httpSearchRespHref == nil || err != nil || (httpSearchRespHref != nil && httpSearchRespHref.StatusCode != 200) || httpSearchRespHref.Body == nil {

		i := 0
		for ; i < 50; i++ {
			time.Sleep(1 * time.Millisecond)
			httpSearchRespHref, err := http.Get(parResultHref.String())
			if httpSearchRespHref == nil || err != nil || (httpSearchRespHref != nil && httpSearchRespHref.StatusCode != 200) || httpSearchRespHref.Body == nil {
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

	dec := mahonia.NewDecoder("utf8")
	rd := dec.NewReader(httpSearchRespHref.Body)
	docHref, err := goquery.NewDocumentFromReader(rd)

	if docHref == nil || err != nil {
		fmt.Println(strHref, "download failed")
		<-goroutine_cnt
		wg.Done()
		return
	}
	defer httpSearchRespHref.Body.Close()

	docHref.Find("#matchData > div").Each(func(i int, s *goquery.Selection) {

		pos := strings.Contains((s.Find("table > tbody > tr:nth-child(1)").Text()), "本场技术统计")
		if pos != false {
			s.Find("table > tbody > tr").Each(func(i int, t *goquery.Selection) {
				pos := strings.Contains(t.Text(), "射门")
				if pos != false {
					if (t.Find("td.bg3").Text() == "射门") || (t.Find("td.bg4").Text() == "射门") ||
						(t.Find("td.bg3").Text() == "射门次数") || (t.Find("td.bg4").Text() == "射门次数") {
						stMatchInfoDB.HTShoot, _ = strconv.ParseInt(t.Find("td:nth-child(2)").Text(), 10, 64)
						stMatchInfoDB.VTShoot, _ = strconv.ParseInt(t.Find("td:nth-child(4)").Text(), 10, 64)

					}
				}

				pos = strings.Contains(t.Text(), "射正")
				if pos != false {
					stMatchInfoDB.HTShooton, _ = strconv.ParseInt(t.Find("td:nth-child(2)").Text(), 10, 64)
					stMatchInfoDB.VTShooton, _ = strconv.ParseInt(t.Find("td:nth-child(4)").Text(), 10, 64)
				}

				pos = strings.Contains(t.Text(), "角球")
				if pos != false {
					if (t.Find("td.bg3").Text() == "角球") || (t.Find("td.bg4").Text() == "角球") {
						stMatchInfoDB.HTTotalCorner, _ = strconv.ParseInt(t.Find("td:nth-child(2)").Text(), 10, 64)
						stMatchInfoDB.VTTotalCorner, _ = strconv.ParseInt(t.Find("td:nth-child(4)").Text(), 10, 64)
					}
				}

				pos = strings.Contains(t.Text(), "半场角球")
				if pos != false {
					stMatchInfoDB.HTHalfCorner, _ = strconv.ParseInt(t.Find("td:nth-child(2)").Text(), 10, 64)
					stMatchInfoDB.VTHalfCorner, _ = strconv.ParseInt(t.Find("td:nth-child(4)").Text(), 10, 64)
				}

				pos = strings.Contains(t.Text(), "红牌")
				if pos != false {
					stMatchInfoDB.HTRed, _ = strconv.ParseInt(t.Find("td:nth-child(2)").Text(), 10, 64)
					stMatchInfoDB.VTRed, _ = strconv.ParseInt(t.Find("td:nth-child(4)").Text(), 10, 64)
				}
			})
		}

		pos = strings.Contains((s.Find("table > tbody > tr:nth-child(1)").Text()), "详细事件")
		if pos != false {
			s.Find("table > tbody > tr").Each(func(i int, t *goquery.Selection) {

				stMatchInfoDB.GoalTime = strings.Replace(stMatchInfoDB.GoalTime, "'", "", -1)

				if stMatchInfoDB.GoalTime == "" {
					s.Find("table > tbody > tr").Each(func(i int, t *goquery.Selection) {
						tempsrc, ok := t.Find("td:nth-child(2) > img").Attr("src")
						if (ok == true) && (strings.Contains(tempsrc, "/images/bf_img/1.png") ||
							strings.Contains(tempsrc, "/images/bf_img/7.png") ||
							strings.Contains(tempsrc, "/images/bf_img/8.png")) {
							stMatchInfoDB.GoalTime += " H" + t.Find("td:nth-child(3)").Text()
						}
						tempsrc, ok = t.Find("td:nth-child(4) > img").Attr("src")
						if (ok == true) && (strings.Contains(tempsrc, "/images/bf_img/1.png") ||
							strings.Contains(tempsrc, "/images/bf_img/7.png") ||
							strings.Contains(tempsrc, "/images/bf_img/8.png")) {
							stMatchInfoDB.GoalTime += " V" + t.Find("td:nth-child(3)").Text()

						}

						stMatchInfoDB.GoalTime = strings.Replace(stMatchInfoDB.GoalTime, "'", "", -1)
					})
				}
				//	fmt.Println("")
			})
		}
	})

	stMatchInfoDB.DetailHref = strHref

	/*插入数据库*/
	strSql := fmt.Sprintf("INSERT into Bet365.leisu_football_data (match_id, name, ht_name, vt_name, ht_total_score, vt_total_score, ht_half_score, vt_half_score, ht_total_corner, vt_total_corner, ht_half_corner, vt_half_corner, ht_shoot, vt_shoot, ht_shoot_on, vt_shoot_on, ht_red, vt_red, asion_odd, corner_odd, num_odd, detail_href, goal_time, match_time) VALUES (\"%s\", \"%s\", \"%s\", \"%s\", %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, \"%s\",\"%s\", \"%s\", \"%s\",\"%s\", \"%s\");",
		stMatchInfoDB.MatchID, stMatchInfoDB.Name, stMatchInfoDB.HTName, stMatchInfoDB.VTName, stMatchInfoDB.HTTotalScore, stMatchInfoDB.VTTotalScore, stMatchInfoDB.HTHalfScore, stMatchInfoDB.VTHalfScore,
		stMatchInfoDB.HTTotalCorner, stMatchInfoDB.VTTotalCorner, stMatchInfoDB.HTHalfCorner, stMatchInfoDB.VTHalfCorner, stMatchInfoDB.HTShoot, stMatchInfoDB.VTShoot, stMatchInfoDB.HTShooton, stMatchInfoDB.VTShooton,
		stMatchInfoDB.HTRed, stMatchInfoDB.VTRed, stMatchInfoDB.Asionodd, stMatchInfoDB.Cornerodd, stMatchInfoDB.Numodd, stMatchInfoDB.DetailHref, stMatchInfoDB.GoalTime, stMatchInfoDB.MatchTime)
	_, err = db.Exec(strSql)
	if err != nil {
		fmt.Println("Error is ", err)

		i := 0
		for ; i < 5; i++ {
			_, err = db.Exec(strSql)
			if err != nil {
				fmt.Println("Error is ", err)
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

	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(100)
	db.Ping()
}

func main() {
	fmt.Println("Start =", time.Now().Format("2006-01-02 15:04:05"))

	InitNumCpu()

	runtime.GOMAXPROCS(runtime.NumCPU())
	if true {
		strUrlKeyword := "http://bf.win007.com/football/hg/Over_"
		//20150401-20151231 20161020-20161231
		dayarray := getdateArray("2018-11-28", "2018-11-30")
		//dayarray = []string{"20150319"}

		for _, value := range dayarray {
			var arrMatchInfoDB []MatchInfoDB

			strUrl := strUrlKeyword + value + ".htm"
			fmt.Print(value, " ")

			timeUse := time.Now() // get current time
			GetBaseInfo(strUrl, value, &arrMatchInfoDB)
			fmt.Print(">->-> 共用时: ", time.Since(timeUse))

			wg.Add(len(arrMatchInfoDB))
			for i := 0; i < len(arrMatchInfoDB); i++ {
				goroutine_cnt <- 1
				strUrl = "http://live.win007.com/detail/" + arrMatchInfoDB[i].MatchID + "cn.htm"
				go GetDetailInfo(strUrl, arrMatchInfoDB[i])
			}
			fmt.Println(">->-> 共用时: ", time.Since(timeUse))
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
