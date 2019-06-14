package godaddy

import (
	"encoding/json"
	"io/ioutil"
)

func InitConfig() *Context {
	var config struct {
		Key    string `json:"key"`
		Secret string `json:"secret"`
	}
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(data, &config); err != nil {
		panic(err)
	}
	return NewGoDaddy(config.Key, config.Secret)
}
