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

func NewScheduleTask(done chan struct{}) iface.ISchedule {
	interval := conf.LogConfObj.TimeInterval
	lc := NewLogConsumer()
	return &ScheduleTask{
		TimeTick: time.NewTicker(interval * time.Second),
		runner:   lc,
		done:     done,
	}
}

func (s *ScheduleTask) StartWorker() {
	for {
		select {
		case <-s.TimeTick.C:
			go s.runner.Start()
		case <- s.done:
			fmt.Println("[STOP]Schedule task is Stop")
			s.TimeTick.Stop()
			return
		}
	}

}

func (s *ScheduleTask) Start () {
	fmt.Println("[START]Running schedule task...")
	s.StartWorker()
}

func (s *ScheduleTask) GetRunner () iface.ILogConsumer {
	return s.runner
}

func (s *ScheduleTask) GetLogChan () chan defs.LogData {
	return s.GetRunner().GetLogChan()
}