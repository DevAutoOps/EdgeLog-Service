package publish_subscribe

import (
	"edgelog/app/global/variable"
	"github.com/streadway/amqp"
)

//  Create a producer 
func CreateProducer() (*producer, error) {
	//  Get configuration information 
	conn, err := amqp.Dial(variable.ConfigYml.GetString("RabbitMq.PublishSubscribe.Addr"))
	exchangeType := variable.ConfigYml.GetString("RabbitMq.PublishSubscribe.ExchangeType")
	exchangeName := variable.ConfigYml.GetString("RabbitMq.PublishSubscribe.ExchangeName")
	queueName := variable.ConfigYml.GetString("RabbitMq.PublishSubscribe.QueueName")
	dura := variable.ConfigYml.GetBool("RabbitMq.PublishSubscribe.Durable")

	if err != nil {
		variable.ZapLog.Error(err.Error())
		return nil, err
	}

	producer := &producer{
		connect:      conn,
		exchangeTyte: exchangeType,
		exchangeName: exchangeName,
		queueName:    queueName,
		durable:      dura,
	}
	return producer, nil
}

//   Define a message queue structure ：PublishSubscribe  Model 
type producer struct {
	connect      *amqp.Connection
	exchangeTyte string
	exchangeName string
	queueName    string
	durable      bool
	occurError   error
}

func (p *producer) Send(data string) bool {

	//  Get a channel 
	ch, err := p.connect.Channel()
	p.occurError = errorDeal(err)
	defer func() {
		_ = ch.Close()
	}()

	//  Claim switch ， In this mode, the producer is only responsible for delivering messages to the switch 
	err = ch.ExchangeDeclare(
		p.exchangeName, // Exchanger name 
		p.exchangeTyte, //fanout pattern ( Sector model  ) 。 solve   release 、 Subscription scenario related issues 
		p.durable,      //durable
		!p.durable,     //autodelete
		false,
		false,
		nil,
	)
	p.occurError = errorDeal(err)

	//  Post message 
	err = ch.Publish(
		p.exchangeName, //  Switch name 
		p.queueName,    // fanout  The mode is empty by default ， Indicates that all subscribed consumers will receive the same message 
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})

	if p.occurError != nil { //   An error occurred ， return  false
		return false
	} else {
		return true
	}
}

// Close manually after sending ， This does not affect send Send data multiple times 
func (p *producer) Close() {
	_ = p.connect.Close()
}

//  Define an error handling function 
func errorDeal(err error) error {
	if err != nil {
		variable.ZapLog.Error(err.Error())
	}
	return err
}
