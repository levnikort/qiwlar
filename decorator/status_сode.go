package decorator

import (
	"github.com/levnikort/qiwlar"
)

func StatusCode(code int) qiwlar.Decorator {
	return func(q *qiwlar.QResponse) interface{} {
		q.StatusCode = code

		return nil
	}
}
