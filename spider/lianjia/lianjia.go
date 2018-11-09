package main

import (
	"container/list"
	"database/sql"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*下载每天数据的数据库结构*/
type HouseInfo struct {
	HouseRecord           string `orm:"column(house_record)"`
	CommunityName         string `orm:"column(community_name)"`
	AreaName              string `orm:"column(area_name)"`
	AllPrice              string `orm:"column(all_price)"`
	UnitPrice             string `orm:"column(unit_price)"`
	HouseType             string `orm:"column(house_type)"`
	Floor                 string `orm:"column(floor)"`
	BuiltArea             string `orm:"column(built_area)"`
	FamilyStructure       string `orm:"column(family_structure)"`
	InnerArea             string `orm:"column(inner_area)"`
	ArchitecturalType     string `orm:"column(architectural_type)"`
	HouseOrientation      string `orm:"column(house_orientation)"`
	BuildingStructure     string `orm:"column(building_structure)"`
	DecorationSituation   string `orm:"column(decoration_situation)"`
	LadderProportion      string `orm:"column(ladder_proportion)"`
	HeatingMode           string `orm:"column(heating_mode)"`
	EquippedWithElevators string `orm:"column(equipped_with_elevators)"`
	YearsOfPropertyRights string `orm:"column(years_of_property_rights)"`
	ListingTime           string `orm:"column(listing_time)"`
	TradingRights         string `orm:"column(trading_rights)"`
	LastTransaction       string `orm:"column(last_transaction)"`
	HousingUse            string `orm:"column(housing_use)"`
	HousingYears          string `orm:"column(housing_years)"`
	PropertyRightsBelong  string `orm:"column(property_rights_belong)"`
	MortgageInformation   string `orm:"column(mortgage_information)"`
	SpareParts            string `orm:"column(spare_parts)"`
	AddTime               string `orm:"column(add_time)"`
	WebHref               string `orm:"column(web_href)"`
}

/*下载历史成交数据的数据库结构*/
type HouseInfoHistory struct {
	HouseRecord           string `orm:"column(house_record)"`
	CommunityName         string `orm:"column(community_name)"`
	DealTotalPrice        string `orm:"column(deal_total_price)"`
	DealUnitPrice         string `orm:"column(deal_unit_price)"`
	DealInfo              string `orm:"column(deal_info)"`
	DealTime              string `orm:"column(deal_time)"`
	HouseType             string `orm:"column(house_type)"`
	Floor                 string `orm:"column(floor)"`
	BuiltArea             string `orm:"column(built_area)"`
	FamilyStructure       string `orm:"column(family_structure)"`
	InnerArea             string `orm:"column(inner_area)"`
	ArchitecturalType     string `orm:"column(architectural_type)"`
	HouseOrientation      string `orm:"column(house_orientation)"`
	BuildingStructure     string `orm:"column(building_structure)"`
	DecorationSituation   string `orm:"column(decoration_situation)"`
	LadderProportion      string `orm:"column(ladder_proportion)"`
	HeatingMode           string `orm:"column(heating_mode)"`
	EquippedWithElevators string `orm:"column(equipped_with_elevators)"`
	YearsOfPropertyRights string `orm:"column(years_of_property_rights)"`
	ListingTime           string `orm:"column(listing_time)"`
	TradingRights         string `orm:"column(trading_rights)"`
	LastTransaction       string `orm:"column(last_transaction)"`
	HousingUse            string `orm:"column(housing_use)"`
	HousingYears          string `orm:"column(housing_years)"`
	PropertyRightsBelong  string `orm:"column(property_rights_belong)"`
	MortgageInformation   string `orm:"column(mortgage_information)"`
	SpareParts            string `orm:"column(spare_parts)"`
	AddTime               string `orm:"column(add_time)"`
	WebHref               string `orm:"column(web_href)"`
}

var wg sync.WaitGroup
var goroutine_cnt = make(chan int, 50) /*最大协程数量*/
var db *sql.DB

const (
	PAGE int = 100
)

/*得到此房源的详细信息*/
func GetDetailInfo(strHref string) {
	var stHouseInfo HouseInfo

	stHouseInfo.WebHref = strHref
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
	if docHref == nil {
		<-goroutine_cnt
		wg.Done()
		return
	}

	/*获取房源的所有介绍，根据不同类型进行存储*/
	lstContentHref := docHref.Find("div.introContent")
	lstContentHref.Each(func(i int, contentSelectionHref *goquery.Selection) {
		if i == 0 {

			lstLabel := contentSelectionHref.Find("span.label").Parent()

			lstLabel.Each(func(i int, contentSelectionLabel *goquery.Selection) {
				strLabel := contentSelectionLabel.Text()
				strLabel = strings.Replace(strLabel, " ", "", -1)
				strLabel = strings.Replace(strLabel, "\n", "", -1)
				strLabel = strings.Replace(strLabel, "\r", "", -1)
				strName := strLabel[0:12]
				strLabel = strLabel[12:len(strLabel)]

				switch strName {
				case "房屋户型":
					stHouseInfo.HouseType = strLabel
				case "所在楼层":
					stHouseInfo.Floor = strLabel
				case "建筑面积":
					stHouseInfo.BuiltArea = strLabel
				case "户型结构":
					stHouseInfo.FamilyStructure = strLabel
				case "套内面积":
					stHouseInfo.InnerArea = strLabel
				case "建筑类型":
					stHouseInfo.ArchitecturalType = strLabel
				case "房屋朝向":
					stHouseInfo.HouseOrientation = strLabel
				case "建筑结构":
					stHouseInfo.BuildingStructure = strLabel
				case "装修情况":
					stHouseInfo.DecorationSituation = strLabel
				case "梯户比例":
					stHouseInfo.LadderProportion = strLabel

				case "供暖方式":
					stHouseInfo.HeatingMode = strLabel
				case "配备电梯":
					stHouseInfo.EquippedWithElevators = strLabel
				case "产权年限":
					stHouseInfo.YearsOfPropertyRights = strLabel
				case "挂牌时间":
					stHouseInfo.ListingTime = strLabel
				case "交易权属":
					stHouseInfo.TradingRights = strLabel
				case "上次交易":
					stHouseInfo.LastTransaction = strLabel
				case "房屋用途":
					stHouseInfo.HousingUse = strLabel
				case "房屋年限":
					stHouseInfo.HousingYears = strLabel
				case "产权所属":
					stHouseInfo.PropertyRightsBelong = strLabel
				case "抵押信息":
					stHouseInfo.MortgageInformation = strLabel
				case "房本备件":
					stHouseInfo.SpareParts = strLabel
				default:
					//fmt.Println(strName, "can not find!")
				}

			})

		}

		//fmt.Println(contentSelectionHref.Find("span.label").Parent().Text())

	})

	/*获取总价*/
	strPrice := docHref.Find("div.price ").Children().Eq(0).Text()
	stHouseInfo.AllPrice = strPrice + "万"

	/*获取单价*/
	strUnitPriceValue := docHref.Find("span.unitPriceValue").Text()
	stHouseInfo.UnitPrice = strUnitPriceValue

	/*获取小区名字*/
	strCommunityName := docHref.Find("div.communityName").Children().Eq(2).Text()
	strCommunityName = strings.Replace(strCommunityName, " ", "", -1)
	strCommunityName = strings.Replace(strCommunityName, "\n", "", -1)
	strCommunityName = strings.Replace(strCommunityName, "\r", "", -1)
	stHouseInfo.CommunityName = strCommunityName

	/*获取所属区域*/
	strAreaName := docHref.Find("div.areaName").Text()
	strAreaName = strings.Replace(strAreaName, " ", "-", -1)
	strAreaName = strings.Replace(strAreaName, "\n", "", -1)
	strAreaName = strings.Replace(strAreaName, "\r", "", -1)
	stHouseInfo.AreaName = strAreaName

	/*获取所属区域*/
	strHouseRecord := docHref.Find("div.houseRecord").Children().Eq(1).Text()
	strHouseRecord = strings.Replace(strHouseRecord, "举报", "", -1)
	stHouseInfo.HouseRecord = strHouseRecord

	/*数据采集日期*/
	stHouseInfo.AddTime = time.Now().Format("2006-01-02")

	/*日期为空或者非日期时用2000-01-01代替*/
	if stHouseInfo.ListingTime == "" || stHouseInfo.ListingTime[0:1] != "2" {
		stHouseInfo.ListingTime = "2000-01-01"
	}

	if stHouseInfo.LastTransaction == "" || stHouseInfo.ListingTime[0:1] != "2" {
		stHouseInfo.LastTransaction = "2000-01-01"
	}

	/*插入数据库*/
	strSql := `INSERT into LianjiaDB.tbl_lj_house (house_record, community_name, area_name, all_price, unit_price, house_type, floor, built_area, family_structure, 
	 		inner_area, architectural_type, house_orientation, building_structure, decoration_situation, ladder_proportion, heating_mode, equipped_with_elevators,
	 		years_of_property_rights, listing_time, trading_rights, last_transaction, housing_use, housing_years, property_rights_belong, mortgage_information, spare_parts, add_time, web_href) VALUES (` +
		"\"" + stHouseInfo.HouseRecord + "\"," +
		"\"" + stHouseInfo.CommunityName + "\"," +
		"\"" + stHouseInfo.AreaName + "\"," +
		"\"" + stHouseInfo.AllPrice + "\"," +
		"\"" + stHouseInfo.UnitPrice + "\"," +
		"\"" + stHouseInfo.HouseType + "\"," +
		"\"" + stHouseInfo.Floor + "\"," +
		"\"" + stHouseInfo.BuiltArea + "\"," +
		"\"" + stHouseInfo.FamilyStructure + "\"," +
		"\"" + stHouseInfo.InnerArea + "\"," +
		"\"" + stHouseInfo.ArchitecturalType + "\"," +
		"\"" + stHouseInfo.HouseOrientation + "\"," +
		"\"" + stHouseInfo.BuildingStructure + "\"," +
		"\"" + stHouseInfo.DecorationSituation + "\"," +
		"\"" + stHouseInfo.LadderProportion + "\"," +
		"\"" + stHouseInfo.HeatingMode + "\"," +
		"\"" + stHouseInfo.EquippedWithElevators + "\"," +
		"\"" + stHouseInfo.YearsOfPropertyRights + "\"," +
		"\"" + stHouseInfo.ListingTime + "\"," +
		"\"" + stHouseInfo.TradingRights + "\"," +
		"\"" + stHouseInfo.LastTransaction + "\"," +
		"\"" + stHouseInfo.HousingUse + "\"," +
		"\"" + stHouseInfo.HousingYears + "\"," +
		"\"" + stHouseInfo.PropertyRightsBelong + "\"," +
		"\"" + stHouseInfo.MortgageInformation + "\"," +
		"\"" + stHouseInfo.SpareParts + "\"," +
		"\"" + stHouseInfo.AddTime + "\"," +
		"\"" + stHouseInfo.WebHref + "\");"

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
			<-goroutine_cnt
			wg.Done()
			return
		}
	}
	<-goroutine_cnt

	wg.Done()
}

