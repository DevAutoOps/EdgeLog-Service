package consts

//  Constants defined here ， Usually with error codes + Error description composition ， Generally used for interface return 
const (
	//  The process was ended 
	ProcessKilled string = " Received signal ， The process was ended "
	//  Form validator prefix 
	ValidatorPrefix              string = "Form_Validator_"
	ValidatorParamsCheckFailCode int    = -400300
	ValidatorParamsCheckFailMsg  string = " Parameter verification failed "

	// An error occurred in the server code 
	ServerOccurredErrorCode int    = -500100
	ServerOccurredErrorMsg  string = " A code execution error occurred inside the server , "

	// token relevant 
	JwtTokenOK            int    = 200100           //token Effective 
	JwtTokenInvalid       int    = -400100          // invalid token
	JwtTokenExpired       int    = -400101          // overdue token
	JwtTokenFormatErrCode int    = -400102          // Submitted  token  Format error 
	JwtTokenFormatErrMsg  string = " Submitted  token  Format error " // Submitted  token  Format error 

	//SnowFlake  Snowflake algorithm 
	StartTimeStamp = int64(1483228800000) // Start time cut  (2017-01-01)
	MachineIdBits  = uint(10)             // machine id Number of digits occupied 
	SequenceBits   = uint(12)             // Number of bits occupied by the sequence 
	//MachineIdMax   = int64(-1 ^ (-1 << MachineIdBits)) // Maximum machines supported id quantity 
	SequenceMask   = int64(-1 ^ (-1 << SequenceBits)) //
	MachineIdShift = SequenceBits                     // machine id Shift left 
	TimestampShift = SequenceBits + MachineIdBits     // Time stamp shift left 

	// CURD  Common business status code 
	CurdStatusOkCode         int    = 200
	CurdStatusOkMsg          string = "Success"
	CurdCreatFailCode        int    = -400200
	CurdCreatFailMsg         string = " Failed to add "
	CurdUpdateFailCode       int    = -400201
	CurdUpdateFailMsg        string = " Update failed "
	CurdDeleteFailCode       int    = -400202
	CurdDeleteFailMsg        string = " Deletion failed "
	CurdSelectFailCode       int    = -400203
	CurdSelectFailMsg        string = " Query no data "
	CurdRegisterFailCode     int    = -400204
	CurdRegisterFailMsg      string = " login has failed "
	CurdLoginFailCode        int    = -400205
	CurdLoginFailMsg         string = " Login failed "
	CurdRefreshTokenFailCode int    = -400206
	CurdRefreshTokenFailMsg  string = " Refresh Token fail "

	// File upload 
	FilesUploadFailCode            int    = -400250
	FilesUploadFailMsg             string = " File upload failed ,  Error getting uploaded file !"
	FilesUploadMoreThanMaxSizeCode int    = -400251
	FilesUploadMoreThanMaxSizeMsg  string = " The long transfer file exceeds the maximum value set by the system , Maximum value allowed by the system （M）："
	FilesUploadMimeTypeFailCode    int    = -400252
	FilesUploadMimeTypeFailMsg     string = " file mime Type not allowed "

	//websocket
	WsServerNotStartCode int    = -400300
	WsServerNotStartMsg  string = "websocket  The service is not turned on ， Please open in the configuration file ， Correlation path ：config/config.yml"
	WsOpenFailCode       int    = -400301
	WsOpenFailMsg        string = "websocket open Phase initialization basic parameters failed "

	// Verification Code 
	CaptchaGetParamsInvalidMsg    string = " Get verification code ： The submitted verification code parameter is invalid , Please check the verification code ID And whether the file name suffix is complete "
	CaptchaGetParamsInvalidCode   int    = -400350
	CaptchaCheckParamsInvalidMsg  string = " Verification code ： The submitted parameter is invalid ， Please check  【 Verification Code ID、 Verification code value 】  Whether the key name at the time of submission is consistent with the configuration item "
	CaptchaCheckParamsInvalidCode int    = -400351
	CaptchaCheckOkMsg             string = " Verification code passed "
	//CaptchaCheckOkCode            int    = 200
	CaptchaCheckFailCode int    = -400355
	CaptchaCheckFailMsg  string = " Verification code verification failed "

	PortCrossBorderMsg string = "Port range out of bounds"
	IpNotFoundCode     int    = -401002
	IpNotFoundMsg      string = "The IP address of this machine was not found"
)

const (
	OsWindows string = "windows"
	OsLinux   string = "linux"
	OsMac     string = "darwin"
	OsOther   string = "other"
)

const (
	MonitorCpuUsage         uint8 = iota //CPU utilization
	MonitorCpuBaseUsage                  //Basic CPU utilization
	MonitorCpuLoad                       //Average CPU load
	MonitorMemUsed                       //Memory usage
	MonitorMemRate                       //Memory utilization
	MonitorDiskPartRate                  //Hard disk partition utilization
	MonitorDiskReadFlow                  //Hard disk read traffic(KB/s)
	MonitorDiskWriteFlow                 //Hard disk write traffic(KB/s)
	MonitorDiskIoWait                    //Hard disk IO wait time(ms)
	MonitorDiskIoService                 //Hard disk IO service time(ms)
	MonitorDiskIoBusy                    //Hard disk IO busy ratio(%)
	MonitorDiskPartTotal                 //Total capacity of hard disk partition(GB)
	MonitorDiskPartUsage                 //Used amount of hard disk partition(GB)
	MonitorNginxConcurrency              //Single nginx concurrency
	MonitorNetworkReception              // Network Reception
	MonitorNetworkSending                // Network sending
)

const (
	StatusNode int8 = iota
	StatusNginx
)

const CollectingMonitorInterval = 10

const (
	HostCpuThreshold             string = "host_cpu_threshold"
	HostMemoryThreshold          string = "host_memory_threshold"
	HostDiskThreshold            string = "host_disk_threshold"
	HostDefaultThresholdValueStr string = "90"
	HostDefaultThresholdValue    int    = 90
)

var codeMsg = map[int]string{
	//Grayscale
	CurdStatusOkCode:   "Success",
	CurdSelectFailCode: "Query failed",
	CurdCreatFailCode:  "Add failed",
	CurdUpdateFailCode: "Update failed",
	CurdDeleteFailCode: "Delete failed",
}

const (
	Success                 string = "Success"
	Fail                    string = "Fail"
	UpdateSuccess           string = "Update success"
	UpdateFail              string = "Update fail"
	UploadSuccess           string = "Upload success"
	UploadFail              string = "Upload fail"
	InitPermissionDeniedMsg string = "Please use the root account and password to operate"
	PermissionDenied        string = "Permission denied"
	AgentConfigName         string = "config.yml"
)
