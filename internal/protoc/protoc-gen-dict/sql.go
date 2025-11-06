package main

import (
	"fmt"
	"path"
	"reflect"
	"strings"

	"github.com/code-crafters-lab/ccl/internal/protoc/protoc-gen-dict/sql"
	"google.golang.org/protobuf/compiler/protogen"
)

func sqlGenerator(dictionaries []*dict, plugin *protogen.Plugin, file *protogen.File, opts option) {
	if !opts.sql {
		return
	}
	filename := fmt.Sprintf("%s.dict.sql", file.GeneratedFilenamePrefix)
	// 如果存在 java 生成选项，将 sql 提到一级目录
	if opts.java {
		filename = fmt.Sprintf("%s.dict.sql", path.Base(file.GeneratedFilenamePrefix))
	}
	g := plugin.NewGeneratedFile(filename, file.GoImportPath)
	generatePreamble(g, "-- ", file, plugin)
	for i, dictionary := range dictionaries {
		g.P()
		g.P(fmt.Sprintf("-- %d. %s", i+1, dictionary.Name))
		// 字典插入
		d := sql.InsertBuilder("t_dict", "").PrimaryKey("id")
		d.Columns("id", "code", "label", "value", "sort", "remark", "created_time", "modified_time")
		d.CreatedAtColumn("created_time")
		d.UpdatedAtColumn("modified_time")
		_ = d.Data([]map[string]interface{}{
			dictionary.toMap(),
		})
		_, _ = d.Build()
		g.P(d.String())

		// 字典项插入
		di := sql.InsertBuilder("t_dict", "").PrimaryKey("id")
		di.Columns("id", "pid", "code", "label", "value", "sort", "remark", "created_time", "modified_time")
		di.CreatedAtColumn("created_time")
		di.UpdatedAtColumn("modified_time")
		_ = di.Data([]map[string]interface{}{
			dictionary.toMap(),
		})
		_, _ = di.Build()
		g.P(di.String())
	}
}

// 生成 sql
func generateSQL(dictionaries []*dict, plugin *protogen.Plugin, file *protogen.File, opts option) {
	if !opts.sql {
		return
	}
	filename := fmt.Sprintf("%s.dict.help.sql", file.GeneratedFilenamePrefix)
	//if opts.java {
	//	filename = fmt.Sprintf("%s.dict.sql", path.Base(file.GeneratedFilenamePrefix))
	//}
	g := plugin.NewGeneratedFile(filename, file.GoImportPath)
	generatePreamble(g, "-- ", file, plugin)

	for i, dictionary := range dictionaries {
		g.P()

		// 备注
		var (
			zeroItem   *dictItem
			result     []string
			value2name []string
			values     []string
		)
		for _, item := range dictionary.Items {
			if item.IsZero() {
				zeroItem = item
				continue
			}
			result = append(result, fmt.Sprintf("%s:%s", item.Value, item.Name))
			if dictionary.ValueType == ValueTypeNumber {
				value2name = append(value2name, fmt.Sprintf(`when %s then '%s'`, item.Value, item.Name))
			} else {
				value2name = append(value2name, fmt.Sprintf(`when '%s' then '%s'`, item.Value, item.Name))
			}
			values = append(values, fmt.Sprintf(`"%s"`, item.Value))
		}
		if zeroItem != nil {
			result = append(result, fmt.Sprintf("%s:%s", zeroItem.Value, zeroItem.Name))
			if dictionary.ValueType == ValueTypeNumber {
				value2name = append(value2name, fmt.Sprintf(`else %s`, zeroItem.Value))
			} else {
				value2name = append(value2name, fmt.Sprintf(`else '%s'`, zeroItem.Name))
			}
			values = append(values, fmt.Sprintf(`"%s"`, zeroItem.Value))
		}
		g.P(fmt.Sprintf("-- %d. %s %s", i+1, dictionary.Name, strings.Join(result, " ")))
		g.P(fmt.Sprintf("-- enum values: %s", strings.Join(values, ",")))
		g.P("select")
		g.P(fmt.Sprintf(`case '%s'`, dictionary.Code))
		g.P(strings.Join(value2name, "\n"))
		g.P("end")
		g.P("from dual;")

	}
}

func resolveDictValues(dict dict) string {
	vals := []any{dict.ID, dict.FullCode, dict.Name, dict.ValueType, dict.Description}
	values := make([]string, len(vals))
	for i, v := range vals {
		values[i] = resolveValue(reflect.ValueOf(v), nil)
	}
	return strings.Join(values, ", ")
}

func resolveItemValues(dictId uint64, item dictItem) string {
	vals := []any{item.ID, dictId, item.Code, item.Name, item.Value, item.Sort, item.Description}
	values := make([]string, len(vals))
	for i, v := range vals {
		values[i] = resolveValue(reflect.ValueOf(v), nil)
	}
	return strings.Join(values, ", ")
}

func resolveValue(value reflect.Value, defaultVale interface{}) string {
	var result string
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
