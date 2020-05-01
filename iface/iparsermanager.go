package iface

import "mayipashu/defs"

type IParserManager interface {
	//开启任务处理池
	RunWorkerPool()

	//回收任务池
	StopWorkerPool()

	ReturnLogChan() 	chan defs.LogData
	ReturnDoneChan() 	chan struct{}

	//正在运行的Parser数量
	Len()               int

	//添加解析器到切片，管理解析器使用切片[]parser，方面管理worker
	//AppendParser()

	//动态管理任务池
	//动态减少、增加 任务池
	//每个Parser中的select都有一个time.After，当长时间不使用一个Parser就会触发这个计时器，
	//触发了计时器会通过管道发送自己的pid到ParserManager，Manager就会关闭这个parser但是不回收资源，将多余的Parser的ID放在SuspendedParserPID[]uint中进行管控。
	//ParserManager中也会有一个计时器，如果长时间没有解析器发送自己的pid，就证明此时所有任务池处于满载状态，需要增加Parser，
	//增加Parser首先查看SuspendedParserPID[]int 中是否存在已经暂停的Parser，如果SuspendedParserPID[]int中为控则创建新的Parser
	//增加的Parser数量受到配置文件maxParser的管控
	DynamicManagementTaskPool()
}
