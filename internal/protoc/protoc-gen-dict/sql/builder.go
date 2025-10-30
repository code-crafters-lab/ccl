package sql

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
	"time"
)

var (
	errTableName         = errors.New("table name is null")
	errInsertNullColumns = errors.New("insert columns is null")
	errInsertNullData    = errors.New("insert statements must have at least one set of values or select clause")
)

type Builder[T any] interface {
	Build() (T, error)
	String() string
	Bytes() []byte
	Write(writer io.Writer) error
}

type Insert struct {
	buf        *bytes.Buffer
	table      string
	columns    []Column
	primaryKey *string
	uniqueKeys [][]string
	data       []map[string]interface{}

	createdAtColumn string     // 创建时间字段
	createTime      *time.Time // 创建时间
	updatedAtColumn string     // 更新时间字段
	updateTime      *time.Time // 更新时间

	buildTime time.Time
	comment   string
}

func InsertBuilder(tableName, comment string) *Insert {
	return &Insert{
		buf:             &bytes.Buffer{},
		table:           tableName,
		columns:         []Column{},
		createdAtColumn: "created_at",
		createTime:      nil,
		updatedAtColumn: "updated_at",
		updateTime:      nil,
		buildTime:       time.Now(),
		comment:         comment,
	}
}

func (sql *Insert) Table(name string) *Insert {
	if name != "" {
		sql.table = name
	}
	return sql
}

func (sql *Insert) Columns(name ...string) *Insert {
	for _, n := range name {
		tName := strings.TrimSpace(n)
		if tName == "" {
			continue
		}
		sql.columns = append(sql.columns, NewColumn(tName))
	}
	return sql
}

func (sql *Insert) ColumnDefaultValue(values ...DefaultValueFn) *Insert {
	for _, vf := range values {
		for _, column := range sql.columns {
			if ok := vf(column); ok {
				break
			}
		}
	}
	return sql
}

func (sql *Insert) CreatedAtColumn(columnName string) *Insert {
	if columnName != "" {
		sql.createdAtColumn = columnName
	}
	return sql
}

func (sql *Insert) UpdatedAtColumn(columnName string) *Insert {
	if columnName != "" {
		sql.updatedAtColumn = columnName
	}
	return sql
}

func (sql *Insert) PrimaryKey(key string) *Insert {
	if key != "" {
		sql.primaryKey = &key
	}
	return sql
}

func (sql *Insert) UniqueKey(keys ...string) *Insert {
	var k []string
	for _, key := range keys {
		tKey := strings.TrimSpace(key)
		if tKey == "" {
			continue
		}
		k = append(k, tKey)
	}
	if len(k) > 0 {
		sql.uniqueKeys = append(sql.uniqueKeys, k)
	}
	return sql
}

func (sql *Insert) Data(data []map[string]interface{}) *Insert {
	sql.data = append(sql.data, data...)
	return sql
}

func (sql *Insert) Build() (*Insert, error) {
	err := sql.verify()
	if err != nil {
		return sql, err
	}

	sql.buf.Reset()

	// 数据插入列
	columns := sql.resolveColumns()
	// duplicate 更新列
	updateColumns := sql.resolveUpdateColumns(columns)

	//update := ""
	//if len(updateColumns) > 0 {
	//	update = "或更新"
	//}
	//comment := fmt.Sprintf("-- %s %s %s%s \n", sql.comment, sql.table, "数据新增", update)
	//sql.writeString(comment)
	sql.writeString(fmt.Sprintf("INSERT INTO %s (%s)\n", sql.table, strings.Join(columns, ", ")))

	// 循环添加数据
	for i, item := range sql.data {
		sql.addValue(item, i, len(sql.data))
	}

	// DUPLICATE KEY UPDATE
	sql.updateOnDuplicate(updateColumns)

	// 结束语句
	sql.writeString(";\n")

	return sql, nil
}

func (sql *Insert) String() string {
	return sql.buf.String()
}

