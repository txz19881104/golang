package main

import (
	"cnblogs/cors"
	"cnblogs/models"
	_ "cnblogs/routers"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/Tang-RoseChild/mahonia"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var goroutine_cnt = make(chan int, 100) /*最大协程数量*/

func main() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:   []string{"Content-Length", "Access-Control-Allow-Origin"},
	}))

	go GetFootballAnalyse(5 * 24 * 60 * 60)

	go GetBasketBallLive(10)

	go GetFootBallLive(10)

	beego.Run()
}

func GetBasketBallLive(nTime int) {
	timer(GetBasketBallLiveInfo, nTime)
}

func GetFootBallLive(nTime int) {
	timer(GetFootBallLiveInfo, nTime)
}

func GetFootballAnalyse(nTime int) {
	GetFootballAnalyseInfo()
	timer(GetFootballAnalyseInfo, nTime)
}

func timer(timer func(), nTime int) {
	ticker := time.NewTicker(time.Duration(nTime) * time.Second)
	for {
		select {
		case <-ticker.C:
			timer()
		}
	}
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
		strDataTime := time.Unix(i, 0).Format("20060102") //设置时间戳 使用模板格式化为日期字符串
		//strDataTime = strings.Replace(strDataTime, "-", "", -1)
		dateArray = append(dateArray, strDataTime)
	}

	return dateArray
}

