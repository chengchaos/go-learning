package helper

import (
	"database/sql"
	"github.com/chengchaos/go-learning/config"
	"log"
	"testing"
)

type Config struct {
	Admin string
	Database Database
}

func Test_OpenMySQLDB(t *testing.T) {


	yamlConfig := &Config{}

	err := config.ReadYamlConfig("../conf/config.yml", yamlConfig)

	if err != nil {
		log.Panic(err)
	}

	helper := NewDBHelper()
	helper.InitProperties(&yamlConfig.Database.Mysql)

	size := 10
	channels := make([]chan int, size)

	for i := 0; i < size; i++ {
		ch := make(chan int)
		channels[i] = ch
		go func() {
			db := helper.GetDB()
			t.Log("db =>", db)
			ch<- 1
		}()
	}

	for i, ch := range channels {
		x := <-ch
		t.Logf("%d => %d\n", i, x)
	}

	t.Log("end")

}


func CreateHelper() *DBHelper{

	yamlConfig := &Config{}

	err := config.ReadYamlConfig("../conf/config.yml", yamlConfig)

	if err != nil {
		log.Panic(err)
	}

	mysql := &yamlConfig.Database.Mysql

	helper := NewDBHelper()
	helper.InitProperties(mysql)

	return helper
}


func Test_InsertUserinfo(t *testing.T) {

	helper := CreateHelper()

	prop := InsertProps{
		InsertSql: `insert into userinfo (username, departname, created) values (? , ?, ?)`,
		Callback: func(stmt *sql.Stmt) (sql.Result, error) {
			return stmt.Exec("chengchao", "研发部", "2020-11-01 23:33:33")
		},
	}

	result, err := helper.Insert(prop)

	if err != nil {
		return
	}

	id, err := result.LastInsertId()

	if err != nil {
		log.Panic(err)
	}

	log.Println("insert id =>", id)
}

func Test_SelectOne(t *testing.T) {
	helper := CreateHelper()

	helper.SelectOne()
}