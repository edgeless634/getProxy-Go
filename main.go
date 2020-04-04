package main

import (
	"fmt"
	"net/http"
	"net/url"
	htmlquery "github.com/antchfx/xquery/html"
	"io/ioutil"
	"io"
	"strings"
	"time"
	"log"
	"sync"
	"math/rand"
	"os"
)

//MaxPages 最大页数
const MaxPages = 20
//MaxTestedProxySubNum TestProxySub的最大协程数
const MaxTestedProxySubNum = 256
var allUserAgent = [5]string{ //allUserAgent 全部UA
	"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.94 Safari/537.36", 
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.93 Safari/537.36",
	"Mozilla/5.0 (Linux; Android 4.0.4; Galaxy Nexus Build/IMM76B) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.133 Mobile Safari/535.19",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 8_0_2 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12A366 Safari/600.1.4",
	"Mozilla/5.0 (Android; Mobile; rv:14.0) Gecko/14.0 Firefox/14.0",
}

func removeSpace(s string) string {
	/*
	移除字符串前后的tab和回车
	*/
	l, r := 0, len(s)-1
	for l < len(s) && (s[l] == '\t' || s[l] == '\n'){
		l++
	}
	for r >= 0 && (s[r] == '\t' || s[r] == '\n') {
		r--
	}
	return s[l:r+1]
}

func getClient(proxy string) *http.Client {
	/*
	获得一个自动配置的http.Client

	不使用代理时应该让proxy = ""
	*/
	if proxy != "" {
		urli := url.URL{}
		urlproxy, _ := urli.Parse(proxy)
		c := http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlproxy),
			},
			Timeout: time.Second * 10,
		}
		return &c
	}
	return &http.Client{Timeout: time.Second * 10}
}

func getHTML(url, proxy string) string {
	/*
	获得一个url的body
	*/
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Add("User-Agent", allUserAgent[rand.Intn(5)])
	client := getClient(proxy)
	resp, err := client.Do(req)
    if err != nil {
        log.Fatalln(err)
    }
    defer resp.Body.Close()
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil && data == nil {
        log.Fatalln(err)
    }
    return fmt.Sprintf("%s", data)
}


//以下是n个免费代理网站的爬虫
func getProxy1(allIPChannel chan string, wg *sync.WaitGroup){
	html := getHTML("https://www.kuaidaili.com/free/", "")
	root, _ := htmlquery.Parse(strings.NewReader(html))
	tr := htmlquery.Find(root, "//*[@id='list']/table/tbody/tr")
	for _, row := range tr {
		item := htmlquery.Find(row, ".//td")
		ip := htmlquery.InnerText(item[0])
		port := htmlquery.InnerText(item[1])
		p := ip + ":" + port
		allIPChannel <- p
	}
	wg.Done()
}

func getProxy2(allIPChannel chan string, wg *sync.WaitGroup){
	for i := 1; i <= MaxPages; i++ {
		urlNow := ""
		if i == 1{
			urlNow = "http://ip.yqie.com/proxygaoni/index.htm"
		} else {
			urlNow = fmt.Sprintf("http://ip.yqie.com/proxygaoni/index_%d.htm", i)
		}
		html := getHTML(urlNow, "")
		root, _ := htmlquery.Parse(strings.NewReader(html))
		//每一行
		tr := htmlquery.Find(root, "//*[@id='GridViewOrder']/tbody/tr[position()>=2]")
		for _, row := range tr {
			item := htmlquery.Find(row, ".//td")
			ip := htmlquery.InnerText(item[1]) //ip
			port := htmlquery.InnerText(item[2]) //port
			p := ip + ":" + port
			allIPChannel <- p
		}
	}
	wg.Done()
}

func getProxy3(allIPChannel chan string, wg *sync.WaitGroup){
	for i := 1; i <= MaxPages; i++ {
		html := getHTML(fmt.Sprintf("http://www.66ip.cn/%d.html", i), "")
		root, _ := htmlquery.Parse(strings.NewReader(html))
		tr := htmlquery.Find(root, "//*[@id='main']/div/div[1]/table/tbody/tr[position()>=2]")
		for _, row := range tr {
			item := htmlquery.Find(row, ".//td")
			ip := htmlquery.InnerText(item[0])
			port := htmlquery.InnerText(item[1])
			p := ip + ":" + port
			allIPChannel <- p
		}
	}
	wg.Done()
}