func (sql *Insert) Bytes() []byte {
	return sql.buf.Bytes()
}

func (sql *Insert) Write(writer io.Writer) error {
	_, err := sql.buf.WriteTo(writer)
	if err != nil {
		return err
	}
	return nil
}

func (sql *Insert) verify() error {
	if sql.table == "" {
		return errTableName
	}
	if len(sql.columns) < 1 {
		return errInsertNullColumns
	}
	if len(sql.data) < 1 {
		return errInsertNullData
	}
	return nil
}

func (sql *Insert) writeString(s string) {
	sql.buf.WriteString(s)
}

func (sql *Insert) resolveColumns() []string {
	columns := make([]string, len(sql.columns))
	for i, c := range sql.columns {
		columns[i] = c.GetName()
	}
	return columns
}

func (sql *Insert) addValue(item map[string]interface{}, index, total int) {
	var (
		prefix  = ""
		sep     = ""
		newline = ""
	)
	if index == 0 {
		prefix = "VALUES"
	}
	if index < total-1 {
		sep = ","
		newline = "\n"
	}
	values := sql.resolveInsertValues(item)
	sql.writeString(fmt.Sprintf("%6s (%s)%s%s", prefix, values, sep, newline))
}

func (sql *Insert) resolveInsertValues(item map[string]interface{}) string {
	cvs := make(columnValues, len(sql.columns))
	for i, c := range sql.columns {
		cname := c.GetName()
		defaultValue := c.GetDefaultValue()
		if sql.createdAtColumn == cname || sql.updatedAtColumn == cname {
			if sql.createTime != nil {
				defaultValue = sql.createTime.Format(time.DateTime)
			} else {
				defaultValue = sql.buildTime.Format(time.DateTime)
			}
		}
		cvs[i] = &column{c.GetName(), defaultValue, item[cname]}
	}
	return cvs.resolveValues()
}

const (
	UpdateSql = "ON DUPLICATE KEY UPDATE"
)

func (sql *Insert) updateOnDuplicate(updateColumns []string) {
	updateLen := len(updateColumns)
	if updateLen > 0 {
		sql.writeString("\n")
		sql.writeString("AS alias\n")
	}
	maxLen := 0
	for _, updateColumn := range updateColumns {
		maxLen = int(math.Max(float64(maxLen), float64(len(updateColumn))))
	}

	format := fmt.Sprintf("%%%ds %%-%ds = %%s%%s%%s", len(UpdateSql), maxLen)
	for i, updateColumn := range updateColumns {
		prefix := ""
		sep := ""
		newline := ""
		if i == 0 {
			prefix = UpdateSql
		} else {
			prefix = ""
		}
		if i < updateLen-1 {
			sep = ","
			newline = "\n"
		}
		updateValue := fmt.Sprintf("alias.%s", updateColumn)
		if sql.updatedAtColumn == updateColumn {
			if sql.updateTime != nil {
				updateValue = fmt.Sprintf("'%s'", sql.updateTime.Format(time.DateTime))
			} else {
				updateValue = fmt.Sprintf("'%s'", sql.buildTime.Format(time.DateTime))
			}
		}
		sql.writeString(fmt.Sprintf(format, prefix, updateColumn, updateValue, sep, newline))
	}
}

func (sql *Insert) resolveUpdateColumns(columns []string) []string {
	var (
		result []string
		keyMap = make(map[string]bool)
	)
	// 主键
	if sql.primaryKey != nil && *sql.primaryKey != "" {
		keyMap[*sql.primaryKey] = true
	}
	// 唯一键
	for _, ks := range sql.uniqueKeys {
		for _, k := range ks {
			keyMap[k] = true
		}
	}
	// 不存在则唯一约束则直接返回
	if len(keyMap) == 0 {
		return nil
	}

	for _, c := range columns {
		// 更新语句忽略创建时间字段
		if sql.createdAtColumn == c {
			continue
		}
		if _, exists := keyMap[c]; !exists {
			result = append(result, c)
		}
	}
	return result
}
