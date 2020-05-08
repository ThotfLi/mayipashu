package parser

import (
	"fmt"
	"mayipashu/conf"
	"mayipashu/defs"
	"mayipashu/iface"
	"time"
)

const (
	WAIT    = 2
	RUNNING = 1
	STOP    = 0
)

type ParserManager struct {
	//任务池从logC中读取要处理的数据
	logC chan defs.LogData

	//用来结束任务的chan
	done    chan struct{}
	parsers []iface.IParser
	//SuspendedParserPID chan int  //存放所有被挂起的Parser的PID
	Suspended chan int //当会Parser没有工作，就将自己的pid发送到这个chan
	//ringQueue是循环队列，其中存放暂停中的ParserID，可以对单一进
	ringQueue iface.IRingQueue
	count     int //已管理的Parser 数量
	status    int //管理器状态
}

func NewParserManager(logC chan defs.LogData) iface.IParserManager {
	done := make(chan struct{})
	Supend := make(chan int, conf.LogConfObj.MaxTaskPool)
	return &ParserManager{
		logC:      logC,
		done:      done,
		parsers:   make([]iface.IParser, conf.LogConfObj.MaxTaskPool),
		Suspended: Supend,
		ringQueue: NewRingQueue(uint(conf.LogConfObj.MaxGogroutineNumber)),
		count:     0,
	}
}

//运行任务池，开启N个worker
//每开启一个worker将worker添加到管理器
func (pm *ParserManager) RunWorkerPool() {

	//加载动态任务池管理器
	go pm.DynamicManagementTaskPool()

	//conf.LogConfObj.MinTaskPool
	for i := 0; i < int(conf.LogConfObj.MinTaskPool); i++ {
		done := make(chan struct{}, 1)
		p := NewOneParser(pm.ReturnLogChan(), done, i, pm.ReturnSuspendedChan())

		pm.count += 1
		//将Parser添加到管理器
		pm.appendParser(p)
		go p.RunOneParser()
	}
}

func (pm *ParserManager) StopOneWorker(pid int) {
	(pm.parsers[pid]).StopParser()
}

func (pm *ParserManager) exitOneWorker(pid int) {
	//在RingQueue中有些Parser是处于暂停状态，有些是运行状态所以需要做判断避免重复暂停

	if (pm.parsers[pid]).GetState() != suspend {
		(pm.parsers[pid]).StopParser()
	}

	(pm.parsers[pid]).Close()
}

//关闭所有worker并回收资源
func (pm *ParserManager) StopAllWorker() {
	var qLst []int

	if !pm.returnRingQueue().IsEmpty() {
		qLst = pm.returnRingQueue().ShowQueue()

		//退出所有Worker
		for _, v := range qLst {
			pm.exitOneWorker(v)
		}
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
func (pm *ParserManager) DynamicManagementTaskPool() {
	for {
		select {
		//空闲，减少一个worker
		case pid := <-pm.ReturnSuspendedChan():

			if int(conf.LogConfObj.MinTaskPool) < pm.ParserCount() {
				//向Queue中添加一个pid
				pm.returnRingQueue().AddQueue(pid)

				//暂停pid代表的Parser
				pm.StopOneWorker(pid)
			}
			//满载，添加一个worker
		case <-time.After(conf.LogConfObj.AfterTimeAddWorker * time.Second):
			//如果队列中有暂停的worker 那么直接从队列中拿到pid并激活
			pid, ok := pm.returnRingQueue().GetQueue()
			if ok {
				pm.parsers[pid].RunOneParser()
				continue
			}
			//队列中为空则创建新的worker，但是要保证worker数量小于conf中的MaxTaskPool
			if pm.ParserCount() < int(conf.LogConfObj.MaxTaskPool) {
				done := make(chan struct{}, 1)

				newWorker := NewOneParser(pm.ReturnLogChan(), done, pm.ParserCount(), pm.ReturnSuspendedChan())
				//将worker添加到管理器
				pm.appendParser(newWorker)
				pm.count += 1
				go newWorker.RunOneParser()
			}
		//关闭动态管理任务池
		case <-pm.done:
			return
		}
	}
}

func (pm *ParserManager) ReturnSuspendedChan() chan int {
	return pm.Suspended
}

func (pm *ParserManager) returnRingQueue() iface.IRingQueue {
	return pm.ringQueue
}

func (pm *ParserManager) StopManager() {
	fmt.Println("[STOP] Stop Manager")
	//如果循环列表中不为空，则拿出所有的PID，一个个STOP
	pm.done <- struct{}{}
	close(pm.done)

	pm.parsers = pm.parsers[:pm.count]
	for _, v := range pm.parsers {
		pm.exitOneWorker(v.GetPID())
	}
	pm.parsers = nil
	close(pm.ReturnSuspendedChan())
}

//正在运行的worker数量
func (pm *ParserManager) ParserCount() int {
	return pm.count
}
