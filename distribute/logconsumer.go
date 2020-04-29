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

//日志消费者
type LogConsumer struct{
	logC chan defs.LogData
	s    iface.IServer
}

func NewLogConsumer()iface.ILogConsumer {
	//return &LogConsumer{logC:make(chan defs.LogData,conf.LogConfObj.LogConsumerChanNumer,s)}
	return &LogConsumer{
		logC: make(chan defs.LogData,conf.LogConfObj.LogConsumerChanNumer),
	}
}

func (l *LogConsumer) GetLogToChan ()  {
	f,err := os.Open(conf.LogConfObj.LogFilePath)
	if err != nil {
		fmt.Println("【ERROR】",err)
		panic("")
	}

	r := bufio.NewReader(f)

	//从日志文件中按行获取单条数据

	fmt.Println("【START】Continue to read the log")
	for {
		line,err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if conf.LogConfObj.TimeInterval <= 0 {
					l.s.Stop()
				}
				//没有日志了结束定时任务
				return
			}
			fmt.Println("【ERROR】",err)
			panic("")
		}

		var newLogData defs.LogData
		newLogData.Data = line

		//发送数据到日志处理goroutine
		l.logC <- newLogData
	}

}

func (l *LogConsumer) GetLogChan() chan defs.LogData {
	return l.logC
}

func (l *LogConsumer) Close () {
	close(l.GetLogChan())
}

func (l *LogConsumer) Start(){
	l.GetLogToChan()
}

func (l *LogConsumer) SetServeObject (s iface.IServer) {
	l.s = s
}

