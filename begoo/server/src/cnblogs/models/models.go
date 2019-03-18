package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type EntertainmentDBName struct {
	ID         int    `orm:"column(pk_id)"`
	Name       string `orm:"column(name)"`
	Authority  int    `orm:"column(authority)"`
	WatchNum   int    `orm:"column(watch_count)"`
	ChapterNum int    `orm:"column(chapter_count)"`
	IsFinish   bool   `orm:"column(is_finish)"`
	Website    string `orm:"column(website)"`
	Introduce  string `orm:"column(introduce)"`
	Author     string `orm:"column(author)"`
	Img        string `orm:"column(cover_img_src)"`
	Type       string `orm:"column(type)"`
	Time       string `orm:"column(add_time)"`
}

type ComicChapter struct {
	ID          int    `orm:"column(pk_id)"`
	ChapterNum  int    `orm:"column(chapter_num)"`
	ChapterName string `orm:"column(chapter_name)"`
	PicNum      int    `orm:"column(pic_count)"`
	Dept_ID     int    `orm:"column(fk_id)"`
}

type ComicImgSrc struct {
	ID              int    `orm:"column(pk_id)"`
	PageNum         int    `orm:"column(page_num)"`
	ComicID         int    `orm:"column(fk_comic_id)"`
	ComicChapterNum int    `orm:"column(fk_comic_chapter_id)"`
	ImgSrc          string `orm:"column(img_src)"`
}

type FictionChapter struct {
	ID          int    `orm:"column(pk_id)"`
	ChapterNum  int    `orm:"column(chapter_num)"`
	ChapterName string `orm:"column(chapter_name)"`
	ContentUrl  string `orm:"column(content_src)"`
	Dept_ID     int    `orm:"column(fk_id)"`
}

type UserCookie struct {
	UserName    string `orm:"column(user_name)"`
	NameID      int    `orm:"column(name_id)"`
	ChapterName string `orm:"column(chapter_name)"`
	ReadNum     int    `orm:"column(cur_read_id)"`
	ReadUrl     string `orm:"column(cur_read_src)"`
	Dept_User   string `orm:"column(fk_user)"`
}

type User struct {
	UserName      string `orm:"column(user_name)"`
	UserAuthority string `orm:"column(user_authority)"`
	UserSetting   string `orm:"column(user_setting)"`
}

type BasketballMatchInfo struct {
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

type BasketballAnalyseInfo struct {
	MatchCount    int64
	HTFirstScore  int64
	HTSecondScore int64
	HTThirdScore  int64
	HTFourthScore int64
	HTFirstBig    float64
	HTSecondBig   float64
	HTThirdBig    float64
	HTFourthBig   float64
	HTMatchBig    int64
	VTFirstScore  int64
	VTSecondScore int64
	VTThirdScore  int64
	VTFourthScore int64
	VTFirstBig    float64
	VTSecondBig   float64
	VTThirdBig    float64
	VTFourthBig   float64
	VTMatchBig    int64
}

type FootballMatchInfo struct {
	ID            int    `orm:"column(pk_id)"`
	MatchID       string `orm:"column(match_id)"`
	Name          string `orm:"column(name)"`
	HTName        string `orm:"column(ht_name)"`
	VTName        string `orm:"column(vt_name)"`
	HTTotalScore  int64  `orm:"column(ht_total_score)"`
	VTTotalScore  int64  `orm:"column(vt_total_score)"`
	HTHalfScore   int64  `orm:"column(ht_half_score)"`
	VTHalfScore   int64  `orm:"column(vt_half_score)"`
	HTTotalCorner int64  `orm:"column(ht_total_corner)"`
	VTTotalCorner int64  `orm:"column(vt_total_corner)"`
	HTHalfCorner  int64  `orm:"column(ht_half_corner)"`
	VTHalfCorner  int64  `orm:"column(vt_half_corner)"`
	HTShoot       int64  `orm:"column(ht_shoot)"`
	VTShoot       int64  `orm:"column(vt_shoot)"`
	HTShooton     int64  `orm:"column(ht_shoot_on)"`
	VTShooton     int64  `orm:"column(vt_shoot_on)"`
	HTRed         int64  `orm:"column(ht_red)"`
	VTRed         int64  `orm:"column(vt_red)"`
	Asionodd      string `orm:"column(asion_odd)"`
	Cornerodd     string `orm:"column(corner_odd)"`
	Numodd        string `orm:"column(num_odd)"`
	DetailHref    string `orm:"column(detail_href)"`
	GoalTime      string `orm:"column(goal_time)"`
	MatchTime     string `orm:"column(match_time)"`
	HGoalTime     string
	VGoalTime     string
	EventSta      string
	NowMatchTime  int
	HTNameAlias   string
	VTNameAlias   string
	HTStatic      string
	VTStatic      string
}

type FootballGoalStatics struct {
	MatchTotal         int64                            /*比赛总数*/
	GoalTotal          int64                            /*进球总数*/
	UpGoalCount        int64                            /*上半场进球总数*/
	DownGoalCount      int64                            /*下半场进球总数*/
	UpHaveGoalMatch    int64                            /*上半场有球比赛总数*/
	DownHaveGoalMatch  int64                            /*下半场有球比赛总数*/
	AllHaveGoalMatch   int64                            /*全场有球比赛总数*/
	HaveCornerMatch    int64                            /*可统计有角球比赛总数*/
	CornerTotal        int64                            /*角球总数*/
	TotalGoalTime      int64                            /*可统计进球时间进球总数*/
	TotalGoalTimeMatch int64                            /*可统计进球时间比赛总数*/
	ShootCount         int64                            /*射门总数*/
	ShootOnCount       int64                            /*射正总数*/
	TotalShootMatch    int64                            /*射门有效比赛数*/
	MapFootballGoal    map[string]FootballGoalNumByTime /*进球时间分布*/
	MapUpEffectDown    map[string]FootballGoalEffect    /*上半场对下半场影响*/
}

type FootballGoalNumByTime struct {
	GoalNum      int64
	GoalMatchNum int64
	HaveGoal     bool
}

type FootballGoalEffect struct {
	GoalHaveGoalMatch  int64
	GoalTotalGoalMatch int64
}

var (
	MapBasketballMatchInfo   = make(map[string]([]BasketballMatchInfo))
	MapFootballMatchInfo     = make(map[string]([]FootballMatchInfo))
	MapFootballOverMatchInfo = make(map[string]([]FootballMatchInfo))
	MapFootballGoalStatics   = make(map[string](FootballGoalStatics))

	//ArrMatchInfo = []MatchInfo{}
)

func init() {
	fmt.Println("数据库连接中！")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:passwd@tcp(149.129.79.183:3306)/EntertainmentDB?charset=utf8")
	orm.SetMaxIdleConns("default", 1000)
	orm.SetMaxOpenConns("default", 1000)

	fmt.Println("数据库连接成功！")
}
