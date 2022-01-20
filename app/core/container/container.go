package container

import (
	"edgelog/app/global/my_errors"
	"edgelog/app/global/variable"
	"log"
	"strings"
	"sync"
)

//  Define a global key value pair storage container 

var sMap sync.Map

//  Create a container factory 
func CreateContainersFactory() *containers {
	return &containers{}
}

//  Define a container structure 
type containers struct {
}

//  1. Register the code to the container as a key value pair 
func (c *containers) Set(key string, value interface{}) (res bool) {

	if _, exists := c.KeyIsExists(key); exists == false {
		sMap.Store(key, value)
		res = true
	} else {
		//  Program startup phase ，zaplog  uninitialized ， Use system log Print the exception log that occurs at startup 
		if variable.ZapLog == nil {
			log.Fatal(my_errors.ErrorsContainerKeyAlreadyExists + ", Please solve the problem of duplicate key names , Correlation key ：" + key)
		} else {
			//  Program startup initialization complete 
			variable.ZapLog.Warn(my_errors.ErrorsContainerKeyAlreadyExists + ",  Correlation key ：" + key)
		}
	}
	return
}

//  2. delete 
func (c *containers) Delete(key string) {
	sMap.Delete(key)
}

//  3. Pass key ， Get value from container 
func (c *containers) Get(key string) interface{} {
	if value, exists := c.KeyIsExists(key); exists {
		return value
	}
	return nil
}

//  4.  Determine whether the key is registered 
func (c *containers) KeyIsExists(key string) (interface{}, bool) {
	return sMap.Load(key)
}

//  Delete the contents registered in the container according to the prefix of the key 
func (c *containers) FuzzyDelete(keyPre string) {
	sMap.Range(func(key, value interface{}) bool {
		if keyname, ok := key.(string); ok {
			if strings.HasPrefix(keyname, keyPre) {
				sMap.Delete(keyname)
			}
		}
		return true
	})
}
