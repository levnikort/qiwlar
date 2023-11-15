package decorator

import (
	"github.com/julienschmidt/httprouter"
	"github.com/levnikort/qiwlar"
)

func Params() qiwlar.Decorator {
	return func(q *qiwlar.QResponse) interface{} {
		q.Results = append(q.Results, httprouter.ParamsFromContext(q.R.Context()))

		return nil
	}
}
