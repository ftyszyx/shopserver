package libs

//事件

import (
	"container/list"
)

type Handler interface {
	invoke()
}

type HandlerFunc func()

func (f HandlerFunc) invoke() {
	f()
}

type Publish struct {
	handles *list.List
}

func NewEvent() *Publish {
	p := new(Publish)
	p.handles = list.New()
	return p
}

func (p *Publish) AddLister(f HandlerFunc) {
	p.handles.PushBack(f)
}

func (p *Publish) Dispactch() {
	for f := p.handles.Front(); f != nil; f = f.Next() {
		f.Value.(HandlerFunc).invoke()
	}
}
