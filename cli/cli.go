package cli

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
)

type CliFlag struct {
	Name     string
	Usage    string
	DefValue *string
	CliFlags *[]CliFlag
}

type ValueFlag struct {
	Value reflect.Value
}

func (this *ValueFlag) String() string {
	return this.Value.Kind().String()
}

func (this *ValueFlag) Set(v string) error {
	return nil
}

func Parse(i interface{}, cliFlag *CliFlag) error {
	ptrRef := reflect.ValueOf(i)

	if ptrRef.IsNil() || ptrRef.Kind() != reflect.Ptr {
		return fmt.Errorf("Expect a structure pointer type instead of %s",
			ptrRef.Kind().String())
	}

	valueOfStruct := ptrRef.Elem()
	if valueOfStruct.Kind() != reflect.Struct {
		return fmt.Errorf("Expect a structure type instead of %s",
			valueOfStruct.Kind().String())
	}

	return parseValue(valueOfStruct, cliFlag)
}

func parseValue(v reflect.Value, cliParent *CliFlag) error {
	cliFlags := []CliFlag{}
	typeOfStruct := v.Type()

	for i := 0; i < v.NumField(); i++ {
		valueOfField := v.Field(i)
		kindOfField := valueOfField.Kind()
		structOfField := typeOfStruct.Field(i)

		if kindOfField == reflect.Ptr {
			if !valueOfField.IsNil() && valueOfField.CanSet() {
				cliFlag := createCliFlagFromTag(structOfField.Tag)
				if cliFlag != nil {
					cliFlags = append(cliFlags, *cliFlag)
				} else {
					cliFlag = cliParent
				}
				Parse(valueOfField.Interface(), cliFlag)
			} else {
				continue
			}
		} else if kindOfField == reflect.Struct {
			cliFlag := createCliFlagFromTag(structOfField.Tag)
			if cliFlag != nil {
				cliFlags = append(cliFlags, *cliFlag)
			} else {
				cliFlag = cliParent
			}
			parseValue(valueOfField, cliFlag)
		}

		if cliFlag := installFlag(valueOfField, structOfField); cliFlag != nil {
			cliFlags = append(cliFlags, *cliFlag)
		}
	}

	if len(cliFlags) > 0 {
		cliParent.CliFlags = &cliFlags
	}
	return nil
}

func installFlag(v reflect.Value, f reflect.StructField) *CliFlag {
	cliFlag := createCliFlagFromTag(f.Tag)
	if cliFlag != nil {
		vFlag := ValueFlag{Value: v}
		flag.Var(&vFlag, cliFlag.Name, cliFlag.Usage)
	}
	fmt.Printf("Installed flag: %d", cliFlag.Name)
	return cliFlag
}

func createCliFlagFromTag(tag reflect.StructTag) *CliFlag {
	cliTag, ok := tag.Lookup("cli")
	if !ok || cliTag == "" {
		return nil
	}

	cliFlag := CliFlag{}
	firstSpace := strings.Index(cliTag, " ")
	cliFlag.Name = cliTag
	if firstSpace > 0 {
		cliFlag.Name = cliTag[0:firstSpace]
		cliFlag.Usage = cliTag[firstSpace+1:]
	}

	defValue, ok := tag.Lookup("default")
	if !ok {
		cliFlag.DefValue = &defValue
	}

	return &cliFlag
}
