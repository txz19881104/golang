package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
	//"reflect"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"strings"
)

type Score struct {
	Score        float64 `json:"score"`
	Strict_match bool    `json:"strict_match"`
}

type Content struct {
	Id     int    `json:"id"`
	Scores Score  `json:"score"`
	Title  string `json:"title"`
}

type Data struct {
	Author []string  `json:"author"`
	Topic  []Content `json:"topic"`
}

type Search struct {
	Code       int    `json:"code"`
	SearchData Data   `json:"data"`
	Message    string `json:"message"`
}

type ComicChapter struct {
	Num    int
	Name   string
	Count  int
	ImgSrc []string
}

type Comic struct {
	Name         string
	WatchNum     int
	Website      string
	ChapterNum   int
	IsFinish     int
	Introduce    string
	Author       string
	Img          string
	Type         string
	Time         string
	ComicChapter []ComicChapter
}

var wg sync.WaitGroup
var db *sql.DB

func GetHref(strHref string, stComicChapter *ComicChapter) {
	//lock.Lock() // 上锁
	httpHref, err := http.Get(strHref)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		defer httpHref.Body.Close()

		doc, _ := goquery.NewDocumentFromReader(httpHref.Body)
		//fmt.Println(reflect.TypeOf(doc))
		lstImg := doc.Find("img.kklazy")
		nImgLength := lstImg.Length()
		stComicChapter.ImgSrc = make([]string, nImgLength)

		lstImg.Each(func(i int, contentSelection *goquery.Selection) {
			src, _ := contentSelection.Attr("data-kksrc")
			stComicChapter.ImgSrc[i] = src
		})
		stComicChapter.Count = nImgLength

	}

	wg.Done()
	//lock.Unlock() // 解锁
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Error is ", err)
		os.Exit(-1)
	}
}

func InitNumCpu() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	db, _ = sql.Open("mysql", "txz:passwd@tcp(198.13.54.7:3306)/EntertainmentDB?charset=utf8")
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)
	db.Ping()
}

func Insert(lstChapter ComicChapter, nId int64, strKeyword string) {

	strSql := fmt.Sprintf(`INSERT INTO %s.%s(ChapterNum, ChapterName, PicNum, Dept_ID) VALUES (%d, "%s", %d, %d);`,
		"EntertainmentDB", "ComicChapter", lstChapter.Num+1, lstChapter.Name, lstChapter.Count, nId)
	_, err := db.Exec(strSql)
	checkErr(err)

	strFilePath := "/mnt/TecentCloud/" + strKeyword + "/" + lstChapter.Name + "/"

	for nNum, strImgSrc := range lstChapter.ImgSrc {
		strName := fmt.Sprintf("%d.jpg", nNum)
		strSrc := "https://txz-1256783950.cos.ap-beijing.myqcloud.com/Comics/" + strKeyword + "/" + lstChapter.Name + "/" + strName

		saveImages(strImgSrc, strFilePath, strName)
		if nNum == 0 {
			strSql = fmt.Sprintf(`INSERT INTO %s.%s(Page_Num, Comic_ID, Comic_ChapterNum, Img_src) VALUES (%d, %d, %d, "%s")`,
				"EntertainmentDB", "ComicImgSrc", nNum, nId, lstChapter.Num+1, strSrc)
		} else {
			strSql = fmt.Sprintf(`%s,(%d, %d, %d, "%s")`, strSql, nNum, nId, lstChapter.Num+1, strSrc)
		}
	}
	strSql += ";"
	_, err = db.Exec(strSql)
	checkErr(err)

	wg.Done()
}

func ConnectDatabase(mapComic map[string]Comic) {

	for strKey, stComic := range mapComic {
		fmt.Println(strKey)
		strSql := fmt.Sprintf(`INSERT INTO %s.%s(Name, Website, ChapterNum, IsFinish, Introduce, Author, Img, Type, Time) VALUES 
			("%s", "%s", %d, %d, "%s", "%s", "%s", "%s", "%s");`,
			"EntertainmentDB", "ComicName", stComic.Name, stComic.Website, stComic.ChapterNum, stComic.IsFinish, stComic.Introduce, stComic.Author, stComic.Img, stComic.Type, stComic.Time)
		stmt, err := db.Exec(strSql)
		checkErr(err)

		nId, err := stmt.LastInsertId()
		checkErr(err)

		fmt.Println(len(stComic.ComicChapter))
		wg.Add(len(stComic.ComicChapter))
		for _, lstChapter := range stComic.ComicChapter {
			go Insert(lstChapter, nId, stComic.Name)
		}
		wg.Wait()
	}

}

