package util

import (
	"reflect"
	"strconv"
)

func SetValueWithBool(v reflect.Value, boolValue string) error {
	value, err := strconv.ParseBool(boolValue)
	if err != nil {
		return err
	}

	v.SetBool(value)
	return nil
}

func SetValueWithFloatX(v reflect.Value, floatValue string, bitSize int) error {
	value, err := strconv.ParseFloat(floatValue, bitSize)
	if err != nil {
		return err
	}

	v.SetFloat(value)
	return nil
}

func SetValueWithIntX(v reflect.Value, intValue string, bitSize int) error {
	value, err := strconv.ParseInt(envValue, 10, bitSize)
	if err != nil {
		return err
	}

	v.SetInt(value)
	return nil
}

func SetValueWithUintX(v reflect.Value, envValue string, bitSize int) error {
	value, err := strconv.ParseUint(envValue, 10, bitSize)
	if err != nil {
		return err
	}

	v.SetUint(value)
	return nil
}

/*
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
}*/
