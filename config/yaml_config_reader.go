package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

/*
 * Read Yaml file to `target`
 * depends gopkg.in/yaml.v2
 */
func ReadYamlConfig(path string, target interface{}) (err error){

	contentsBytes, err := ioutil.ReadFile(path)

	if err != nil {
		log.Printf("to read yaml file (%s) but got an err => %v\n", path, err)
		return
	}

	err = yaml.Unmarshal(contentsBytes, target)

	if err != nil {
		log.Printf("unmarshal yaml but got an err => %v\n", err)
	}
	return
}
