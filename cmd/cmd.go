package cmd

import (
	"flag"
	"fmt"
	"reflect"
	"strings"

	"github.com/eschao/config/util"
)

type anyValue struct {
	any reflect.Value
}

func newAnyValue(v Value) *anyValue {
	return anyValue{any: v}
}

func (this *anyValue) String() string {
	return this.any.Kind().String()
}

func (this *anyValue) Set(v string) error {
	kind := this.any.Kind()
	switch kind {
	case reflect.Bool:

	}
	return nil
}

type Command struct {
	Name        string
	FlagSet     *flag.FlagSet
	Usage       string
	SubCommands map[string]*Command
}

func New(name string) *Command {
	cmd := Command{
		Name:        name,
		FlagSet:     flag.NewFlagSet(name, flag.ExitOnError),
		SubCommands: make(map[string]*Command),
	}
	return &cmd
}

func (this *Command) Init(i interface{}) error {
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

	return this.parseValue(valueOfStruct)
}

func (this *Command) parseValue(v reflect.Value) error {
	typeOfStruct := v.Type()

	for i := 0; i < v.NumField(); i++ {
		valueOfField := v.Field(i)
		kindOfField := valueOfField.Kind()
		structOfField := typeOfStruct.Field(i)

		if kindOfField == reflect.Ptr {
			if !valueOfField.IsNil() && valueOfField.CanSet() {
				cmd := this.createCliFlagSet(structOfField.Tag)
				if err := cmd.Init(valueOfField.Interface()); err != nil {
					return err
				}
			}
		} else if kindOfField == reflect.Struct {
			cmd := this.createCliFlagSet(structOfField.Tag)
			if err := cmd.parseValue(valueOfField); err != nil {
				return err
			}
		} else {
			this.addFlag(valueOfField, structOfField)
		}
	}

	return nil
}

func (this *Command) addFlag(v reflect.Value, f reflect.StructField) {
	cmdTag, ok := f.Tag.Lookup("cmd")
	if !ok || cmdTag == "" {
		return
	}

	firstSpace := strings.Index(cmdTag, " ")
	name := cmdTag
	usage := ""
	if firstSpace > 0 {
		name = cmdTag[0:firstSpace]
		usage = cmdTag[firstSpace+1:]
	}

	//defValue, ok := f.Tag.Lookup("default")
	vFlag := ValueFlag{Value: v}
	this.FlagSet.Var(&vFlag, name, usage)
	//fmt.Printf("[%s]: Added Flag: %s\n", this.Name, name)
}

func (this *Command) createCliFlagSet(tag reflect.StructTag) *Command {
	cmdTag, ok := tag.Lookup("cmd")
	if !ok || cmdTag == "" {
		return this
	}

	cmd := Command{SubCommands: make(map[string]*Command)}
	firstSpace := strings.Index(cmdTag, " ")
	name := cmdTag
	usage := ""
	if firstSpace > 0 {
		name = cmdTag[0:firstSpace]
		usage = cmdTag[firstSpace+1:]
	}

	cmd.Name = name
	cmd.FlagSet = flag.NewFlagSet(name, flag.ExitOnError)
	cmd.Usage = usage
	this.SubCommands[name] = &cmd
	return &cmd
}

func (this *Command) Parse(i interface{}, args []string) error {

	return nil
}
