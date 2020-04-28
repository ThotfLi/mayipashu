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
}
