package tools

import (
	"encoding/json"
	"fmt"
)

func InterfaceToStr(modify interface{}) string {
	modityStr := ""
	if modityStr == "null" || modityStr == "nil" {
		modityStr = ""
	}
	modityJson, err := json.Marshal(modify)
	if err == nil {
		modityStr = fmt.Sprintf("%s", modityJson)
	}
	return modityStr
}
