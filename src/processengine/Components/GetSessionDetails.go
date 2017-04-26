package Components

import (
	// "encoding/json"
	"io/ioutil"
	"processengine/logger"
	"runtime"
	// "strings"
	// "strconv"
)

func GetSessionDetails(sessionId, sessiontype string) (sessionDetailObj SessionTranDetails) {

	logger.Log_PE("~~ Getting details for Session ID: "+sessionId+" Type: "+sessiontype, logger.Information, "Y29tLnNtb290aGZsb3cuaW8tMDAwMA")

	var slash = ""
	if runtime.GOOS == "windows" {
		slash = "\\"
	} else {
		slash = "/"
	}

	// get log file according to its type
	var logfilename = ""
	var isValid bool = false
	logfilename = sessiontype + "_" + sessionId + ".log"
	filepath := "Logs" + slash + logfilename
	displayMessage := ""
	logContent := ""
	logger.Log_PE("File path: "+filepath, logger.Debug, "Y29tLnNtb290aGZsb3cuaW8tMDAwMA")

	if sessiontype == "PE" {
		isValid = true
		logContent = GetLogFile(filepath)
		displayMessage = "Process Engine trace details were retrieved successfully for the given Session ID."
	} else if sessiontype == "WF" {
		isValid = true
		logContent = GetLogFile(filepath)
		displayMessage = "Workflow trace details were retrieved successfully for the given Session ID."
	} else if sessiontype == "ACT" {
		isValid = true
		logContent = GetLogFile(filepath)
		displayMessage = "Activity trace details were retrieved successfully for the given Session ID."
	} else if sessiontype == "ALL" {
		isValid = true

		PElog := GetLogFile("Logs" + slash + "PE_" + sessionId + ".log")
		WFlog := GetLogFile("Logs" + slash + "WF_" + sessionId + ".log")
		ACTlog := GetLogFile("Logs" + slash + "ACT_" + sessionId + ".log")

		logContent = PElog + "\n\n" + WFlog + "\n\n" + ACTlog

		displayMessage = "All the trace details were retrived successfully for the given SessionID."
	} else {
		isValid = false
		displayMessage = "The SessionType you provided is not as any of these formats PE,WF,ACT or ALL"
	}

	if isValid == true {
		if logContent == "" {
			displayMessage = "The file you requested is not available in the Logs directory."
		}
		sessionDetailObj.SessionID = sessionId
		sessionDetailObj.SessionType = sessiontype
		sessionDetailObj.SessionDetails = logContent
		sessionDetailObj.Message = displayMessage
	} else {
		sessionDetailObj.SessionID = sessionId
		sessionDetailObj.SessionType = sessiontype
		sessionDetailObj.SessionDetails = ""
		sessionDetailObj.Message = displayMessage
	}

	logger.Log_PE("~~ Session detail request is completed", logger.Information, "Y29tLnNtb290aGZsb3cuaW8tMDAwMA")
	return
}

func GetLogFile(filepath string) string {
	var returnTrace []byte
	activityTrace, init_error := ioutil.ReadFile(filepath)
	if init_error == nil {
		returnTrace = activityTrace
	} else {
		returnTrace = []byte("")
	}
	return string(returnTrace)
}
