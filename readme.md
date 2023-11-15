# Qiwlar framework

A handy Web API development tool inspired by [NestJS](https://nestjs.com/).
Uses [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) as a router

## Установка

```shell
go get github.com/levnikort/qiwlar
```

## Модули

Модуль состоит из четырех полей:

1. `Imports` - поле необходимо для подулючения других модулей, и доступа к их экспортируемым элементам;
2. `Providers` - 
3. `Controllers` -
4. `Exports` - 

### Создание модуля

```go
type AppModule struct {
  Imports     func(foo.FooModule{})
  Providers   func(AppProvider)
  Controllers func(AppController)
  Exports     func(AppProvider, foo.FooModule{})
}
```

or

```go
type AppModule struct {
  Controllers func(AppController)
}
```

```go
type AppModule struct {
  Controllers func(AppController)
}

func (am AppModule) UseRegister() qiwlar.Provides {
  return qiwlar.Provides{
    reflect.TypeOf(foo.FooModule{}): {
      Name: "FOO"
      Value: map[string]interface{}{"opt": false} 
    }
  }
}
```

## Поставщики

### Создание поставщика

```go
type AppProvider struct {}
func (ap *AppProvider) Bind() {}
```

```go
type AppProvider struct {
  fp *foo.FooProvider
}

func (ap *AppProvider) Bind(fp *foo.FooProvider) {
  ap.fp = fp
}
```

```go
type AppProvider struct {
  fp *foo.FooProvider
}

func (ap *AppProvider) Bind(fp *foo.FooProvider) {
  ap.fp = fp

  db.DBConnect("Data")
}
```

```go
type FooProvider struct {
  settings *FooSettings
}

func (ap *AppProvider) Bind(s qiwlar.Provide) {
  ap.settings = s.Value.(*FooSettings)
}
```

## Контроллеры

### Создание контролерра