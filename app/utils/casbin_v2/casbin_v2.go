package casbin_v2

import (
	"edgelog/app/global/my_errors"
	"edgelog/app/global/variable"
	"errors"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"strings"
	"time"
)

// establish  casbin Enforcer( Actuator )
func InitCasbinEnforcer() (*casbin.SyncedEnforcer, error) {
	var tmpDbConn *gorm.DB
	var Enforcer *casbin.SyncedEnforcer
	switch strings.ToLower(variable.ConfigGormv2Yml.GetString("Gormv2.UseDbType")) {
	case "mysql":
		if variable.GormDb == nil {
			return nil, errors.New(my_errors.ErrorCasbinCanNotUseDbPtr)
		}
		tmpDbConn = variable.GormDb
	default:
	}

	prefix := variable.ConfigYml.GetString("Casbin.TablePrefix")
	tbName := variable.ConfigYml.GetString("Casbin.TableName")

	a, err := gormadapter.NewAdapterByDBUseTableName(tmpDbConn, prefix, tbName)
	if err != nil {
		return nil, errors.New(my_errors.ErrorCasbinCreateAdaptFail)
	}
	modelConfig := variable.ConfigYml.GetString("Casbin.ModelConfig")

	if m, err := model.NewModelFromString(modelConfig); err != nil {
		return nil, errors.New(my_errors.ErrorCasbinNewModelFromStringFail + err.Error())
	} else {
		if Enforcer, err = casbin.NewSyncedEnforcer(m, a); err != nil {
			return nil, errors.New(my_errors.ErrorCasbinCreateEnforcerFail)
		}
		_ = Enforcer.LoadPolicy()
		AutoLoadSeconds := variable.ConfigYml.GetDuration("Casbin.AutoLoadPolicySeconds")
		Enforcer.StartAutoLoadPolicy(time.Second * AutoLoadSeconds)
		return Enforcer, nil
	}
}
