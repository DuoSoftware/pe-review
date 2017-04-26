package Components

import (
	perf "duov6.com/common"
	"io/ioutil"
	"processengine/logger"
	"runtime"
	"strconv"
	"time"
)

//import "strings"
//import "strconv"

var ProcessEngineStartTime time.Time

func GetVersionDetails(sessionId string) (verionResult VersionResponse) {

	logger.Log_PE("~~ Getting Process engine version details", logger.Information, sessionId)
	//var verionResult = new(VersionResponse)

	var slash = ""
	if runtime.GOOS == "windows" {
		slash = "\\"
	} else {
		slash = "/"
	}

	//var versionInfo []VersionInfo

	PEversionFile, init_error := ioutil.ReadFile(".." + slash + "WFGOcommitHistory.txt")
	if init_error != nil {
		logger.Log_PE("Error occured on Unmarshal.", logger.Error, sessionId)
		logger.Log_PE("Error - "+init_error.Error(), logger.Error, sessionId)
	}

	SFversionFile, init_error := ioutil.ReadFile(".." + slash + "DPDcommitHistory.txt")
	if init_error != nil {
		logger.Log_PE("Error occured on Unmarshal.", logger.Error, sessionId)
		logger.Log_PE("Error - "+init_error.Error(), logger.Error, sessionId)
	}
	// run a bash and check if there are any updates on git repository
	verionResult.IsUptodate = true
	verionResult.PEVersionDetails = string(PEversionFile)
	verionResult.SFVersionDetails = string(SFversionFile)
	verionResult.Engine_Details = GetVersionDetailsJson()
	logger.Log_PE("~~ Version details retrived succesfully", logger.Information, sessionId)
	return
}

func GetVersionDetailsJson() map[string]interface{} {
	cpuUsage := strconv.Itoa(int(perf.GetProcessorUsage()))
	cpuCount := strconv.Itoa(runtime.NumCPU())

	versionData := make(map[string]interface{})
	versionData["API Name"] = "Smooth Flow Process Engine"
	versionData["API Version"] = "3.0.10f"

	changeLogs := make(map[string]interface{})

	changeLogs["3.0.10"] = "Fixed basic WF-BuiltFlow with hibernate, Added Jira Trigger and auth methods."
	changeLogs["3.0.9"] = "Added toggle Logstash method, Added Parameter for RAM : DSF-441"
	changeLogs["3.0.8"] = "converted BuildFlows msgs to logger."
	changeLogs["3.0.7"] = "Added logger to Components"
	changeLogs["3.0.6"] = "Added configurable domain based logger"
	changeLogs["3.0.5"] = "Added basic logstash server"
	changeLogs["3.0.4"] = "Improved Disk Log function"
	changeLogs["3.0.3"] = "Added delete activity from Disk. ( DSF-343 )"
	changeLogs["3.0.2"] = "Added SecurityToken for all methods, Added CORS support"
	changeLogs["3.0.1"] = "Added versioning document, Fixed Display Workflow details on the downloaded Executable"

	versionData["Change Logs"] = changeLogs

	gitMap := make(map[string]string)
	gitMap["Type"] = "git"
	gitMap["URL"] = "https://github.com/DuoSoftware/WF-GO/"
	versionData["Repository"] = gitMap

	statMap := make(map[string]string)
	statMap["CPU"] = cpuUsage + " (percentage)"
	statMap["CPU Cores"] = cpuCount
	nowTime := time.Now()
	elapsedTime := nowTime.Sub(ProcessEngineStartTime)
	statMap["Time Started"] = ProcessEngineStartTime.UTC().Add(330 * time.Minute).Format(time.RFC1123)
	statMap["Time Elapsed"] = elapsedTime.String()
	versionData["Metrics"] = statMap

	authorMap := make(map[string]string)
	authorMap["Name"] = "Duo Software Pvt Ltd"
	authorMap["URL"] = "http://www.duosoftware.com/"
	versionData["Project Author"] = authorMap

	return versionData
}
