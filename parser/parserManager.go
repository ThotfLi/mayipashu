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
	//SuspendedParserPID chan int  //存放所有被挂起的Parser的PID
	Suspended          chan int //当会Parser没有工作，就将自己的pid发送到这个chan
	ringQueue          iface.IRingQueue

}

//运行任务池，开启N个worker
//每开启一个worker将worker添加到管理器
func (pm *ParserManager) RunWorkerPool() {
	for i := 0; i < int(conf.LogConfObj.MinTaskPool); i++ {
		done := make(chan struct{},1)
		p := NewOneParser(pm.ReturnLogChan(), done, i,pm.ReturnSuspendedChan())

		//将Parser添加到管理器
		pm.appendParser(p)
		go p.RunOneParser()
	}
}

func (pm *ParserManager) StopOneWorker(pid int) {
	(pm.parsers[pid]).StopParser()
}

//关闭所有worker并回收资源
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

//动态管理任务池
//监听SuspendedChan 有worker将自己的id发进来就判断是否需要暂停
//如果当前worker过少就不暂停这个任务池
func (pm *ParserManager) DynamicManagementTaskPool () {
	for{
		select {
		case pid :=<- pm.ReturnSuspendedChan():
			if int(conf.LogConfObj.MinTaskPool) < pm.Len() {
				//向Queue中添加一个pid
				pm.returnRingQueue().AddQueue(pid)

				//暂停pid代表的Parser
				pm.StopOneWorker(pid)
			}

			//满载时添加worker
		case <- time.After(1*time.Second) :
			//如果队列中有暂停的worker 那么直接从队列中拿到pid并激活
			pid,ok := pm.returnRingQueue().GetQueue()
			if ok {
				pm.parsers[pid].RunOneParser()
				continue
			}
			//队列中为空则创建新的worker，但是要保证worker数量小于conf中的MaxTaskPool
			if pm.Len() < int(conf.LogConfObj.MaxTaskPool) {
				done := make(chan struct{},1)

				newWorker := NewOneParser(pm.ReturnLogChan(), done, pm.Len(),pm.ReturnSuspendedChan())
				//将worker添加到管理器
				pm.appendParser(newWorker)

				go newWorker.RunOneParser()
			}
		}
	}
}

func (pm *ParserManager) ReturnSuspendedChan() chan int {
	return pm.Suspended
}

//返回正在运行的worker数量
func (pm *ParserManager) Len () int {
	return len(pm.parsers) - len(pm.ReturnLogChan())
}

func (pm *ParserManager) returnRingQueue () iface.IRingQueue {
	return pm.ringQueue
}
