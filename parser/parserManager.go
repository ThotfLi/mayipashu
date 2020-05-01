package parser

import (
	"mayipashu/conf"
	"mayipashu/defs"
	"mayipashu/iface"
	"time"
)

type ParserManager struct {
	//任务池从logC中读取要处理的数据
	logC 	chan defs.LogData

	//用来结束任务的chan
	done    chan struct{}
	parsers []iface.IParser
	SuspendedParserPID []uint  //存放所有被挂起的Parser的PID
	Suspended          chan int //当会Parser没有工作，就将自己的pid发送到这个chan

}

func (pm *ParserManager) RunWorkerPool() {
	for i := 0; i < int(conf.LogConfObj.MaxGogroutineNumber); i++ {
		done := make(chan struct{},1)
		p := NewOneParser(pm.ReturnLogChan(), done, i,pm.ReturnSuspendedChan())

		//将Parser添加到管理器
		pm.appendParser(p)
		go p.RunOneParser()
	}
}

func (pm *ParserManager) StopWorkerPool() {
}

//关闭所有解析器并回收资源
func (pm *ParserManager) StopAll() {
	for _,i := range pm.parsers {
		i.StopParser()
		i.Close()
	}
}

func (pm *ParserManager) ReturnLogChan() chan defs.LogData {
	return pm.logC
}

func (pm *ParserManager) ReturnDoneChan() chan struct{} {
	return pm.done
}

func (pm *ParserManager) appendParser(p iface.IParser) {
	pm.parsers[p.GetPID()] = p
}

func (pm *ParserManager) DynamicManagementTaskPool () {
	for{
		select {
		case pid :=<- pm.ReturnSuspendedChan():
			if int(conf.LogConfObj.MinTaskPool) < pm.Len() {
				//向Suspend中添加一个pid
				//暂停pid代表的Parser
			}

		case <- time.After(1*time.Second) :

		}
	}
}

func (pm *ParserManager) ReturnSuspendedChan() chan int {
	return pm.Suspended
}

func (pm *ParserManager) addSuspendParser() {
	//循环队列添加
	pm.
}

func (pm *ParserManager) Len () int {
	return len(pm.parsers)
}
