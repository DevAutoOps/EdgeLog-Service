package redis_factory

import (
	"edgelog/app/core/event_manage"
	"edgelog/app/global/my_errors"
	"edgelog/app/global/variable"
	"edgelog/app/utils/yml_config"
	"edgelog/app/utils/yml_config/ymlconfig_interf"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"time"
)

var redisPool *redis.Pool
var configYml ymlconfig_interf.YmlConfigInterf

//  Package at the bottom of the program ，init  The execution of the initialized code segment takes precedence over the upper layer code ， Therefore, the global configuration item variable cannot be used to read configuration items here 
func init() {
	configYml = yml_config.CreateYamlFactory()
	redisPool = initRedisClientPool()
}
func initRedisClientPool() *redis.Pool {
	redisPool = &redis.Pool{
		MaxIdle:     configYml.GetInt("Redis.MaxIdle"),                        // Maximum idle 
		MaxActive:   configYml.GetInt("Redis.MaxActive"),                      // Maximum active number 
		IdleTimeout: configYml.GetDuration("Redis.IdleTimeout") * time.Second, // Maximum idle connection wait time ， After this time ， Idle connections will be closed 
		Dial: func() (redis.Conn, error) {
			// Corresponding here redis ip And port number 
			conn, err := redis.Dial("tcp", configYml.GetString("Redis.Host")+":"+configYml.GetString("Redis.Port"))
			if err != nil {
				variable.ZapLog.Error(my_errors.ErrorsRedisInitConnFail + err.Error())
				return nil, err
			}
			auth := configYml.GetString("Redis.Auth") // Set by configuration item redis password 
			if len(auth) >= 1 {
				if _, err := conn.Do("AUTH", auth); err != nil {
					_ = conn.Close()
					variable.ZapLog.Error(my_errors.ErrorsRedisAuthFail + err.Error())
				}
			}
			_, _ = conn.Do("select", configYml.GetInt("Redis.IndexDb"))
			return conn, err
		},
	}
	//  take redis Close event for ， Register in global event unified Manager ， Unified destruction upon program exit 
	event_manage.CreateEventManageFactory().Set(variable.EventDestroyPrefix+"Redis", func(args ...interface{}) {
		_ = redisPool.Close()
	})
	return redisPool
}

//   Get a connection from the connection pool redis connect 
func GetOneRedisClient() *RedisClient {
	maxRetryTimes := configYml.GetInt("Redis.ConnFailRetryTimes")
	var oneConn redis.Conn
	for i := 1; i <= maxRetryTimes; i++ {
		oneConn = redisPool.Get()
		if oneConn.Err() != nil {
			//variable.ZapLog.Error("Redis： Network interruption , Start reconnection in progress ..." , zap.Error(oneConn.Err()))
			if i == maxRetryTimes {
				variable.ZapLog.Error(my_errors.ErrorsRedisGetConnFail, zap.Error(oneConn.Err()))
				return nil
			}
			// If there is a short jitter in the network ， After brief dormancy ， Support automatic reconnection 
			time.Sleep(time.Second * configYml.GetDuration("Redis.ReConnectInterval"))
		} else {
			break
		}
	}
	return &RedisClient{oneConn}
}

//  Define a redis Client structure 
type RedisClient struct {
	client redis.Conn
}

//  by redis-go  Client encapsulation unified operation function entry 
func (r *RedisClient) Execute(cmd string, args ...interface{}) (interface{}, error) {
	return r.client.Do(cmd, args...)
}

//  Release connection to connection pool 
func (r *RedisClient) ReleaseOneRedisClient() {
	_ = r.client.Close()
}

//   Encapsulates several data type conversion functions 

//bool  Type conversion 
func (r *RedisClient) Bool(reply interface{}, err error) (bool, error) {
	return redis.Bool(reply, err)
}

//string  Type conversion 
func (r *RedisClient) String(reply interface{}, err error) (string, error) {
	return redis.String(reply, err)
}

//strings  Type conversion 
func (r *RedisClient) Strings(reply interface{}, err error) ([]string, error) {
	return redis.Strings(reply, err)
}

//Float64  Type conversion 
func (r *RedisClient) Float64(reply interface{}, err error) (float64, error) {
	return redis.Float64(reply, err)
}

//int  Type conversion 
func (r *RedisClient) Int(reply interface{}, err error) (int, error) {
	return redis.Int(reply, err)
}

//int64  Type conversion 
func (r *RedisClient) Int64(reply interface{}, err error) (int64, error) {
	return redis.Int64(reply, err)
}

//uint64  Type conversion 
func (r *RedisClient) Uint64(reply interface{}, err error) (uint64, error) {
	return redis.Uint64(reply, err)
}

//Bytes  Type conversion 
func (r *RedisClient) Bytes(reply interface{}, err error) ([]byte, error) {
	return redis.Bytes(reply, err)
}
