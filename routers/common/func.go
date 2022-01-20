package common

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StructToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func PrintStruct(a interface{}) {
	typ := reflect.TypeOf(a)
	val := reflect.ValueOf(a)
	kd := val.Kind()
	if kd != reflect.Struct {
		fmt.Println("expect struct")
		return
	}
	num := val.NumField()
	for i := 0; i < num; i++ {
		tagVal := typ.Field(i).Tag.Get("json")
		if tagVal != "" {
			fmt.Printf("%v:", tagVal)
		}
		fmt.Printf("%v ", val.Field(i))
	}
	fmt.Println()
}

func InStrList(target string, array []string) bool {
	for _, element := range array {
		if target == element {
			return true
		}
	}
	return false
}

func StandardManageList(c *gin.Context, f func() (*gorm.DB, interface{}, error)) {
	db, result, err := f()
	if err != nil {
		Error(c, err)
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		Error(c, err)
		return
	}
	pageNum, err := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	if err != nil {
		Error(c, err)
		return
	}
	var count int64
	if err = db.Debug().Count(&count).
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(result).Error; err != nil {
		Error(c, err)
		return
	}
	Ok(c, PageRes{List: result, Count: int(count), PageIndex: pageNum, PageSize: pageSize})
}

func StructToMapFilter(s interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	elem := reflect.ValueOf(s).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		if relType.Field(i).Name == "Model" {
			continue
		}
		m[relType.Field(i).Name] = elem.Field(i).Interface()
	}
	return m
}

func NotImplemented(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 500,
		"msg":  "not implemented",
	})
	c.Abort()
}

func Error(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"code": 500,
		"msg":  err.Error(),
	})
	c.Abort()
}

func Ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

type PageRes struct {
	List      interface{} `json:"rows"`
	Count     int         `json:"total"`
	PageIndex int         `json:"pageNum"`
	PageSize  int         `json:"pageSize"`
}