func GetFootballAnalyseInfo() {
	var (
		strSql string
		nNum   int64
		err    error
	)
	o := orm.NewOrm()
	var stMatchInfo []models.FootballMatchInfo
	arrNum := getNumStr(0, 90, 15)

	dayarray := getdateArray("2017-01-01", "2018-12-31")
	for _, value := range dayarray {
		stMatchInfo = []models.FootballMatchInfo{}
		strSql = "SELECT * FROM Bet365.leisu_football_data where match_time like \"" + value + "\";"
		fmt.Println(strSql)

		nNum, err = o.Raw(strSql).QueryRows(&stMatchInfo)
		if err != nil {
			fmt.Println(err, nNum)
			return
		}

		for i := 0; i < len(stMatchInfo); i++ {
			nHTScore := stMatchInfo[i].HTTotalScore - stMatchInfo[i].HTHalfScore
			nVTScore := stMatchInfo[i].VTTotalScore - stMatchInfo[i].VTHalfScore

			_, ok := models.MapFootballGoalStatics[stMatchInfo[i].Name]
			if !ok {
				stGoalStatics := models.FootballGoalStatics{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, map[string]models.FootballGoalNumByTime{}, map[string]models.FootballGoalEffect{}}
				stGoalStatics.MapUpEffectDown = map[string]models.FootballGoalEffect{"0-0": models.FootballGoalEffect{0, 0}, "0-1": models.FootballGoalEffect{0, 0},
					"0-2": models.FootballGoalEffect{0, 0}, "0-3": models.FootballGoalEffect{0, 0}, "1-1": models.FootballGoalEffect{0, 0},
					"1-2": models.FootballGoalEffect{0, 0}, "1-3": models.FootballGoalEffect{0, 0}, "2-2": models.FootballGoalEffect{0, 0},
					"2-3": models.FootballGoalEffect{0, 0}, "3-3": models.FootballGoalEffect{0, 0}}

				for j := 0; j < len(arrNum)-1; j++ {
					strKey := fmt.Sprintf("%d-%d", arrNum[j], arrNum[j+1])
					stGoalStatics.MapFootballGoal[strKey] = models.FootballGoalNumByTime{0, 0, false}
				}

				models.MapFootballGoalStatics[stMatchInfo[i].Name] = stGoalStatics
			}
			stGoalStatics := models.MapFootballGoalStatics[stMatchInfo[i].Name]

			/*比赛总数*/
			stGoalStatics.MatchTotal++

			/*上半场进球总数*/
			stGoalStatics.UpGoalCount = stGoalStatics.UpGoalCount + stMatchInfo[i].HTHalfScore + stMatchInfo[i].VTHalfScore

			/*下半场进球总数*/
			stGoalStatics.DownGoalCount = stGoalStatics.DownGoalCount + nHTScore + nVTScore

			/*进球总数*/
			stGoalStatics.GoalTotal = stGoalStatics.GoalTotal + stMatchInfo[i].HTTotalScore + stMatchInfo[i].VTTotalScore

			/*上半场有球比赛总数*/
			if stMatchInfo[i].HTHalfScore+stMatchInfo[i].VTHalfScore > 0 {
				stGoalStatics.UpHaveGoalMatch++
			}

			/*下半场有球比赛总数*/
			if nHTScore+nVTScore > 0 {
				stGoalStatics.DownHaveGoalMatch++
			}

			/*全场有球比赛总数*/
			if stMatchInfo[i].HTTotalScore+stMatchInfo[i].VTTotalScore > 0 {
				stGoalStatics.AllHaveGoalMatch++
			}

			if stMatchInfo[i].HTTotalCorner > 0 && stMatchInfo[i].VTTotalCorner > 0 {
				/*可统计有角球比赛总数*/
				stGoalStatics.HaveCornerMatch++
				/*角球总数*/
				stGoalStatics.CornerTotal = stGoalStatics.CornerTotal + stMatchInfo[i].HTTotalCorner + stMatchInfo[i].VTTotalCorner
			}

			if stMatchInfo[i].HTShoot > 0 && stMatchInfo[i].VTShoot > 0 {
				/*射门有效比赛数*/
				stGoalStatics.TotalShootMatch++
				/*射门总数*/
				stGoalStatics.ShootCount = stGoalStatics.ShootCount + stMatchInfo[i].HTShoot + stMatchInfo[i].VTShoot
				/*射正总数*/
				stGoalStatics.ShootOnCount = stGoalStatics.ShootOnCount + stMatchInfo[i].HTShooton + stMatchInfo[i].VTShooton
			}

			if len(stMatchInfo[i].GoalTime) > 0 && (stMatchInfo[i].HTTotalScore+stMatchInfo[i].VTTotalScore) != 0 {
				/*可统计进球时间比赛总数*/
				stGoalStatics.TotalGoalTimeMatch++

				arrGoalTime := strings.Split(stMatchInfo[i].GoalTime[1:], " ")
				for j := 0; j < len(arrGoalTime); j++ {

					strTime := arrGoalTime[j][1:]
					nTime, _ := strconv.Atoi(strTime)

					/*可统计进球时间比赛总数*/
					stGoalStatics.TotalGoalTime++

					for k := 0; k < len(arrNum)-1; k++ {
						if nTime > arrNum[k] && nTime <= arrNum[k+1] {
							strKey := fmt.Sprintf("%d-%d", arrNum[k], arrNum[k+1])
							_, ok = stGoalStatics.MapFootballGoal[strKey]

							stFootballGoalNumByTime := stGoalStatics.MapFootballGoal[strKey]
							stFootballGoalNumByTime.GoalNum++
							if !stFootballGoalNumByTime.HaveGoal {
								stFootballGoalNumByTime.GoalMatchNum++
								stFootballGoalNumByTime.HaveGoal = true
							}
							stGoalStatics.MapFootballGoal[strKey] = stFootballGoalNumByTime
						}
					}
				}

				for key, _ := range stGoalStatics.MapFootballGoal {
					stFootballGoalNumByTime := stGoalStatics.MapFootballGoal[key]
					stFootballGoalNumByTime.HaveGoal = false
					stGoalStatics.MapFootballGoal[key] = stFootballGoalNumByTime
				}

			} else if (stMatchInfo[i].HTTotalScore + stMatchInfo[i].VTTotalScore) == 0 {
				stGoalStatics.TotalGoalTimeMatch++
			}

			/*上半场比分对下半场影响*/
			for key, _ := range stGoalStatics.MapUpEffectDown {
				arrUpHV := strings.Split(key, "-")
				nUpScore, _ := strconv.Atoi(arrUpHV[0])
				nDownScore, _ := strconv.Atoi(arrUpHV[1])

				if (stMatchInfo[i].HTHalfScore == int64(nUpScore) && stMatchInfo[i].VTHalfScore == int64(nDownScore)) ||
					(stMatchInfo[i].HTHalfScore == int64(nDownScore) && stMatchInfo[i].VTHalfScore == int64(nUpScore)) {
					stMapUpEffectDown := stGoalStatics.MapUpEffectDown[key]
					if nHTScore != 0 || nVTScore != 0 {

						stMapUpEffectDown.GoalHaveGoalMatch++
						stMapUpEffectDown.GoalTotalGoalMatch++
					} else {
						stMapUpEffectDown.GoalTotalGoalMatch++
					}
					stGoalStatics.MapUpEffectDown[key] = stMapUpEffectDown
				}
			}

			models.MapFootballGoalStatics[stMatchInfo[i].Name] = stGoalStatics
		}

	}
	/*
		for key, value := range models.MapFootballGoalStatics {
			fmt.Println(key, value)
		}
	*/
	fmt.Println("Get Football Statics Finish")
}

