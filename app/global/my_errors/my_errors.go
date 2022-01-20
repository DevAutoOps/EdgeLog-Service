package my_errors

const (
	// System part 
	ErrorsContainerKeyAlreadyExists string = " The key is already registered in the container "
	ErrorsPublicNotExists           string = "public  directory does not exist "
	ErrorsConfigYamlNotExists       string = "config.yml  Profile does not exist "
	ErrorsConfigGormNotExists       string = "gorm_v2.yml  Profile does not exist "
	ErrorsStorageLogsNotExists      string = "storage/logs  directory does not exist "
	ErrorsConfigInitFail            string = " Error initializing configuration file "
	ErrorsSoftLinkCreateFail        string = " Automatic creation of soft connection failed , Please run the client as an administrator ( The development environment is goland etc. ， Production environment check command executor permission )"
	ErrorsSoftLinkDeleteFail        string = " Failed to delete soft connection "

	ErrorsFuncEventAlreadyExists   string = " Failed to register function class event ， The key name has already been registered "
	ErrorsFuncEventNotRegister     string = " The function corresponding to the key name was not found "
	ErrorsFuncEventNotCall         string = " The registered function cannot execute correctly "
	ErrorsBasePath                 string = " Failed to initialize the project root directory "
	ErrorsNoAuthorization          string = "token Authentication failed ， Please pass token Authorization interface reacquire token,"
	ErrorsParseTokenFail           string = " analysis token fail "
	ErrorsGormInitFail             string = "Gorm  Database driven 、 Connection initialization failed "
	ErrorsCasbinNoAuthorization    string = "Casbin  Authentication failed ， Please check in the background  casbin  Set parameters "
	ErrorsGormNotInitGlobalPointer string = "%s  Database global variable pointer not initialized ， Please enter the configuration file  Gormv2.yml  set up  Gormv2.%s.IsInitGolobalGormMysql = 1,  And ensure that the database configuration is correct  \n"
	//  Database part 
	ErrorsDbDriverNotExists   string = " Database driver type does not exist , Currently supported database types ：mysql、sqlserver、postgresql， You must submit the database type ："
	ErrorsDialectorDbInitFail string = "gorm dialector  initialization failed ,dbType:"

	//redis part 
	ErrorsRedisInitConnFail string = " initialization redis Connection pool failed "
	ErrorsRedisAuthFail     string = "Redis Auth  Authentication failed ， Password error "
	ErrorsRedisGetConnFail  string = "Redis  Failed to get a connection from the connection pool ， Maximum number of retries exceeded "
	//  Verifier error 
	ErrorsValidatorNotExists      string = " Non existent verifier "
	ErrorsValidatorBindParamsFail string = " Validator binding parameters failed "
	//token part 
	ErrorsTokenInvalid      string = " invalid token"
	ErrorsTokenNotActiveYet string = "token  Not activated "
	ErrorsTokenMalFormed    string = "token  Incorrect format "

	//snowflake
	ErrorsSnowflakeGetIdFail string = " obtain snowflake only ID An error occurred in the process "

	// File upload 
	ErrorsFilesUploadOpenFail string = " fail to open file ， details ："
	ErrorsFilesUploadReadFail string = " read file 32 Byte failure ， details ："

	//Agent correlation
	ErrorCommunicationEnd       string = "communication end"
	ErrorDataLength             string = "data length error"
	ErrorEmptyArray             string = "array is empty"
	ErrorMonitorDataPrefixError string = "monitor data prefix error"

	//Log correlation
	ErrorNotFoundParam          string = "not found param"
	ErrorInsufficientParameters string = "Insufficient number of parameters"
	ErrorNotParameters          string = "The parameter is null"
	ErrorResultEmpty            string = "The result is empty"

	// Taos
	ErrorTaosDiskspace string = "System out of disk space"
)
