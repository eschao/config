package cli

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/eschao/config/utils"
)

type anyValue struct {
	any reflect.Value
}

func newAnyValue(v reflect.Value) *anyValue {
	return &anyValue{any: v}
}

func (this *anyValue) String() string {
	kind := this.any.Kind()
	switch kind {
	case reflect.Bool:
		return strconv.FormatBool(this.any.Bool())
	case reflect.String:
		return this.any.String()
	case reflect.Int8,
		reflect.Int16,
		reflect.Int,
		reflect.Int32,
		reflect.Int64:
		return strconv.FormatInt(this.any.Int(), 10)
	case reflect.Uint8,
		reflect.Uint16,
		reflect.Uint,
		reflect.Uint32,
		reflect.Uint64:
		return strconv.FormatUint(this.any.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(this.any.Float(), 'E', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(this.any.Float(), 'E', -1, 64)
	}
	return fmt.Sprintf("unsupport type %s", kind.String())
}

func (this *anyValue) Set(v string) error {
	kind := this.any.Kind()
	switch kind {
	case reflect.String:
		this.any.SetString(v)
	case reflect.Float32:
		return utils.SetValueWithFloatX(this.any, v, 32)
	case reflect.Float64:
		return utils.SetValueWithFloatX(this.any, v, 64)
	case reflect.Int8:
		return utils.SetValueWithIntX(this.any, v, 8)
	case reflect.Int16:
		return utils.SetValueWithIntX(this.any, v, 16)
	case reflect.Int, reflect.Int32:
		return utils.SetValueWithIntX(this.any, v, 32)
	case reflect.Int64:
		return utils.SetValueWithIntX(this.any, v, 64)
	case reflect.Uint8:
		return utils.SetValueWithUintX(this.any, v, 8)
	case reflect.Uint16:
		return utils.SetValueWithUintX(this.any, v, 16)
	case reflect.Uint, reflect.Uint32:
		return utils.SetValueWithUintX(this.any, v, 32)
	case reflect.Uint64:
		return utils.SetValueWithUintX(this.any, v, 64)
	default:
		return fmt.Errorf("Can't support type %s", kind.String())
	}

	return nil
}

var errorHandling = flag.ExitOnError

type UsageFunc func(*Command) func()

var usageHandler UsageFunc = nil

type Command struct {
	Name        string
	FlagSet     *flag.FlagSet
	Usage       string
	SubCommands map[string]*Command
}

func New(name string) *Command {
	cmd := Command{
		Name:        name,
		FlagSet:     flag.NewFlagSet(name, errorHandling),
		SubCommands: make(map[string]*Command),
	}

	return &cmd
}

func NewWith(name string, errHandling flag.ErrorHandling,
	usageHandling UsageFunc) *Command {
	errorHandling = errHandling
	usageHandler = usageHandling

	cmd := Command{
		Name:        name,
		FlagSet:     flag.NewFlagSet(name, errorHandling),
		SubCommands: make(map[string]*Command),
	}

	if usageHandler != nil {
		cmd.FlagSet.Usage = usageHandler(&cmd)
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
	var err error

	for i := 0; i < v.NumField() && err == nil; i++ {
		valueOfField := v.Field(i)
		kindOfField := valueOfField.Kind()
		structOfField := typeOfStruct.Field(i)

		if kindOfField == reflect.Ptr {
			if !valueOfField.IsNil() && valueOfField.CanSet() {
				cmd := this.createSubCommand(structOfField.Tag)
				err = cmd.Init(valueOfField.Interface())
			}
		} else if kindOfField == reflect.Struct {
			cmd := this.createSubCommand(structOfField.Tag)
			err = cmd.parseValue(valueOfField)
		} else {
			err = this.addFlag(valueOfField, structOfField)
		}
	}

	return err
}

func (this *Command) addFlag(v reflect.Value, f reflect.StructField) error {
	cmdTag, ok := f.Tag.Lookup("cmd")
	if !ok || cmdTag == "" {
		return nil
	}

	firstSpace := strings.Index(cmdTag, " ")
	name := cmdTag
	usage := ""
	if firstSpace > 0 {
		name = cmdTag[0:firstSpace]
		usage = cmdTag[firstSpace+1:]
	}

	kind := v.Kind()
	switch kind {
	case reflect.Bool:
		this.FlagSet.BoolVar((*bool)(unsafe.Pointer(v.UnsafeAddr())), name,
			false, usage)
		return nil
	case reflect.String,
		reflect.Int8,
		reflect.Int16,
		reflect.Int,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64:
		anyValue := newAnyValue(v)
		this.FlagSet.Var(anyValue, name, usage)
	default:
		return fmt.Errorf("Can't support type %s", kind.String())
	}

	return nil
}

func (this *Command) createSubCommand(tag reflect.StructTag) *Command {
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
	cmd.FlagSet = flag.NewFlagSet(name, errorHandling)
	cmd.Usage = usage

	if usageHandler != nil {
		cmd.FlagSet.Usage = usageHandler(&cmd)
	}

	this.SubCommands[name] = &cmd
	return &cmd
}

func (this *Command) Parse(args []string) error {
	if err := this.FlagSet.Parse(args); err != nil {
		return err
	}

	unprocessed := this.FlagSet.Args()
	if len(unprocessed) < 1 {
		return nil
	}

	if this.SubCommands == nil {
		return fmt.Errorf("Command: %s is unsupport", unprocessed[0])
	}

	cmd := this.SubCommands[unprocessed[0]]
	if cmd == nil {
		return fmt.Errorf("Command: %s is unsupport", unprocessed[0])
	}

	return cmd.Parse(unprocessed[1:])
}
