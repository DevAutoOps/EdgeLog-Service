package hello_world

import (
	"edgelog/app/global/variable"
	"github.com/streadway/amqp"
	"time"
)

func CreateConsumer() (*consumer, error) {
	//  Get configuration information 

	conn, err := amqp.Dial(variable.ConfigYml.GetString("RabbitMq.HelloWorld.Addr"))
	queueName := variable.ConfigYml.GetString("RabbitMq.HelloWorld.QueueName")
	dura := variable.ConfigYml.GetBool("RabbitMq.HelloWorld.Durable")
	chanNumber := variable.ConfigYml.GetInt("RabbitMq.HelloWorld.ConsumerChanNumber")
	reconnectInterval := variable.ConfigYml.GetDuration("RabbitMq.HelloWorld.OffLineReconnectIntervalSec")
	retryTimes := variable.ConfigYml.GetInt("RabbitMq.HelloWorld.RetryCount")

	if err != nil {
		//log.Println(err.Error())
		return nil, err
	}
	consumer := &consumer{
		connect:                     conn,
		queueName:                   queueName,
		durable:                     dura,
		chanNumber:                  chanNumber,
		connErr:                     conn.NotifyClose(make(chan *amqp.Error, 1)),
		offLineReconnectIntervalSec: reconnectInterval,
		retryTimes:                  retryTimes,
	}
	return consumer, nil
}

//   Define a message queue structure ：helloworld  Model 
type consumer struct {
	connect                     *amqp.Connection
	queueName                   string
	durable                     bool
	chanNumber                  int
	occurError                  error
	connErr                     chan *amqp.Error
	callbackForReceived         func(receivedData string) //    Disconnection reconnection ， Internal use of structure 
	offLineReconnectIntervalSec time.Duration
	retryTimes                  int
	callbackOffLine             func(err *amqp.Error) //    Disconnection reconnection ， Internal use of structure 
}

//  receive 、 Processing messages 
func (c *consumer) Received(callbackFunDealSmg func(receivedData string)) {
	defer func() {
		_ = c.connect.Close()
	}()
	//  Assign the callback function address to the structure variable ， Used for disconnection and reconnection 
	c.callbackForReceived = callbackFunDealSmg

	blocking := make(chan bool)

	for i := 1; i <= c.chanNumber; i++ {
		go func(chanNo int) {
			ch, err := c.connect.Channel()
			c.occurError = errorDeal(err)
			defer func() {
				_ = ch.Close()
			}()

			queue, err := ch.QueueDeclare(
				c.queueName,
				c.durable,
				!c.durable,
				false,
				false,
				nil,
			)

			c.occurError = errorDeal(err)

			msgs, err := ch.Consume(
				queue.Name,
				queue.Name, //   Consumer marking ， Please make sure that it is unique in one message channel 
				true,       // Automatic response confirmation 
				false,      // Private queue ，false Multiple identities are allowed  consumer  Post messages to the queue ，true  Indicates exclusive 
				false,      //RabbitMQ I won't support it noLocal sign 。
				false,      //  If the queue is already declared on the server ， Set to  true ， Otherwise set to  false；
				nil,
			)
			c.occurError = errorDeal(err)

			for msg := range msgs {
				//  Processing messages through callbacks 
				callbackFunDealSmg(string(msg.Body))
			}

		}(i)
	}

	<-blocking

}

// Consumer end ， Drop line reconnection monitor 
func (c *consumer) OnConnectionError(callbackOfflineErr func(err *amqp.Error)) {
	c.callbackOffLine = callbackOfflineErr
	go func() {
		select {
		case err := <-c.connErr:
			var i = 1
			for i = 1; i <= c.retryTimes; i++ {
				//  Automatic reconnection mechanism 
				time.Sleep(c.offLineReconnectIntervalSec * time.Second)
				conn, err := CreateConsumer()
				if err != nil {
					continue
				} else {
					go func() {
						c.connErr = conn.connect.NotifyClose(make(chan *amqp.Error, 1))
						go conn.OnConnectionError(c.callbackOffLine)
						conn.Received(c.callbackForReceived)
					}()
					break
				}
			}
			if i > c.retryTimes {
				callbackOfflineErr(err)
			}
		}
	}()

}
