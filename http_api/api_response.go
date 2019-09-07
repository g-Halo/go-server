package http_api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
)

type Err struct {
	Code int
	Text string
}

type APIHandler func(http.ResponseWriter, *http.Request, httprouter.Params) (interface{}, error)
type Decorator func(APIHandler) APIHandler

func Decorate(f APIHandler, ds ...Decorator) httprouter.Handle {
	decorated := f
	for _, decorate := range ds {
		decorated = decorate(decorated)
	}
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		decorated(w, req, ps)
	}
}

func PlainText(f APIHandler) APIHandler {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
		code := 200
		data, err := f(w, req, ps)
		if err != nil {
			data = err.Error()
		}
		switch d := data.(type) {
		case string:
			w.WriteHeader(code)
			io.WriteString(w, d)
		case []byte:
			w.WriteHeader(code)
			w.Write(d)
		default:
			panic(fmt.Sprintf("unknown response type %T", data))
		}
		return nil, nil
	}
}