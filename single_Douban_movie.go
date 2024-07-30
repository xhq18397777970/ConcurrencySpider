package Concurrency_Spider

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 定义数据库常量
const (
	USERNAME = "root"
	PASSWORD = "1097966132"
	HOST     = "127.0.0.1"
	PORT     = "3306"
	DBNAME   = "test"
)

// 定义数据库全局变量
var DB *sql.DB

type movieData struct {
	Title    string
	Director string
	Picture  string
	Actor    string
	Score    string
	Year     string
	Quote    string
}

func Start() {
	InitDB()
	for i := 0; i <= 10; i++ {
		fmt.Printf("正在爬取第%d页数据\n", i+1)
		Spider(strconv.Itoa(i * 25))
	}

}

func Spider(page string) {
	//1、发送请求
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/top250?start="+page, nil)
	if err != nil {
		fmt.Println("req err", err)
	}
	//防止浏览器反扒，设置请求头（只要加user-agent即可）
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败", err)
	}
	defer resp.Body.Close()

	//2、解析网页
	//得到请求树
	docDetail, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("解析失败", err)
	}

	//3、获取节点信息

	//#content > div > div.article > ol > li
	//#content > div > div.article > ol > li:nth-child(1) > div > div.info > div.hd > a > span:nth-child(1)
	//#content > div > div.article > ol > li:nth-child(1) > div > div.pic > a > img

	//Find要拿*goquery.Selecttion对象中的文本，需要加上text.否则返回的是列表

	docDetail.Find("#content > div > div.article > ol > li").
		Each(func(i int, s *goquery.Selection) {
			var data movieData

			title := s.Find(" div > div.info > div.hd > a > span:nth-child(1)").Text()
			//此时需要得到图片的src，而这个获取的是img标签
			img := s.Find(" div > div.pic > a > img")
			//
			imgTemp, ok := img.Attr("src")

			info := s.Find("div > div.info > div.bd > p:nth-child(1)").Text()
			scope := s.Find(" div > div.info > div.bd > div > span.rating_num").Text()
			quote := s.Find(" div > div.info > div.bd > div > span:nth-child(4)").Text()

			if ok {

				director, actor, year := InfoSpite(info)
				data.Actor = actor
				data.Director = director
				data.Picture = imgTemp
				data.Year = year
				data.Score = scope
				data.Title = title
				data.Quote = quote

			}
			if InsertData(data) {

			} else {
				fmt.Println("插入失败")
				return
			}

			// fmt.Println("data ", data)

		})

	fmt.Println("插入成功!")
	return

	//4、保存数据

}
func InfoSpite(info string) (director, actor, year string) {
	//regexp包 --正则表达式
	yearRe, _ := regexp.Compile(`(\d+)`)
	year = yearRe.FindString(info)

	directorRe, _ := regexp.Compile(`导演:(.*)`)
	director = directorRe.FindString(info)

	actorRe, _ := regexp.Compile(`主演: (.*)`)
	actor = actorRe.FindString(info)
	return
}

// 初始化数据库
func InitDB() {

	path := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(", HOST, ":", PORT, ")/", DBNAME, "?charset=utf8"}, "")
	// sql.Open()返回2个参数：*sql.DB数据库连接对象，错误对象
	//注意，不要再将DB新赋值，此时为=而不是：=
	DB, _ = sql.Open("mysql", path)

	DB.SetConnMaxLifetime(10)
	DB.SetMaxIdleConns(5)

	if err := DB.Ping(); err != nil {
		fmt.Println("ping database fail")
		return
	}
	fmt.Println("connect database success!")

}
func InsertData(moviedata movieData) bool {

	sqlInsert := "INSERT INTO movie_data(Title,Director,Picture,Actor,Year,Score,Quote)VALUES(?,?,?,?,?,?,?)"
	_, err := DB.Exec(sqlInsert, moviedata.Title, moviedata.Director, moviedata.Picture, moviedata.Actor, moviedata.Year, moviedata.Score, moviedata.Quote)
	if err != nil {
		fmt.Println("exec error", err)
	}
	//表的ID为自增字段，当不传入id时，会自动传入id值，
	//打印最后一条的id
	// id, err := res.LastInsertId()
	// if err != nil {
	// 	fmt.Println("id get erro", err)
	// }
	// fmt.Println("id", id)
	return true
}
