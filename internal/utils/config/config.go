package config

import "encoding/json"

func ConfigFromStruct(config interface{}, configStruct interface{}) {
	jsonM, _ := json.Marshal(config)
	json.Unmarshal(jsonM, configStruct)
}
