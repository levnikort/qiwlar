package decorator

import "github.com/levnikort/qiwlar"

func CustomFormat(formatter func(qiwlar.AnswerScheme) interface{}) qiwlar.Decorator {
	return func(q *qiwlar.QResponse) interface{} {
		q.CustomFormat = formatter

		return nil
	}
}
