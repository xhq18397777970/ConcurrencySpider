package concurrentyspider

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	USERNAME = "root"
	PASSWORD = "1097966132"
	PORT     = "3306"
	DBNAME   = "test"
	HOST     = "127.0.0.1"
)

var (
	DB *sql.DB

	TaskChannel = make(chan string, 10)

	page int    = 4
	url  string = "https://movie.douban.com/top250?start="

	maxWorker int = 4
)

type Movie_data struct {
	Year     string
	Director string
	Actor    string
	Quote    string
	Title    string
	Score    string
	Img      string
}

func Start() {
	var wg sync.WaitGroup
	now := time.Now()
	InitDB()
	distributeTask(TaskChannel)
	//range应该从通道接收数据，而不是发送数据
	// for i := range TaskChannel {
	// 	fmt.Println(i)
	// }
	for i := 0; i < maxWorker; i++ {
		wg.Add(1)
		go Spider(&wg, TaskChannel)
	}
	//要加上wait（）来等待goroutine执行完毕函数才返回，之前错误就是由于goroutine还没执行成功，主函数退出
	wg.Wait()
	defer func() { fmt.Printf("数据爬取完毕，插入数据库完成,耗时：%s\n", time.Since(now)) }()
}

func Spider(wg *sync.WaitGroup, ch chan string) {
	//wg.Done()应该在工作开始之前调用，确保所有工作完成后执行（都可以）
	defer wg.Done()
	for task := range ch {
		//1、创建请求
		fmt.Println("task", task)
		client := http.Client{}
		req, err := http.NewRequest("GET", task, nil)
		if err != nil {
			fmt.Println("request error", err)
		}

		req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("client do error", err)
		}
		// defer resp.Body.Close()

		//3、解析网页数据
		docDetail, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			fmt.Println("analyse page error", err)
		}

		docDetail.Find("#content > div > div.article > ol > li").
			//s取标签地址的值，得到标签的value
			Each(func(i int, s *goquery.Selection) {
				// #content > div > div.article > ol > li:nth-child(1) > div > div.info > div.hd > a > span:nth-child(1)
				title := s.Find("div.info > div.hd > a > span:nth-child(1)").Text()
				info := s.Find("div > div.info > div.bd > p:nth-child(1)").Text()
				quote := s.Find("div > div.info > div.bd > div > span:nth-child(4)").Text()
				score := s.Find("div > div.info > div.bd > div > span.rating_num").Text()
				imgTemp := s.Find("div > div.pic > a > img")
				img, _ := imgTemp.Attr("src")
				year, director, actor := InfoSpite(info)

				var moviedata Movie_data
				// var moviedata = new(Movie_data)  定义的结构体，只是为了构造该数据结构，直接写入数据库，不需要修改其定义的值，使用new初始化结构体，得到的是结构体指针
				moviedata.Actor = actor
				moviedata.Img = img
				moviedata.Director = director
				moviedata.Quote = quote
				moviedata.Score = score
				moviedata.Title = title
				moviedata.Year = year

				fmt.Println(moviedata)
				if InsertData(moviedata) {
					fmt.Println("插入成功")
				} else {
					fmt.Println("插入失败")
				}

			})

	}
	// wg.Done() 在此处执行也可以

}

// 4、数据处理
func InfoSpite(info string) (year, director, actor string) {
	// year,director,actor:=
	yearRe, _ := regexp.Compile(`(\d+)`)
	directorRe, _ := regexp.Compile(`导演: (.*)`)
	actorRe, _ := regexp.Compile(`主演: (.*)`)

	year = yearRe.FindString(info)
	director = directorRe.FindString(info)
	actor = actorRe.FindString(info)
	return
}

// 5、存储数据库
func InitDB() {
	path := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(", HOST, ":", PORT, ")/", DBNAME, "?charset=utf8"}, "")
	var err error
	DB, err = sql.Open("mysql", path)
	if err != nil {
		fmt.Println("open mysql error", err)
	}
	DB.SetConnMaxLifetime(10)
	DB.SetMaxIdleConns(5)

	if err := DB.Ping(); err != nil {
		fmt.Println("DB ping fail", err)
		return
	}
	fmt.Println("connect databse successful!")
}

func InsertData(data Movie_data) bool {
	sqlInsert := "INSERT INTO movie_data (Title,Director,Picture,Actor,Year,Score,Quote)VALUES(?,?,?,?,?,?,?)"
	_, err := DB.Exec(sqlInsert, data.Title, data.Director, data.Img, data.Actor, data.Year, data.Score, data.Quote)
	if err != nil {
		fmt.Println("insert to tables error ", err)
		return false
	} else {
		return true
	}
}

// func distributeTask(ch chan<- string) {   此时全局变量定义的是chan string双向通道，而不是单向
func distributeTask(ch chan string) {
	// for i := range page {  page不是一个可迭代的对象，
	for i := 0; i < page; i++ {
		ch <- url + strconv.Itoa(i*25)
	}
	close(ch)
}
