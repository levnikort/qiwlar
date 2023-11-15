package decorator

import "github.com/levnikort/qiwlar"

func Header(key, value string) qiwlar.Decorator {
	return func(q *qiwlar.QResponse) interface{} {
		q.W.Header().Set(key, value)
		q.Results = append(q.Results, q.W.Header())

		return nil
	}
}
