package config

import (
	"github.com/chengchaos/go-learning/helper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"testing"
)


type Config struct {
	Admin string
	Database helper.Database
}

func Test_ReadConfig(t *testing.T) {

	path := "../conf/config.yml"
	fileAsBytes, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatalf("read yaml file but ...%s\n", err)
	}

	config := &Config{}
	err = yaml.Unmarshal(fileAsBytes, config)

	if err != nil {
		log.Fatalf("unmarshal err : %s\n", err)
	}

	log.Println("config =>", config)
	log.Println("config admin =>", config.Admin)
	db := config.Database
	log.Println("config db => ", db)

	mysql := db.Mysql
	log.Println("mysql =>", mysql)

}

func Test_ReadConfig2(t *testing.T) {

	path := "../conf/config.yml"
	themap := make(map[string]interface{})
	err := ReadYamlConfig(path, &themap)
	if err != nil {
		log.Fatalf("read yaml config but get some error => %v\n", err)
	}

	log.Println("themap =>", themap)
	log.Println("admin =>", themap["admin"])
	db0 := themap["database"]
	db := db0.(map[interface{}]interface{})
	mysql0 := db["mysql"]
	mysql := mysql0.(map[interface{}]interface{})

	log.Println("host =>", mysql["host"])
	log.Println("post =>", mysql["port"])
	log.Println("username =>", mysql["username"])
	log.Println("password =>", mysql["password"])
}