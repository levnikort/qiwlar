package qiwlar

import (
	"reflect"
)

type module struct {
	store       *store
	isActive    bool
	meta        reflect.Value
	imports     []*module
	providers   []interface{}
	controllers []interface{}
	exports     []interface{}
}

func (m *module) collect(buildModule func(interface{}) *module) []interface{} {
	if hookRegister := m.meta.MethodByName("UseRegister"); hookRegister.IsValid() {
		p := hookRegister.Call(nil)[0].Interface().(Provides)

		for key, value := range p {
			nestedModule := buildModule(reflect.New(key).Elem().Interface())

			if value.Name != "" {
				nestedModule.providers = append(nestedModule.providers, value)
			}

			m.imports = append(m.imports, nestedModule)
		}
	}

	if imports := m.meta.FieldByName("Imports"); imports.IsValid() {
		for i := 0; i < imports.Type().NumIn(); i++ {
			nestedModule := m.store.find(imports.Type().In(i))

			if nestedModule == nil {
				nestedModule = buildModule(reflect.New(imports.Type().In(i)).Elem().Interface())
				m.store.add(nestedModule)
			}

			m.imports = append(m.imports, nestedModule)
		}
	}

	if providers := m.meta.FieldByName("Providers"); providers.IsValid() {
		for i := 0; i < providers.Type().NumIn(); i++ {
			p := reflect.New(providers.Type().In(i)).Interface()
			m.providers = append(m.providers, p)
		}
	}

	if controllers := m.meta.FieldByName("Controllers"); controllers.IsValid() {
		for i := 0; i < controllers.Type().NumIn(); i++ {
			c := reflect.New(controllers.Type().In(i)).Interface()
			m.controllers = append(m.controllers, c)
		}
	}

	if exports := m.meta.FieldByName("Exports"); exports.IsValid() {
	ExportsLoop:
		for i := 0; i < exports.Type().NumIn(); i++ {
			for _, nestedModule := range m.imports {
				if exports.Type().In(i) == nestedModule.meta.Type() {
					m.exports = append(m.exports, nestedModule)
					continue ExportsLoop
				}
			}

			for _, provider := range m.providers {
				if exports.Type().In(i) == reflect.TypeOf(provider).Elem() {
					m.exports = append(m.exports, provider)
					continue ExportsLoop
				}
			}

			panic("invalid exports")
		}
	}

	return nil
}