//下载图片
func saveImages(strImgUrl string, strFilePath string, strFileName string) {
	//检查目录是否存在
	file, err := os.Stat(strFilePath)
	if err != nil || !file.IsDir() {
		dir_err := os.Mkdir(strFilePath, os.ModePerm)
		if dir_err != nil {
			fmt.Println("create dir failed")
			os.Exit(1)
		}
	}

	strFile := strFilePath + strFileName
	fmt.Println("download file ", strFile)
	exists := checkExists(strFile)
	if exists {
		fmt.Println(strFile, " is exists")
		return
	}

	response, err := http.Get(strImgUrl)
	checkErr(err)

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	checkErr(err)

	image, err := os.Create(strFile)
	checkErr(err)

	defer image.Close()
	image.Write(data)
}

func checkExists(strDownloadPath string) bool {
	_, err := os.Stat(strDownloadPath)
	return err == nil
}

func main() {
	fmt.Println("Start =", time.Now().Format("2006-01-02 15:04:05"))

	InitNumCpu()
	strKeyword := "西行纪" //"航海王（海贼王）"
	//download_path := "/home/txz/download/海贼王"
	strUrlKeyword := "http://www.kuaikanmanhua.com/web/topic/search?keyword=" + strKeyword

	parResult, _ := url.Parse(strUrlKeyword)

	parResult.RawQuery = parResult.Query().Encode()
	httpSearchResp, err := http.Get(parResult.String())

	checkErr(err)

	defer httpSearchResp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(httpSearchResp.Body)
	//fmt.Println(reflect.TypeOf(doc))

	strMsg := doc.Text()
	var stSearchResult Search
	err = json.Unmarshal([]byte(strMsg), &stSearchResult)
	checkErr(err)

	for _, topic := range stSearchResult.SearchData.Topic {
		strUrlKeywordContent := "http://www.kuaikanmanhua.com/web/topic/" + strconv.Itoa(topic.Id)
		httpRespKeywordContent, err := http.Get(strUrlKeywordContent)
		checkErr(err)

		defer httpRespKeywordContent.Body.Close()
		doc, _ := goquery.NewDocumentFromReader(httpRespKeywordContent.Body)
		strAuthor := doc.Find("div.author-nickname").Text()
		strIntroduce := doc.Find("div.switch-content").Text()
		// 去除空格
		strIntroduce = strings.Replace(strIntroduce, " ", "", -1)
		// 去除换行符
		strIntroduce = strings.Replace(strIntroduce, "\n", "", -1)

		// 去除"
		strIntroduce = strings.Replace(strIntroduce, "\"", "", -1)

		strImg, _ := doc.Find("img.kk-img").Attr("src")

		lstBook := doc.Find("a.article-img")

		nLength := lstBook.Length()
		nChpterNum := nLength
		lstComicChapter := make([]ComicChapter, nLength)

		//lock := &sync.Mutex{}

		wg.Add(nChpterNum)
		lstBook.Each(func(i int, contentSelection *goquery.Selection) {
			nLength--
			//var stComicChapter ComicChapter
			name, _ := contentSelection.Attr("title")
			lstComicChapter[nLength].Name = name
			lstComicChapter[nLength].Num = nLength

			href, _ := contentSelection.Attr("href")
			strHref := "http://www.kuaikanmanhua.com" + href

			go GetHref(strHref, &lstComicChapter[nLength])

			//lstComicChapter[nLength] = stComicChapter

		})

		wg.Wait()
		/*
			for {
				lock.Lock() // 上锁
				c := counter
				lock.Unlock() // 解锁

				runtime.Gosched() // 出让时间片

				if c >= nChpterNum {
					break
				}
			}*/

		strSrc := "https://txz-1256783950.cos.ap-beijing.myqcloud.com/Comics/" + strKeyword + "/" + "封面.jpg"
		strFilePath := "/mnt/TecentCloud/" + strKeyword + "/"
		saveImages(strImg, strFilePath, "封面.jpg")

		strTime := time.Now().Format("2006-01-02 15:04:05")

		var stComic Comic
		mapSql := make(map[string]Comic)
		stComic.Name = topic.Title
		stComic.WatchNum = 0
		stComic.Website = "http://www.kuaikanmanhua.com"
		stComic.ChapterNum = nChpterNum
		stComic.IsFinish = 0
		stComic.Introduce = strIntroduce
		stComic.Author = strAuthor
		stComic.Img = strSrc
		stComic.Type = "热血"
		stComic.Time = strTime
		stComic.ComicChapter = lstComicChapter
		mapSql[topic.Title] = stComic

		ConnectDatabase(mapSql)

	}
	fmt.Println("End =", time.Now().Format("2006-01-02 15:04:05"))
}
