package helper

import "database/sql"

type Mysql struct {
	Host string
	Port string
	Username string
	Password string
	Schema string
}
type Database struct {
	Mysql Mysql
}

type InsertCallback func(stmt *sql.Stmt) (sql.Result, error)

type InsertProps struct {
	InsertSql string
	Callback InsertCallback
}
