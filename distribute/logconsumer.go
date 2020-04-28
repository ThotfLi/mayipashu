package distribute

import (
	"bufio"
	"fmt"
	"io"
	"mayipashu/defs"
	"mayipashu/iface"
	"mayipashu/conf"
	"os"
	"time"
)

//从log日志文件中获取日志，发送到logchan

//日志消费者
type LogConsumer struct{
	logC chan defs.LogData
}

func NewLogConsumer()iface.ILogConsumer {
	return LogConsumer{logC:make(chan defs.LogData,conf.LogConfObj.LogConsumerChanNumer)}
}

func (l LogConsumer) GetLogToChan ()  {
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
				//log日志文件中没有内容暂停5秒继续获取
				time.Sleep(5*time.Second)
				continue
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

func (l LogConsumer) GetLogChan() chan defs.LogData {
	return l.logC
}

func (l LogConsumer) Close () {
	close(l.GetLogChan())
}