package job

import (
	"github.com/g-Halo/go-server/internal/job/conf"
	"github.com/g-Halo/go-server/pkg/logger"
	"github.com/g-Halo/go-server/pkg/mq"
	"github.com/nsqio/go-nsq"
)

type Job struct {
	MQ *mq.Server
	Conf *conf.Config
	Consumer *Consumer
}

type Consumer struct {

}

func New(conf *conf.Config) *Job {
	consumer := &Consumer{}
	return &Job{
		Conf: conf,
		Consumer: consumer,
	}
}

func (j *Job) ConsumerHandle() {
	cfg := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(j.Conf.Nsq.Topic, j.Conf.Nsq.Channel, cfg)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debug("subscribe the topic in nsq")

	// Set the number of messages that can be in flight at any given time
	// you'll want to set this number as the default is only 1. This can be
	// a major concurrency knob that can change the performance of your application.
	consumer.ChangeMaxInFlight(200)

	consumer.AddHandler(j.Consumer)

	err = consumer.ConnectToNSQLookupd(j.Conf.Nsq.LookUpAddress)

	<- consumer.StopChan
}

func (j *Job) Close() error {
	return nil
}

func (c *Consumer) HandleMessage(message *nsq.Message) error {
	logger.Debug("consumer message: %s", string(message.Body))

	return nil
}
