package gorm_v2

//  Database parameter configuration ， structural morphology 
//  It is used to solve the deployment of complex business scenarios connected to multiple servers  mysql、sqlserver、postgresql  database 
//  See unit test for specific usage (test/gormv2_test.go) file ，TestCustomeParamsConnMysql  Function code snippet 

type ConfigParams struct {
	Write ConfigParamsDetail
	Read  ConfigParamsDetail
}
type ConfigParamsDetail struct {
	Host     string
	DataBase string
	Port     int
	Prefix   string
	User     string
	Pass     string
	Charset  string
}
