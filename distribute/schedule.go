package distribute

import (
	"fmt"
	"mayipashu/conf"
	"mayipashu/defs"
	"mayipashu/iface"
	"time"
)

//定时任务

type ScheduleTask struct {
	TimeTick *time.Ticker
	runner   iface.ILogConsumer
	done     chan struct{}
}

func NewScheduleTask() iface.ISchedule {
	interval := conf.LogConfObj.TimeInterval
	lc := NewLogConsumer()
	return &ScheduleTask{
		TimeTick: time.NewTicker(interval * time.Second),
		runner:   lc,
		done:     make(chan struct{}),
	}
}

func (s *ScheduleTask) StartWorker() {
	for {
		select {
		case <-s.TimeTick.C:
			if s.runner.GetStatus() == 0 {
				go s.runner.Start()
			}
		case <- s.done:
			fmt.Println("[STOP]Schedule task is Stop")
			s.TimeTick.Stop()
			return
		}
	}

}

func (s *ScheduleTask) Start () {
	fmt.Println("[START]Running schedule task...")
	go s.runner.Start()
	s.runner.SetStatus(1)
	s.StartWorker()
}

func (s *ScheduleTask) GetRunner () iface.ILogConsumer {
	return s.runner
}

func (s *ScheduleTask) GetLogChan () chan defs.LogData {
	return s.GetRunner().GetLogChan()
}

func (s *ScheduleTask) StopSchedule () {
	s.runner.Close()
	s.done <- struct{}{}
}