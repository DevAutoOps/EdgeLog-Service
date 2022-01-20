package event_manage

import (
	"edgelog/app/global/my_errors"
	"edgelog/app/global/variable"
	"strings"
	"sync"
)

//  Define a global event storage variable ， This module is only responsible for storage   key  =>  function  ，  The function is slightly weaker than that of the container ， But the call is simpler 、 convenient 、 quick 
var sMap sync.Map

//  Create an event management factory 
func CreateEventManageFactory() *eventManage {

	return &eventManage{}
}

//  Define an event management structure 
type eventManage struct {
}

//  1. Registration event 
func (e *eventManage) Set(key string, keyFunc func(args ...interface{})) bool {
	// judge key Whether there are events under 
	if _, exists := e.Get(key); exists == false {
		sMap.Store(key, keyFunc)
		return true
	} else {
		variable.ZapLog.Info(my_errors.ErrorsFuncEventAlreadyExists + " ,  Related key name ：" + key)
	}
	return false
}

// 2. Get event 
func (e *eventManage) Get(key string) (interface{}, bool) {
	if value, exists := sMap.Load(key); exists {
		return value, exists
	}
	return nil, false
}

//  3. Execution event 
func (e *eventManage) Call(key string, args ...interface{}) {
	if valueInterface, exists := e.Get(key); exists {
		if fn, ok := valueInterface.(func(args ...interface{})); ok {
			fn(args...)
		} else {
			variable.ZapLog.Error(my_errors.ErrorsFuncEventNotCall + ",  Key name ：" + key + ",  Related functions cannot be called ")
		}

	} else {
		variable.ZapLog.Error(my_errors.ErrorsFuncEventNotRegister + ",  Key name ：" + key)
	}
}

//  4. Delete event 
func (e *eventManage) Delete(key string) {
	sMap.Delete(key)
}

//  5. According to the prefix of the key ， Fuzzy call .  Use with caution .
func (e *eventManage) FuzzyCall(keyPre string) {

	sMap.Range(func(key, value interface{}) bool {
		if keyName, ok := key.(string); ok {
			if strings.HasPrefix(keyName, keyPre) {
				e.Call(keyName)
			}
		}
		return true
	})
}
