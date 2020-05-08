package iface

import "mayipashu/defs"

//定义了定时任务接口
//当log日志文件中没有数据时，就会关闭任务，在一段时间后再次开启
type ISchedule interface {
	//任务本体
	StartWorker()

	//开启定时任务
	Start()

	//返回定时任务
	GetRunner () ILogConsumer

	//返回日志获取通道，要交给任务解析goroutine读取任务
	GetLogChan () chan defs.LogData

	//停止定时任务器
	StopSchedule ()
}
