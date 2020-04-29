package iface

import "mayipashu/defs"

type  IServer interface {
	Start()

	Stop()

	//返回日志获取通道，要交给任务解析goroutine读取任务
	GetLogDataChan () chan defs.LogData
}
