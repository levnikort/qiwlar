package decorator

import (
	"net/http"

	"github.com/levnikort/qiwlar"
)

func Get(path string) qiwlar.Decorator {
	return func(qr *qiwlar.QResponse) interface{} {
		qr.Start = true
		qr.StatusCode = 200
		qr.Method = http.MethodGet
		qr.Path = path

		return nil
	}
}

func Post(path string) qiwlar.Decorator {
	return func(qr *qiwlar.QResponse) interface{} {
		qr.Start = true
		qr.StatusCode = 201
		qr.Method = http.MethodPost
		qr.Path = path

		return nil
	}
}

func Put(path string) qiwlar.Decorator {
	return func(qr *qiwlar.QResponse) interface{} {
		qr.Start = true
		qr.StatusCode = 200
		qr.Method = http.MethodPut
		qr.Path = path

		return nil
	}
}

func Delete(path string) qiwlar.Decorator {
	return func(qr *qiwlar.QResponse) interface{} {
		qr.Start = true
		qr.StatusCode = 200
		qr.Method = http.MethodDelete
		qr.Path = path

		return nil
	}
}

func Head(path string) qiwlar.Decorator {
	return func(qr *qiwlar.QResponse) interface{} {
		qr.Start = true
		qr.StatusCode = 200
		qr.Method = http.MethodHead
		qr.Path = path

		return nil
	}
}

func Options(path string) qiwlar.Decorator {
	return func(qr *qiwlar.QResponse) interface{} {
		qr.Start = true
		qr.StatusCode = 200
		qr.Method = http.MethodOptions
		qr.Path = path

		return nil
	}
}

func Patch(path string) qiwlar.Decorator {
	return func(qr *qiwlar.QResponse) interface{} {
		qr.Start = true
		qr.StatusCode = 200
		qr.Method = http.MethodPatch
		qr.Path = path

		return nil
	}
}
