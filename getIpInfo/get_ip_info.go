package proxy

//获取当前机器的公网ip地址需要通过第三方服务，因为本地ip地址是私有的，并不能直接在互联网上使用

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

type IpaddrGeoInfo struct {
	As          string  `json:"as"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Isp         string  `json:"isp"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Org         string  `json:"org"`
	Query       string  `json:"query"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	Status      string  `json:"status"`
	Timezone    string  `json:"timezone"`
	Zip         string  `json:"zip"`
}

func Private_Ip() {
	// 获取本地IP地址
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error getting interface addresses:", err)
		return
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		// 打印IPv4地址
		if ip.To4() != nil {
			fmt.Printf("[局域网ip:  %s]\n", ip.String())
		}
	}
}
func Public_Ip() {
	// 使用ipify服务获取公网IP地址
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		fmt.Println("Error fetching IP address:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	ipBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 打印IP地址
	fmt.Println("Public IP Address:", string(ipBytes))

	// 如果你还想获取IP地址的地理位置信息，可以使用ip-api服务
	ipinfoBytes, err := http.Get("http://ip-api.com/json/?fields=61439")
	if err != nil {
		fmt.Println("ip地理信息查询接口失效", err)
	}

	info, err := ioutil.ReadAll(ipinfoBytes.Body)
	defer ipinfoBytes.Body.Close()
	var geoinfo IpaddrGeoInfo
	_ = json.Unmarshal(info, &geoinfo)

	if err != nil {
		fmt.Println("Error fetching IP geolocation:", err)
		return
	}
	fmt.Printf("[公网Ip: %s]\n", geoinfo.Query)
	fmt.Printf("[国家:%s %s]\n", geoinfo.Country, geoinfo.City)
	fmt.Printf("[经度: %f 纬度: %f]\n", geoinfo.Lon, geoinfo.Lat)

	// 这里没有实现处理IP-API响应的代码，你需要根据API的响应格式来解析地理位置信息
}
