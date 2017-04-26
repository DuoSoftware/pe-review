package Components

import (
	"encoding/json"
	"os"
	"os/exec"
	"processengine/context"
	"processengine/logger"
	"runtime"
)

// method used to test the workflow with dummy data
func TestWorkflow(TestData TestWorkflowInvoke, flowName, sessionId string, isWorkflow bool) (wfResponse context.TestWorkflowResponse) {

	var slash = ""
	var executableName = ""
	if runtime.GOOS == "windows" {
		slash = "\\"
		executableName = flowName + ".exe"
	} else {
		slash = "/"
		executableName = flowName
	}

	pwd, _ := os.Getwd()

	// setting the default set of variables to he method
	wfResponse.Status = true
	wfResponse.ErrorDetails = ""
	if isWorkflow {
		wfResponse.Message = "Starting workflow test invoke"
		logger.Log_PE("~~ Starting Workflow test", logger.Information, sessionId)
	} else {
		wfResponse.Message = "Starting activity test invoke"
		logger.Log_PE("~~ Starting Activity test", logger.Information, sessionId)
	}
	// starting the process

	logger.Log_PE("Received data: ", logger.Debug, sessionId)
	logger.Log_PE(TestData.InArguments, logger.Debug, sessionId)
	logger.Log_PE("Executable Name: "+executableName, logger.Debug, sessionId)

	// check for the published workflow existence
	executableFound := IsExists(executableName)
	if executableFound == false {
		msg := ""
		if isWorkflow {
			msg = "The executable for given workflow is not found. Make sure that you have already published it before testing it."
		} else {
			msg = "The executable for given activity is not found. Make sure that you have already build it before testing it."
		}
		logger.Log_PE(msg, logger.Error, sessionId)
		wfResponse.Message = msg
		wfResponse.Status = false
		wfResponse.ErrorCode = 0
		return wfResponse
	} else {
		logger.Log_PE("Runtime Platform: Environment - "+string(runtime.GOOS), logger.Debug, sessionId)

		convertedStruct, _ := json.Marshal(TestData.InArguments)
		logger.Log_PE("Arguments : "+string(convertedStruct), logger.Debug, sessionId)
		//WFParameters := string(convertedStruct)

		var out []byte
		var err error
		if runtime.GOOS == "windows" {
			out, err = exec.Command(pwd+slash+executableName, string(TestData.InArguments), "CMD", "start").CombinedOutput()
		} else {
			out, err = exec.Command(pwd+slash+executableName, string(TestData.InArguments), "sh").CombinedOutput()
		}

		// check if there is an error or converting to Argument list or not
		if err != nil {
			msg := ""
			if isWorkflow {
				msg = "There was an error when testing the workflow."
			} else {
				msg = "There was an error when testing the activity."
			}
			logger.Log_PE(msg, logger.Error, sessionId)
			wfResponse.Message = msg
			wfResponse.ErrorDetails = string(out)
			wfResponse.Status = false
			wfResponse.ErrorCode = 0
			return wfResponse
		}

		logger.Log_PE("Starting response Unmarshal process.", logger.Debug, sessionId)

		res := context.ReturnData{}
		jerr := json.Unmarshal([]byte(string(out)), &res)
		// check if there is an error or converting to JSON or not
		if jerr != nil {
			msg := "Error occured on Unmarshalling result from the workflow."
			logger.Log_PE(msg, logger.Error, sessionId)
			logger.Log_PE("Result from workflow:", logger.Error, sessionId)
			logger.Log_PE(string(out), logger.Error, sessionId)
			wfResponse.Message = msg
			wfResponse.ErrorDetails = string(out)
			wfResponse.Status = false
			wfResponse.ErrorCode = 0
			return wfResponse
		} else {
			wfResponse.ResponseData = res
		}

		//print the valid data on screen
		wfResponse.Status = true
		if isWorkflow {
			wfResponse.Message = "Workflow has completed executing with an output."
		} else {
			wfResponse.Message = "Activity has completed executing with an output."
		}
		logger.Log_PE(executableName+" has executed with an output.", logger.Debug, sessionId)
		logger.Log_PE(string(out), logger.Debug, sessionId)
		logger.Log_PE("Please check the WorkFlowLog for more details.", logger.Debug, sessionId)
		if isWorkflow {
			logger.Log_PE("Test Workflow Invoke completed!", logger.Information, sessionId)
		} else {
			logger.Log_PE("Test Activity Invoke completed!", logger.Information, sessionId)
		}

		wfResponse.Message = "Test workflow invoke was successful with an output."
		wfResponse.Status = true
		wfResponse.ErrorDetails = ""
		wfResponse.ErrorCode = 1
	}

	// when the process is complete it passes the response back to the front
	return wfResponse
}
