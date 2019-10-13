package server

type context struct {
	chatS *ChatS
}

func NewContext(s *ChatS) *context {
	return &context{
		chatS: s,
	}
}
