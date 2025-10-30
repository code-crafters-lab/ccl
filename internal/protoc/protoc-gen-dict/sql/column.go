package sql

import (
	"fmt"
	"reflect"
	"strings"
)

type Column interface {
	GetName() string                   // 字段名称
	GetDefaultValue() interface{}      // 获取默认值
	SetDefaultValue(value interface{}) // 设置默认值
}

type column struct {
	name         string
	defaultValue interface{}
	rawValue     interface{}
}

func NewColumn(name string, opts ...ColumnOption) Column {
	c := &column{name: name}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *column) GetName() string {
	return c.name
}

func (c *column) GetDefaultValue() interface{} {
	return c.defaultValue
}

func (c *column) SetDefaultValue(value interface{}) {
	c.defaultValue = &value
}

type ColumnOption func(*column)

func WithDefaultValue(defaultValue any) ColumnOption {
	return func(c *column) {
		c.defaultValue = &defaultValue
	}
}

func WithValue(value any) ColumnOption {
	return func(c *column) {
		c.rawValue = value
	}
}

type DefaultValueFn = func(Column) bool

func WithColumnDefault(name string, value any) DefaultValueFn {
	return func(column Column) bool {
		if column.GetName() == name {
			column.SetDefaultValue(value)
			return true
		}
		return false
	}
}

type columnValues []*column

func (cvs columnValues) resolve() []string {
	values := make([]string, len(cvs))
	for i, cv := range cvs {
		values[i] = resolveValue(reflect.ValueOf(cv.rawValue), cv.defaultValue)
	}
	return values
}

func (cvs columnValues) resolveValues() string {
	return strings.Join(cvs.resolve(), ", ")
}

func resolveValue(value reflect.Value, defaultVale interface{}) string {
	var (
		result string
	)
	switch value.Kind() {
	case reflect.Invalid:
		def := reflect.ValueOf(defaultVale)
		if def.Kind() == reflect.Invalid {
			result = "null"
		} else {
			result = resolveValue(def, nil)
		}
	case reflect.String:
		// todo 增加特殊字符串转义处理函数配置
		v := strings.ReplaceAll(strings.TrimSpace(value.String()), "'", "\\'")
		if v == "" {
			result = "null"
		} else {
			result = fmt.Sprintf("'%s'", v)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = fmt.Sprintf("%d", value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result = fmt.Sprintf("%d", value.Uint())
	case reflect.Float32, reflect.Float64:
		result = fmt.Sprintf("%f", value.Float())
	default:
		result = "null"
	}
	return fmt.Sprintf("%s", result)
}
