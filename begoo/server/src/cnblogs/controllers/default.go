package controllers

import (
	"cnblogs/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	//"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	Success    = 0
	Failed     = -1
	TokenError = "TokenError"
)

type MainController struct {
	beego.Controller
}

type Name struct {
	beego.Controller
}

type NameKeyword struct {
	beego.Controller
}

type DataType struct {
	beego.Controller
}

type Chapter struct {
	beego.Controller
}

type ChapterNum struct {
	beego.Controller
}

type ChapterContent struct {
	beego.Controller
}

type Register struct {
	beego.Controller
}

type Login struct {
	beego.Controller
}

type ChangeSetting struct {
	beego.Controller
}

type UserCookie struct {
	beego.Controller
}

type GetCookie struct {
	beego.Controller
}

type GetLiveBasketBall struct {
	beego.Controller
}

type GetBasketBallAnalyse struct {
	beego.Controller
}

type GetLiveFootBall struct {
	beego.Controller
}

type GetFootBallAnalyse struct {
	beego.Controller
}

type GetFootBallTeamAnalyse struct {
	beego.Controller
}

type CheckUserToken struct {
	beego.Controller
}

var filterUser = func(ctx *context.Context) {

	token := ctx.Input.Header("Token")

	bFind := CheckToken(token)

	//验证Token是否合法
	if !bFind {
		fmt.Println("token check failed")
		ctx.ResponseWriter.Write([]byte(TokenError))
		//http.Error(ctx.ResponseWriter, "Token verification not pass", http.StatusBadRequest)
		return
	}

	fmt.Println("Find token:", bFind)
}

func init() {
	//访问接口前验证token
	//beego.InsertFilter("/Live|Finish/*", beego.BeforeRouter, filterUser)
	beego.InsertFilter("/Api/Sports/*", beego.BeforeRouter, filterUser)
	beego.InsertFilter("/Api/Check/*", beego.BeforeRouter, filterUser)
}

type Token struct {
	Token string `json:"token"`
}

// 校验token是否有效
func CheckToken(token string) (b bool) {

	t, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(beego.AppConfig.String("jwtkey")), nil
	})

	if err != nil {
		fmt.Println("转换为jwt claims失败.", err)
		return false
	}

	conn, err := redis.Dial("tcp", "www.tangxinzhu.com:6379")
	if err != nil {
		fmt.Println("connect redis error :", err)
		return
	}
	conn.Do("AUTH", "passwd")
	defer conn.Close()

	strGettoken, err := redis.String(conn.Do("GET", t.Claims.(jwt.MapClaims)["nameid"]))
	if err != nil {
		return false
	} else {
		fmt.Printf("Get %s: %s \n", t.Claims.(jwt.MapClaims)["nameid"], strGettoken)
	}

	return true
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (this *CheckUserToken) Get() {
	token := this.Ctx.Input.Header("Token")

	t, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(beego.AppConfig.String("jwtkey")), nil
	})

	strUser := t.Claims.(jwt.MapClaims)["nameid"]

	conn, err := redis.Dial("tcp", "www.tangxinzhu.com:6379")
	if err != nil {
		fmt.Println("connect redis error :", err)
		this.Data["json"] = map[string]interface{}{"status": Success}
		this.ServeJSON()
		return
	}
	conn.Do("AUTH", "passwd")
	defer conn.Close()

	_, err = conn.Do("expire", strUser, 7*24*60*60) //10秒过期
	if err != nil {
		fmt.Println("set expire error: ", err)
	}

	this.Data["json"] = map[string]interface{}{"status": Success}
	this.ServeJSON()
}

func (this *Name) Get() {
	var (
		strSql string
		nNum   int64
		err    error
	)

	//orm.Debug = true
	o := orm.NewOrm()

	strAuthority := this.Ctx.Input.Param(":authority")
	strName := this.Ctx.Input.Param(":name")

	strSql = fmt.Sprintf("SELECT * FROM EntertainmentDB.tbl_%s_name WHERE authority <= %s;", strName, strAuthority)

	var stEntertainmentDBNames []models.EntertainmentDBName
	nNum, err = o.Raw(strSql).QueryRows(&stEntertainmentDBNames)

	if err == nil && nNum > 0 {
		this.Data["json"] = map[string]interface{}{"status": Success, "data": stEntertainmentDBNames}
	} else {
		this.Data["json"] = map[string]interface{}{"status": Failed, "err": err}
	}
	this.ServeJSON()
	//this.StopRun()
}

