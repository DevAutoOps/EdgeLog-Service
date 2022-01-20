package yml_config

import (
	"edgelog/app/core/container"
	"edgelog/app/global/my_errors"
	"edgelog/app/global/variable"
	"edgelog/app/utils/yml_config/ymlconfig_interf"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"time"
)

//  because  vipver  The package itself has a for file change events bug， Related events are called back twice 
//  It has not been completely solved for many years ， dependent  issue  detailed list ：https://github.com/spf13/viper/issues?q=OnConfigChange
//  Set an internal global variable ， Record the point in time when the profile changes ， If the difference between two callback events is less than 1 second ， We think it is the second callback event ， Instead of modifying the configuration file manually 
//  This avoids  vipver  This bag bug

var lastChangeTime time.Time

func init() {
	lastChangeTime = time.Now()
}

//  Create a yaml Profile Factory 
//  The parameter is set to the file name of the variable parameter ， In this way, the parameters do not need to be passed ， If more than one is passed ， We only take the first parameter as the configuration file name 
func CreateYamlFactory(fileName ...string) ymlconfig_interf.YmlConfigInterf {

	yamlConfig := viper.New()
	//  Directory of configuration file 
	yamlConfig.AddConfigPath(variable.BasePath + "/config")
	//  File name to read , Default to ：config
	if len(fileName) == 0 {
		yamlConfig.SetConfigName("config")
	} else {
		yamlConfig.SetConfigName(fileName[0])
	}
	// Set profile type ( suffix ) by  yml
	yamlConfig.SetConfigType("yml")

	if err := yamlConfig.ReadInConfig(); err != nil {
		log.Fatal(my_errors.ErrorsConfigInitFail + err.Error())
	}

	return &ymlConfig{
		yamlConfig,
	}
}

type ymlConfig struct {
	viper *viper.Viper
}

// Listening for file changes 
func (y *ymlConfig) ConfigFileChangeListen() {
	y.viper.OnConfigChange(func(changeEvent fsnotify.Event) {
		if time.Now().Sub(lastChangeTime).Seconds() >= 1 {
			if changeEvent.Op.String() == "WRITE" {
				y.clearCache()
				lastChangeTime = time.Now()
			}
		}
	})
	y.viper.WatchConfig()
}

//  Determine whether the phase key has been cached 
func (y *ymlConfig) keyIsCache(keyName string) bool {
	if _, exists := container.CreateContainersFactory().KeyIsExists(variable.ConfigKeyPrefix + keyName); exists {
		return true
	} else {
		return false
	}
}

//  Cache key values 
func (y *ymlConfig) cache(keyName string, value interface{}) bool {
	return container.CreateContainersFactory().Set(variable.ConfigKeyPrefix+keyName, value)
}

//  Get cached value by key 
func (y *ymlConfig) getValueFromCache(keyName string) interface{} {
	return container.CreateContainersFactory().Get(variable.ConfigKeyPrefix + keyName)
}

//  Clear the configuration item information that has been changed 
func (y *ymlConfig) clearCache() {
	container.CreateContainersFactory().FuzzyDelete(variable.ConfigKeyPrefix)
}

//  allow  clone  A structure with the same function 
func (y *ymlConfig) Clone(fileName string) ymlconfig_interf.YmlConfigInterf {
	//  There is a deep copy here ， Need attention ， Avoid the impact of copied structure operations on the original structure 
	var ymlC = *y
	var ymlConfViper = *(y.viper)
	(&ymlC).viper = &ymlConfViper

	(&ymlC).viper.SetConfigName(fileName)
	if err := (&ymlC).viper.ReadInConfig(); err != nil {
		variable.ZapLog.Error(my_errors.ErrorsConfigInitFail, zap.Error(err))
	}
	return &ymlC
}

// Get  A raw value 
func (y *ymlConfig) Get(keyName string) interface{} {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName)
	} else {
		value := y.viper.Get(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetString
func (y *ymlConfig) GetString(keyName string) string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(string)
	} else {
		value := y.viper.GetString(keyName)
		y.cache(keyName, value)
		return value
	}

}

// GetBool
func (y *ymlConfig) GetBool(keyName string) bool {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(bool)
	} else {
		value := y.viper.GetBool(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt
func (y *ymlConfig) GetInt(keyName string) int {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int)
	} else {
		value := y.viper.GetInt(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt32
func (y *ymlConfig) GetInt32(keyName string) int32 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int32)
	} else {
		value := y.viper.GetInt32(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt64
func (y *ymlConfig) GetInt64(keyName string) int64 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int64)
	} else {
		value := y.viper.GetInt64(keyName)
		y.cache(keyName, value)
		return value
	}
}

// float64
func (y *ymlConfig) GetFloat64(keyName string) float64 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(float64)
	} else {
		value := y.viper.GetFloat64(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetDuration
func (y *ymlConfig) GetDuration(keyName string) time.Duration {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(time.Duration)
	} else {
		value := y.viper.GetDuration(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetStringSlice
func (y *ymlConfig) GetStringSlice(keyName string) []string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).([]string)
	} else {
		value := y.viper.GetStringSlice(keyName)
		y.cache(keyName, value)
		return value
	}
}