/*得到每页所有房源列表*/
func GetAllHouse(strAreaUrl string, lstHouseInfo *list.List) {
	parResult, _ := url.Parse(strAreaUrl)

	parResult.RawQuery = parResult.Query().Encode()

	httpSearchResp, err := http.Get(parResult.String())

	if err != nil {

		i := 0
		for ; i < 5; i++ {
			httpSearchResp, err = http.Get(parResult.String())
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

	defer httpSearchResp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(httpSearchResp.Body)

	lstContent := doc.Find("li.clear")
	//wg.Add(lstContent.Length())

	/*得到此页所有房源*/
	lstContent.Each(func(i int, contentSelection *goquery.Selection) {
		strHref, _ := contentSelection.Find("a").Attr("href")
		lstHouseInfo.PushBack(strHref)

		//getDetailInfo(strHref)

		/*
			strInfo := contentSelection.Find("a").Eq(1).Text()
			fmt.Println(strInfo)

			strAddress := contentSelection.Find("div.address").Text()
			fmt.Println(strAddress)

			strFlood := contentSelection.Find("div.flood").Text()
			fmt.Println(strFlood)

			strFollowInfo := contentSelection.Find("div.followInfo").Text()
			fmt.Println(strFollowInfo)

			strTag := contentSelection.Find("div.tag").Text()
			fmt.Println(strTag)

			strPriceInfo := contentSelection.Find("div.priceInfo").Text()
			fmt.Println(strPriceInfo)
		*/

	})
	//wg.Wait()
	<-goroutine_cnt

	wg.Done()
}

/*得到此房源的历史成交详细信息*/
func GetDetailInfoHistory(strHref string) {
	var stHouseInfo HouseInfoHistory

	stHouseInfo.WebHref = strHref
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
	lstContentHref := docHref.Find("div.introContent")
	lstContentHref.Each(func(i int, contentSelectionHref *goquery.Selection) {
		if i == 0 {
			lstLabel := contentSelectionHref.Find("span.label").Parent()

			lstLabel.Each(func(i int, contentSelectionLabel *goquery.Selection) {
				strLabel := contentSelectionLabel.Text()
				strLabel = strings.Replace(strLabel, " ", "", -1)
				strLabel = strings.Replace(strLabel, "\n", "", -1)
				strLabel = strings.Replace(strLabel, "\r", "", -1)
				strName := strLabel[0:12]
				strLabel = strLabel[12:len(strLabel)]

				switch strName {
				case "房屋户型":
					stHouseInfo.HouseType = strLabel
				case "所在楼层":
					stHouseInfo.Floor = strLabel
				case "建筑面积":
					stHouseInfo.BuiltArea = strLabel
				case "户型结构":
					stHouseInfo.FamilyStructure = strLabel
				case "套内面积":
					stHouseInfo.InnerArea = strLabel
				case "建筑类型":
					stHouseInfo.ArchitecturalType = strLabel
				case "房屋朝向":
					stHouseInfo.HouseOrientation = strLabel
				case "建筑结构":
					stHouseInfo.BuildingStructure = strLabel
				case "装修情况":
					stHouseInfo.DecorationSituation = strLabel
				case "梯户比例":
					stHouseInfo.LadderProportion = strLabel

				case "供暖方式":
					stHouseInfo.HeatingMode = strLabel
				case "配备电梯":
					stHouseInfo.EquippedWithElevators = strLabel
				case "产权年限":
					stHouseInfo.YearsOfPropertyRights = strLabel
				case "挂牌时间":
					stHouseInfo.ListingTime = strLabel
				case "交易权属":
					stHouseInfo.TradingRights = strLabel
				case "上次交易":
					stHouseInfo.LastTransaction = strLabel
				case "房屋用途":
					stHouseInfo.HousingUse = strLabel
				case "房屋年限":
					stHouseInfo.HousingYears = strLabel
				case "产权所属":
					stHouseInfo.PropertyRightsBelong = strLabel
				case "抵押信息":
					stHouseInfo.MortgageInformation = strLabel
				case "房本备件":
					stHouseInfo.SpareParts = strLabel
				case "链家编号":
					stHouseInfo.HouseRecord = strLabel
				default:
					//fmt.Println(strName, "can not find!")
				}

			})

		}

		//fmt.Println(contentSelectionHref.Find("span.label").Parent().Text())

	})

	/*交易总价*/
	strDealTotalPrice := docHref.Find("div.price").Children().Eq(0).Text()
	stHouseInfo.DealTotalPrice = strDealTotalPrice

	/*交易单价*/
	strDealUnitPriceValue := docHref.Find("div.price").Children().Eq(1).Text()
	stHouseInfo.DealUnitPrice = strDealUnitPriceValue

	/*小区名字*/
	strCommunityName := docHref.Find("div.house-title").Children().Eq(0).Text()
	strCommunityName = strings.Split(strCommunityName, " ")[0]
	strCommunityName = strings.Replace(strCommunityName, "\n", "", -1)
	strCommunityName = strings.Replace(strCommunityName, "\r", "", -1)
	stHouseInfo.CommunityName = strCommunityName

	/*交易过程信息*/
	strDealInfo := docHref.Find("div.msg").Eq(0).Text()
	stHouseInfo.DealInfo = strDealInfo

	/*交易时间*/
	strDealTime := docHref.Find("div.house-title").Children().Children().Eq(0).Text()
	strDealTime = strings.Split(strDealTime, " ")[0]
	stHouseInfo.DealTime = strDealTime

	if stHouseInfo.ListingTime == "" || stHouseInfo.ListingTime[0:1] != "2" {
		stHouseInfo.ListingTime = "2000-01-01"
	}

	if stHouseInfo.LastTransaction == "" || stHouseInfo.ListingTime[0:1] != "2" {
		stHouseInfo.LastTransaction = "2000-01-01"
	}

	if stHouseInfo.DealTime == "" || stHouseInfo.DealTime[0:1] != "2" {
		stHouseInfo.DealTime = "2000-01-01"
	}

	/*插入数据库*/
	strSql := `INSERT into LianjiaDB.tbl_lj_house_history (house_record, community_name, deal_total_price, deal_unit_price, deal_info, deal_time, house_type, floor, built_area, family_structure,
		 		inner_area, architectural_type, house_orientation, building_structure, decoration_situation, ladder_proportion, heating_mode, equipped_with_elevators,
		 		years_of_property_rights, listing_time, trading_rights, last_transaction, housing_use, housing_years, property_rights_belong, mortgage_information, spare_parts, web_href) VALUES (` +
		"\"" + stHouseInfo.HouseRecord + "\"," +
		"\"" + stHouseInfo.CommunityName + "\"," +
		"\"" + stHouseInfo.DealTotalPrice + "\"," +
		"\"" + stHouseInfo.DealUnitPrice + "\"," +
		"\"" + stHouseInfo.DealInfo + "\"," +
		"\"" + stHouseInfo.DealTime + "\"," +
		"\"" + stHouseInfo.HouseType + "\"," +
		"\"" + stHouseInfo.Floor + "\"," +
		"\"" + stHouseInfo.BuiltArea + "\"," +
		"\"" + stHouseInfo.FamilyStructure + "\"," +
		"\"" + stHouseInfo.InnerArea + "\"," +
		"\"" + stHouseInfo.ArchitecturalType + "\"," +
		"\"" + stHouseInfo.HouseOrientation + "\"," +
		"\"" + stHouseInfo.BuildingStructure + "\"," +
		"\"" + stHouseInfo.DecorationSituation + "\"," +
		"\"" + stHouseInfo.LadderProportion + "\"," +
		"\"" + stHouseInfo.HeatingMode + "\"," +
		"\"" + stHouseInfo.EquippedWithElevators + "\"," +
		"\"" + stHouseInfo.YearsOfPropertyRights + "\"," +
		"\"" + stHouseInfo.ListingTime + "\"," +
		"\"" + stHouseInfo.TradingRights + "\"," +
		"\"" + stHouseInfo.LastTransaction + "\"," +
		"\"" + stHouseInfo.HousingUse + "\"," +
		"\"" + stHouseInfo.HousingYears + "\"," +
		"\"" + stHouseInfo.PropertyRightsBelong + "\"," +
		"\"" + stHouseInfo.MortgageInformation + "\"," +
		"\"" + stHouseInfo.SpareParts + "\"," +
		"\"" + stHouseInfo.WebHref + "\");"

	_, err = db.Exec(strSql)
	if err != nil {

		i := 0
		for ; i < 5; i++ {
			_, err = db.Exec(strSql)
			if err != nil {
				continue
			} else {
				break
			}
		}
		if i == 5 {
			fmt.Println("Error is ", err)
			fmt.Println(strSql)
			<-goroutine_cnt
			wg.Done()
			return
		}
	}

	<-goroutine_cnt

	wg.Done()
}

/*得到每页所有房源列表*/
func GetAllHouseHistory(strAreaUrl string, lstHouseInfo *list.List) {
	parResult, _ := url.Parse(strAreaUrl)

	parResult.RawQuery = parResult.Query().Encode()

	httpSearchResp, err := http.Get(parResult.String())

	if err != nil {

		i := 0
		for ; i < 5; i++ {
			httpSearchResp, err = http.Get(parResult.String())
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

	defer httpSearchResp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(httpSearchResp.Body)

	lstContent := doc.Find("ul.listContent").Children()

	//wg.Add(lstContent.Length())

	lstContent.Each(func(i int, contentSelection *goquery.Selection) {
		strHref, _ := contentSelection.Find("a").Attr("href")
		lstHouseInfo.PushBack(strHref)

	})
	<-goroutine_cnt

	wg.Done()
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

	lstHouseInfo := list.New()

	runtime.GOMAXPROCS(runtime.NumCPU())

	strKey := flag.String("download", "today", "datatype")
	flag.Parse()

	strUrlKeyword := "https://bj.lianjia.com/ershoufang/"

	if *strKey == "history" {
		strUrlKeyword = "https://bj.lianjia.com/chengjiao/"
	}

	/*按照每平米一个区间来获取房源（链家每种搜素最多可以显示3000套房源），尽可能获取所有房源*/
	wg.Add(200 * PAGE)
	for j := 0; j < 200; j++ {
		for i := 1; i <= PAGE; i++ {
			strUrl := strUrlKeyword
			if i > 1 {
				strUrl = strUrl + "pg" + strconv.Itoa(i)
			}

			strArea := "ba" + strconv.Itoa(j) + "ea" + strconv.Itoa((j + 1)) + "ep100000/"
			if *strKey == "history" {
				strArea = "ba" + strconv.Itoa(j) + "ea" + strconv.Itoa((j + 1))
			}

			strAreaUrl := strUrl + strArea

			goroutine_cnt <- 1
			if *strKey == "history" {
				go GetAllHouseHistory(strAreaUrl, lstHouseInfo)
			} else {
				go GetAllHouse(strAreaUrl, lstHouseInfo)
			}
		}

	}
	wg.Wait()

	/*循环每套房源，获取详细信息，并存入数据库*/
	wg.Add(lstHouseInfo.Len())
	for e := lstHouseInfo.Front(); e != nil; e = e.Next() {
		goroutine_cnt <- 1

		if *strKey == "history" {
			go GetDetailInfoHistory(e.Value.(string))
		} else {
			go GetDetailInfo(e.Value.(string))
		}
	}
	wg.Wait()

	fmt.Println("End =", time.Now().Format("2006-01-02 15:04:05"))
}
