package qiwlar

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type router struct {
	httprouter.Router
	panicHandler func(*QResponse, interface{})
}

func newRouter() *router {
	r := &router{
		Router: httprouter.Router{
			RedirectTrailingSlash:  true,
			RedirectFixedPath:      true,
			HandleMethodNotAllowed: true,
			HandleOPTIONS:          true,
		},
		panicHandler: func(qr *QResponse, r interface{}) {
			qr.W.Header().Set("Content-Type", "application/json; charset=utf-8")
			qr.W.WriteHeader(http.StatusInternalServerError)

			d, _ := json.Marshal(map[string]interface{}{
				"status_code": 500,
				"errors":      "internal server error",
				"data":        nil,
			})

			fmt.Fprint(qr.W, string(d))
		},
	}

	return r
}

func (r *router) add(endPoint EndPoint, prefix string) {
	qr := new(QResponse)

	for _, decorator := range endPoint.Decorators {
		decorator(qr)

		if qr.Start {
			qr.Path = pattern(prefix + "/" + qr.Path)
			break
		}
	}

	r.Handle(qr.Method, qr.Path, func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		if len(p) > 0 {
			ctx := req.Context()
			ctx = context.WithValue(ctx, httprouter.ParamsKey, p)
			req = req.WithContext(ctx)
		}

		handler(w, req, &endPoint, r.panicHandler)
	})
}

type AnswerScheme struct {
	StatusCode int         `json:"status_code"`
	Errors     interface{} `json:"errors"`
	Data       interface{} `json:"data"`
}

func launchDecorators(qr *QResponse, d []Decorator) interface{} {
	for _, decorator := range d {
		if err := decorator(qr); err != nil {
			return err
		}
	}

	return nil
}

func endPointResult(qr *QResponse, ep *EndPoint) interface{} {
	fnValue := reflect.ValueOf(ep.Fn)
	params := []reflect.Value{}

	if fnValue.Type().NumIn() > 0 {
		for i := 0; i < fnValue.Type().NumIn(); i++ {
			for _, result := range qr.Results {
				if reflect.TypeOf(result) != fnValue.Type().In(i) {
					continue
				}

				params = append(params, reflect.ValueOf(result))
			}
		}
	}

	return fnValue.Call(params)[0].Interface()
}

func handler(w http.ResponseWriter, r *http.Request, ep *EndPoint, ph func(*QResponse, interface{})) {
	qr := &QResponse{R: r, W: w}

	defer func(qr *QResponse) {
		if r := recover(); r != nil {
			ph(qr, r)
		}
	}(qr)

	as := AnswerScheme{}

	if err := launchDecorators(qr, ep.Decorators); err != nil {
		as.Errors = err
	}

	if as.Errors == nil {
		result := endPointResult(qr, ep)

		if err, ok := result.(Exception); ok {
			as.Errors = err.Data
			qr.StatusCode = err.StatusCode
		} else if err, ok := result.(error); ok {
			qr.StatusCode = 400
			as.Errors = err.Error()
		} else {
			as.Data = result
		}
	}

	var result interface{}
	as.StatusCode = qr.StatusCode

	if qr.CustomFormat != nil {
		result = qr.CustomFormat(as)
		b, _ := json.Marshal(&result)
		result = string(b)
	}

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	if !strings.Contains(w.Header().Get("Content-Type"), "application/json") {
		if as.Data != nil {
			result = as.Data
		} else {
			result = as.Errors
		}
	}

	if result == nil {
		result = &as
		b, _ := json.Marshal(result)
		result = string(b)
	}

	w.WriteHeader(as.StatusCode)
	fmt.Fprint(w, result)
}

func pattern(path string) string {
	p := ""
	words := strings.Split(
		strings.ReplaceAll(path, "/", " "),
		" ",
	)

	for _, word := range words {
		if word == "" {
			continue
		}

		p += word + "/"
	}

	if len(p) > 0 {
		if p[len(p)-1] == '/' {
			p = p[:len(p)-1]
		}
	} else {
		return "/"
	}

	if p[0] != '/' {
		p = "/" + p
	}

	return p
}
