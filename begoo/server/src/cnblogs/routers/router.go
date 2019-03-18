package routers

import (
	"cnblogs/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	/*娱乐：查找*/
	beego.Router("/Api/Entertainment/:authority/Name/:name", &controllers.Name{})
	beego.Router("/Api/Entertainment/:authority/SearchResult/:name/All/:keyword", &controllers.NameKeyword{})
	beego.Router("/Api/Entertainment/:authority/SearchResult/:name/Type/:datatype", &controllers.DataType{})

	/*娱乐：漫画及小说*/
	beego.Router("/Api/Entertainment/:name/Chapter/:id/:begin", &controllers.Chapter{})
	beego.Router("/Api/Entertainment/:name/:id/ChapterByNum/:num", &controllers.ChapterNum{})
	beego.Router("/Api/Entertainment/:name/:id/Chapter/:num", &controllers.ChapterContent{})

	/*用户配置*/
	beego.Router("/Api/Check/UserToken", &controllers.CheckUserToken{})
	beego.Router("/Api/Register", &controllers.Register{})
	beego.Router("/Api/Login", &controllers.Login{})
	beego.Router("/Api/User/ChangeSetting", &controllers.ChangeSetting{})
	beego.Router("/Api/User/Cookie", &controllers.UserCookie{})
	beego.Router("/Api/User/Type/:type/User/:user/NameID/:nameid", &controllers.GetCookie{})

	/*运动：篮球*/
	beego.Router("/Api/Sports/BasketBall/Live", &controllers.GetLiveBasketBall{})
	beego.Router("/Api/Sports/Basketball/Finish/Analyse/:Name/:HorV/:TeamName", &controllers.GetBasketBallAnalyse{})

	/*运动：足球*/
	beego.Router("/Api/Sports/FootBall/Live", &controllers.GetLiveFootBall{})
	beego.Router("/Api/Sports/FootBall/Finish/Analyse/:Name", &controllers.GetFootBallAnalyse{})
	beego.Router("/Api/Sports/FootBall/Finish/TeamAnalyse/:Name/:HorV/:TeamName/:NameAlias", &controllers.GetFootBallTeamAnalyse{})
}
