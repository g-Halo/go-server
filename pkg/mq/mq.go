package mq

type Server interface {
	Push() error
}

type MQ struct {
	Server Server
}

func New() *MQ {
	return &MQ{
		Server: NewNsq(),
	}
}
