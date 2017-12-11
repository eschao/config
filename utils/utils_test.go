package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type data struct {
	BoolValue    bool
	Int8Value    int8
	IntValue     int
	Uint16Value  uint16
	UintValue    uint
	Float32Value float32
	Names        []string
}

func TestSetValueWithBool(t *testing.T) {
	d := data{}
	v := reflect.ValueOf(&d).Elem().FieldByName("BoolValue")

	assert := assert.New(t)
	assert.NoError(SetValueWithBool(v, "true"))
	assert.Equal(true, d.BoolValue)
}

func TestSetValueWithFloat32(t *testing.T) {
	d := data{}
	v := reflect.ValueOf(&d).Elem().FieldByName("Float32Value")

	assert := assert.New(t)
	assert.NoError(SetValueWithFloatX(v, "123.456", 32))
	assert.Equal(float32(123.456), d.Float32Value)
}

func TestSetValueWithInt8(t *testing.T) {
	d := data{}
	v := reflect.ValueOf(&d).Elem().FieldByName("Int8Value")

	assert := assert.New(t)
	assert.NoError(SetValueWithIntX(v, "10", 8))
	assert.Equal(int8(10), d.Int8Value)
}

func TestSetValueWithInt(t *testing.T) {
	d := data{}
	v := reflect.ValueOf(&d).Elem().FieldByName("IntValue")

	assert := assert.New(t)
	assert.NoError(SetValueWithIntX(v, "10000", 32))
	assert.Equal(10000, d.IntValue)
}

func TestSetValueWithUint16(t *testing.T) {
	d := data{}
	v := reflect.ValueOf(&d).Elem().FieldByName("Uint16Value")

	assert := assert.New(t)
	assert.NoError(SetValueWithUintX(v, "100", 16))
	assert.Equal(uint16(100), d.Uint16Value)
}

func TestSetValueWithUint(t *testing.T) {
	d := data{}
	v := reflect.ValueOf(&d).Elem().FieldByName("UintValue")

	assert := assert.New(t)
	assert.NoError(SetValueWithUintX(v, "2000", 32))
	assert.Equal(uint(2000), d.UintValue)
}

func TestSetValueWithSlice(t *testing.T) {
	d := data{}
	v := reflect.ValueOf(&d).Elem().FieldByName("Names")

	assert := assert.New(t)
	assert.NoError(SetValueWithSlice(v, "xx:yy:zz", ":"))
	assert.Equal(3, len(d.Names))
	assert.Equal("xx", d.Names[0])
	assert.Equal("yy", d.Names[1])
	assert.Equal("zz", d.Names[2])
}
