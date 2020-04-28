package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var uaList = []string{
	"Mpzilla/5.0",
	"IE/10",
	"Chorme/34.2",
	"Safari/2.12",
}

// 资源结构体，存放模拟url的生成配置
type resource struct {
	url string   // 用户访问页面的url
	target string    // 用户访问页面的具体资源对象，id表示
	start int    // 资源数的起始量
	end int    // 资源数的结束量
}

/*
* 获取初始化的url生成配置，以便生成用户访问日志数据
 */
func ruleResource() []resource {
	var res []resource
	r1 := resource{
		url:    "http://localhost:8888",
		target: "",
		start: 0,
		end: 0,
	}
	r2 := resource{
		url:    "http://localhost:8888/list/{$id}.html",
		target: "{$id}",
		start:  1,
		end:    21,
	}
	r3 := resource{
		url:    "http://localhost:8888/movie/{$id}.html",
		target: "{$id}",
		start:  1,
		end:    12924,
	}
	res = append(res,r1,r2,r3)
	return res
}

/*
* 根据url生成配置，构造用户访问的url，返回生成好的切片
 */
func buildUrl(res []resource) []string {
	var list []string
	for _,resItem := range res {
		if len(resItem.target) == 0 {    // 如果访问的是首页
			list = append(list,resItem.url)
		} else {
			for i:=resItem.start;i<=resItem.end;i++{
				urlStr := strings.Replace(resItem.url,resItem.target,strconv.Itoa(i),-1)
				list = append(list,urlStr)
			}
		}
	}
	return list
}

/*
* 生成具体的日志数据
 */
func makeLog( current, refer, ua string) string {
	// 生成url query部分字符串
	u := url.Values{}
	u.Set( "time","1")
	u.Set( "url",current)
	u.Set( "refer",refer)
	u.Set( "ua",ua)
	paramsStr := u.Encode()

	logTemplate := "127.0.0.1 - - [02/Mar/2020:00:32:1 +0800] \"OPTIONS /dig?{$paramsStr} HTTP/1.1\" 200 43 \"-\" \"{$ua}\" \"-\" "
	// 替换掉模板的 $paramsStr 和 ua 部分
	log := strings.Replace(logTemplate,"{$paramsStr}",paramsStr,-1)
	log = strings.Replace(log,"{$ua}",ua,-1)
	return log
}

/*
* 获取随机数
 */
func randInt(min, max int) int {
	r := rand.New( rand.NewSource( time.Now().UnixNano()))
	if min > max {
		return max
	}
	return r.Intn(max-min)+min
}


func main() {
	// 通过命令行收集参数 total-要创建的日志行数，filepath-要保存的日志文件路径
	total := flag.Int("total",100,"rows be created")
	filePath := flag.String("filePath","./dig.log","file path")
	flag.Parse()

	// 需要构造出真实的网站url集合
	res := 	ruleResource()
	list := buildUrl( res )

	// 随机取currentUrl，referUrl，ua，循环拼接 total 行日志
	logStr := ""
	for i := 0; i < *total; i++{
		currentUrl := list[ randInt(0, len(list)-1) ]
		referUrl := list[ randInt(0, len(list)-1) ]
		ua := uaList[ randInt(0, len(uaList)-1) ]
		logStr = logStr + makeLog( currentUrl, referUrl, ua ) + "\n"
		//ioutil.WriteFile(*filePath,[]byte(logStr),0644)
	}
	// 写日志
	fd,_ := os.OpenFile(*filePath,os.O_RDWR|os.O_APPEND|os.O_CREATE,0644)
	fd.Write([]byte( logStr ))
	fd.Close()

	//按照要求，生成total行日志内容，源自上面的这个集合
	fmt.Println("done.\n")
}
