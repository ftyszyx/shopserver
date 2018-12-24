package db

import (
	"sync"
)

//DBManger 数据库集
var DBManger = &dbmangerType{dict: make(map[string]*DBInfo)}

type dbmangerType struct {
	mux  sync.RWMutex
	dict map[string]*DBInfo //所有的数据库
}

func (self *dbmangerType) add(name string, info *DBInfo) (added bool) {
	self.mux.Lock()
	defer self.mux.Unlock()
	if _, ok := self.dict[name]; !ok {
		self.dict[name] = info
		added = true
	}
	return
}

func (self *dbmangerType) get(name string) (info *DBInfo, ok bool) {
	self.mux.RLock()
	defer self.mux.RUnlock()
	info, ok = self.dict[name]
	return
}

func (self *dbmangerType) getDefault() (info *DBInfo) {
	info, _ = self.get("default")
	return
}
