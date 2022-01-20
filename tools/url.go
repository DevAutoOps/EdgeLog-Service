package tools

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// obtain URL Medium batch id And analyze
func IdsStrToIdsIntGroup(key string, c *gin.Context) []uint {
	return IdsStrToIdsIntGroupStr(c.Param(key))
}

type Ids struct {
	Ids []int
}

func IdsStrToIdsIntGroupStr(keys string) []uint {
	IDS := make([]uint, 0)
	ids := strings.Split(keys, ",")
	for i := 0; i < len(ids); i++ {
		ID := StrToUint(ids[i])
		IDS = append(IDS, ID)
	}
	return IDS
}
