package parser

import (
	"fmt"
	"mayipashu/conf"
	"mayipashu/defs"
	"mayipashu/iface"
	"net/url"
	"regexp"
	"time"
)

//开启多个goroutine 从logC中读取数据进行分析
//done 用来关闭goroutine
//每个解析器都有一个pid，方面以后扩展功能，对单一解析器进行操作
//statistician 是可扩展统计器

//每个Parser 状态用state 表示 0==退出，1==运行，2==暂停运行 3==创建了但尚未运行
const (
	stop    = 0
	running = 1
	suspend = 2   //暂停
	wait    = 3   //默认值
)
type Parser struct {
	logC chan defs.LogData
	done chan struct{}
	pid  int
	//statistician iface.IStatistician   //统计管理器
	Suspended    chan<- int //如果有Parser没有任务就将自己的pid放入这个chan中
	state	     int
}

func NewOneParser(logChan chan defs.LogData, done chan struct{},pid int,Suspended chan<- int) iface.IParser {
	return &Parser{
		logC: logChan,
		done: done,
		pid:  pid,
		Suspended:Suspended,
		state:wait,
	}
}

func (p *Parser) RunOneParser() {
	//将当前Parser改为运行状态
	p.state = running
	fmt.Println("[START] Running a New Parser PID:",p.GetPID())
	for {
		select {
		case d, ok := <-p.logC:
			//当前Parser处于状态1时，才正常进行log解析
			if !ok {
				return
			}
			if p.state == running {
				digdata := p.parseData(d)
				if digdata.Url == ""  {
					//数据处理异常舍弃处理下一个数据
					continue
				}
				p.WriteUrlDataToStatistician(digdata)
				continue
			}
			//下面这段代码是为了防止将要停止的Parser继续拿到数据，防止数据丢失
			//当前Parser已经处于非正常运行状态，将拿到的数据放回到logC中并暂停当前Parser
			p.logC <- d
			return

		case <-p.done:
			p.state = suspend
			return

		case <- time.After(conf.LogConfObj.AfterTimeDelWorker*time.Second):
			//将当前Parser状态改为暂停状态
			p.state = 1
			//告知Parser管理器当前Parser处于空闲状态可以暂停
			p.Suspended <- p.pid
			fmt.Println("[INFO] 向Manager 发送暂停请求 PID：",p.GetPID())
		}
	}
}

func (p *Parser) parseData(data defs.LogData) defs.DigData {
	/*
		{127.0.0.1 - - [02/Mar/2020:00:32:1 +0800] "OPTIONS /
		dig?refer=http%3A%2F%2Flocalhost%3A8888%2Fmovie%2F12824.html&time=1&ua=Chorme%2F34.2&url=http%3A%2F%2Flocalhost%3A8888%2Fmovie%2F11257.html
		HTTP/1.1" 200 43 "-" "Chorme/34.2" "-"
		}
	*/
	r := regexp.MustCompile(`(dig.*?)\s`)
	a := r.FindAllString(data.Data, -1)
	//dig?refer=http%3A%2F%2Flocalhost%3A8888%2Fmovie%2F12824.html&time=1&ua=Chorme%2F34.2&url=http%3A%2F%2Flocalhost%3A8888%2Fmovie%2F11257.html
	if len(a) == 0 {
		return defs.DigData{}
	}
	urlInfo, err := url.Parse("http://localhost/?" + a[0])
	if err != nil {
		fmt.Println("[ERROR]url.Parse is Error:", err)
		return defs.DigData{}
	}
	d := urlInfo.Query()
	return	defs.DigData{d.Get("time"),
						d.Get("url"),
						d.Get("refer"),
						d.Get("ua"),}

}

func (p *Parser) StopParser() {
	if p.state == running {
		p.done <- struct{}{}
		p.state = suspend
		fmt.Println("[SUSPEND]暂停一个Parser PID:", p.GetPID())
	}
}

func (p *Parser) WriteUrlDataToStatistician (data defs.DigData) {
	fmt.Println(data)
}

func (p *Parser) GetPID() int {
	return p.pid
}

func (p *Parser) Close () {
	if p.state == running || p.state == suspend {
		p.StopParser()
		p.state = stop
		close(p.done)
		fmt.Println("[STOP]关闭一个Parser PID ：",p.GetPID())
	}

}

func (p *Parser) Activation () {
	fmt.Println("[START]重启Parser  PID：",p.GetPID())
	go p.RunOneParser()
}

func (p *Parser) GetState() int {
	return p.state
}