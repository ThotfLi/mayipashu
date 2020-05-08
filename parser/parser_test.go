package parser

import (
	"fmt"
	"mayipashu/defs"
	"mayipashu/distribute"
	"testing"
	"time"
)

func TestMain(m *testing.M){
	m.Run()
}


func TestWrker(t *testing.T){
	//t.Run("newObj",testNewOneParser)
	t.Run("newManager",testManager)
	//t.Run("SuspendTest",testSuspend)
}

func testNewOneParser(t *testing.T) {
	tlogc := make(chan defs.LogData,10)
	done := make(chan struct{})
	pid := 1
	suppendC := make(chan int,10)

	one := NewOneParser(tlogc,done,pid,suppendC)
	go one.RunOneParser()
	//{127.0.0.1 - - [02/Mar/2020:00:32:1 +0800] "OPTIONS /dig?refer=http%3A%2F%2Flocalhost%3A8888%2Fmovie%2F8328.html&time=1&ua=Mpzilla%2F5.0&url=http%3A%2F%2Flocalhost%3A8888%2Fmovie%2F8328.html HTTP/1.1" 200 43 "-" "Mpzilla/5.0" "-"
	//}
	newData := defs.LogData{Data:"OPTIONS /dig?refer=http%3A%2F%2Flocalhost%3A8888%2Fmovie%2F8328.html&time=1&ua=Mpzilla%2F5.0&url=http%3A%2F%2Flocalhost%3A8888%2Fmovie%2F8328.html HTTP/1.1" + "-"+ "Mpzilla/5.0"+ "-"}

	tlogc <- newData
	time.Sleep(2*time.Second)
	fmt.Println("暂停Parser")
	one.StopParser()
	time.Sleep(2 * time.Second)
	fmt.Println("激活Parser")
	one.Activation()
	tlogc <- newData
	time.Sleep(2*time.Second)
	one.StopParser()
}

func testManager(t *testing.T){
	l := distribute.NewLogConsumer()

	m := NewParserManager(l.GetLogChan())
	go m.RunWorkerPool()
	go l.Start()
	time.Sleep(100*time.Second)
	m.StopManager()

}

func testSuspend(t *testing.T) {
	tlogc := make(chan defs.LogData,10)
	done := make(chan struct{})
	pid := 1
	suppendC := make(chan int,10)

	one := NewOneParser(tlogc,done,pid,suppendC)
	go one.RunOneParser()
	time.Sleep(3*time.Second)
}