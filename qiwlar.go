package qiwlar

import (
	"fmt"
	"net/http"
	"reflect"
)

type Qiwlar struct {
	store  *store
	router *router
}

func New(metaModule interface{}) *Qiwlar {
	q := &Qiwlar{
		store:  &store{make(map[reflect.Type]*module)},
		router: newRouter(),
	}

	appModule := q.buildModule(metaModule)
	q.resolveDependencies(appModule)

	return q
}

func (q *Qiwlar) SetPanicHandler(ph func(qr *QResponse, r interface{})) {
	q.router.panicHandler = ph
}

func (q *Qiwlar) Run(port string) {
	err := http.ListenAndServe(port, q.router)

	if err != nil {
		panic(err)
	}
}

func (q *Qiwlar) buildModule(metaModule interface{}) *module {
	realModule := module{meta: reflect.ValueOf(metaModule), store: q.store}
	realModule.collect(q.buildModule)

	return &realModule
}

func (q *Qiwlar) resolveDependencies(m *module) {
	if m.isActive {
		return
	}

	m.isActive = true

	for _, nestedModule := range m.imports {
		q.resolveDependencies(nestedModule)
	}

	for _, provider := range m.providers {
		if _, ok := provider.(Provide); ok {
			continue
		}

		bind := reflect.ValueOf(provider).MethodByName("Bind")
		params := []reflect.Value{}

	ILoop:
		for i := 0; i < bind.Type().NumIn(); i++ {
			for _, p := range m.providers {
				if bind.Type().In(i) == reflect.TypeOf(p) {
					params = append(params, reflect.ValueOf(p))
					continue ILoop
				}
			}

			for _, nestedModule := range m.imports {
				result := recursiveSearchDep(nestedModule, bind.Type().In(i))

				if result.IsValid() {
					params = append(params, result)
					continue ILoop
				}
			}
		}

		bind.Call(params)
	}

	for _, controller := range m.controllers {
		cValue := reflect.ValueOf(controller)
		bind := cValue.MethodByName("Bind")
		build := cValue.MethodByName("Build")
		params := []reflect.Value{}

	CLoop:
		for i := 0; i < bind.Type().NumIn(); i++ {
			for _, p := range m.providers {
				if bind.Type().In(i) == reflect.TypeOf(p) {
					params = append(params, reflect.ValueOf(p))
					continue CLoop
				}
			}

			for _, nestedModule := range m.imports {
				result := recursiveSearchDep(nestedModule, bind.Type().In(i))

				if result.IsValid() {
					params = append(params, result)
					continue CLoop
				}
			}
		}

		bind.Call(params)
		build.Call(nil)
		ctrl := reflect.ValueOf(controller).Elem().FieldByName("Controller").Interface().(Controller)

		for _, ep := range ctrl.EndPoints() {
			q.router.add(ep, ctrl.prefix)
		}
	}
}

func recursiveSearchDep(m *module, paramType reflect.Type) reflect.Value {
	modules := []*module{}

	for _, exportElem := range m.exports {
		if exportM, ok := exportElem.(*module); ok {
			modules = append(modules, exportM)
			continue
		}

		if paramType == reflect.TypeOf(exportElem) {
			return reflect.ValueOf(exportElem)
		}
	}

	for _, nestedModule := range modules {
		result := recursiveSearchDep(nestedModule, paramType)

		if result.IsValid() {
			return result
		}
	}

	return reflect.ValueOf(nil)
}

type Provides map[reflect.Type]Provide

type Exception struct {
	StatusCode int
	Data       interface{}
}

func (e *Exception) Error() string {
	return fmt.Sprintf("[ %v ], %v", e.StatusCode, e.Data)
}

type Provide struct {
	Name  string
	Value interface{}
}

type store struct {
	modules map[reflect.Type]*module
}

func (s *store) find(mType reflect.Type) *module {
	if fm, ok := s.modules[mType]; ok {
		return fm
	}

	return nil
}

func (s *store) add(m *module) {
	s.modules[m.meta.Type()] = m
}
