package routers

import (
	"cnblogs/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.Router("/api/:authority/Name/:name", &controllers.Name{})
	beego.Router("/api/:authority/SearchResult/:name/All/:keyword", &controllers.NameKeyword{})
	beego.Router("/api/:authority/SearchResult/:name/Type/:datatype", &controllers.DataType{})

	beego.Router("/api/Name/:name/Chapter/:id/:begin", &controllers.Chapter{})
	beego.Router("/api/Name/:name/:id/ChapterByNum/:num", &controllers.ChapterNum{})
	beego.Router("/api/Name/:name/:id/Chapter/:num", &controllers.ChapterContent{})

	beego.Router("/api/Register", &controllers.Register{})
	beego.Router("/api/Login", &controllers.Login{})
	beego.Router("/api/ChangeSetting", &controllers.ChangeSetting{})
	beego.Router("/api/UserCookie", &controllers.UserCookie{})
	beego.Router("/api/Type/:type/User/:user/NameID/:nameid", &controllers.GetCookie{})

	beego.Router("/api/Live/BasketBall", &controllers.GetLiveBasketBall{})
	beego.Router("/api/Finish/BasketballAnalyse/:Name/:HorV/:TeamName", &controllers.GetBasketBallAnalyse{})

	beego.Router("/api/Live/FootBall", &controllers.GetLiveFootBall{})
	beego.Router("/api/Finish/FootballAnalyse/:Name", &controllers.GetFootBallAnalyse{})
	beego.Router("/api/Finish/FootBallTeamAnalyse/:Name/:HorV/:TeamName/:NameAlias", &controllers.GetFootBallTeamAnalyse{})
}