func getNumStr(nStart int, nEnd int, nSeq int) []int {
	var arrNum []int
	for nStart := 0; nStart <= nEnd; nStart += nSeq {
		//temps := fmt.Sprintf("%d-%d", nStart, nStart+nSeq)
		arrNum = append(arrNum, nStart)
	}
	return arrNum
}

func getMatchNames(matchname map[string]string, strall string) {

	//strall = "美职业^21^1^^41^美国^52!日职乙^284^1^^46^日本^58!巴西乙^358^1^^39^巴西^50!立陶甲^217^1^^29^立陶宛^28!哥伦甲秋^250^1^^63^哥伦比亚^69!威超^135^1^^30^威尔士^29!西杯^81^1^^3^西班牙^3!委內超秋^391^1^^82^委內瑞拉^87!哥斯甲春^504^1^^85^哥斯达黎加^93!巴西杯^186^1^^39^巴西^50!印度超^1367^1^^101^印度^90!韩足总^468^1^^47^南韩^59!印尼超^1122^1^^109^印度尼西亚^89!巴圣杯^1358^1^^39^巴西^50!国际友谊^1366^1^^52^国际赛事^191!意丙1A^142^0^^2^意大利^2!德地区北^1416^0^^4^德国^4!欧女杯^534^0^^53^欧洲赛事^192!女欧U17^526^0^^53^欧洲赛事^192!斯亚杯^208^0^^35^斯洛文尼亚^34!塞尔杯^211^0^^37^塞尔维亚^36!女瑞典杯^464^0^^10^瑞典^10!土乙白^1431^0^^32^土耳其^31!萨尔超春^511^0^^86^萨尔瓦多^95!格鲁甲^563^0^^96^格鲁吉亚^46!洪都甲春^577^0^^83^洪都拉斯^88!沙地甲^745^0^^57^沙地阿拉伯^64!俄乙中^1433^0^^18^俄罗斯^18!智利乙^611^0^^42^智利^53!港后备^805^0^^45^中国^56!球会友谊^41^0^^52^国际赛事^191!斯伐丙^431^0^^24^斯洛伐克^24!土丙A^1432^0^^32^土耳其^31!西室内足^587^0^^3^西班牙^3!瑞典女甲^1430^0^^10^瑞典^10!匈U19^893^0^^31^匈牙利^30!立陶乙^897^0^^29^立陶宛^28!马尔女甲^994^0^^93^马耳他^43!摩尔乙^1014^0^^90^摩尔多瓦^40!白俄女超^1086^0^^28^白俄罗斯^27!罗乙^1423^0^^23^罗马尼亚^23!世女美预^1121^0^^54^美洲赛事^194!乌兹甲^1160^0^^103^乌兹别克^94!捷女甲^1191^0^^21^捷克^21!印果阿超^1411^0^^101^印度^90!希丙南^1424^0^^22^希腊^22!美超^1457^0^^41^美国^52!阿丙曼特^1461^0^^38^阿根廷^49!科特超^1479^0^^124^科特迪瓦^124!巴马甲^1494^0^^151^巴拿马^151!委內乙^1501^0^^82^委內瑞拉^87!巴拉后备^1531^0^^66^巴拉圭^72!新加国联^1544^0^^48^新加坡^60!德堡州联^1558^0^^4^德国^4!巴马后备^1568^0^^54^美洲赛事^194!尼拉甲^1586^0^^54^美洲赛事^194!塞尔女甲^1604^0^^37^塞尔维亚^36!北爱后备^1606^0^^17^北爱尔兰^17!俄丙^1607^0^^18^俄罗斯^18!爱尔高联^1620^0^^16^爱尔兰^16!阿后备^1634^0^^38^阿根廷^49!乌拉后备^1635^0^^40^乌拉圭^51!斯女联^1640^0^^35^斯洛文尼亚^34!意丁^1650^0^^2^意大利^2!印米佐超^1682^0^^101^印度^90!罗女超^1686^0^^23^罗马尼亚^23!塞尔联^1689^0^^37^塞尔维亚^36!科索沃超^1750^0^^148^科索沃^149!巴地区^1760^0^^39^巴西^50!巴克甲^1761^0^^39^巴西^50!玻利乙^1775^0^^128^斯里兰卡^128!洪都拉杯^1778^0^^83^洪都拉斯^88!塞浦女甲^1826^0^^36^塞浦路斯^35!希业余杯^1829^0^^22^希腊^22!塔吉克杯^1855^0^^136^塔吉克斯坦^137!尼加青联^1885^0^^54^美洲赛事^194!爱沙丁^1894^0^^89^爱沙尼亚^39!哥伦U20^1897^0^^63^哥伦比亚^69!德地杯^1903^0^^4^德国^4!危地杯^1954^0^^84^危地马拉^91!菲律宾杯^1966^0^^150^菲律宾^150!巴圣青联^1389^0^^39^巴西^50!波青联^500^0^^19^波兰^19"
	tempstrarr := strings.Split(strall, "!")
	for _, str := range tempstrarr {
		temp1 := strings.Split(str, "^")
		if len(temp1) > 1 {
			matchname[temp1[1]] = temp1[0]
		}
	}

}

