package main

import (
	"database/sql"
	_ "encoding/xml"
	"fmt"
	_ "github.com/PuerkitoBio/goquery"
	_ "github.com/Tang-RoseChild/mahonia"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*下载每天数据的数据库结构*/
type MatchInfo struct {
	ID          int     `orm:"column(pk_id)"`
	Name        string  `orm:"column(name)"`
	HTName      string  `orm:"column(ht_name)"`
	HTFirst     int64   `orm:"column(ht_first_score)"`
	HTSecond    int64   `orm:"column(ht_second_score)"`
	HTThird     int64   `orm:"column(ht_third_score)"`
	HTFourth    int64   `orm:"column(ht_fourth_score)"`
	VTName      string  `orm:"column(vt_name)"`
	VTFirst     int64   `orm:"column(vt_first_score)"`
	VTSecond    int64   `orm:"column(vt_second_score)"`
	VTThird     int64   `orm:"column(vt_third_score)"`
	VTFourth    int64   `orm:"column(vt_fourth_score)"`
	ExpectTotal float64 `orm:"column(expect_total)"`
	ExpectDiff  float64 `orm:"column(expect_diff)"`
	MatchTime   string  `orm:"column(match_time)"`
}

var wg sync.WaitGroup
var goroutine_cnt = make(chan int, 100) /*最大协程数量*/
var db *sql.DB

var strUrlKeyword = "http://bf.win007.com/nba_date.aspx?date="
var strUrlEx = "http://bf.win007.com/nba_Oddsdata.aspx?date="

func DecodeToGBK(text string) (string, error) {

	dst := make([]byte, len(text)*2)
	tr := simplifiedchinese.GBK.NewDecoder()
	nDst, _, err := tr.Transform(dst, []byte(text), true)
	if err != nil {
		return text, err
	}

	return string(dst[:nDst]), nil
}

func GetDetailInfo(strData string) {
	parResultHref, _ := url.Parse(strUrlKeyword + strData)

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
			<-goroutine_cnt
			wg.Done()
			return
		}
	}

	if httpSearchRespHref.StatusCode != http.StatusOK {
		fmt.Println(httpSearchRespHref.StatusCode)
		<-goroutine_cnt
		wg.Done()
		return
	}
	defer httpSearchRespHref.Body.Close()

	body, err := ioutil.ReadAll(httpSearchRespHref.Body)

	strBody, _ := DecodeToGBK(string(body))

	mapMatchInfo := make(map[string]*MatchInfo)
	arrMatch := strings.Split(strBody, "<![CDATA[")
	for i := 0; i < len(arrMatch)-1; i++ {
		arrOneMatch := strings.Split(arrMatch[i], "^")
		if len(arrOneMatch) < 44 {
			continue
		}

		id := arrOneMatch[0]

		/*查看元素在集合中是否存在 */
		_, ok := mapMatchInfo[id] /*如果确定是真实的,则存在,否则不存在 */

		if ok {
			fmt.Println(id, "已经存在！！")
			continue
		}

		var stMatchInfo MatchInfo

		stMatchInfo.Name = strings.Split(arrOneMatch[1], ",")[0] + "_" + arrOneMatch[36]
		stMatchInfo.MatchTime = strings.Replace(arrOneMatch[4], "<br>", " ", -1)
		stMatchInfo.HTName = strings.Split(arrOneMatch[8], ",")[0]
		stMatchInfo.VTName = strings.Split(arrOneMatch[10], ",")[0]
		stMatchInfo.HTFirst, _ = strconv.ParseInt(arrOneMatch[13], 10, 64)
		stMatchInfo.VTFirst, _ = strconv.ParseInt(arrOneMatch[14], 10, 64)
		stMatchInfo.HTSecond, _ = strconv.ParseInt(arrOneMatch[15], 10, 64)
		stMatchInfo.VTSecond, _ = strconv.ParseInt(arrOneMatch[16], 10, 64)
		stMatchInfo.HTThird, _ = strconv.ParseInt(arrOneMatch[17], 10, 64)
		stMatchInfo.VTThird, _ = strconv.ParseInt(arrOneMatch[18], 10, 64)
		stMatchInfo.HTFourth, _ = strconv.ParseInt(arrOneMatch[19], 10, 64)
		stMatchInfo.VTFourth, _ = strconv.ParseInt(arrOneMatch[20], 10, 64)
		stMatchInfo.MatchTime = arrOneMatch[42] + "年" + stMatchInfo.MatchTime

		if len(stMatchInfo.MatchTime) < 23 {
			continue
		} else if len(stMatchInfo.MatchTime) > 23 {
			stMatchInfo.MatchTime = stMatchInfo.MatchTime[0:23]
		}

		mapMatchInfo[id] = &stMatchInfo

	}

	parEx, _ := url.Parse(strUrlEx + strData + "&id=3")

	parEx.RawQuery = parEx.Query().Encode()
	httpRespEx, err := http.Get(parEx.String())
	if err != nil {

		i := 0
		for ; i < 5; i++ {
			httpRespEx, err = http.Get(parEx.String())
			if err != nil {
				continue
			} else {
				break
			}
		}
		if i == 5 {
			fmt.Println("Error is ", err)
			<-goroutine_cnt
			wg.Done()
			return
		}
	}

	if httpRespEx.StatusCode != http.StatusOK {
		fmt.Println(httpRespEx.StatusCode)
		<-goroutine_cnt
		wg.Done()
		return
	}
	defer httpRespEx.Body.Close()

	body, err = ioutil.ReadAll(httpRespEx.Body)
	strBody, _ = DecodeToGBK(string(body))
	arrMatch = strings.Split(strBody, "<m>")

	for i := 0; i < len(arrMatch)-1; i++ {
		arrOneMatch := strings.Split(arrMatch[i], ",")
		if len(arrOneMatch) < 8 {
			continue
		}

		id := arrOneMatch[0]

		/*查看元素在集合中是否存在 */
		stMatchInfo, ok := mapMatchInfo[id] /*如果确定是真实的,则存在,否则不存在 */

		if !ok {
			continue
		}

		stMatchInfo.ExpectTotal, _ = strconv.ParseFloat(arrOneMatch[4], 64)
		stMatchInfo.ExpectDiff, _ = strconv.ParseFloat(arrOneMatch[1], 64)

		strSql := fmt.Sprintf("INSERT into Bet365.leisu_basketball_data (name, ht_name, ht_first_score, ht_second_score, ht_third_score, ht_fourth_score, vt_name, vt_first_score, vt_second_score, vt_third_score, vt_fourth_score, expect_total, expect_diff, match_time) VALUES (\"%s\", \"%s\", %d, %d, %d, %d, \"%s\", %d, %d, %d, %d, %f, %f, \"%s\");",
			stMatchInfo.Name, stMatchInfo.HTName, stMatchInfo.HTFirst, stMatchInfo.HTSecond, stMatchInfo.HTThird, stMatchInfo.HTFourth, stMatchInfo.VTName, stMatchInfo.VTFirst,
			stMatchInfo.VTSecond, stMatchInfo.VTThird, stMatchInfo.VTFourth, stMatchInfo.ExpectTotal, stMatchInfo.ExpectDiff, stMatchInfo.MatchTime)

		_, err = db.Exec(strSql)
		if err != nil {

			i := 0
			for ; i < 5; i++ {
				_, err = db.Exec(strSql)
				if err != nil {
					time.Sleep(100 * time.Millisecond)
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

	//var stMatchInfo MatchInfo
	/*
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
				<-goroutine_cnt
				wg.Done()
				return
			}
		}

		defer httpSearchRespHref.Body.Close()

		dec := mahonia.NewDecoder("gbk")
		rd := dec.NewReader(httpSearchRespHref.Body)

		fmt.Println(httpSearchRespHref.Body)

		docHref, _ := goquery.NewDocumentFromReader(rd)
		fmt.Println(docHref.Text())

		arrMatch := strings.Split(docHref.Text(), "]]>")
		for i := 0; i < len(arrMatch)-1; i++ {
			arrOneMatch := strings.Split(arrMatch[i], "^")
			fmt.Println(arrOneMatch)
			fmt.Println("******************************")
		}
	*/
	/*获取房源的所有介绍，根据不同类型进行存储*/

	//stMatchInfo.ExpectTotal
	//fmt.Println(contentSelectionHref.Children().Eq(0).Text())

	//fmt.Println(contentSelectionHref.Find("span.label").Parent().Text())

	/*插入数据库
	strSql := fmt.Sprintf("INSERT into Bet365.leisu_basketball_data (name, ht_name, ht_first_score, ht_second_score, ht_third_score, ht_fourth_score, vt_name, vt_first_score, vt_second_score, vt_third_score, vt_fourth_score, expect_total, match_time) VALUES (\"%s\", \"%s\", %s, %s, %s, %s, \"%s\", %s, %s, %s, %s, %d, \"%s\");",
		stMatchInfo.Name, stMatchInfo.HTName, stMatchInfo.HTFirst, stMatchInfo.HTSecond, stMatchInfo.HTThird, stMatchInfo.HTFourth, stMatchInfo.VTName, stMatchInfo.VTFirst,
		stMatchInfo.VTSecond, stMatchInfo.VTThird, stMatchInfo.VTFourth, stMatchInfo.ExpectTotal, stMatchInfo.MatchTime)

	_, err = db.Exec(strSql)
	if err != nil {

		i := 0
		for ; i < 5; i++ {
			_, err = db.Exec(strSql)
			if err != nil {
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				break
			}
		}
		if i == 5 {
			fmt.Println("Error is ", err)
			fmt.Println(strSql)
		}
	}*/

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
		//strDataTime = strings.Replace(strDataTime, "-", "", -1)
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

		dayarray := getdateArray("2010-01-01", "2018-10-31")
		wg.Add(len(dayarray))
		for _, value := range dayarray {
			fmt.Println(value)
			goroutine_cnt <- 1
			go GetDetailInfo(value)
		}
		wg.Wait()
	}

	fmt.Println("End =", time.Now().Format("2006-01-02 15:04:05"))
}
