package http_api

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/yigger/go-server/logger"
	"github.com/yigger/go-server/model"
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
		case map[string]interface{}:
			data, err := json.Marshal(data)
			if err != nil {
				panic(fmt.Sprintf("response json %T", data))
			}
			w.WriteHeader(code)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write([]byte(string(data)))
		default:
			panic(fmt.Sprintf("unknown response type %T", data))
		}
		return nil, nil
	}
}

func MiddlewareHandler(mdFunc func(string) (*model.User, bool), fn func(http.ResponseWriter, *http.Request, *model.User) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		user, _ := mdFunc(token)
		if user == nil {
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized"))
		} else {
			// common setting
			w.Header().Set("content-type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(200)
			// 从回调函数获取 data
			data, err := fn(w, r, user)

			if err != nil {
				data = err.Error()
			}

			switch d := data.(type) {
			case string:
				io.WriteString(w, d)
			case []byte:
				w.Write(d)
			case map[string]interface{}:
				data, err := json.Marshal(data)
				if err != nil {
					logger.Errorf("response json %T", data)
				}
				w.Write([]byte(string(data)))
			default:
				logger.Errorf("unknown response type %T", data)
			}
		}
	}
}