func gettimeminus(strStartTime string, strEndTime string) int {
	if strEndTime == "" {
		return 0
	}
	//获取本地location
	strTimeLayout := "20060102150405"                                    //转化所需模板
	loc, _ := time.LoadLocation("Local")                                 //重要：获取时区
	theTime, _ := time.ParseInLocation(strTimeLayout, strStartTime, loc) //使用模板在对应时区转化为time.time类型
	nStartTime := theTime.Unix()
	theTime, _ = time.ParseInLocation(strTimeLayout, strEndTime, loc) //使用模板在对应时区转化为time.time类型
	nEndTime := theTime.Unix()

	nMinusTime := nEndTime - nStartTime

	return time.Unix(nMinusTime, 0).Hour()*60 + time.Unix(nMinusTime, 0).Minute()
}

func getMatchDetail(mapDetail map[string]string, str []string) {
	/*
		[ 0:比赛id 1：联赛 2：比赛状态 3：开始时间 4：上半场为实际开始时间，下半场为下半场开始时间 5：主队 6：客队  7：主队进球 8：客队进球 9：主队半场进球 10：客队半场进球
		11：主队红卡数 12：客队红卡数 13：主队黄卡 14：客队黄卡 15：亚指初盘     16 17 18      19：主队排名 20：客队排名 21：底行备注    22 23 24    25：主队角球 26：客队角球
		27 28     29：主队半场角球 30：客队半场角球 31：大小初盘 32：角球初盘 33：比赛国家 共计33 ]
	*/

	mapDetail["比赛id"] = str[0]
	mapDetail["联赛"] = str[1]
	mapDetail["状态"] = str[2]
	mapDetail["时间1"] = str[3]
	mapDetail["时间2"] = str[4]
	mapDetail["主队"] = str[5]
	mapDetail["客队"] = str[6]
	mapDetail["主队进球"] = str[7]
	mapDetail["客队进球"] = str[8]
	mapDetail["主队半场进球"] = str[9]
	mapDetail["客队半场进球"] = str[10]
	mapDetail["主队红卡数"] = str[11]
	mapDetail["客队红卡数"] = str[12]
	mapDetail["主队黄卡数"] = str[13]
	mapDetail["客队黄卡数"] = str[14]
	mapDetail["亚指初盘"] = str[15]
	mapDetail["16"] = str[16]
	mapDetail["17"] = str[17]
	mapDetail["18"] = str[18]
	mapDetail["主队排名"] = str[19]
	mapDetail["客队排名"] = str[20]
	mapDetail["底行备注"] = str[21]
	mapDetail["22"] = str[22]
	mapDetail["23"] = str[23]
	mapDetail["24"] = str[24]
	mapDetail["主队角球"] = str[25]
	mapDetail["客队角球"] = str[26]
	mapDetail["27"] = str[27]
	mapDetail["28"] = str[28]
	mapDetail["主队半场角球"] = str[29]
	mapDetail["客队半场角球"] = str[30]
	mapDetail["大小初盘"] = str[31]
	mapDetail["角球初盘"] = str[32]
	mapDetail["比赛国家"] = str[33]

	//fmt.Println(mapDetail)
}