func (this *NameKeyword) Get() {
	var (
		strSql string
		nNum   int64
		err    error
	)

	orm.Debug = true
	o := orm.NewOrm()

	strAuthority := this.Ctx.Input.Param(":authority")
	strName := this.Ctx.Input.Param(":name")
	strKeyword := this.Ctx.Input.Param(":keyword")

	strSql = fmt.Sprintf("SELECT * FROM EntertainmentDB.tbl_%s_name WHERE authority <= %s AND name LIKE \"%%%s%%\";", strName, strAuthority, strKeyword)

	var stEntertainmentDBNames []models.EntertainmentDBName
	nNum, err = o.Raw(strSql).QueryRows(&stEntertainmentDBNames)
	if err == nil && nNum > 0 {
		this.Data["json"] = map[string]interface{}{"status": Success, "data": stEntertainmentDBNames}
	} else {
		this.Data["json"] = map[string]interface{}{"status": Failed, "err": err}
	}

	this.ServeJSON()
	this.StopRun()
}

func (this *DataType) Get() {
	var (
		strSql string
		nNum   int64
		err    error
	)

	orm.Debug = true
	o := orm.NewOrm()

	strAuthority := this.Ctx.Input.Param(":authority")
	strName := this.Ctx.Input.Param(":name")
	strType := this.Ctx.Input.Param(":datatype")

	strSql = fmt.Sprintf("SELECT * FROM EntertainmentDB.tbl_%s_name WHERE authority <= %s AND type LIKE \"%%%s%%\";", strName, strAuthority, strType)

	var stEntertainmentDBNames []models.EntertainmentDBName
	nNum, err = o.Raw(strSql).QueryRows(&stEntertainmentDBNames)
	if err == nil && nNum > 0 {
		this.Data["json"] = map[string]interface{}{"status": Success, "data": stEntertainmentDBNames}
	} else {
		this.Data["json"] = map[string]interface{}{"status": Failed, "err": err}
	}

	this.ServeJSON()
	this.StopRun()
}

func (this *Chapter) Get() {
	var (
		strSql string
		nNum   int64
		err    error
	)

	orm.Debug = true
	o := orm.NewOrm()

	strName := this.Ctx.Input.Param(":name")
	strID := this.Ctx.Input.Param(":id")
	strBeginNum := this.Ctx.Input.Param(":begin")

	nBeginNum, _ := strconv.Atoi(strBeginNum)
	strEndNum := strconv.Itoa(nBeginNum + 400)

	var nMaxChapter []int
	if "comic" == strName {
		strSql = fmt.Sprintf("SELECT MAX(chapter_num) FROM EntertainmentDB.tbl_comic_chapter  WHERE fk_id = %s;", strID)

		nNum, err = o.Raw(strSql).QueryRows(&nMaxChapter)
		fmt.Println("nMaxChapter", nMaxChapter[0])
		if err != nil || nNum <= 0 {
			nMaxChapter[0] = 0
		}

		strSql = fmt.Sprintf("SELECT * FROM EntertainmentDB.tbl_comic_chapter WHERE fk_id = %s and chapter_num > %s  and chapter_num <= %s ORDER BY chapter_num ASC;;", strID, strBeginNum, strEndNum)
		var stComicChapters []models.ComicChapter
		nNum, err = o.Raw(strSql).QueryRows(&stComicChapters)
		if err == nil && nNum > 0 {
			this.Data["json"] = map[string]interface{}{"status": Success, "data": stComicChapters, "max_chapter": nMaxChapter[0]}
		} else {
			this.Data["json"] = map[string]interface{}{"status": Failed, "err": err}
		}

	} else if "fiction" == strName {
		strSql = fmt.Sprintf("SELECT MAX(chapter_num) FROM EntertainmentDB.tbl_fiction_chapter  WHERE fk_id = %s;", strID)

		nNum, err = o.Raw(strSql).QueryRows(&nMaxChapter)
		fmt.Println("nMaxChapter", nMaxChapter[0])
		if err != nil || nNum <= 0 {
			nMaxChapter[0] = 0
		}

		strSql = fmt.Sprintf("SELECT * FROM EntertainmentDB.tbl_fiction_chapter WHERE fk_id = %s and chapter_num > %s  and chapter_num <= %s ORDER BY chapter_num ASC;;", strID, strBeginNum, strEndNum)
		var stFictionChapter []models.FictionChapter
		nNum, err = o.Raw(strSql).QueryRows(&stFictionChapter)
		if err == nil && nNum > 0 {
			this.Data["json"] = map[string]interface{}{"status": Success, "data": stFictionChapter, "max_chapter": nMaxChapter[0]}
		} else {
			this.Data["json"] = map[string]interface{}{"status": Failed, "err": err}
		}
	}

	this.ServeJSON()
	this.StopRun()
}

