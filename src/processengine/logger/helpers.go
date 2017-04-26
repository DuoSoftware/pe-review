package logger

import (
	"duov6.com/cebadapter"
	"duov6.com/common"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func GetDomainBySessionID(sessionID string) (domain string) {
	tokens := strings.Split(common.DecodeFromBase64(sessionID), "-")
	domain = tokens[0]
	return
}

func VerifyIsSmoothFlowEngine() bool {
	if cebadapter.GetAgent() != nil {
		if strings.Contains(cebadapter.GetAgent().Client.GetAgentName(), "ProcessEngine@") {
			return true
		}
	}
	//either not connected to CEB and not ProcessEngine at all.
	return false
}

func GetMessageInString(data interface{}) string {
	Lable := ""
	if reflect.TypeOf(data).String() == "string" {
		Lable = data.(string)
	} else {
		byteArray, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err.Error())
			return err.Error()
		}
		Lable = string(byteArray)
	}
	return Lable
}

func GetLogNamesByCategory(logCategory int, sessionID string) []string {
	var logNames []string

	if logCategory == Default {
		logNames = make([]string, 1)
		logNames[0] = "ProcessEngineLog.log"
	} else if logCategory == ProcessEngine {
		logNames = make([]string, 2)
		logNames[0] = "ProcessEngineLog.log"
		logNames[1] = "PE_" + sessionID + ".log"
	} else if logCategory == Activity {
		logNames = make([]string, 2)
		logNames[0] = "ActivityLog.log"
		logNames[1] = "ACT_" + sessionID + ".log"
	} else if logCategory == WorkFlow {
		logNames = make([]string, 2)
		logNames[0] = "WorkflowLog.log"
		logNames[1] = "WF_" + sessionID + ".log"
	}
	return logNames
}

func GetLogTypeByID(logType int) (typeInstring string) {
	if logType == Information {
		typeInstring = "Information"
	} else if logType == Error {
		typeInstring = "Error"
	} else if logType == Debug {
		typeInstring = "Debug"
	} else if logType == Splash {
		typeInstring = "Splash"
	} else if logType == Blank {
		typeInstring = "Blank"
	} else if logType == Warning {
		typeInstring = "Warning"
	}
	return
}

func GetCatTypeByID(catType int) (typeInstring string) {
	if catType == Default {
		typeInstring = "Default Log"
	} else if catType == ProcessEngine {
		typeInstring = "ProcessEngine Log"
	} else if catType == WorkFlow {
		typeInstring = "WorkFlow Log"
	} else if catType == Activity {
		typeInstring = "Activity Log"
	}
	return
}
