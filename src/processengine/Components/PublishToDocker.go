package Components

import (
	"os"
	"os/exec"
	"processengine/context"
	"processengine/logger"
	"runtime"
)

// method used to test the workflow with dummy data
func PublishToDocker(wfdetails DockerDeployement) (flowResult context.FlowResult) {

	sessionId := wfdetails.SessionID

	var slash = ""
	if runtime.GOOS != "linux" {
		// currently docker is supported only for linux, since then this should allow installing docker only for linux servers
		flowResult.Status = false
		flowResult.Message = "The server is not compatible for Docker depolyement."
		flowResult.SessionID = sessionId
		flowResult.FlowName = wfdetails.WFName
		return flowResult
	} else {
		slash = "/"
	}

	pwd, _ := os.Getwd()

	// initiating default variables
	flowResult.Status = true
	flowResult.Message = "Starting docker publisher"
	// starting the process
	logger.Log_PE("~~ Starting docker publisher.", logger.Information, sessionId)

	// pending tasks to do
	/*
		- get the Executable file location
		- get the docker file location (should be inside the newly generated docker file)
	*/
	FromFileLocation := pwd + slash + "BuiltExecutables" + slash + wfdetails.WFName + slash
	ToFileLocation := "/home/smoothflow/executables/"

	logger.Log_PE("From: "+FromFileLocation, logger.Debug, sessionId)
	logger.Log_PE("To: "+ToFileLocation, logger.Debug, sessionId)
	logger.Log_PE("WFName: "+wfdetails.WFName, logger.Debug, sessionId)
	logger.Log_PE("Port: "+wfdetails.Port, logger.Debug, sessionId)
	logger.Log_PE("RAM: "+wfdetails.RAM, logger.Debug, sessionId)
	logger.Log_PE("CPU: "+wfdetails.CPU, logger.Debug, sessionId)

	//RAM := "--memory=\"" + wfdetails.RAM + "\""
	//CPU := "--cpus=" + wfdetails.CPU

	RAM := wfdetails.RAM
	CPU := wfdetails.CPU

	out, err := exec.Command(pwd+slash+"publishToDocker.sh", FromFileLocation, ToFileLocation, wfdetails.WFName, wfdetails.Port, RAM, CPU, "sh").Output()
	// check if there is an error or not
	if err != nil {
		msg := "Error on publishing docker"
		flowResult.Status = false
		flowResult.Message = msg
		logger.Log_PE(msg, logger.Debug, sessionId)
		logger.Log_PE(string(out), logger.Error, sessionId)
		logger.Log_PE(err.Error(), logger.Error, sessionId)
	} else {
		msg := "Docker publishment succesfull."
		flowResult.Status = true
		flowResult.Message = msg
		logger.Log_PE(string(out), logger.Debug, sessionId)
		logger.Log_PE(msg, logger.Information, sessionId)
	}
	// when the process is complete it passes the response back to the front
	return flowResult
}
