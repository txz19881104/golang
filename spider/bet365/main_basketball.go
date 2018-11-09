package main

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/text/encoding/simplifiedchinese"
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
	ID          int    `orm:"column(pk_id)"`
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

/*下载每天数据的数据库结构*/
type MatchInfoDB struct {
	ID          int    `orm:"column(pk_id)"`
	Name        string `orm:"column(name)"`
	HTName      string `orm:"column(ht_name)"`
	HTFirst     int    `orm:"column(ht_first_score)"`
	HTSecond    int    `orm:"column(ht_second_score)"`
	HTThird     int    `orm:"column(ht_third_score)"`
	HTFourth    int    `orm:"column(ht_fourth_score)"`
	VTName      string `orm:"column(vt_name)"`
	VTFirst     int    `orm:"column(vt_first_score)"`
	VTSecond    int    `orm:"column(vt_second_score)"`
	VTThird     int    `orm:"column(vt_third_score)"`
	VTFourth    int    `orm:"column(vt_fourth_score)"`
	ExpectTotal int    `orm:"column(expect_total)"`
	MatchTime   string `orm:"column(match_time)"`
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

	docHref, _ := goquery.NewDocumentFromReader(httpSearchRespHref.Body)

	/*获取房源的所有介绍，根据不同类型进行存储*/
	lstContentHref := docHref.Find("div.list-right")
	lstContentHref.Each(func(i int, contentSelectionHref *goquery.Selection) {
		strHead := contentSelectionHref.Find("span.h2-title").Find("font.alt-top").Text()
		stMatchInfo.Name = strHead

		stMatchInfo.HTName = contentSelectionHref.Find("div.tbody-one").Find("span.name").Find("span.o-hidden").Text()
		stMatchInfo.HTFirst = strings.Replace(contentSelectionHref.Find("div.tbody-one").Find("div.race-node").Find("div.float-left").Eq(0).Text(), " ", "", -1)
		stMatchInfo.HTSecond = strings.Replace(contentSelectionHref.Find("div.tbody-one").Find("div.race-node").Find("div.float-left").Eq(1).Text(), " ", "", -1)
		stMatchInfo.HTThird = strings.Replace(contentSelectionHref.Find("div.tbody-one").Find("div.race-node").Find("div.float-left").Eq(2).Text(), " ", "", -1)
		stMatchInfo.HTFourth = strings.Replace(contentSelectionHref.Find("div.tbody-one").Find("div.race-node").Find("div.float-left").Eq(3).Text(), " ", "", -1)

		if stMatchInfo.HTFourth == "" {
			stMatchInfo.HTThird = "0"
			stMatchInfo.HTFourth = "0"
		}

		if stMatchInfo.HTThird != "-" && stMatchInfo.HTFourth != "-" {
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
			if stMatchInfo.VTFourth == "" {
				stMatchInfo.VTThird = "0"
				stMatchInfo.VTFourth = "0"
			}

			stMatchInfo.MatchTime = strData
			//stMatchInfo.ExpectTotal
			//fmt.Println(contentSelectionHref.Children().Eq(0).Text())

			//fmt.Println(contentSelectionHref.Find("span.label").Parent().Text())

			/*插入数据库*/
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
		strUrlKeyword := "https://live.leisu.com/lanqiu/wanchang?date="

		dayarray := getdateArray("2018-10-03", "2018-10-25")
		wg.Add(len(dayarray))
		for _, value := range dayarray {
			strUrl := strUrlKeyword + value
			fmt.Println(strUrl)
			goroutine_cnt <- 1
			go GetDetailInfo(strUrl, value)
		}
		wg.Wait()
	} else {
		strSql := `SELECT * FROM Bet365.leisu_basketball_data where match_time like "201%" and name like "中国NBL 21岁以下";`
		rows, _ := db.Query(strSql)
		var fTotalNum float64 = 0
		var nDiffGT20, nDiff20To15, nDiff15To10, nDiff10To5, nDiff5To0, nDiff0ToN5, nDiffN5ToN10, nDiffN10ToN15, nDiffN15ToN20, nDiffLTN20 int = 0, 0, 0, 0, 0, 0, 0, 0, 0, 0

		nStatics := 1

		switch nStatics {
		case 1:
			for rows.Next() {
				var stMatchInfo MatchInfoDB

				_ = rows.Scan(&stMatchInfo.ID, &stMatchInfo.Name, &stMatchInfo.HTName, &stMatchInfo.HTFirst, &stMatchInfo.HTSecond, &stMatchInfo.HTThird, &stMatchInfo.HTFourth,
					&stMatchInfo.VTName, &stMatchInfo.VTFirst, &stMatchInfo.VTSecond, &stMatchInfo.VTThird, &stMatchInfo.VTFourth, &stMatchInfo.ExpectTotal, &stMatchInfo.MatchTime)

				nHTTotal := stMatchInfo.HTFirst + stMatchInfo.HTSecond + stMatchInfo.HTThird + stMatchInfo.HTFourth
				nVTTotal := stMatchInfo.VTFirst + stMatchInfo.VTSecond + stMatchInfo.VTThird + stMatchInfo.VTFourth

				nTotalScore := nHTTotal + nVTTotal
				if stMatchInfo.ExpectTotal != 0 && nTotalScore < 300 && (stMatchInfo.HTThird+stMatchInfo.HTFourth) != 0 && (stMatchInfo.VTThird+stMatchInfo.VTFourth) != 0 {

					nDiff := nTotalScore - stMatchInfo.ExpectTotal
					if nDiff >= 20 {
						nDiffGT20++
					} else if nDiff >= 15 && nDiff < 20 {
						nDiff20To15++
					} else if nDiff >= 10 && nDiff < 15 {
						nDiff15To10++
					} else if nDiff >= 5 && nDiff < 10 {
						nDiff10To5++
					} else if nDiff >= 0 && nDiff < 5 {
						nDiff5To0++
					} else if nDiff >= -5 && nDiff < 0 {
						nDiff0ToN5++
					} else if nDiff >= -10 && nDiff < -5 {
						nDiffN5ToN10++
					} else if nDiff >= -15 && nDiff < -10 {
						nDiffN10ToN15++
					} else if nDiff >= -20 && nDiff < -15 {
						nDiffN15ToN20++
					} else if nDiff < -20 {
						nDiffLTN20++
					}

					fTotalNum++

				}

			}

			fmt.Println("总样本数：", fTotalNum)
			fmt.Printf("( 20至90 ) : %d(%.2f)\n( 20至15 ) : %d(%.2f)\n( 15至10 ) : %d(%.2f)\n( 10至5  ) : %d(%.2f)\n(  5至0  ) : %d(%.2f)\n(  0至-5 ) : %d(%.2f)\n( -5至-10) : %d(%.2f)\n(-10至-15) : %d(%.2f)\n(-15至-20) : %d(%.2f)\n(-20至-90) : %d(%.2f)\n",
				nDiffGT20, ((float64(nDiffGT20) / fTotalNum) * 100),
				nDiff20To15, ((float64(nDiff20To15) / fTotalNum) * 100),
				nDiff15To10, ((float64(nDiff15To10) / fTotalNum) * 100),
				nDiff10To5, ((float64(nDiff10To5) / fTotalNum) * 100),
				nDiff5To0, ((float64(nDiff5To0) / fTotalNum) * 100),
				nDiff0ToN5, ((float64(nDiff0ToN5) / fTotalNum) * 100),
				nDiffN5ToN10, ((float64(nDiffN5ToN10) / fTotalNum) * 100),
				nDiffN10ToN15, ((float64(nDiffN10ToN15) / fTotalNum) * 100),
				nDiffN15ToN20, ((float64(nDiffN15ToN20) / fTotalNum) * 100),
				nDiffLTN20, ((float64(nDiffLTN20) / fTotalNum) * 100))

		case 2:
			for rows.Next() {
				var stMatchInfo MatchInfoDB

				_ = rows.Scan(&stMatchInfo.ID, &stMatchInfo.Name, &stMatchInfo.HTName, &stMatchInfo.HTFirst, &stMatchInfo.HTSecond, &stMatchInfo.HTThird, &stMatchInfo.HTFourth,
					&stMatchInfo.VTName, &stMatchInfo.VTFirst, &stMatchInfo.VTSecond, &stMatchInfo.VTThird, &stMatchInfo.VTFourth, &stMatchInfo.ExpectTotal, &stMatchInfo.MatchTime)

				nHTTotalHalf := stMatchInfo.HTFirst + stMatchInfo.HTSecond
				nVTTotalHalf := stMatchInfo.VTFirst + stMatchInfo.VTSecond
				nHTTotal := stMatchInfo.HTFirst + stMatchInfo.HTSecond + stMatchInfo.HTThird + stMatchInfo.HTFourth
				nVTTotal := stMatchInfo.VTFirst + stMatchInfo.VTSecond + stMatchInfo.VTThird + stMatchInfo.VTFourth

				nTotalScoreHalf := nHTTotalHalf + nVTTotalHalf
				nTotalScore := nHTTotal + nVTTotal
				if stMatchInfo.ExpectTotal != 0 && nTotalScore < 300 && (stMatchInfo.HTThird+stMatchInfo.HTFourth) != 0 && (stMatchInfo.VTThird+stMatchInfo.VTFourth) != 0 {

					var nExpectTotal int = stMatchInfo.ExpectTotal / 2
					if (nTotalScoreHalf - nExpectTotal) < -15 {

						nDiff := nTotalScore - stMatchInfo.ExpectTotal

						//fmt.Println(nExpectTotal, nTotalScoreHalf, (nTotalScoreHalf - nExpectTotal), nDiff)

						if nDiff >= 20 {
							nDiffGT20++
						} else if nDiff >= 15 && nDiff < 20 {
							nDiff20To15++
						} else if nDiff >= 10 && nDiff < 15 {
							nDiff15To10++
						} else if nDiff >= 5 && nDiff < 10 {
							nDiff10To5++
						} else if nDiff >= 0 && nDiff < 5 {
							nDiff5To0++
						} else if nDiff >= -5 && nDiff < 0 {
							//fmt.Println(stMatchInfo.Name, stMatchInfo.HTName, stMatchInfo.VTName, stMatchInfo.MatchTime)
							nDiff0ToN5++
						} else if nDiff >= -10 && nDiff < -5 {
							nDiffN5ToN10++
						} else if nDiff >= -15 && nDiff < -10 {
							nDiffN10ToN15++
						} else if nDiff >= -20 && nDiff < -15 {
							nDiffN15ToN20++
						} else if nDiff < -20 {
							nDiffLTN20++
						}

						fTotalNum++
					}

				}

			}

			fmt.Println("总样本数：", fTotalNum)
			fmt.Printf("( 20至90 ) : %d(%.2f)\n( 20至15 ) : %d(%.2f)\n( 15至10 ) : %d(%.2f)\n( 10至5  ) : %d(%.2f)\n(  5至0  ) : %d(%.2f)\n(  0至-5 ) : %d(%.2f)\n( -5至-10) : %d(%.2f)\n(-10至-15) : %d(%.2f)\n(-15至-20) : %d(%.2f)\n(-20至-90) : %d(%.2f)\n",
				nDiffGT20, ((float64(nDiffGT20) / fTotalNum) * 100),
				nDiff20To15, ((float64(nDiff20To15) / fTotalNum) * 100),
				nDiff15To10, ((float64(nDiff15To10) / fTotalNum) * 100),
				nDiff10To5, ((float64(nDiff10To5) / fTotalNum) * 100),
				nDiff5To0, ((float64(nDiff5To0) / fTotalNum) * 100),
				nDiff0ToN5, ((float64(nDiff0ToN5) / fTotalNum) * 100),
				nDiffN5ToN10, ((float64(nDiffN5ToN10) / fTotalNum) * 100),
				nDiffN10ToN15, ((float64(nDiffN10ToN15) / fTotalNum) * 100),
				nDiffN15ToN20, ((float64(nDiffN15ToN20) / fTotalNum) * 100),
				nDiffLTN20, ((float64(nDiffLTN20) / fTotalNum) * 100))
		}
	}
	fmt.Println("End =", time.Now().Format("2006-01-02 15:04:05"))
}
