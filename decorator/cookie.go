package decorator

import (
	"net/http"

	"github.com/levnikort/qiwlar"
)

func Cookie(c *http.Cookie) qiwlar.Decorator {
	return func(q *qiwlar.QResponse) interface{} {
		if c != nil {
			http.SetCookie(q.W, c)
		}

		q.Results = append(q.Results, q.R.Cookies())

		return nil
	}
}
