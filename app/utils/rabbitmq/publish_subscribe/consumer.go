package publish_subscribe

import (
	"edgelog/app/global/variable"
	"github.com/streadway/amqp"
	"time"
)

func CreateConsumer() (*consumer, error) {
	//  Get configuration information 
	conn, err := amqp.Dial(variable.ConfigYml.GetString("RabbitMq.PublishSubscribe.Addr"))
	exchangeType := variable.ConfigYml.GetString("RabbitMq.PublishSubscribe.ExchangeType")
	exchangeName := variable.ConfigYml.GetString("RabbitMq.PublishSubscribe.ExchangeName")
	queueName := variable.ConfigYml.GetString("RabbitMq.PublishSubscribe.QueueName")
	dura := variable.ConfigYml.GetBool("RabbitMq.PublishSubscribe.Durable")
	chanNumber := variable.ConfigYml.GetInt("RabbitMq.PublishSubscribe.ConsumerChanNumber")
	reconnectInterval := variable.ConfigYml.GetDuration("RabbitMq.PublishSubscribe.OffLineReconnectIntervalSec")
	retryTimes := variable.ConfigYml.GetInt("RabbitMq.PublishSubscribe.RetryCount")

	if err != nil {
		return nil, err
	}

	consumer := &consumer{
		connect:                     conn,
		exchangeTyte:                exchangeType,
		exchangeName:                exchangeName,
		queueName:                   queueName,
		durable:                     dura,
		chanNumber:                  chanNumber,
		connErr:                     conn.NotifyClose(make(chan *amqp.Error, 1)),
		offLineReconnectIntervalSec: reconnectInterval,
		retryTimes:                  retryTimes,
	}

	return consumer, nil
}

//   Define a message queue structure ：PublishSubscribe  Model 
type consumer struct {
	connect                     *amqp.Connection
	exchangeTyte                string
	exchangeName                string
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
				"", //  fanout  Mode set to   empty   that will do 
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

		}(i)
	}

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
