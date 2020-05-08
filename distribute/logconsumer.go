package distribute

import (
	"bufio"
	"fmt"
	"io"
	"mayipashu/conf"
	"mayipashu/defs"
	"mayipashu/iface"
	"os"
)

//从log日志文件中获取日志，发送到logchan
//status = 1  任务执行中  status = 0 任务暂停中
//日志消费者
const (
	RUNNING = 1
	SUSPEND = 0 //暂停
)

type LogConsumer struct {
	logC   chan defs.LogData
	s      iface.IServer
	status int
	done   chan struct{}
}

func NewLogConsumer() iface.ILogConsumer {
	//return &LogConsumer{logC:make(chan defs.LogData,conf.LogConfObj.LogConsumerChanNumer,s)}
	return &LogConsumer{
		logC:   make(chan defs.LogData, conf.LogConfObj.LogConsumerChanNumer),
		status: 0,
		done:   make(chan struct{}),
	}
}

func (l *LogConsumer) GetLogToChan() {
	f, err := os.Open(conf.LogConfObj.LogFilePath)
	defer f.Close()
	if err != nil {
		fmt.Println("【ERROR】", err)
		panic("")
	}

	r := bufio.NewReader(f)

	//从日志文件中按行获取单条数据

	fmt.Println("【START】Continue to read the log")
Loop:
	for {
		select {
		case <-l.getDoneC():
			break Loop
		default:
			line, err := r.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					//没开启定时任务
					if conf.LogConfObj.TimeInterval <= 0 {
						l.s.Stop()
					}
					//没有日志了结束定时任务
					l.status = 0
					return
				}
				fmt.Println("【ERROR】", err)

			}

			var newLogData defs.LogData
			newLogData.Data = line

			//发送数据到日志处理goroutine
			l.logC <- newLogData
		}
	}

}

func (l *LogConsumer) GetLogChan() chan defs.LogData {
	return l.logC
}

func (l *LogConsumer) Close() {
	l.getDoneC() <- struct{}{}
	close(l.GetLogChan())
}

func (l *LogConsumer) Start() {
	l.GetLogToChan()
}

func (l *LogConsumer) SetServeObject(s iface.IServer) {
	l.s = s
}

func (l *LogConsumer) GetStatus() int {
	return l.status
}

func (l *LogConsumer) SetStatus(status int) {
	l.status = status
}

func (l *LogConsumer) getDoneC() chan struct{} {
	return l.done
}
