package main

import (
	"ccl/authn"
	"ccl/db/ent"
	"fmt"

	"entgo.io/ent/dialect"
	_ "github.com/go-sql-driver/mysql"
)

var (
	client *ent.Client
	err    error
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", "root", "root!@@&", "localhost", "3306", "cds_infra",
		"charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai")
	client, err = ent.Open(dialect.MySQL, dsn)
	err = authn.NewAuthorizationServer(client).Run()
	if err != nil {
		return
	}
}
