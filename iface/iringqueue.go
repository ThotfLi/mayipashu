package iface

//循环队列
type IRingQueue interface {
	//向队列添加一个元素
	AddQueue(n int) error
	//从队列获取一个元素
	GetQueue() (int,bool)
	//队列是否为空满了
	IsFull() bool
	//队列是否为空
	IsEmpty() bool
	//返回队列中的所有值
	ShowQueue() []int
	//返回当前队列长度
	Lenght() int
}