func getMatchDetil(matchnamemap map[string]string, strMatch string, arrNowdata *[]models.FootballMatchInfo, arrOverdata *[]models.FootballMatchInfo) {

	/*
		str = "1572799^250^3^20181018070000^20181018080429^拉伊奎达德^托利马体育^0^1^0^1^0^0^0^2^0.25^^^1^1^2^^^0^0^3^4^1^^3^4^2^9^63!" +
			"1572801^250^2^20181018070000^^巴兰基亚青年^帕特里奥坦斯^2^0^2^0^0^0^1^1^1.25^^^1^5^15^^^0^0^3^4^1^^3^4^2.25^9^63!" +
			"1601932^611^3^20181018070000^20181018080043^巴列彻^塞雷那^2^0^2^0^0^0^0^0^0.5^^^0^13^15^^^0^0^1^1^1^^1^0^2.25^9.5^42!" +
			"1629401^1358^3^20181018070000^20181018080238^皮拉西卡巴^伊图诺^0^0^0^0^0^0^0^1^0.25^^^0^4^8^^^0^0^1^0^1^^1^0^2.25^10^39!" +
			"1500472^21^1^20181018073000^20181018073823^奥兰多城^西雅图音速^0^2^^^0^0^0^0^-0.5^^^1^11^5^^^1^0^2^1^1^^2^1^3.25^10^41!" +
			"1633112^1121^1^20181018080000^20181018080000^加拿大女足^美国女足^0^0^^^0^0^0^0^-1.5^^^0^5^1^^^0^0^0^0^0^^0^0^3^^54"
	*/
	var stFootballInfo models.FootballMatchInfo

	arrMatchInfo := strings.Split(strMatch, "!")
	for i := 0; i < len(arrMatchInfo); i++ {
		arrMatchDetail := strings.Split(arrMatchInfo[i], "^")
		if len(arrMatchDetail) < 34 {
			continue
		}

		mapDetail := make(map[string]string)
		getMatchDetail(mapDetail, arrMatchDetail)

		nMinuteMinusES := gettimeminus(mapDetail["时间1"], mapDetail["时间2"])
		nMinuteMinusNS := gettimeminus(mapDetail["时间1"], time.Now().Format("20060102150405"))
		tempn := nMinuteMinusNS - nMinuteMinusES
		temps := ""

		strMatchSta := mapDetail["状态"]
		if !((strMatchSta == "1") || (strMatchSta == "2") || (strMatchSta == "3") || (strMatchSta == "-1")) {
			continue
		}

		switch strMatchSta {
		case "1":
			temps = fmt.Sprintf("%d'", tempn)
			stFootballInfo.NowMatchTime = tempn
		case "2":
			temps = "中"
			stFootballInfo.NowMatchTime = 65
		case "3":
			tempn += 45

			temps = fmt.Sprintf("%d'", tempn)
			if mapDetail["主队半场进球"] == "" {
				stFootballInfo.HTHalfScore = 0
			} else {
				stFootballInfo.HTHalfScore, _ = strconv.ParseInt(mapDetail["主队半场进球"], 10, 64)
			}

			if mapDetail["客队半场进球"] == "" {
				stFootballInfo.VTHalfScore = 0
			} else {
				stFootballInfo.VTHalfScore, _ = strconv.ParseInt(mapDetail["客队半场进球"], 10, 64)
			}
			stFootballInfo.NowMatchTime = tempn + 20
		case "-1":
			temps = "完"
		default:
		}

		stFootballInfo.EventSta = temps //比赛状态

		stFootballInfo.MatchID = mapDetail["比赛id"]          //比赛ID
		stFootballInfo.Name = matchnamemap[mapDetail["联赛"]] //联赛名称

		if mapDetail["主队进球"] == "" {
			stFootballInfo.HTTotalScore = 0
		} else {
			stFootballInfo.HTTotalScore, _ = strconv.ParseInt(mapDetail["主队进球"], 10, 64)
		}

		if mapDetail["客队进球"] == "" {
			stFootballInfo.VTTotalScore = 0
		} else {
			stFootballInfo.VTTotalScore, _ = strconv.ParseInt(mapDetail["客队进球"], 10, 64)
		}

		if temps == "中" || temps == "完" || strMatchSta == "1" {
			stFootballInfo.HTHalfScore = stFootballInfo.HTTotalScore
			stFootballInfo.VTHalfScore = stFootballInfo.VTTotalScore
		}
		//stFootballInfo.Eventdate = mapDetail["时间1"][6:8] + "日 "
		//stFootballInfo.Eventtime = fmt.Sprintf("%s:%s ", mapDetail["时间1"][8:10], mapDetail["时间1"][10:12])
		//stFootballInfo.Sortdate = mapDetail["时间1"]

		stFootballInfo.HTName = mapDetail["主队"]
		stFootballInfo.VTName = mapDetail["客队"]

		if mapDetail["主队红卡数"] == "" {
			stFootballInfo.HTRed = 0
		} else {
			stFootballInfo.HTRed, _ = strconv.ParseInt(mapDetail["主队红卡数"], 10, 64)
		}

		if mapDetail["客队红卡数"] == "" {
			stFootballInfo.VTRed = 0
		} else {
			stFootballInfo.VTRed, _ = strconv.ParseInt(mapDetail["客队红卡数"], 10, 64)
		}

		stFootballInfo.Asionodd = mapDetail["亚指初盘"]
		stFootballInfo.Cornerodd = mapDetail["角球初盘"]
		stFootballInfo.Numodd = mapDetail["大小初盘"]

		if strMatchSta == "-1" {
			*arrOverdata = append(*arrOverdata, stFootballInfo)
		} else {
			*arrNowdata = append(*arrNowdata, stFootballInfo)
		}

	}
}

