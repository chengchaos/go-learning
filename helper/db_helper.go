package helper

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type DBHelper struct {
	once      sync.Once
	mysqlDb   *sql.DB
	mysqlProp *Mysql
}

func NewDBHelper() *DBHelper {
	return &DBHelper{}
}

func CloseCloser(closer io.Closer) {

	err := closer.Close()
	if err != nil {
		log.Println("close Closer but got an error =>", err)
	}
}



func (helper *DBHelper) openMySQLDB0() {

	log.Println("Go Go Go ..")
	prop := helper.mysqlProp
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s%s)/%s?charset=utf8",
		prop.Username, prop.Password, prop.Host, prop.Port, prop.Schema)

	log.Println("dataSourceName =>", dataSourceName)

	db, err := sql.Open("mysql", dataSourceName)

	if err != nil {
		log.Panicf("Open MySQL but got an error => %v\n", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(300)

	helper.mysqlDb = db

}

func (helper *DBHelper) InitProperties(prop *Mysql) *DBHelper {
	helper.mysqlProp = prop
	return helper
}

func (helper *DBHelper) GetDB() *sql.DB {
	helper.once.Do(helper.openMySQLDB0)
	return helper.mysqlDb
}

func (helper *DBHelper) Insert(props InsertProps) (sql.Result, error) {
	db := helper.GetDB()
	stmt, err := db.Prepare(props.InsertSql)
	if err != nil {
		return nil, err
	}
	defer CloseCloser(stmt)
	return props.Callback(stmt)
}


func (helper *DBHelper) SelectOne() {

	selectSql := "SELECT * FROM userinfo WHERE uid > ? ORDER BY uid DESC LIMIT 2  OFFSET 0"
	db := helper.GetDB()

	stmt, err := db.Prepare(selectSql)

	if err != nil {
		log.Panic(err)
	}
	defer CloseCloser(stmt)

	var id int64 = 0
	rows, err :=stmt.Query(id)

	if err != nil {
		log.Panic(err)
	}
	defer CloseCloser(rows)

	columns, err := rows.Columns()

	if err != nil {
		log.Panic(err)
	}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		log.Panic(err)
	}

	log.Println("columns =>", columns)
	for _, columnType := range columnTypes{
		/*
			name string

			hasNullable       bool
			hasLength         bool
			hasPrecisionScale bool

			nullable     bool
			length       int64
			databaseType string
			precision    int64
			scale        int64
			scanType     reflect.Type
		 */
		log.Println("Name =>", columnType.Name())
		log.Println("DatabaseTypeName =>", columnType.DatabaseTypeName())

		log.Println("ScanType =>", columnType.ScanType())

		if precision, scale, ok  := columnType.DecimalSize(); ok {
			log.Println("DecimalSize =>",precision, scale )
		}

		if length, ok := columnType.Length(); ok {
			log.Println("Length =>", length)
		}

		if nullable, ok := columnType.Nullable(); ok {
			log.Println("Nullable =>", nullable)
		}

		log.Println("=======================")
	}

	recordLen := len(columns)


	for rows.Next() {
		record := make([][]byte, recordLen)
		recordRef := make([]interface{}, recordLen)
		for i, _ := range record {
			recordRef[i] = &record[i]
		}

		rows.Scan(recordRef...)

		log.Println("record =>", record)
		log.Println("recordRef =>", recordRef)

		for i, name := range columns {
			v := string(record[i])
			log.Printf( "%s => %s",name, v)
		}

	}



}
