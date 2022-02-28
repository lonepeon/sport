package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type FieldType int

const (
	FieldTypeInt FieldType = iota
	FieldTypeString
	FieldTypeStrings
	FieldTypeFloat64
)

type Field struct {
	Name    string
	Type    FieldType
	Int     int
	String  string
	Strings []string
	Float64 float64
}

func (f Field) ToZapField() zap.Field {
	switch f.Type {
	case FieldTypeInt:
		return zap.Int(f.Name, f.Int)
	case FieldTypeString:
		return zap.String(f.Name, f.String)
	case FieldTypeStrings:
		return zap.Strings(f.Name, f.Strings)
	case FieldTypeFloat64:
		return zap.Float64(f.Name, f.Float64)
	default:
		panic(fmt.Sprintf("unsupported field type: %d", f.Type))
	}
}

func String(name string, value string) Field {
	return Field{Name: name, String: value, Type: FieldTypeString}
}

func Strings(name string, value []string) Field {
	return Field{Name: name, Strings: value, Type: FieldTypeStrings}
}

func Int(name string, value int) Field {
	return Field{Name: name, Int: value, Type: FieldTypeInt}
}

func Float64(name string, value float64) Field {
	return Field{Name: name, Float64: value, Type: FieldTypeFloat64}
}
