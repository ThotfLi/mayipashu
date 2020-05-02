package parser

import (
	"errors"
	"mayipashu/iface"
)

type RingQueue struct {
	q       []int
	font    int  //指向第一个元素
	real    int  //指向最后一个元素的后一个位置
	maxSize uint //队列支持存储的最大值
}

func NewRingQueue(maxsize uint) iface.IRingQueue {
	return &RingQueue{
		q:       make([]int, maxsize),
		font:    0,
		real:    0,
		maxSize: maxsize,
	}
}

func (eue *RingQueue) AddQueue(n int) error {
	if eue.IsFull() {
		return errors.New("Queue is full")
	}

	eue.returnQueue()[eue.returnLastIndex()] = n
	eue.real = (eue.returnLastIndex() + 1) % eue.returnMaxSize()
	return nil
}
func (eue *RingQueue) GetQueue() (int, bool) {
	if eue.IsEmpty() {
		return 0, false
	}

	n := eue.returnQueue()[eue.returnOneIndex()]
	eue.font = (eue.returnOneIndex() + 1) % eue.returnMaxSize()
	return n, true
}
func (eue *RingQueue) IsFull() bool {
	//满了
	if eue.returnLastIndex()+1 == eue.returnOneIndex() {
		return true
	}
	return false
}

func (eue *RingQueue) IsEmpty() bool {
	//为空
	if eue.returnLastIndex() == eue.returnOneIndex() {
		return true
	}

	return false
}
func (eue *RingQueue) ShowQueue() []int {
	newlist := make([]int, eue.returnMaxSize())
	copy(newlist, eue.returnQueue())
	return newlist
}

func (eue *RingQueue) returnLastIndex() int {
	return eue.real
}

func (eue *RingQueue) returnOneIndex() int {
	return eue.font
}

func (eue *RingQueue) returnMaxSize() int {
	return int(eue.maxSize)
}

func (eue *RingQueue) returnQueue() []int {
	return eue.q
}

func (eue *RingQueue) Lenght () int {
	return len(eue.q)
}