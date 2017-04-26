package Components

import (
	"os"
	"os/exec"
	"processengine/context"
	"processengine/logger"
	"runtime"
)

// method used to remove docker related files from the backend
func RemoveDocker(wfdetails DockerDeployement) (flowResult context.FlowResult) {

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
	flowResult.Message = "Starting docker removal"
	// starting the process
	logger.Log_PE("~~ Starting docker removal.", logger.Information, sessionId)

	dockerfullname := wfdetails.WFName + ":latest"

	logger.Log_PE("Docker name: "+dockerfullname, logger.Debug, sessionId)

	out, err := exec.Command(pwd+slash+"removeDocker", dockerfullname, "sh").Output()
	// check if there is an error or not
	if err != nil {
		msg := "Error on removing docker"
		flowResult.Status = false
		flowResult.Message = msg
		logger.Log_PE(msg, logger.Error, sessionId)
		logger.Log_PE(string(out), logger.Error, sessionId)
		logger.Log_PE(err.Error(), logger.Error, sessionId)
	} else {
		msg := "Docker removal succesfull."
		flowResult.Status = true
		flowResult.Message = msg
		logger.Log_PE(msg, logger.Debug, sessionId)
		logger.Log_PE(string(out), logger.Debug, sessionId)
	}
	// when the process is complete it passes the response back to the front
	return flowResult
}
