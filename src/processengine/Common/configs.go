package Common

import (
	"encoding/json"
	"io/ioutil"
)

var settingsConfig map[string]string

/*get value of a key from the config*/
func GetConfigValue(key string) (output string) {

	if settingsConfig == nil {
		settingsConfig = make(map[string]string)

		content, err := ioutil.ReadFile("settings.config")
		if err == nil {
			_ = json.Unmarshal(content, &settingsConfig)
		}

		output = settingsConfig[key]
	} else {
		output = settingsConfig[key]
	}

	return
}

func GetConfig() (object map[string]string) {
	if settingsConfig == nil {
		settingsConfig = make(map[string]string)

		content, err := ioutil.ReadFile("settings.config")
		if err == nil {
			_ = json.Unmarshal(content, &settingsConfig)
		}
		object = settingsConfig
	} else {
		object = settingsConfig
	}

	return
}

func SaveConfig(object map[string]string) {
	if settingsConfig == nil {
		settingsConfig = make(map[string]string)
	}
	settingsConfig = object
	byteArray, _ := json.Marshal(object)
	_ = ioutil.WriteFile("settings.config", byteArray, 0666)
}
