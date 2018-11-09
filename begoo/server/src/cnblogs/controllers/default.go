package controllers

import (
	"cnblogs/models"
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

	this.Data["json"] = map[string]interface{}{"status": 1, "data": models.ArrMatchInfo}

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

	var stMatchInfo []models.MatchInfo

	strSql = "SELECT * FROM Bet365.leisu_basketball_data where " + strHV + "_name like \"%" + strTeamName + "%\" and name like \"%" + strName + "%\" and match_time like \"2018%\" ORDER BY match_time desc LIMIT 20; "
	nNum, err = o.Raw(strSql).QueryRows(&stMatchInfo)

	if err == nil && nNum > 0 {
		this.Data["json"] = map[string]interface{}{"status": 1, "data": stMatchInfo}
	} else {
		this.Data["json"] = map[string]interface{}{"status": -1, "err": err}
	}

	this.ServeJSON()
	this.StopRun()
}
