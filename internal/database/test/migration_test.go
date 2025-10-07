package test

import (
	"ccl/db/ent"
	"ccl/db/ent/hook"
	_ "ccl/db/ent/runtime"
	"context"
	"fmt"
	"log"
	"testing"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/go-sql-driver/mysql"
)

var (
	client *ent.Client
	tx     *ent.Tx
	err    error
	ctx    = context.Background()
)

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", "root", "root!@@&", "localhost", "3306", "cds_infra",
		"charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai")
	client, err = ent.Open(dialect.MySQL, dsn)
	// 全局 hook
	//client.Use(hooks.PasswordHook)
	client.User.Use(hook.PasswordHook)
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	//defer client.Close()
}

func Test_Migration(t *testing.T) {
	// Run the auto migration tool.
	if err := client.Debug().Schema.Create(ctx,
		//schema.WithGlobalUniqueID(true),
		schema.WithDropIndex(true),
		schema.WithDropColumn(true),
	); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}

func WithTx(fn func(tx *ent.Tx) error) {
	tx, err = client.Tx(ctx)
	if err != nil {
		log.Fatalf("failed starting a transaction: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
		}
	}()

	if err = fn(tx); err != nil {
		if err = tx.Rollback(); err != nil {
			log.Fatalf("rolling back transaction: %v", err)
		}
		return
	}
	if err = tx.Commit(); err != nil {
		log.Fatalf("committing transaction: %v", err)
	}
}
