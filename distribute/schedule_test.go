package distribute

import (
	"mayipashu/conf"
	"testing"
	"time"
)

func testNewScheduleTask(t *testing.T) {
	conf.InitLogConf()
	a := make(chan struct{},1)
	ls := NewScheduleTask(a)
	go ls.Start()
	time.Sleep(3*time.Second)
	a <- struct{}{}
}

func TestMain(m *testing.M){
	m.Run()
}

func TestWork(t *testing.T){
	//t.Run("scheduletask",testNewScheduleTask)
	t.Run("abc",testLogConsumer_GetLogToChan)
}
