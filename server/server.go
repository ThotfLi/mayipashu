package server

import (
	"fmt"
	"mayipashu/defs"
	"mayipashu/distribute"
	"mayipashu/iface"
	"mayipashu/conf"
	"mayipashu/parser"
	"os"
)

type Server struct{
	logConsumer  iface.ILogConsumer
	tickTask     iface.ISchedule
	done         chan struct{}  //退出循环任务chan

}
func NewServer() iface.IServer {
	go conf.InitLogConf()

	//没开启定时任务
	if conf.LogConfObj.TimeInterval <= 0 {
		logCons := distribute.NewLogConsumer()
		return &Server{logConsumer:logCons}
	}
	//开启了定时任务
	done := make(chan struct{},1)
	tickTask := distribute.NewScheduleTask(done)
	return &Server{
		logConsumer: nil,
		tickTask:    tickTask,
		done:done,
	}

}

func (s *Server) Start () {
	s.Serve()
}

func (s *Server) Stop () {
	fmt.Println("[ERROR] Server is stop")
	//退出定时任务
	s.done <- struct{}{}

	//回收chan
	s.logConsumer.Close()

	os.Exit(-1)
}

func (s *Server) Serve () {
	//开启日志消费者
	if s.logConsumer == nil {
		//开启了定时任务
		s.tickTask.GetRunner().SetServeObject(s)
		go s.tickTask.Start()
	}else {
		s.logConsumer.SetServeObject(s)
		go s.logConsumer.Start()
	}

	p := parser.NewOneParser(s.GetLogDataChan(),1)
	p.RunOneParser()

}

func (s *Server)GetLogDataChan () chan defs.LogData {
	if s.logConsumer == nil {
		return s.tickTask.GetLogChan()
	}
	return s.logConsumer.GetLogChan()
}
