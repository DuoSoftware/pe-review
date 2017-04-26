package Common

import (
	"duov6.com/common"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"io/ioutil"
	"net"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
)

var CEB_URL string
var OBJECTSTORE_URL string
var AUTH_URL string
var LOGSTASH_URL string

func VerifyCEBAgentConfig() (config map[string]interface{}) {
	config = make(map[string]interface{})
	content, err := ioutil.ReadFile("agent.config")
	if err != nil {
		//Agent File not Available
		config["cebUrl"] = "ceb.smoothflow.io:5000"
		config["authUrl"] = "auth.smoothflow.io:3048"
		config["logstashUrl"] = "logstash.smoothflow.io:5044"
		config["objUrl"] = "obj.smoothflow.io:3000"
		config["canMonitorOutput"] = true
		config["ResourceClass"] = (runtime.GOOS + "_" + runtime.GOARCH)
		byteArray, _ := json.Marshal(config)
		_ = ioutil.WriteFile("agent.config", byteArray, 0666)
	} else {
		_ = json.Unmarshal(content, &config)
		if config["ResourceClass"] == nil {
			config["ResourceClass"] = (runtime.GOOS + "_" + runtime.GOARCH)
		}
	}

	CEB_URL = config["cebUrl"].(string)
	AUTH_URL = config["authUrl"].(string)
	LOGSTASH_URL = config["logstashUrl"].(string)
	OBJECTSTORE_URL = config["objUrl"].(string)

	return
}

func CheckConnectionToCEB() (err error) {
	config := make(map[string]interface{})
	content, _ := ioutil.ReadFile("agent.config")
	_ = json.Unmarshal(content, &config)
	host := config["cebUrl"].(string)
	_, err = net.Dial("tcp", host)
	return
}

func AuthenticateSecurityToken(securityToken string) (status bool, session map[string]interface{}) {
	session = make(map[string]interface{})

	url := "http://"
	config := make(map[string]interface{})
	content, err := ioutil.ReadFile("agent.config")
	if err != nil {
		status = false
		return
	} else {
		_ = json.Unmarshal(content, &config)
		if config["authUrl"] == nil {
			//set default
			config["authUrl"] = "auth.smoothflow.io:3048"
		}
		url += config["authUrl"].(string) + "/GetSession/" + securityToken + "/Nil"
	}

	var bodyBytes []byte
	err, bodyBytes = common.HTTP_GET(url, nil, false)
	if err != nil {
		status = false
	} else {
		status = true
		_ = json.Unmarshal(bodyBytes, &session)
		fmt.Println(session)
	}
	return
}

func VerifyDependencies() (status bool) {
	status = true
	var err error

	packageNames := [...]string{"python", "autoconf", "python-gpiozero", "python-pkg-resources", "python-picamera"}

	IsAllPackagesInstalled := true

	for x := 0; x < len(packageNames); x++ {
		byteArray, err1 := exec.Command("sh", "-c", ("dpkg -l | grep " + packageNames[x])).Output()

		if err1 != nil {
			fmt.Println(err1.Error())
			status = false
			break
		} else {
			if strings.EqualFold(strings.TrimSpace(string(byteArray)), "") {
				status = false
				break
			}
		}
	}

	if !IsAllPackagesInstalled {
		_, err = exec.Command("sh", "-c", "sudo apt-get update").Output()

		for x := 0; x < len(packageNames); x++ {
			_, err = exec.Command("sh", "-c", ("sudo apt-get install -y " + packageNames[x])).Output()
			if err != nil {
				fmt.Println(err.Error())
				status = false
				break
			}
		}
	}

	return
}

func VerifyGPIOCapability() (status bool) {
	status = false
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	if strings.Contains(goos, "linux") && strings.Contains(goarch, "arm") {
		status = true
	} else {
		status = false
	}
	return
}

func GetDomainBySessionID(sessionID string) (domain string) {
	tokens := strings.Split(common.DecodeFromBase64(sessionID), "-")
	domain = tokens[0]
	return
}

func ConvertInterfaceArrayToObjectArray(objs interface{}) []map[string]interface{} {
	s := reflect.ValueOf(objs)
	var interfaceList []map[string]interface{}
	interfaceList = make([]map[string]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		obj := s.Index(i).Interface()
		v := reflect.ValueOf(obj)
		k := v.Kind()
		var newMap map[string]interface{}

		if k != reflect.Map {
			newMap = structs.Map(obj)
		} else {
			newMap = obj.(map[string]interface{})
		}

		interfaceList[i] = newMap
	}
	return interfaceList
}
