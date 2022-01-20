package routing

import (
	"edgelog/app/global/variable"
	"github.com/streadway/amqp"
	"time"
)

func CreateConsumer() (*consumer, error) {
	//  Get configuration information 
	conn, err := amqp.Dial(variable.ConfigYml.GetString("RabbitMq.Routing.Addr"))
	exchangeType := variable.ConfigYml.GetString("RabbitMq.Routing.ExchangeType")
	exchangeName := variable.ConfigYml.GetString("RabbitMq.Routing.ExchangeName")
	queueName := variable.ConfigYml.GetString("RabbitMq.Routing.QueueName")
	dura := variable.ConfigYml.GetBool("RabbitMq.Routing.Durable")
	reconnectInterval := variable.ConfigYml.GetDuration("RabbitMq.Routing.OffLineReconnectIntervalSec")
	retryTimes := variable.ConfigYml.GetInt("RabbitMq.Routing.RetryCount")

	if err != nil {
		return nil, err
	}

	consumer := &consumer{
		connect:                     conn,
		exchangeTyte:                exchangeType,
		exchangeName:                exchangeName,
		queueName:                   queueName,
		durable:                     dura,
		connErr:                     conn.NotifyClose(make(chan *amqp.Error, 1)),
		offLineReconnectIntervalSec: reconnectInterval,
		retryTimes:                  retryTimes,
	}
	return consumer, nil
}

//   Define a message queue structure ：Routing  Model 
type consumer struct {
	connect                     *amqp.Connection
	exchangeTyte                string
	exchangeName                string
	queueName                   string
	durable                     bool
	occurError                  error
	connErr                     chan *amqp.Error
	routeKey                    string                    //    Disconnection reconnection ， Internal use of structure 
	callbackForReceived         func(receivedData string) //    Disconnection reconnection ， Internal use of structure 
	offLineReconnectIntervalSec time.Duration
	retryTimes                  int
	callbackOffLine             func(err *amqp.Error) //    Disconnection reconnection ， Internal use of structure 
}

//  receive 、 Processing messages 
func (c *consumer) Received(routeKey string, callbackFunDealSmg func(receivedData string)) {
	defer func() {
		_ = c.connect.Close()
	}()
	//  Assign the callback function address to the structure variable ， Used for disconnection and reconnection 
	c.routeKey = routeKey
	c.callbackForReceived = callbackFunDealSmg

	blocking := make(chan bool)

	go func(key string) {

		ch, err := c.connect.Channel()
		c.occurError = errorDeal(err)
		defer func() {
			_ = ch.Close()
		}()

		//  statement exchange Switch 
		err = ch.ExchangeDeclare(
			c.exchangeName, //exchange name
			c.exchangeTyte, //exchange kind
			c.durable,      // Is the data persistent 
			!c.durable,     // When all connections are disconnected ， Delete switch 
			false,
			false,
			nil,
		)
		//  Declaration queue 
		queue, err := ch.QueueDeclare(
			c.queueName,
			c.durable,
			!c.durable,
			false,
			false,
			nil,
		)
		c.occurError = errorDeal(err)

		// Queue binding 
		err = ch.QueueBind(
			queue.Name,
			key, //  routing  pattern , The producer will deliver the message to the switch route_key，  Consumers match different key Get message 、 handle 
			c.exchangeName,
			false,
			nil,
		)
		c.occurError = errorDeal(err)

		msgs, err := ch.Consume(
			queue.Name, //  Queue name 
			"",         //   Consumer marking ， Please make sure that it is unique in one message channel 
			true,       // Automatic response confirmation ， Set here as false， Manual confirmation 
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

	}(routeKey)

	<-blocking

}

// Consumer end ， Error callback after disconnection and reconnection failure 
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
						conn.Received(c.routeKey, c.callbackForReceived)
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
