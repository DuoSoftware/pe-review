package Common

import (
	"github.com/fatih/color"
	"io/ioutil"
)

func Log(message, sessionid string) {
	logname := "PE_" + sessionid + ".log"
	PublishLog(logname, message)
	PublishLog("ProcessEngineLog.log", message)
}

func LogHighlighted(message, sessionid string) {
	font := color.New(color.FgBlack)
	colorscheme := font.Add(color.BgWhite)
	colorscheme.Println(message)
	logname := "PE_" + sessionid + ".log"
	PublishLog(logname, message)
	PublishLog("ProcessEngineLog.log", message)
}

func LogWFMSG(message, sessionid string) {
	logname := "WF_" + sessionid + ".log"
	PublishLog(logname, message)
	PublishLog("WorkflowLog.log", message)
}

func LogACT(message, sessionid string) {
	logname := "ACT_" + sessionid + ".log"
	PublishLog(logname, message)
	PublishLog("ActivityLog.log", message)
}

func GetActLog(sessionID string) string {
	var returnTrace []byte
	activityTrace, init_error := ioutil.ReadFile("ACT_" + sessionID)
	if init_error == nil {
		returnTrace = activityTrace
		removeFile("ACT_" + sessionID)
	} else {
		returnTrace = []byte("")
	}
	return string(returnTrace)
}

func GetWFLog(sessionID string) string {
	var returnTrace []byte
	activityTrace, init_error := ioutil.ReadFile("WF_" + sessionID)
	if init_error == nil {
		returnTrace = activityTrace
		removeFile("WF_" + sessionID)
	} else {
		returnTrace = []byte("")
	}

	return string(returnTrace)
}

func GetPELog(sessionID string) string {
	var returnTrace []byte
	activityTrace, init_error := ioutil.ReadFile("PE_" + sessionID)
	if init_error == nil {
		returnTrace = activityTrace
		removeFile("PE_" + sessionID)
	} else {
		returnTrace = []byte("")
	}
	return string(returnTrace)
}
