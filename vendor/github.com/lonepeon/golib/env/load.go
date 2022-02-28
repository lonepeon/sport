package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type TagValue struct {
	variableName string
	options      map[string]string
}

func (t TagValue) VariableValue() string {
	return strings.TrimSpace(os.Getenv(t.variableName))
}

func (t TagValue) OptionWithDefaultValue(name string, defaultValue string) string {
	if value, ok := t.options[name]; ok {
		return value
	}

	return defaultValue
}

func newTagValue(rawTag string) (TagValue, error) {
	nameAndOptions := strings.Split(strings.TrimSpace(rawTag), ",")

	name := nameAndOptions[0]
	rawOptions := nameAndOptions[1:]
	options := make(map[string]string)
	for i := range rawOptions {
		kindAndValue := strings.Split(strings.TrimSpace(rawOptions[i]), "=")
		if len(kindAndValue) != 2 {
			return TagValue{}, fmt.Errorf("can't parse env tag option \"%s\"", rawOptions[i])
		}

		options[kindAndValue[0]] = kindAndValue[1]
	}

	return TagValue{variableName: name, options: options}, nil
}

func Load(input interface{}) error {
	inputType := reflect.TypeOf(input)
	if inputType.Kind() != reflect.Ptr {
		return fmt.Errorf("expected input to be a pointer")
	}

	inputPtrType := inputType.Elem()
	if inputPtrType.Kind() != reflect.Struct {
		return fmt.Errorf("expected input to be a pointer to a struct")
	}
	inputValue := reflect.ValueOf(input).Elem()

	for i := 0; i < inputPtrType.NumField(); i++ {
		if err := parseField(inputPtrType, inputValue, i); err != nil {
			return err
		}
	}

	return nil
}

func parseField(inputPtrType reflect.Type, inputValue reflect.Value, i int) error {
	inputFieldType := inputPtrType.Field(i)
	inputFieldValue := inputValue.Field(i)

	if !inputFieldValue.CanSet() {
		return nil
	}

	tag, err := newTagValue(inputFieldType.Tag.Get("env"))
	if err != nil {
		return err
	}

	value := tag.VariableValue()
	if value == "" {
		value = tag.OptionWithDefaultValue("default", "")
	}

	if value == "" {
		if tag.OptionWithDefaultValue("required", "false") == "true" {
			return fmt.Errorf("environment variable %s is required", tag.variableName)
		}
		return nil
	}

	if err := setFieldValue(inputFieldType, inputFieldValue, value, tag); err != nil {
		return fmt.Errorf("can't set variable %s to %s field: %v", tag.variableName, inputFieldType.Type.Kind(), err)
	}

	return nil
}

func setFieldValue(inputFieldType reflect.StructField, inputFieldValue reflect.Value, value string, tag TagValue) error {
	switch inputFieldType.Type.Kind() {
	case reflect.String:
		inputFieldValue.SetString(value)
	case reflect.Int:
		x, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("can't parse to int")
		}

		inputFieldValue.SetInt(x)
	case reflect.Slice:
		separator := tag.OptionWithDefaultValue("sep", ",")
		switch inputFieldType.Type.Elem().Kind() {
		case reflect.String:
			inputFieldValue.Set(reflect.ValueOf(strings.Split(value, separator)))
		default:
			return fmt.Errorf("unsupported slice value type %s", inputFieldType.Type.Elem().Kind())
		}
	default:
		return fmt.Errorf("unsupported value type %s", inputFieldType.Type.Kind())
	}

	return nil
}
