package decorator

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/levnikort/qiwlar"
)

var Validate = func() func(interface{}) qiwlar.Decorator {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return func(dto interface{}) qiwlar.Decorator {
		nDto := reflect.New(reflect.TypeOf(dto)).Interface()

		return func(q *qiwlar.QResponse) interface{} {
			b, err := io.ReadAll(q.R.Body)

			if err != nil {
				q.StatusCode = 400
				q.Results = append(q.Results, nDto)
				return "failed to read request body"
			}

			err = json.Unmarshal(b, nDto)

			if err != nil {
				q.StatusCode = 400
				q.Results = append(q.Results, nDto)
				return "unable to convert"
			}

			err = validate.Struct(nDto)

			if err != nil {
				q.StatusCode = 400
				q.Results = append(q.Results, nDto)

				if _, ok := err.(*validator.InvalidValidationError); ok {
					fmt.Println(err)
					return err
				}

				var errors []map[string]interface{}

				for _, err := range err.(validator.ValidationErrors) {
					errors = append(errors, map[string]interface{}{
						"field": err.Field(),
						"tag":   err.Tag(),
						"param": err.Param(),
					})
				}

				return errors
			}

			q.Results = append(q.Results, reflect.ValueOf(nDto).Elem().Interface())

			return nil
		}
	}
}()