func getProxy4(allIPChannel chan string, wg *sync.WaitGroup){
	html := getHTML("http://m.feizhuip.com/Index/article/id/470.html", "")
	root, _ := htmlquery.Parse(strings.NewReader(html))
	tr := htmlquery.Find(root, "/html/body/div[3]/div[2]/div/table/tbody/tr")
	for _, row := range tr {
		item := htmlquery.Find(row, ".//td")
		ip := htmlquery.InnerText(item[0])
		port := htmlquery.InnerText(item[1])
		p := ip + ":" + port
		allIPChannel <- p
	}
	wg.Done()
}

func getProxy5(allIPChannel chan string, wg *sync.WaitGroup){
	html := getHTML("http://www.xiladaili.com/gaoni/", "")
	root, _ := htmlquery.Parse(strings.NewReader(html))
	tr := htmlquery.Find(root, "/html/body/div/div[3]/div[2]/table/tbody/tr")
	for _, row := range tr {
		item := htmlquery.Find(row, ".//td")
		p := htmlquery.InnerText(item[0])
		allIPChannel <- p
	}
	wg.Done()
}

func getProxy6(allIPChannel chan string, wg *sync.WaitGroup){
	html := getHTML("http://www.89ip.cn/", "")
	root, _ := htmlquery.Parse(strings.NewReader(html))
	tr := htmlquery.Find(root, "//tbody/tr")
	for _, row := range tr {
		item := htmlquery.Find(row, ".//td")
		ip := removeSpace(htmlquery.InnerText(item[0]))
		port := removeSpace(htmlquery.InnerText(item[1]))
		p := ip + ":" + port
		allIPChannel <- p
	}
	wg.Done()
}

func getProxy(allIPChannel chan string){
	/*
	启动并管理n个网站的爬虫
	*/
	wg := sync.WaitGroup{}
	//启动n个不同网站的爬虫
	wg.Add(6)
	go getProxy1(allIPChannel, &wg)
	go getProxy2(allIPChannel, &wg)
	go getProxy3(allIPChannel, &wg)
	go getProxy4(allIPChannel, &wg)
	go getProxy5(allIPChannel, &wg)
	go getProxy6(allIPChannel, &wg)
	//等待爬虫完成
	wg.Wait()
	close(allIPChannel)
}

func testProxy(ipPort string) bool {
	/*
	测试代理ipPort是否有效

	ipPort : "xxx.xxx.xxx.xxx:xxxx"
	*/
	url := "https://www.baidu.com"
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3776.0 Safari/537.36")
	client := getClient("http://" + ipPort)
	resp, err := client.Do(req)
    if err != nil {
        return false
    }
    defer resp.Body.Close()
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil && data == nil {
        return false
    }
    return true
}

func testedProxySub(allIPChannel, channel chan string, wg *sync.WaitGroup) {
	/*
	把allIPChannel中的可用代理传入channel
	*/
	for ip := range allIPChannel{
		success := 0
		for i := 1; i <= 5; i++ {
			if testProxy(ip) {
				success++
			}
		}
		if success >= 2 {
			channel <- ip
		}
	}
	wg.Done()
}

func testedProxy(channel chan string) {
	/*
	向channel传入可用的代理
	*/
	wg := sync.WaitGroup{}
	allIPChannel := make(chan string, 10)
	go getProxy(allIPChannel)
	for i := 1; i <= MaxTestedProxySubNum; i++ {
		wg.Add(1)
		go testedProxySub(allIPChannel, channel, &wg)
	}
	wg.Wait()
	close(channel)
}

func main() {
	channel := make(chan string, 10)
	f, err := os.OpenFile("ips.txt", os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	go testedProxy(channel)
	for i := range channel {
		fmt.Println(i)
		_, err := io.WriteString(f, i + "\n")
		if err != nil {
			log.Fatalln(err)
		}
	}
}