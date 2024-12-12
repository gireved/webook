package main

type List interface {
	Add(index int, val any)
	Append(val any) error
	Delete(index int) error
}

type LinkedList struct {
	head node
	//Head node
}

func (l *LinkedList) Append(val any) error {
	//TODO implement me
	panic("implement me")
}

func (l *LinkedList) Delete(index int) error {
	//TODO implement me
	panic("implement me")
}

func (l *LinkedList) Add(index int, val any) {
	// 实现这个方法
}

type node struct {
	//next node
	next *node
}

func UseListV1() {
	l := &LinkedList{}
	l.Add(1, 123)
	l.Add(1, "123")
	l.Add(1, nil)
}
