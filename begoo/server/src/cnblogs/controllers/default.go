package controllers

import (
	"cnblogs/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
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

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
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

	if "Comic" == strName {
		strSql = "SELECT * FROM EntertainmentDB.tbl_comic_name" + " WHERE authority <= " + strAuthority + ";"
	} else if "Fiction" == strName {
		strSql = "SELECT * FROM EntertainmentDB.tbl_fiction_name" + " WHERE authority <= " + strAuthority + ";"
	}

	var stEntertainmentDBNames []models.EntertainmentDBName
	nNum, err = o.Raw(strSql).QueryRows(&stEntertainmentDBNames)

	if err == nil && nNum > 0 {
		this.Data["json"] = map[string]interface{}{"data": stEntertainmentDBNames}
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

	if "Comic" == strName {
		strSql = "SELECT * FROM EntertainmentDB.tbl_comic_name" + " WHERE authority <= " + strAuthority + " AND name LIKE \"%" + strKeyword + "%\";"
	} else if "Fiction" == strName {
		strSql = "SELECT * FROM EntertainmentDB.tbl_fiction_name" + " WHERE authority <= " + strAuthority + " AND name LIKE \"%" + strKeyword + "%\";"
	}

	var stEntertainmentDBNames []models.EntertainmentDBName
	nNum, err = o.Raw(strSql).QueryRows(&stEntertainmentDBNames)
	if err == nil && nNum > 0 {
		this.Data["json"] = map[string]interface{}{"data": stEntertainmentDBNames}
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

	if "Comic" == strName {
		strSql = "SELECT * FROM EntertainmentDB.tbl_comic_name" + " WHERE authority <= " + strAuthority + " AND type LIKE \"%" + strType + "%\";"
	} else if "Fiction" == strName {
		strSql = "SELECT * FROM EntertainmentDB.tbl_fiction_name" + " WHERE authority <= " + strAuthority + " AND type LIKE \"%" + strType + "%\";"
	}

	var stEntertainmentDBNames []models.EntertainmentDBName
	nNum, err = o.Raw(strSql).QueryRows(&stEntertainmentDBNames)
	if err == nil && nNum > 0 {
		this.Data["json"] = map[string]interface{}{"data": stEntertainmentDBNames}
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

	if "Comic" == strName {
		strSql = "SELECT * FROM EntertainmentDB.tbl_comic_chapter WHERE fk_id=" + strID + " ORDER BY chapter_num ASC;"
		var stComicChapters []models.ComicChapter
		nNum, err = o.Raw(strSql).QueryRows(&stComicChapters)
		if err == nil && nNum > 0 {
			this.Data["json"] = map[string]interface{}{"data": stComicChapters}
		}
	} else if "Fiction" == strName {
		strSql = "SELECT * FROM EntertainmentDB.tbl_fiction_chapter WHERE fk_id	=" + strID + " ORDER BY chapter_num ASC;"
		var stFictionChapter []models.FictionChapter
		nNum, err = o.Raw(strSql).QueryRows(&stFictionChapter)
		if err == nil && nNum > 0 {
			this.Data["json"] = map[string]interface{}{"data": stFictionChapter}
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

	if "Comic" == strName {
		strSql = "SELECT * FROM EntertainmentDB.tbl_comic_chapter WHERE fk_id = " + strID + " AND chapter_num = " + strNum + ";"
		var stComicChapters []models.ComicChapter
		nNum, err = o.Raw(strSql).QueryRows(&stComicChapters)
		if err == nil && nNum > 0 {
			this.Data["json"] = map[string]interface{}{"data": stComicChapters}
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

	if "Comic" == strName {
		strSql = "SELECT * FROM EntertainmentDB.tbl_comic_img WHERE fk_comic_id = " + strID + " AND fk_comic_chapter_id = " + strNum + ";"
		var stComicImgSrcs []models.ComicImgSrc
		nNum, err = o.Raw(strSql).QueryRows(&stComicImgSrcs)
		if err == nil && nNum > 0 {
			this.Data["json"] = map[string]interface{}{"status": 1, "data": stComicImgSrcs}
		} else {
			this.Data["json"] = map[string]interface{}{"status": -1, "err": err}
		}
	} else if "Fiction" == strName {
		strSql = "SELECT * FROM EntertainmentDB.tbl_fiction_chapter WHERE fk_id=" + strID + " AND chapter_num = " + strNum + ";"
		var stFictionChapter []models.FictionChapter
		nNum, err = o.Raw(strSql).QueryRows(&stFictionChapter)
		if err == nil && nNum > 0 {
			this.Data["json"] = map[string]interface{}{"status": 1, "data": stFictionChapter}
		} else {
			this.Data["json"] = map[string]interface{}{"status": -1, "err": err}
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
	strSql = "INSERT INTO EntertainmentDB.user_info (user_name, user_password) VALUES (" + "\"" + strUser + "\",\"" + strPasswd + "\");"
	rows, err := o.Raw(strSql).Exec()

	if err == nil {
		o.Commit()
		this.Data["json"] = map[string]interface{}{"status": 1, "rows": rows}
	} else {
		o.Rollback()
		this.Data["json"] = map[string]interface{}{"status": -1, "err": err}
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

	strSql = "SELECT * FROM EntertainmentDB.user_info WHERE user_name = " + "\"" + strUser + "\"" + " AND user_password = " + "\"" + strPasswd + "\";"
	var stUser []models.User
	nNum, err = o.Raw(strSql).QueryRows(&stUser)
	if err == nil && nNum > 0 {
		this.Data["json"] = map[string]interface{}{"status": 1, "rows": stUser}
	} else {
		this.Data["json"] = map[string]interface{}{"status": -1, "err": err}
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
	strSql = "UPDATE EntertainmentDB.user_info SET user_setting = " + "\"" + strSetting + "\"" + " WHERE user_name = " + "\"" + strUser + "\";"

	rows, err := o.Raw(strSql).Exec()
	if err == nil {
		o.Commit()
		this.Data["json"] = map[string]interface{}{"status": 1, "rows": rows}
	} else {
		o.Rollback()
		this.Data["json"] = map[string]interface{}{"status": -1, "err": err}
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

	strWhereinfo := " WHERE cookie_type = " + "\"" + strCookieType + "\"" + " AND name_id = " + strNameID + " AND fk_user = " + "\"" + strUser + "\";"

	strSql = "SELECT pk_id FROM EntertainmentDB.user_cookie" + strWhereinfo
	var nID []int
	nNum, err = o.Raw(strSql).QueryRows(&nID)
	if err == nil && nNum > 0 {
		if strCookieType == "Fiction" {
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
		if strCookieType == "Fiction" {
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
		this.Data["json"] = map[string]interface{}{"status": 1, "rows": rows}
	} else {
		o.Rollback()
		this.Data["json"] = map[string]interface{}{"status": -1, "err": err}
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
	strSql = "SELECT * FROM EntertainmentDB.user_cookie WHERE cookie_type = " + "\"" + strType + "\"" + " AND name_id = " + strNameID + " AND fk_user = " + "\"" + strUser + "\";"

	var stUserCookies []models.UserCookie
	nNum, err = o.Raw(strSql).QueryRows(&stUserCookies)
	if err == nil && nNum > 0 {
		this.Data["json"] = map[string]interface{}{"status": 1, "data": stUserCookies}
	} else {
		this.Data["json"] = map[string]interface{}{"status": -1, "err": err}
	}

	this.ServeJSON()
	this.StopRun()
}

func (this *GetLiveBasketBall) Get() {

	this.Data["json"] = map[string]interface{}{"status": 1, "data": models.MapBasketballMatchInfo}

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

	strSql = "SELECT * FROM Bet365.leisu_basketball_data where " + strHV + "_name like \"%" + strTeamName + "%\" and name like \"%" + strName + "%\" and match_time like \"2018%\" ORDER BY match_time desc LIMIT 20; "
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

	this.Data["json"] = map[string]interface{}{"status": 1, "data": models.MapFootballMatchInfo}

	this.ServeJSON()
	this.StopRun()
}

func (this *GetFootBallAnalyse) Get() {
	strName := this.Ctx.Input.Param(":Name")

	_, ok := models.MapFootballGoalStatics[strName]
	fmt.Println(ok, strName)
	if ok {
		this.Data["json"] = map[string]interface{}{"status": 1, "analyse": models.MapFootballGoalStatics[strName]}
		fmt.Println(models.MapFootballGoalStatics[strName])
	} else {
		this.Data["json"] = map[string]interface{}{"status": -1}
	}

	this.ServeJSON()
	this.StopRun()
}
