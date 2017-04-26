package File_RemoveTxtFile

import "processengine/context"
import "runtime"
import "processengine/logger"
import "os"

// this method can be used to remove text files stored in the location
func Invoke(FlowData map[string]interface{}) (flowResult map[string]interface{}, activityResult *context.ActivityContext) {

	logger.Log_ACT("Starting txt file read request.", logger.Debug, FlowData["InSessionID"].(string))
	//creating new instance of context.ActivityContext
	var activityContext = new(context.ActivityContext)

	//creating new instance of context.ActivityError
	var activityError context.ActivityError

	//setting activityError proprty values
	activityError.Encrypt = false
	activityError.ErrorString = "exception"
	activityError.Forward = false
	activityError.SeverityLevel = context.Info

	// getting the details from the input arguments
	FileName := FlowData["InNamespace"].(string) + "_" + FlowData["FileName"].(string) + ".txt"

	// check if its running in windows or linux
	var slash = ""
	if runtime.GOOS == "windows" {
		slash = "\\"
	} else {
		slash = "/"
	}

	// check if the Activity folder exists or not, create one if not available
	pwd, _ := os.Getwd()
	txtFilePath := pwd + slash + "TextFiles"
	_, txtFilePatherr := os.Stat(txtFilePath)
	if txtFilePatherr != nil {
		// create folder in the given path and permissions
		os.Mkdir(txtFilePath, 0777)
	} else {
		logger.Log_ACT("Text File folder already available.", logger.Debug, FlowData["InSessionID"].(string))
	}

	// prepare the file path of the file to be removed
	concatinatedFilePath := txtFilePath + slash + FileName

	// remove the file from the location
	err := os.Remove(concatinatedFilePath)
	if err != nil {
		str := "Text file was deleted successful."
		FlowData["Response"] = str
		logger.Log_ACT(str, logger.Debug, FlowData["InSessionID"].(string))
		activityContext.ActivityStatus = true
		activityContext.Message = str
	} else {
		str := "There was no text file available."
		FlowData["Response"] = str
		logger.Log_ACT(str, logger.Debug, FlowData["InSessionID"].(string))
		activityContext.ActivityStatus = false
		activityContext.Message = str
	}

	activityContext.ErrorState = activityError

	logger.Log_ACT("Finishing txt file write request.", FlowData["InSessionID"].(string))
	// return the data once completed
	return FlowData, activityContext
}
