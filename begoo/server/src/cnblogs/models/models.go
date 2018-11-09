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
	UserPassword  string `orm:"column(user_password)"`
	UserAuthority string `orm:"column(user_authority)"`
	UserSetting   string `orm:"column(user_setting)"`
}

/*下载每天数据的数据库结构*/
type MatchInfo struct {
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

var (
	ArrMatchInfo = []MatchInfo{}
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "txz:passwd@tcp(198.13.54.7:3306)/EntertainmentDB?charset=utf8", 1000, 2000)
	fmt.Println("数据库连接成功！")
}