func getShoot(item *models.FootballMatchInfo) string {

	strHref := "http://www.win0168.com/detail/" + item.MatchID + "sb.htm"
	item.DetailHref = strHref

	parResultHref, _ := url.Parse(strHref)

	parResultHref.RawQuery = parResultHref.Query().Encode()
	httpSearchRespHref, err := http.Get(parResultHref.String())

	if httpSearchRespHref == nil || err != nil || (httpSearchRespHref != nil && httpSearchRespHref.StatusCode != 200) || httpSearchRespHref.Body == nil {

		/*i := 0
		for ; i < 2; i++ {
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
			return getErrorLine()
		}*/
		fmt.Println(strHref, "download failed")
		<-goroutine_cnt
		wg.Done()
		return getErrorLine()
	}

	dec := mahonia.NewDecoder("utf-8")
	rd := dec.NewReader(httpSearchRespHref.Body)
	docHref, err := goquery.NewDocumentFromReader(rd)

	if docHref == nil || err != nil {
		fmt.Println(strHref, "download failed")
		<-goroutine_cnt
		wg.Done()
		return getErrorLine()
	}
	defer httpSearchRespHref.Body.Close()

	docHref.Find("#matchData > div").Each(func(i int, s *goquery.Selection) {

		pos := strings.Contains((s.Find("table > tbody > tr:nth-child(1)").Text()), "本场技术统计")
		if pos != false {
			s.Find("table > tbody > tr").Each(func(i int, t *goquery.Selection) {
				pos := strings.Contains(t.Text(), "射门")
				if pos != false {
					if (t.Find("td.bg3").Text() == "射门") || (t.Find("td.bg4").Text() == "射门") {
						temp := strings.Split(t.Text(), "射门")
						if len(temp) > 1 {
							item.HTShoot, _ = strconv.ParseInt(temp[0], 10, 64)
							item.VTShoot, _ = strconv.ParseInt(temp[1], 10, 64)
						}
					}
				}
				pos = strings.Contains(t.Text(), "射正")
				if pos != false {
					temp := strings.Split(t.Text(), "射正")
					if len(temp) > 1 {
						item.HTShooton, _ = strconv.ParseInt(temp[0], 10, 64)
						item.VTShooton, _ = strconv.ParseInt(temp[1], 10, 64)
					}
				}

				pos = strings.Contains(t.Text(), "角球")
				if pos != false {
					if (t.Find("td.bg3").Text() == "角球") || (t.Find("td.bg4").Text() == "角球") {
						temp := strings.Split(t.Text(), "角球")
						if len(temp) > 1 {
							item.HTTotalCorner, _ = strconv.ParseInt(temp[0], 10, 64)
							item.VTTotalCorner, _ = strconv.ParseInt(temp[1], 10, 64)
						}
					}
				}
				pos = strings.Contains(t.Text(), "半场角球")
				if pos != false {
					temp := strings.Split(t.Text(), "半场角球")
					if len(temp) > 1 {
						item.HTHalfCorner, _ = strconv.ParseInt(temp[0], 10, 64)
						item.VTHalfCorner, _ = strconv.ParseInt(temp[1], 10, 64)
					}
				}

			})
		}

		pos = strings.Contains((s.Find("table > tbody > tr:nth-child(1)").Text()), "详细事件")
		if pos != false {

			s.Find("table > tbody > tr").Each(func(i int, t *goquery.Selection) {
				temptitle, ok := t.Find("td:nth-child(2) > img").Attr("title")
				if ok == true {
					//if strings.Contains(temptitle, "球") != false {
					if (temptitle == "入球") || (temptitle == "点球") || (temptitle == "乌龙") {
						//fmt.Println("H:",temptitle,t.Find("td:nth-child(3)").Text())
						item.GoalTime += " H" + t.Find("td:nth-child(3)").Text()
						item.HGoalTime += " H" + t.Find("td:nth-child(3)").Text()
					}
				}
				temptitle, ok = t.Find("td:nth-child(4) > img").Attr("title")
				if ok == true {
					//if strings.Contains(temptitle, "球") != false {
					if (temptitle == "入球") || (temptitle == "点球") || (temptitle == "乌龙") {
						//fmt.Println("G:",temptitle,t.Find("td:nth-child(3)").Text())
						item.GoalTime += " V" + t.Find("td:nth-child(3)").Text()
						item.VGoalTime += " V" + t.Find("td:nth-child(3)").Text()
					}
				}
			})
			//	fmt.Println("")
		}
	})
	//printItem(item)

	<-goroutine_cnt
	wg.Done()

	return ""
}

