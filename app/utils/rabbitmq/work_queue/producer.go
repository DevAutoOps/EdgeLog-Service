package work_queue

import (
	"edgelog/app/global/variable"
	"github.com/streadway/amqp"
)

//  Create a producer 
func CreateProducer() (*producer, error) {
	//  Get configuration information 
	conn, err := amqp.Dial(variable.ConfigYml.GetString("RabbitMq.WorkQueue.Addr"))
	queueName := variable.ConfigYml.GetString("RabbitMq.WorkQueue.QueueName")
	dura := variable.ConfigYml.GetBool("RabbitMq.WorkQueue.Durable")

	if err != nil {
		variable.ZapLog.Error(err.Error())
		return nil, err
	}

	vProducer := &producer{
		connect:   conn,
		queueName: queueName,
		durable:   dura,
	}
	return vProducer, nil
}

//   Define a message queue structure ：helloworld  Model 
type producer struct {
	connect    *amqp.Connection
	queueName  string
	durable    bool
	occurError error
}

func (p *producer) Send(data string) bool {

	//  Get a channel 
	ch, err := p.connect.Channel()
	p.occurError = errorDeal(err)

	defer func() {
		_ = ch.Close()
	}()

	//  Declare message queue 
	_, err = ch.QueueDeclare(
		p.queueName, //  Queue name 
		p.durable,   // Persistent ，false All mode data is in memory ，true Will be saved in erlang Built in database ， But it affects speed 
		!p.durable,  // producer 、 Delete queue when all consumers are disconnected 。 generally speaking ， Data needs to be persistent ， Do not delete ； Non persistent ， Just delete 
		false,       // Private queue ，false Multiple identities are allowed  consumer  Post messages to the queue ，true  Indicates exclusive 
		false,       //  If the queue is already declared on the server ， Set to  true ， Otherwise set to  false；
		nil,         //  Related parameters 
	)
	p.occurError = errorDeal(err)

	//  Post message 
	err = ch.Publish(
		"",          // helloworld 、workqueue  The mode is set to an empty string ， Indicates that the default switch is used 
		p.queueName, //  be careful ： Simple mode  key  Indicates the queue name 
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
