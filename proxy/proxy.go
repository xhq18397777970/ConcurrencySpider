package proxy

import (
	"net/http"
	"net/url"
)

var client http.Client

func SetProxy() {
	//设置代理，并打印当前ip地址及ip所在地，需要向提供ip地址的http api发起请求
	proxyURL := "http://127.0.0.1:7890"
	proxy, _ := url.Parse(proxyURL)
	//在该结构体内定义代理配置方式、TLS配置、请求超时等
	client = http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}

}
