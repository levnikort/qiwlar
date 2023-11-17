package decorator

import "github.com/levnikort/qiwlar"

func Query() qiwlar.Decorator {
	return func(q *qiwlar.QResponse) interface{} {
		q.Results = append(q.Results, q.R.URL.Query())

		return nil
	}
}
