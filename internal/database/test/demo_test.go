package test

import (
	"ccl/db/ent"
	"ccl/db/ent/user"
	"context"
	"fmt"
	"log"
	"testing"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/go-sql-driver/mysql"
)

func Test_Example_Mysql(t *testing.T) {
	ctx := context.Background()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", "root", "root!@@&", "localhost", "3306", "cds_infra",
		"charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai")
	// 1. 创建一个日志记录器，将日志输出到标准输出
	// log.New(os.Stdout, "[ent] ", log.LstdFlags) 会在每条日志前加上 "[ent] " 前缀
	//logger := log.New(os.Stdout, "[ent] ", log.LstdFlags)
	client, err := ent.Open(dialect.MySQL, dsn, ent.Debug())
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(ctx, schema.WithGlobalUniqueID(true)); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	u, err := client.User.Update().Where(user.UsernameEQ("coffee377")).
		SetEmail("coffee377@dingtalk.com").
		Save(ctx)
	if err != nil {
		log.Fatalf("failed creating user: %v", err)
	}
	fmt.Println(u)
}
