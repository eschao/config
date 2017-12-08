package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func Unmarshal(i interface{}) error {
	ptrRef := reflect.ValueOf(i)

	if ptrRef.IsNil() || ptrRef.Kind() != reflect.Ptr {
		return fmt.Errorf("Expect a structure pointer type instead of %s",
			ptrRef.Kind().String())
	}

	valueOfStruct := ptrRef.Elem()
	if valueOfStruct.Kind() != reflect.Struct {
		return fmt.Errorf("Expect a structure pointer type instead of %s",
			valueOfStruct.Kind().String())
	}

	return unmarshal(valueOfStruct)
}

func unmarshal(v reflect.Value) error {
	typeOfStruct := v.Type()
	for i := 0; i < v.NumField(); i++ {
		valueOfField := v.Field(i)
		kindOfField := valueOfField.Kind()
		structOfField := typeOfStruct.Field(i)
		//fmt.Printf("Name: %s, Type: %s\n", structOfField.Name, kindOfField.String())

		// recursively unmarshal if value is ptr type
		if kindOfField == reflect.Ptr {
			if !valueOfField.IsNil() && valueOfField.CanSet() {
				Unmarshal(valueOfField.Interface())
			} else {
				continue
			}
		} else if kindOfField == reflect.Struct {
			unmarshal(valueOfField)
		}

		if err := setFieldValue(valueOfField, structOfField); err != nil {
			return err
		}
	}

	return nil
}

func getEnvValue(envName string, f reflect.StructField) (string, bool) {
	envValue, ok := os.LookupEnv(envName)
	if !ok {
		envValue, ok = f.Tag.Lookup("default")
	}

	return envValue, ok
}

func setFieldValue(v reflect.Value, f reflect.StructField) error {
	envName := f.Tag.Get("env")
	if envName == "" {
		return nil
	}

	envValue, ok := getEnvValue(envName, f)
	if !ok {
		return nil
	}

	if !v.CanSet() {
		return fmt.Errorf("%s: can't be set", f.Name)
	}

	kind := v.Kind()
	name := f.Name
	switch kind {
	case reflect.Bool:
		return setFieldValueWithBool(name, v, envValue)

	case reflect.String:
		v.SetString(envValue)
		return nil

	case reflect.Int8:
		return setFieldValueWithIntX(name, v, envValue, 8)

	case reflect.Int16:
		return setFieldValueWithIntX(name, v, envValue, 16)

	case reflect.Int, reflect.Int32:
		return setFieldValueWithIntX(name, v, envValue, 32)

	case reflect.Int64:
		return setFieldValueWithIntX(name, v, envValue, 64)

	case reflect.Uint8:
		return setFieldValueWithUintX(name, v, envValue, 8)

	case reflect.Uint16:
		return setFieldValueWithUintX(name, v, envValue, 16)

	case reflect.Uint, reflect.Uint32:
		return setFieldValueWithUintX(name, v, envValue, 32)

	case reflect.Uint64:
		return setFieldValueWithUintX(name, v, envValue, 64)

	case reflect.Float32:
		return setFieldValueWithFloatX(name, v, envValue, 32)

	case reflect.Float64:
		return setFieldValueWithFloatX(name, v, envValue, 64)

	case reflect.Slice:
		sp, ok := f.Tag.Lookup("separator")
		if !ok {
			sp = ":"
		}
		return setFieldValueWithSlice(name, v, envValue, sp)

	default:
		return fmt.Errorf("Can't support type: %s", kind.String())
	}
}

func setFieldValueWithBool(name string, v reflect.Value,
	envValue string) error {
	value, err := strconv.ParseBool(envValue)
	if err != nil {
		return fmt.Errorf("%s: can't convert %s to bool value. %s", name, envValue,
			err.Error())
	}

	v.SetBool(value)
	return nil
}

func setFieldValueWithFloatX(name string, v reflect.Value, envValue string,
	bitSize int) error {
	value, err := strconv.ParseFloat(envValue, bitSize)
	if err != nil {
		return fmt.Errorf("%s: can't convert %s to float%d value. %s", name,
			envValue, bitSize, err.Error())
	}

	v.SetFloat(value)
	return nil
}

func setFieldValueWithIntX(name string, v reflect.Value, envValue string,
	bitSize int) error {
	value, err := strconv.ParseInt(envValue, 10, bitSize)
	if err != nil {
		return fmt.Errorf("%s: can't convert %s to int%d value. %s", name,
			envValue, bitSize, err.Error())
	}

	v.SetInt(value)
	return nil
}

func setFieldValueWithUintX(name string, v reflect.Value, envValue string,
	bitSize int) error {
	value, err := strconv.ParseUint(envValue, 10, bitSize)
	if err != nil {
		return fmt.Errorf("%s: can't convert %s to uint%d value. %s", name,
			envValue, bitSize, err.Error())
	}

	v.SetUint(value)
	return nil
}

func setFieldValueWithSlice(name string, v reflect.Value, envValue string,
	separator string) error {
	data := strings.Split(envValue, separator)
	size := len(data)
	if size > 0 {
		slice := reflect.MakeSlice(v.Type(), size, size)
		for i := 0; i < size; i++ {
			ele := slice.Index(i)
			kind := ele.Kind()
			switch kind {
			case reflect.Bool:
				if err := setFieldValueWithBool(name, ele, data[i]); err != nil {
					return err
				}
			case reflect.String:
				ele.SetString(data[i])
			case reflect.Uint8:
				if err := setFieldValueWithUintX(name, ele, data[i], 8); err != nil {
					return err
				}
			case reflect.Uint16:
				if err := setFieldValueWithUintX(name, ele, data[i], 16); err != nil {
					return err
				}
			case reflect.Uint, reflect.Uint32:
				if err := setFieldValueWithUintX(name, ele, data[i], 32); err != nil {
					return err
				}
			case reflect.Uint64:
				if err := setFieldValueWithUintX(name, ele, data[i], 64); err != nil {
					return err
				}
			case reflect.Int8:
				if err := setFieldValueWithIntX(name, ele, data[i], 8); err != nil {
					return err
				}
			case reflect.Int16:
				if err := setFieldValueWithIntX(name, ele, data[i], 16); err != nil {
					return err
				}
			case reflect.Int, reflect.Int32:
				if err := setFieldValueWithIntX(name, ele, data[i], 32); err != nil {
					return err
				}
			case reflect.Int64:
				if err := setFieldValueWithIntX(name, ele, data[i], 64); err != nil {
					return err
				}
			case reflect.Float32:
				if err := setFieldValueWithFloatX(name, ele, data[i], 32); err != nil {
					return err
				}
			case reflect.Float64:
				if err := setFieldValueWithFloatX(name, ele, data[i], 64); err != nil {
					return err
				}
			default:
				return fmt.Errorf("%s: can't support type: %s", name, kind.String())
			}
		}
		v.Set(slice)
	}

	return nil
}
