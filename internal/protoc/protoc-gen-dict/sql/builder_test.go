package sql

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInsertBuilder(t *testing.T) {
	builder := InsertBuilder("test", "")
	builder.Columns("type", "code", "name", "created_at").
		UniqueKey("type", "code")

	builder.Data([]map[string]interface{}{
		{"type": "11", "code": "111"},
		{"code": "222", "type": "22"},
	})
	_, err := builder.Build()
	assert.Nil(t, err)

	file, err := os.OpenFile("test.sql", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		t.Errorf("%v", err)
	}
	//_ = builder.Write(os.Stdout)
	_ = builder.Write(file)
}
