package envconfig

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	timeType     = reflect.TypeOf(time.Time{})
	durationType = reflect.TypeOf(time.Duration(0))
)

func Load(cfg any) error {
	value := reflect.ValueOf(cfg)

	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("argument'type should be pointer. Kind = %s", value.Kind().String())
	}
	value = value.Elem()
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("argument'type should be pointer. Kind = %s", value.Kind().String())
	}

	for i := range value.NumField() {
		field := value.Field(i)
		structField := value.Type().Field(i)

		envName := toUpperSnakeCase(structField.Name)
		envValue, ok := os.LookupEnv(envName)
		if !ok {
			continue
		}

		switch structField.Type.Kind() {
		case reflect.Bool:
			v, err := strconv.ParseBool(envValue)
			if err != nil {
				return fmt.Errorf("fail to the convert environment value to bool. EnvName = %s, EnvValue = %s", envName, envValue)
			}
			field.SetBool(v)
		case reflect.Int32:
			v, err := strconv.ParseInt(envValue, 10, 32)
			if err != nil {
				return fmt.Errorf("fail to the convert environment value to int32. EnvName = %s, EnvValue = %s", envName, envValue)
			}
			field.SetInt(v)
		case reflect.Int64:
			if structField.Type == durationType {
				v, err := time.ParseDuration(envValue)
				if err != nil {
					return fmt.Errorf("fail to the convert environment value to time.Duration. EnvName = %s, EnvValue = %s", envName, envValue)
				}
				field.SetInt(int64(v))
			} else {
				v, err := strconv.ParseInt(envValue, 10, 64)
				if err != nil {
					return fmt.Errorf("fail to the convert environment value to int64. EnvName = %s, EnvValue = %s", envName, envValue)
				}
				field.SetInt(v)
			}
		case reflect.Float64:
			v, err := strconv.ParseFloat(envValue, 64)
			if err != nil {
				return fmt.Errorf("fail to the convert environment value to float64. EnvName = %s, EnvValue = %s", envName, envValue)
			}
			field.SetFloat(v)
		case reflect.String:
			field.SetString(envValue)
		case reflect.Struct:
			if structField.Type == timeType {
				t, err := time.ParseInLocation(time.RFC3339, envValue, time.Local)
				if err != nil {
					return fmt.Errorf("fail to the convert environment value to time.RFC3339. EnvName = %s, EnvValue = %s", envName, envValue)
				}
				field.Set(reflect.ValueOf(t))
			}
		case reflect.Slice:
			strs := strings.Split(envValue, ",")
			switch field.Type().Elem().Kind() {
			case reflect.Bool:
				s := make([]bool, 0, len(strs))
				for _, e := range strs {
					v, err := strconv.ParseBool(e)
					if err != nil {
						return fmt.Errorf("fail to the convert environment value to bool. EnvName = %s, EnvValue = %s", envName, envValue)
					}
					s = append(s, v)
				}
				field.Set(reflect.ValueOf(s))
			case reflect.Int32:
				s := make([]int32, 0, len(strs))
				for _, e := range strs {
					v, err := strconv.ParseInt(e, 10, 32)
					if err != nil {
						return fmt.Errorf("fail to the convert environment value to int32. EnvName = %s, EnvValue = %s", envName, envValue)
					}
					s = append(s, int32(v))
				}
				field.Set(reflect.ValueOf(s))
			case reflect.Int64:
				if field.Type().Elem() == durationType {
					s := make([]time.Duration, 0, len(strs))
					for _, e := range strs {
						v, err := time.ParseDuration(e)
						if err != nil {
							return fmt.Errorf("fail to the convert environment value to time.Duration. EnvName = %s, EnvValue = %s", envName, envValue)
						}
						s = append(s, v)
					}
					field.Set(reflect.ValueOf(s))
				} else {
					s := make([]int64, 0, len(strs))
					for _, e := range strs {
						v, err := strconv.ParseInt(e, 10, 64)
						if err != nil {
							return fmt.Errorf("fail to the convert environment value to int64. EnvName = %s, EnvValue = %s", envName, envValue)
						}
						s = append(s, v)
					}
					field.Set(reflect.ValueOf(s))
				}
			case reflect.Float64:
				s := make([]float64, 0, len(strs))
				for _, e := range strs {
					v, err := strconv.ParseFloat(e, 64)
					if err != nil {
						return fmt.Errorf("fail to the convert environment value to float64. EnvName = %s, EnvValue = %s", envName, envValue)
					}
					s = append(s, v)
				}
				field.Set(reflect.ValueOf(s))
			case reflect.String:
				field.Set(reflect.ValueOf(strs))
			case reflect.Invalid,
				reflect.Int, reflect.Int8, reflect.Int16,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Uintptr, reflect.Float32, reflect.Complex64, reflect.Complex128,
				reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
				reflect.Pointer, reflect.Slice, reflect.Struct, reflect.UnsafePointer:
			}
		case reflect.Invalid,
			reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Uintptr, reflect.Float32, reflect.Complex64, reflect.Complex128,
			reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
			reflect.Pointer, reflect.UnsafePointer:
		}
	}
	return nil
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toUpperSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = strings.ToUpper(snake)

	snake = strings.ReplaceAll(snake, "MY_SQL", "MYSQL")   // MySQL -> MY_SQL -> MYSQL
	snake = strings.ReplaceAll(snake, "O_AUTH2", "OAUTH2") // OAuth2 -> O_AUTH2 -> OAUTH2
	return snake
}