func (this *ChapterNum) Get() {
	var (
		strSql string
		nNum   int64
		err    error
	)

	orm.Debug = true
	o := orm.NewOrm()

	strName := this.Ctx.Input.Param(":name")
	strID := this.Ctx.Input.Param(":id")
	strNum := this.Ctx.Input.Param(":num")

	if "comic" == strName {
		strSql = fmt.Sprintf("SELECT * FROM EntertainmentDB.tbl_comic_chapter WHERE fk_id = %s AND chapter_num = %s;", strID, strNum)
		var stComicChapters []models.ComicChapter
		nNum, err = o.Raw(strSql).QueryRows(&stComicChapters)
		if err == nil && nNum > 0 {
			this.Data["json"] = map[string]interface{}{"status": Success, "data": stComicChapters}
		} else {
			this.Data["json"] = map[string]interface{}{"status": Failed, "err": err}
		}
	}

	this.ServeJSON()
	this.StopRun()
}

func (this *ChapterContent) Get() {
	var (
		strSql string
		nNum   int64
		err    error
	)

	orm.Debug = true
	o := orm.NewOrm()

	strName := this.Ctx.Input.Param(":name")
	strID := this.Ctx.Input.Param(":id")
	strNum := this.Ctx.Input.Param(":num")

	if "comic" == strName {
		strSql = fmt.Sprintf("SELECT * FROM EntertainmentDB.tbl_comic_img WHERE fk_comic_id = %s AND fk_comic_chapter_id = %s;", strID, strNum)
		var stComicImgSrcs []models.ComicImgSrc
		nNum, err = o.Raw(strSql).QueryRows(&stComicImgSrcs)
		if err == nil && nNum > 0 {
			this.Data["json"] = map[string]interface{}{"status": Success, "data": stComicImgSrcs}
		} else {
			this.Data["json"] = map[string]interface{}{"status": Failed, "err": err}
		}
	} else if "fiction" == strName {
		strSql = fmt.Sprintf("SELECT * FROM EntertainmentDB.tbl_fiction_chapter WHERE fk_id = %s AND chapter_num = %s;", strID, strNum)
		var stFictionChapter []models.FictionChapter
		nNum, err = o.Raw(strSql).QueryRows(&stFictionChapter)
		if err == nil && nNum > 0 {
			this.Data["json"] = map[string]interface{}{"status": Success, "data": stFictionChapter}
		} else {
			this.Data["json"] = map[string]interface{}{"status": Failed, "err": err}
		}
	}

	this.ServeJSON()
	this.StopRun()
}

func (this *Register) Post() {
	var (
		strSql string
		err    error
	)

	orm.Debug = true
	o := orm.NewOrm()

	strUser := this.GetString("user")
	strPasswd := this.GetString("password")

	strSql = fmt.Sprintf("INSERT INTO EntertainmentDB.user_info (user_name, user_password) VALUES ( \"%s\", \"%s\");", strUser, strPasswd)
	rows, err := o.Raw(strSql).Exec()
	if err == nil {
		o.Commit()
		this.Data["json"] = map[string]interface{}{"status": Success, "rows": rows}
	} else {
		o.Rollback()
		this.Data["json"] = map[string]interface{}{"status": Failed, "err": err}
	}

	this.ServeJSON()
	this.StopRun()
}

