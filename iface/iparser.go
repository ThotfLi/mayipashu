package iface

//解析器负责对
type IParser interface {
	//开启一个解析器
	RunOneParser()

	//关闭一个解析器，但是没有回收资源，因为还可能复用
    StopParser()

	//每个解析器都有一个单独的PID，方便解析器管理器控制
	GetPID() int

	//回收解析器资源
	Close()
}
