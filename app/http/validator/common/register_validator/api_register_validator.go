package register_validator

import (
	"edgelog/app/core/container"
	"edgelog/app/global/consts"
	"edgelog/app/http/validator/api/home"
)

//Each business module validator must be registered (initialized), and will be automatically loaded into the container when the program is started
func ApiRegisterValidator() {
	//Create container
	containers := container.CreateContainersFactory()

	//Key registers each module verification in the container in the format of prefix + module + verification action
	var key string

	//Register portal class form parameter validator
	key = consts.ValidatorPrefix + "HomeNews"
	containers.Set(key, home.News{})
}