func (this *Login) Post() {
	var (
		strSql string
		nNum   int64
		err    error
	)

	orm.Debug = true
	o := orm.NewOrm()

	strUser := this.GetString("user")
	strPasswd := this.GetString("password")

	strSql = fmt.Sprintf("SELECT user_name, user_authority, user_setting FROM EntertainmentDB.user_info WHERE user_name = \"%s\" AND user_password = \"%s\";", strUser, strPasswd)
	var stUser []models.User
	nNum, err = o.Raw(strSql).QueryRows(&stUser)
	if err == nil && nNum > 0 {

		claims := make(jwt.MapClaims)
		claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
		claims["iat"] = time.Now().Unix()
		claims["nameid"] = strUser
		claims["User"] = "true"
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(beego.AppConfig.String("jwtkey")))
		if err != nil {
			this.Data["json"] = map[string]interface{}{"status": Failed, "err": err}

		} else {
			fmt.Println("Token:", tokenString)
			this.Data["json"] = map[string]interface{}{"status": Success, "rows": stUser, "Token": tokenString}

			conn, err := redis.Dial("tcp", "www.tangxinzhu.com:6379")
			if err != nil {
				fmt.Println("connect redis error :", err)
				return
			}
			conn.Do("AUTH", "passwd")
			defer conn.Close()

			token, err := redis.String(conn.Do("GET", strUser))
			if err != nil {
				_, err = conn.Do("SET", strUser, tokenString)
				if err != nil {
					fmt.Println("redis set error:", err)
				}

				_, err = conn.Do("expire", strUser, 7*24*60*60) //10秒过期
				if err != nil {
					fmt.Println("set expire error: ", err)
				}
			} else {
				fmt.Printf("Get [%s]: %s \n", strUser, token)
			}
		}

	} else {
		this.Data["json"] = map[string]interface{}{"status": Failed, "err": err}
	}

	this.ServeJSON()
	this.StopRun()
}

func (this *ChangeSetting) Post() {
	var (
		strSql string
		err    error
	)

	orm.Debug = true
	o := orm.NewOrm()

	strUser := this.GetString("user")
	strSetting := this.GetString("setting")

	strSql = fmt.Sprintf("UPDATE EntertainmentDB.user_info SET user_setting = \"%s\" WHERE user_name = \"%s\";", strSetting, strUser)
	rows, err := o.Raw(strSql).Exec()
	if err == nil {
		o.Commit()
		this.Data["json"] = map[string]interface{}{"status": Success, "rows": rows}
	} else {
		o.Rollback()
		this.Data["json"] = map[string]interface{}{"status": Failed, "err": err}
	}

	this.ServeJSON()
	this.StopRun()
}

func (this *UserCookie) Post() {
	var (
		strSql string
		nNum   int64
		err    error
	)

	orm.Debug = true
	o := orm.NewOrm()

	strCookieType := this.GetString("CookieType")
	strNameID := this.GetString("NameID")
	strUser := this.GetString("User")
	strChapterName := this.GetString("ChapterName")
	strReadNum := this.GetString("ReadNum")

	strWhereinfo := fmt.Sprintf(" WHERE cookie_type = \"%s\" AND name_id = %s AND fk_user = \"%s\";", strCookieType, strNameID, strUser)
	strSql = fmt.Sprintf("SELECT pk_id FROM EntertainmentDB.user_cookie %s", strWhereinfo)

	var nID []int
	nNum, err = o.Raw(strSql).QueryRows(&nID)
	if err == nil && nNum > 0 {
		if strCookieType == "fiction" {
			strReadUrl := this.GetString("ReadUrl")
			strSql = "UPDATE EntertainmentDB.user_cookie SET chapter_name = " + "\"" + strChapterName +
				"\",cur_read_id = " + strReadNum +
				",cur_read_src = " +
				"\"" + strReadUrl + "\"," +
				" fk_user = \"" +
				strUser + "\"" +
				strWhereinfo
		} else {
			strSql = "UPDATE EntertainmentDB.user_cookie SET chapter_name = " + "\"" + strChapterName +
				"\",cur_read_id = " + strReadNum +
				",fk_user = \"" +
				strUser + "\"" +
				strWhereinfo
		}
	} else {
		if strCookieType == "fiction" {
			strReadUrl := this.GetString("ReadUrl")
			strSql = "INSERT into EntertainmentDB.user_cookie (cookie_type, name_id, chapter_name, cur_read_id, cur_read_src, fk_user) VALUES (" +
				"\"" + strCookieType + "\"," +
				strNameID +
				",\"" + strChapterName + "\"," +
				strReadNum + ",\"" +
				strReadUrl + "\",\"" +
				strUser + "\");"
		} else {
			strSql = "INSERT into EntertainmentDB.user_cookie (cookie_type, name_id, chapter_name, cur_read_id, fk_user) VALUES (" +
				"\"" + strCookieType + "\"," +
				strNameID +
				",\"" + strChapterName + "\"," +
				strReadNum + ",\"" +
				strUser + "\");"
		}
	}

	rows, err := o.Raw(strSql).Exec()
	if err == nil {
		o.Commit()
		this.Data["json"] = map[string]interface{}{"status": Success, "rows": rows}
	} else {
		o.Rollback()
		this.Data["json"] = map[string]interface{}{"status": Failed, "err": err}
	}

	this.ServeJSON()
	this.StopRun()
}

