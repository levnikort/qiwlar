package decorator

import (
	"io"

	"github.com/levnikort/qiwlar"
)

type RCBody io.ReadCloser

func Body() qiwlar.Decorator {
	return func(q *qiwlar.QResponse) interface{} {
		var b RCBody = q.R.Body

		q.Results = append(q.Results, b)

		return nil
	}
}