func getErrorLine() string {
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc).Name()
	return "Error: " + file + " -> " + f + " ->line:" + strconv.Itoa(line)
}

func getShootDetail(arrNowDetail *[]models.FootballMatchInfo, arrOverDetail *[]models.FootballMatchInfo) {
	timeUse := time.Now()

	wg.Add(len(*arrNowDetail) + len(*arrOverDetail))
	for i := 0; i < len(*arrNowDetail); i++ {
		goroutine_cnt <- 1
		go getShoot(&(*arrNowDetail)[i])
	}
	for i := 0; i < len(*arrOverDetail); i++ {
		goroutine_cnt <- 1
		go getShoot(&(*arrOverDetail)[i])
	}
	wg.Wait()

	strTemp := fmt.Sprintf("->-> 现场:%-4d 完场新增:%-4d  用时:", len(*arrNowDetail), len(*arrOverDetail))
	fmt.Println(strTemp, time.Since(timeUse))
}

func getMatchs(strHref string) (string, error) {

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

			return "", err
		}
	}

	dec := mahonia.NewDecoder("utf-8")
	rd := dec.NewReader(httpSearchRespHref.Body)
	docHref, err := goquery.NewDocumentFromReader(rd)

	if docHref == nil || err != nil {
		fmt.Println(strHref, "download failed")

		return "", err
	}
	defer httpSearchRespHref.Body.Close()

	return docHref.Text(), nil
}