func (this *GetCookie) Get() {
	var (
		strSql string
		nNum   int64
		err    error
	)

	orm.Debug = true
	o := orm.NewOrm()

	strType := this.Ctx.Input.Param(":type")
	strUser := this.Ctx.Input.Param(":user")
	strNameID := this.Ctx.Input.Param(":nameid")

	strSql = fmt.Sprintf("SELECT * FROM EntertainmentDB.user_cookie WHERE cookie_type = \"%s\"  AND name_id = %s AND fk_user = \"%s\";", strType, strNameID, strUser)

	var stUserCookies []models.UserCookie
	nNum, err = o.Raw(strSql).QueryRows(&stUserCookies)
	if err == nil && nNum > 0 {
		this.Data["json"] = map[string]interface{}{"status": Success, "data": stUserCookies}
	} else {
		this.Data["json"] = map[string]interface{}{"status": Failed, "err": err}
	}

	this.ServeJSON()
	this.StopRun()
}

func (this *GetLiveBasketBall) Get() {

	this.Data["json"] = map[string]interface{}{"status": Success, "data": models.MapBasketballMatchInfo}

	this.ServeJSON()
	this.StopRun()
}

func (this *GetBasketBallAnalyse) Get() {
	var (
		strSql string
		nNum   int64
		err    error
	)

	orm.Debug = true
	o := orm.NewOrm()

	strName := this.Ctx.Input.Param(":Name")
	strHV := this.Ctx.Input.Param(":HorV")
	strTeamName := this.Ctx.Input.Param(":TeamName")

	var stMatchInfo []models.BasketballMatchInfo
	var stAnalyseInfo models.BasketballAnalyseInfo
	strSql = fmt.Sprintf("SELECT * FROM Bet365.leisu_basketball_data where %s_name like \"%%%s%% and name like \"%%%s%% and match_time like \"2018%%\" ORDER BY match_time desc LIMIT 20;", strHV, strTeamName, strName)
	nNum, err = o.Raw(strSql).QueryRows(&stMatchInfo)

	for i := 0; i < len(stMatchInfo); i++ {
		nHTTotal := stMatchInfo[i].HTFirst + stMatchInfo[i].HTSecond + stMatchInfo[i].HTThird + stMatchInfo[i].HTFourth
		nVTTotal := stMatchInfo[i].VTFirst + stMatchInfo[i].VTSecond + stMatchInfo[i].VTThird + stMatchInfo[i].VTFourth
		fHExpectTotal := float64((stMatchInfo[i].ExpectTotal + stMatchInfo[i].ExpectDiff) / 2)
		fVExpectTotal := float64((stMatchInfo[i].ExpectTotal - stMatchInfo[i].ExpectDiff) / 2)
		fHExpect := float64((fHExpectTotal) / 4)
		fVExpect := float64((fVExpectTotal) / 4)

		stAnalyseInfo.MatchCount++
		stAnalyseInfo.HTFirstScore += stMatchInfo[i].HTFirst
		stAnalyseInfo.HTSecondScore += stMatchInfo[i].HTSecond
		stAnalyseInfo.HTThirdScore += stMatchInfo[i].HTThird
		stAnalyseInfo.HTFourthScore += stMatchInfo[i].HTFourth
		stAnalyseInfo.VTFirstScore += stMatchInfo[i].VTFirst
		stAnalyseInfo.VTSecondScore += stMatchInfo[i].VTSecond
		stAnalyseInfo.VTThirdScore += stMatchInfo[i].VTThird
		stAnalyseInfo.VTFourthScore += stMatchInfo[i].VTFourth

		/*主队单节大数量*/
		if float64(stMatchInfo[i].HTFirst) > fHExpect {
			stAnalyseInfo.HTFirstBig++
		}
		if float64(stMatchInfo[i].HTSecond) > fHExpect {
			stAnalyseInfo.HTSecondBig++
		}
		if float64(stMatchInfo[i].HTThird) > fHExpect {
			stAnalyseInfo.HTThirdBig++
		}
		if float64(stMatchInfo[i].HTFourth) > fHExpect {
			stAnalyseInfo.HTFourthBig++
		}

		/*客队单节大数量*/
		if float64(stMatchInfo[i].VTFirst) > fVExpect {
			stAnalyseInfo.VTFirstBig++
		}
		if float64(stMatchInfo[i].VTSecond) > fVExpect {
			stAnalyseInfo.VTSecondBig++
		}
		if float64(stMatchInfo[i].VTThird) > fVExpect {
			stAnalyseInfo.VTThirdBig++
		}
		if float64(stMatchInfo[i].VTFourth) > fVExpect {
			stAnalyseInfo.VTFourthBig++
		}

		/*总数量*/
		if float64(nHTTotal) > fHExpectTotal {
			stAnalyseInfo.HTMatchBig++
		}
		if float64(nVTTotal) > fVExpectTotal {
			stAnalyseInfo.VTMatchBig++
		}
	}

	stAnalyseInfo.HTFirstScore = int64(stAnalyseInfo.HTFirstScore / stAnalyseInfo.MatchCount)
	stAnalyseInfo.HTSecondScore = int64(stAnalyseInfo.HTSecondScore / stAnalyseInfo.MatchCount)
	stAnalyseInfo.HTThirdScore = int64(stAnalyseInfo.HTThirdScore / stAnalyseInfo.MatchCount)
	stAnalyseInfo.HTFourthScore = int64(stAnalyseInfo.HTFourthScore / stAnalyseInfo.MatchCount)

	stAnalyseInfo.VTFirstScore = int64(stAnalyseInfo.VTFirstScore / stAnalyseInfo.MatchCount)
	stAnalyseInfo.VTSecondScore = int64(stAnalyseInfo.VTSecondScore / stAnalyseInfo.MatchCount)
	stAnalyseInfo.VTThirdScore = int64(stAnalyseInfo.VTThirdScore / stAnalyseInfo.MatchCount)
	stAnalyseInfo.VTFourthScore = int64(stAnalyseInfo.VTFourthScore / stAnalyseInfo.MatchCount)

	stAnalyseInfo.HTFirstBig = stAnalyseInfo.HTFirstBig / float64(stAnalyseInfo.MatchCount)
	stAnalyseInfo.HTSecondBig = stAnalyseInfo.HTSecondBig / float64(stAnalyseInfo.MatchCount)
	stAnalyseInfo.HTThirdBig = stAnalyseInfo.HTThirdBig / float64(stAnalyseInfo.MatchCount)
	stAnalyseInfo.HTFourthBig = stAnalyseInfo.HTFourthBig / float64(stAnalyseInfo.MatchCount)

	stAnalyseInfo.VTFirstBig = stAnalyseInfo.VTFirstBig / float64(stAnalyseInfo.MatchCount)
	stAnalyseInfo.VTSecondBig = stAnalyseInfo.VTSecondBig / float64(stAnalyseInfo.MatchCount)
	stAnalyseInfo.VTThirdBig = stAnalyseInfo.VTThirdBig / float64(stAnalyseInfo.MatchCount)
	stAnalyseInfo.VTFourthBig = stAnalyseInfo.VTFourthBig / float64(stAnalyseInfo.MatchCount)

	fmt.Println(stAnalyseInfo)
	if err == nil && nNum > 0 {
		this.Data["json"] = map[string]interface{}{"status": 1, "data": stMatchInfo, "analyse": stAnalyseInfo}
	} else {
		this.Data["json"] = map[string]interface{}{"status": -1, "err": err}
	}

	this.ServeJSON()
	this.StopRun()
}

