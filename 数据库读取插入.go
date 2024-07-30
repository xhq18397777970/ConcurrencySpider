package Concurrency_Spider

import (
	"database/sql"
	"fmt"
	"strings"

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
	Title    string `json:"title"`
	Director string `json:"director"`
	Picture  string `json:"picture"`
	Actor    string `json:"actor"`
	Year     string `json:"year"`
	Score    string `json:"score"`
	Quote    string `json:"quote"`
}

func Start() {
	if err := InitDB(); err != nil {
		fmt.Println("failed to connect to the database", err)
		return
	}
	//不要在initDB函数中调用DB.close（）
	defer DB.Close()
	//1、new：用于实例化（结构体、int、float等类型）用于申请一片内存空间，返回内存空间的地址（指针类型）
	//2、make：创建数据结构（slice、chan、map），分配的空间清零（返回值类型）
	xsk := new(movieData)
	xsk.Title = "向幸福出发"
	xsk.Director = "张艺谋"
	xsk.Picture = "1231.png"
	xsk.Actor = "彭于晏"
	xsk.Score = "10"
	xsk.Quote = "1w"
	xsk.Year = "2012年9月"
	fmt.Println(*xsk)

	if InsertData(*xsk) {
		fmt.Println("插入成功！")
	}

}

func InsertData(moviedata movieData) bool {

	sqlInsert := "INSERT INTO movie_data(Title,Director,Picture,Actor,Year,Score,Quote)VALUES(?,?,?,?,?,?,?)"
	res, err := DB.Exec(sqlInsert, moviedata.Title, moviedata.Director, moviedata.Picture, moviedata.Actor, moviedata.Year, moviedata.Score, moviedata.Quote)
	if err != nil {
		fmt.Println("exec error", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println("id get erro", err)
	}
	fmt.Println("id", id)
	return true
}

func InitDB() error {

	path := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(", HOST, ":", PORT, ")/", DBNAME, "?charset=utf8"}, "")
	var err error
	DB, err = sql.Open("mysql", path)
	if err != nil {
		fmt.Println(err)
	}
	DB.SetConnMaxLifetime(10)
	DB.SetMaxIdleConns(5)
	if err := DB.Ping(); err != nil {
		fmt.Println("ping database fail")
		return err
	}
	fmt.Println("connect database success!")
	return nil
}
func show() {
	sql := "SHOW TABLES"
	rows, err := DB.Query(sql)
	if err != nil {
		fmt.Println("failed to query databse", err)
		return
	}
	defer rows.Close()
	var tableName string
	for rows.Next() {
		//将读取到的表，赋值给tablename内存空间
		if err := rows.Scan(&tableName); err != nil {
			fmt.Println("fail to scan tables", err)
			return
		}
		fmt.Println(tableName)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("error iterating", err)
	}

}
