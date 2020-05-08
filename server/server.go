package server

import (
	"fmt"
	"mayipashu/conf"
	"mayipashu/defs"
	"mayipashu/distribute"
	"mayipashu/iface"
	"mayipashu/parser"
	"os"
)

type Server struct{
	logConsumer  iface.ILogConsumer
	tickTask     iface.ISchedule
	DynamicWorkerPoolManager  iface.IParserManager

}
func NewServer() iface.IServer {
	go conf.InitLogConf()

	//没开启定时任务
	if conf.LogConfObj.TimeInterval <= 0 {
		logCons := distribute.NewLogConsumer()
		return &Server{logConsumer:logCons,
						DynamicWorkerPoolManager:parser.NewParserManager(logCons.GetLogChan())}
	}
	//开启了定时任务
	tickTask := distribute.NewScheduleTask()
	return &Server{
		logConsumer: nil,
		tickTask:    tickTask,

	}

}

func (s *Server) Start () {
	s.Serve()
}

func (s *Server) Stop () {
	fmt.Println("[STOP] Server is stop")
	//退出定时任务
	if s.tickTask != nil {
		s.tickTask.StopSchedule()
	}else {
		s.logConsumer.Close()
	}

	//回收chan

	//退出并回收任务管理器
	s.DynamicWorkerPoolManager.StopManager()

	os.Exit(0)
}

func (s *Server) Serve () {
	//在NewServer中，如果没开定时任务就会在NewSserver中设置LogConsumer
	//开启日志消费者
	if s.logConsumer == nil {
		//开启了定时任务
		s.tickTask.GetRunner().SetServeObject(s)
		go s.tickTask.Start()

		go func() {
			s.DynamicWorkerPoolManager = parser.NewParserManager(s.tickTask.GetLogChan())
			s.DynamicWorkerPoolManager.RunWorkerPool()
		}()

	}else {
		s.logConsumer.SetServeObject(s)
		go s.logConsumer.Start()
	}

}

func (s *Server)GetLogDataChan () chan defs.LogData {
	if s.logConsumer == nil {
		return s.tickTask.GetLogChan()
	}
	return s.logConsumer.GetLogChan()
}