func (this *GetLiveFootBall) Get() {
	fmt.Println("GetLiveFootBall")
	MapGoalStatics := make(map[string](models.FootballGoalStatics))

	for key, _ := range models.MapFootballMatchInfo {
		_, ok := MapGoalStatics[key] /*如果确定是真实的,则存在,否则不存在 */
		if !ok {
			MapGoalStatics[key] = models.MapFootballGoalStatics[key]
		}
	}

	this.Data["json"] = map[string]interface{}{"status": Success, "data": MapGoalStatics}

	this.ServeJSON()
	this.StopRun()
}

func (this *GetFootBallAnalyse) Get() {
	fmt.Println("GetFootBallAnalyse")
	strName := this.Ctx.Input.Param(":Name")

	_, ok := models.MapFootballMatchInfo[strName]
	fmt.Println(ok, strName)
	if ok {
		this.Data["json"] = map[string]interface{}{"status": Success, "data": models.MapFootballMatchInfo[strName], "over": models.MapFootballOverMatchInfo[strName], "analyse": models.MapFootballGoalStatics[strName]}
		fmt.Println(models.MapFootballGoalStatics[strName])
	} else {
		this.Data["json"] = map[string]interface{}{"status": Failed}
	}

	this.ServeJSON()
	this.StopRun()
}

