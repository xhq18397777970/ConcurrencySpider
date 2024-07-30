package rank_king

//start为主函数入口
//该模块为爬取动态js数据（api接口数据）
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type KingRankResp struct {
	Code int64 `json:"code"`
	Data struct {
		Replies []struct {
			Content struct {
				Device  string        `json:"device"`
				JumpURL struct{}      `json:"jump_url"`
				MaxLine int64         `json:"max_line"`
				Members []interface{} `json:"members"`
				Message string        `json:"message"`
				Plat    int64         `json:"plat"`
			} `json:"content"`
			Count  int64 `json:"count"`
			Folder struct {
				HasFolded bool   `json:"has_folded"`
				IsFolded  bool   `json:"is_folded"`
				Rule      string `json:"rule"`
			} `json:"folder"`
			Like    int64 `json:"like"`
			Replies []struct {
				Action  int64 `json:"action"`
				Assist  int64 `json:"assist"`
				Attr    int64 `json:"attr"`
				Content struct {
					Device  string   `json:"device"`
					JumpURL struct{} `json:"jump_url"`
					MaxLine int64    `json:"max_line"`
					Message string   `json:"message"`
					Plat    int64    `json:"plat"`
				} `json:"content"`
				Rcount  int64       `json:"rcount"`
				Replies interface{} `json:"replies"`
			} `json:"replies"`
			Type int64 `json:"type"`
		} `json:"replies"`
	} `json:"data"`
	Message string `json:"message"`
}

var url string = "https://www.bilibili.com/bangumi/play/ss39462"

func start() {
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/v2/reply/wbi/main?oid=312421175&type=1&mode=3&pagination_str=%7B%22offset%22:%22%22%7D&plat=1&web_location=1315875&w_rid=9b957fa58985a205d60595db3280199c&wts=1722183758", nil)
	if err != nil {
		fmt.Println("request error", err)
	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("", err)
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("io error", err)
	}
	var resultList KingRankResp
	//由于go语言不同于python直接json转字典，需要构造与json中对应的struct，再进行操作
	//将响应内容转为struct结构
	_ = json.Unmarshal(bodyText, &resultList)

	//为列表类型
	for _, result := range resultList.Data.Replies {
		fmt.Println("一级评论为:", result.Content.Message)
		for _, reply := range result.Replies {
			fmt.Println("回复为:", reply.Content.Message)
		}
	}
}
