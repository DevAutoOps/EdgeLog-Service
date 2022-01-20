package register_validator

import (
	"edgelog/app/core/container"
	"edgelog/app/global/consts"
	"edgelog/app/http/validator/common/upload_files"
)

//Each business module validator must be registered (initialized), and will be automatically loaded into the container when the program is started
func WebRegisterValidator() {
	//Create container
	containers := container.CreateContainersFactory()

	//Key registers each module verification in the container in the format of prefix + module + verification action
	var key string
	//The Users module form validator is registered in the container in the form of key = value, which facilitates the call in the routing module.

	//File upload
	key = consts.ValidatorPrefix + "UploadFiles"
	containers.Set(key, upload_files.UpFiles{})
}
