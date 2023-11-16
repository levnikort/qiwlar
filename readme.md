# Qiwlar

Этот фреймворк вдохновлен [NestJS](https://nestjs.com/) и предназначен для создания легко масштабируемых
серверных приложений. Он использует стандартный пакет **net/http** для обработки запросов
и [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) для маршрутизации.

## Установка

```shell
go get github.com/levnikort/qiwlar
```

## Модули

**Модуль** — это структура, содержащая хотя бы одно из четырех полей:
1. `Imports` — список импортированных модулей, которые экспортируют необходимые провайдеры
или другие модули.
2. `Providers` — поле, необходимое для регистрации провайдеров, чтобы в дальнейшем передавать
их другим элементам, которые этого требуют.
3. `Controllers` — набор контроллеров, определенных в этом модуле, экземпляры которых
необходимо создать.
4. `Exports` — список провайдеров или модулей, которые будут доступны в других модулях.

```go
type AppModule struct {
  Imports     func(foo.FooModule{})
  Providers   func(AppProvider)
  Controllers func(AppController)
  Exports     func(AppProvider, foo.FooModule{})
}
```

или

```go
type AppModule struct {
  Controllers func(AppController)
}
```


У модуля можно реализовать метод `UseRegister`, необходимый для регистрации
изолированного модуля. Проще говоря, создается абсолютно новая копия модуля, в которую также
можно передать тип `qiwlar.Provide`. Этот тип используется для передачи каких-либо настроек, например, адреса
для подключения к базе данных. `qiwlar.Provide` передается в поле Provider создаваемого модуля.

**Примечание**: Доступ к `qiwlar.Provide` можно получить из метода `Bind()` в провайдере или контроллере.

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

В приложении обязательно должен быть корневой модуль, с которого начинается построение
графа и в дальнейшем внедрение зависимостей.

```go
// Точка входа

func main() {
  q := qiwlar.New(app.AppModule{}).Run(":8080")
}
```

## Провайдеры

**Провайдер** — это структура, которая должна содержать метод `Bind()`.
Основная идея провайдера заключается в том, что он может быть внедрен как зависимость.

Метод `Bind()` в провайдере и контроллере необходим для распределения зависимостей внутри структуры,
а также его можно использовать как функцию инициализации. То есть при создании провайдера или контроллера
можно выполнить какой-либо дополнительный код.

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

  // Выполнение дополнительного кода.
  db.DBConnect("Data")
}
```

## Контроллеры

**Контроллер** — это структура, которая должна включать в себя `qiwlar.Controller` и содержать методы `Bind()`, `Build()`.

Цель контроллера — обрабатывать определенные запросы для приложения. Механизм маршрутизации управляет
тем, какой из контроллеров обрабатывает какие запросы. Зачастую каждый контроллер имеет более одного маршрута,
и разные маршруты могут выполнять разные действия.

Метод `Build()` необходим для начала сборки контроллера.

Метод `Collect()` собирает конечные точки из декораторов и функций.

Метод `Prefix()` добавляет префикс ко всем путям в контроллере.

Метод `ForAll()` применяет декораторы ко всем конечным точкам в контроллере.

```go
type AppController struct {
  qiwlar.Controller
  as *AppService
}

func (ac *AppController) Bind(as *AppService) {
  ac.as = as
}

func (ac *AppController) Build() {
  ac.Collect([]interface{}{
    ac.Prefix("v1"),
    ac.ForAll(d.MyLogger())

    d.Post("test/:id")
    d.Validate(UserDTO{})
    d.Params()
    func(dto UserDTO, params httprouter.Params) string {
      return params.ByName("id")
    }

    d.Get("cats")
    func() any {
      return "no"
    }

    // Вернуть ошибку, с http кодом 500
    d.Get("error")
    func() any {
      return qiwlar.Exception{StatusCode: 500, Data: "server error"}
    }
  })
}
```

## Декораторы

Декоратор — это функция, которая реализует тип `qiwlar.Decorator`.
Декоратор представляет собой, по сути, промежуточное программное обеспечение.

Пример простой реализации декоратора выглядит следующим образом:

```go
func MyLoggerMethod() qiwlar.Decorator {
  return func(q *qiwlar.QResponse) interface{} {
    fmt.Println(q.Method)
    
    return nil
  }
}
```

Если мы хотим, чтобы наш декоратор предоставлял какие-то данные конечной функции:

```go
type CurrentMethod string

func MyLoggerMethod() qiwlar.Decorator {
  return func(q *qiwlar.QResponse) interface{} {
    fmt.Println(q.Method)

    q.Results = append(Results, CurrentMethod(q.Method))

    return nil
  }
}
```

Если в декораторе может произойти ошибка, то мы можем обработать ее следующим образом:
если декоратор предоставляет какие-либо данные, то до выхода из функции и возврата ошибки,
он должен передать в `q.Results` пустой тип данных.

```go
type CurrentMethod string

func MyLoggerMethod() qiwlar.Decorator {
  return func(q *qiwlar.QResponse) interface{} {
    fmt.Println(q.Method)

    err := logic(...)

    if err != nil {
      q.StatusCode = 400
      q.Results = append(q.Results, CurrentMethod(""))
      return err
    }

    q.Results = append(Results, CurrentMethod(q.Method))

    return nil
  }
}
```