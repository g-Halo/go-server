package mq

import (
	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/pkg/logger"
	"github.com/nsqio/go-nsq"
)

type Nsq struct {
	Topic string
	Producer *nsq.Producer
}

func NewNsq() *Nsq {
	cfg := nsq.NewConfig()
	producer, err := nsq.NewProducer(conf.Conf.NSQTopic, cfg)
	if err != nil {
		logger.Error(err)
		return nil
	}

	return &Nsq{
		Topic: conf.Conf.NSQTopic,
		Producer: producer,
	}
}

func (nsq *Nsq) Push() error {
	producer := nsq.Producer
	err := producer.Publish(nsq.Topic, []byte("hello"))

	if err != nil {
		return err
	} else {
		return nil
	}
}