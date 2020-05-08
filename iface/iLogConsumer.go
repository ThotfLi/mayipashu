package iface

import (
	"mayipashu/defs"
)

type ILogConsumer interface {
	//从log文件中获取数据发送到chan7··那你，。、
	GetLogToChan()

	//返回channel
	GetLogChan() chan defs.LogData

	//关闭Chan
	Close ()

	//开始任务
	Start()

	//在当前对象中设置server对象
	SetServeObject (s IServer)

	//任务状态
	GetStatus ()  int
	SetStatus(status int)
}
