package dao

import (
	"edgelog/app/global/variable"
	"errors"
	"fmt"
)

type Config struct {
}

// GetSystemConfig  Get system configuration parameters
func (c *Config) GetSystemConfig(itemName string) (string, error) {
	systemConfig := make(map[string]interface{})
	if err := variable.GormDb.Raw(fmt.Sprintf("SELECT * FROM `config` WHERE item='%s' LIMIT 1", itemName)).
		First(&systemConfig).Error; err != nil {
		return "", err
	}
	if itemValue, ok := systemConfig["value"].(string); ok {
		return itemValue, nil
	}
	return "", errors.New("no found value of " + itemName)
}

func (c *Config) SetSystemConfig(itemName, value string) error {
	if err := variable.GormDb.Table("config").Where("item=?", itemName).
		Update("value", value).Error; err != nil {
		return err
	}
	return nil
}

func (c *Config) AddSystemConfig(itemName, value string) error {
	var err error
	if _, err = c.GetSystemConfig(itemName); err != nil {
		if err = variable.GormDb.
			Exec(fmt.Sprintf("insert into `config` (`item`,`value`) values ('%s','%s')",
				itemName, value)).Error; err != nil {
			return err
		}
	}
	return err
}

func (c *Config) AddAndSetSystemConfig(itemName, value string) error {
	var err error
	if _, err = c.GetSystemConfig(itemName); err != nil {
		if err = variable.GormDb.
			Exec(fmt.Sprintf("insert into `config` (`item`,`value`) values ('%s','%s')",
				itemName, value)).Error; err != nil {
			return err
		}
	} else {
		err = c.SetSystemConfig(itemName, value)
	}
	return err
}

func (c *Config) DeleteSystemConfig(itemName string) error {
	if err := variable.GormDb.
		Exec(fmt.Sprintf("delete from config where item='%s'", itemName)).
		Error; err != nil {
		return err
	}
	return nil
}

func (c *Config) CheckSystemConfigIsExistence(itemName string) bool {
	_, err := c.GetSystemConfig(itemName)
	if err != nil {
		return false
	}
	return true
}
