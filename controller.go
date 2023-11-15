package qiwlar

import (
	"net/http"
)

type Controller struct {
	prefix           string
	commonDecorators []Decorator
	endPoints        []EndPoint
}

func (c *Controller) Prefix(prefix string) struct{} {
	c.prefix = prefix

	return struct{}{}
}

func (c *Controller) ForAll(decorators ...Decorator) struct{} {
	c.commonDecorators = append(c.commonDecorators, decorators...)

	return struct{}{}
}

func (c *Controller) Collect(parts []interface{}) {
	var endPoints []EndPoint
	var endPoint EndPoint

	for _, part := range parts {
		if decorator, ok := part.(Decorator); ok {
			endPoint.Decorators = append(endPoint.Decorators, decorator)
			continue
		}

		if _, ok := part.(struct{}); ok {
			continue
		}

		endPoint.Fn = part
		d := []Decorator{}
		d = append(d, c.commonDecorators...)
		d = append(d, endPoint.Decorators...)
		endPoint.Decorators = d
		endPoints = append(endPoints, endPoint)
		endPoint = EndPoint{}
	}

	c.endPoints = endPoints
}

func (c *Controller) EndPoints() []EndPoint {
	for _, endPoint := range c.endPoints {
		d := []Decorator{}
		d = append(d, endPoint.Decorators...)

		endPoint.Decorators = d
	}

	return c.endPoints
}

type QResponse struct {
	Start        bool
	StatusCode   int
	Method       string
	Path         string
	CustomFormat func(AnswerScheme) interface{}
	Results      []interface{}
	W            http.ResponseWriter
	R            *http.Request
}

type Decorator func(*QResponse) interface{}

type EndPoint struct {
	Decorators []Decorator
	Fn         interface{}
}