/*当前足球比赛抓取*/
func GetFootBallLiveInfo() {
	strMatchTxt, err := getMatchs("http://m.win007.com/phone/Schedule_0_0.txt")

	if err != nil {
		return
	}

	strTempArry := strings.Split(strMatchTxt, "$$")
	if len(strTempArry) < 2 {
		return
	}

	/*
		timeNowHour := time.Now().Hour()
		timeNowMinute := time.Now().Minute()
		timeNow := timeNowHour*60 + timeNowMinute
		if (timeNow > 13*60) && (timeNow < 13*60+10) { //每天下午13:00~13:10清空完场数组 重新累计
			g_Overmatchdata = nil
		}
	*/

	mapMatchName := make(map[string]string)
	getMatchNames(mapMatchName, strTempArry[0])

	var arrNowDetail []models.FootballMatchInfo
	var arrOverDetail []models.FootballMatchInfo
	getMatchDetil(mapMatchName, strTempArry[1], &arrNowDetail, &arrOverDetail)

	getShootDetail(&arrNowDetail, &arrOverDetail)

	models.MapFootballMatchInfo = make(map[string][]models.FootballMatchInfo)
	for i := 0; i < len(arrNowDetail); i++ {
		/*查看元素在集合中是否存在 */
		_, ok := models.MapFootballMatchInfo[arrNowDetail[i].Name] /*如果确定是真实的,则存在,否则不存在 */
		if ok {
			models.MapFootballMatchInfo[arrNowDetail[i].Name] = append(models.MapFootballMatchInfo[arrNowDetail[i].Name], arrNowDetail[i])
		} else {
			ArrMatchInfo := []models.FootballMatchInfo{}
			ArrMatchInfo = append(ArrMatchInfo, arrNowDetail[i])
			models.MapFootballMatchInfo[arrNowDetail[i].Name] = ArrMatchInfo
		}
	}
	//getMatchColors(&arrNowDetail)

}

/*当前篮球比赛抓取*/
func GetBasketBallLiveInfo() {
	strHref := "https://live.leisu.com/lanqiu"
	var stMatchInfo models.BasketballMatchInfo

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
		models.MapBasketballMatchInfo = make(map[string][]models.BasketballMatchInfo)
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

		//models.ArrMatchInfo = append(models.ArrMatchInfo, stMatchInfo)

		/*查看元素在集合中是否存在 */
		_, ok := models.MapBasketballMatchInfo[stMatchInfo.Name] /*如果确定是真实的,则存在,否则不存在 */
		if ok {
			models.MapBasketballMatchInfo[stMatchInfo.Name] = append(models.MapBasketballMatchInfo[stMatchInfo.Name], stMatchInfo)
		} else {
			ArrMatchInfo := []models.BasketballMatchInfo{}
			ArrMatchInfo = append(ArrMatchInfo, stMatchInfo)
			models.MapBasketballMatchInfo[stMatchInfo.Name] = ArrMatchInfo
		}

	})

}
