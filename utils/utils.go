/*
 * Copyright (C) 2017 eschao <esc.chao@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
	value, err := strconv.ParseInt(intValue, 10, bitSize)
	if err != nil {
		return err
	}

	v.SetInt(value)
	return nil
}

func SetValueWithUintX(v reflect.Value, uintValue string, bitSize int) error {
	value, err := strconv.ParseUint(uintValue, 10, bitSize)
	if err != nil {
		return err
	}

	v.SetUint(value)
	return nil
}

func SetValueWithSlice(v reflect.Value, slice string, separator string) error {
	data := strings.Split(slice, separator)
	size := len(data)
	if size > 0 {
		slice := reflect.MakeSlice(v.Type(), size, size)
		for i := 0; i < size; i++ {
			ele := slice.Index(i)
			kind := ele.Kind()
			var err error
			switch kind {
			case reflect.Bool:
				err = SetValueWithBool(ele, data[i])
			case reflect.String:
				ele.SetString(data[i])
			case reflect.Uint8:
				err = SetValueWithUintX(ele, data[i], 8)
			case reflect.Uint16:
				err = SetValueWithUintX(ele, data[i], 16)
			case reflect.Uint, reflect.Uint32:
				err = SetValueWithUintX(ele, data[i], 32)
			case reflect.Uint64:
				err = SetValueWithUintX(ele, data[i], 64)
			case reflect.Int8:
				err = SetValueWithIntX(ele, data[i], 8)
			case reflect.Int16:
				err = SetValueWithIntX(ele, data[i], 16)
			case reflect.Int, reflect.Int32:
				err = SetValueWithIntX(ele, data[i], 32)
			case reflect.Int64:
				err = SetValueWithIntX(ele, data[i], 64)
			case reflect.Float32:
				err = SetValueWithFloatX(ele, data[i], 32)
			case reflect.Float64:
				err = SetValueWithFloatX(ele, data[i], 64)
			default:
				return fmt.Errorf("Can't support type: %s", kind.String())
			}

			if err != nil {
				return err
			}
		}

		v.Set(slice)
	}

	return nil
}