func (this *GetFootBallTeamAnalyse) Get() {
	var (
		strSql     string
		strStatics string
		nNum       int64
		err        error
	)

	orm.Debug = true
	o := orm.NewOrm()

	strHV := this.Ctx.Input.Param(":HorV")
	strName := this.Ctx.Input.Param(":Name")
	strTeamName := this.Ctx.Input.Param(":TeamName")
	strNameAlias := this.Ctx.Input.Param(":NameAlias")

	var stMatchInfo []models.FootballMatchInfo

	strSql = fmt.Sprintf("SELECT * FROM Bet365.leisu_football_data where (%s_name like \"%s\" or %s_name like \"%s\") and name like \"%s\" and match_time like \"2018%%\" ORDER BY match_time desc LIMIT 15;", strHV, strTeamName, strHV, strNameAlias, strName)
	nNum, err = o.Raw(strSql).QueryRows(&stMatchInfo)

	var fTotalGoalTime, f75to90MatchNum, fZeroMatchNum, fUpZeroMatchNum, fDownZeroMatchNum float64 = 0, 0, 0, 0, 0
	for i := 0; i < len(stMatchInfo); i++ {
		nHTScore := stMatchInfo[i].HTTotalScore - stMatchInfo[i].HTHalfScore
		nVTScore := stMatchInfo[i].VTTotalScore - stMatchInfo[i].VTHalfScore

		if (stMatchInfo[i].HTTotalScore + stMatchInfo[i].VTTotalScore) == 0 {
			fTotalGoalTime++
			fZeroMatchNum++
		}

		if (stMatchInfo[i].HTHalfScore + stMatchInfo[i].VTHalfScore) == 0 {
			fUpZeroMatchNum++
		}

		if (nHTScore + nVTScore) == 0 {
			fUpZeroMatchNum++
		}

		arrGoalTime := strings.Split(stMatchInfo[i].GoalTime[1:], " ")
		for j := 0; j < len(arrGoalTime); j++ {

			strTime := arrGoalTime[j][1:]
			nTime, _ := strconv.Atoi(strTime)

			/*可统计进球时间比赛总数*/
			fTotalGoalTime++
			if nTime > 75 && nTime <= 90 {
				f75to90MatchNum++
				break
			}

		}

		fUpZero := (fTotalGoalTime - fUpZeroMatchNum) / fTotalGoalTime
		fDownZero := (fTotalGoalTime - fDownZeroMatchNum) / fTotalGoalTime
		fZero := (fTotalGoalTime - fZeroMatchNum) / fTotalGoalTime
		f75to90 := f75to90MatchNum / fTotalGoalTime
		strStatics = fmt.Sprintf("%.0f/%.2f/%.2f/%.2f/%.2f", fTotalGoalTime, fUpZero, fDownZero, fZero, f75to90)
	}

	if err == nil && nNum > 0 {
		this.Data["json"] = map[string]interface{}{"status": Success, "data": strStatics}
	} else {
		this.Data["json"] = map[string]interface{}{"status": Failed, "err": err}
	}

	this.ServeJSON()
	this.StopRun()
}
