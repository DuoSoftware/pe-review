package File_WriteToTxt

import "processengine/context"
import "runtime"
import "processengine/logger"
import "os"

//import "fmt"

// method used to write to text files
func Invoke(FlowData map[string]interface{}) (flowResult map[string]interface{}, activityResult *context.ActivityContext) {

	logger.Log_ACT("Starting txt file write request.", logger.Debug, FlowData["InSessionID"].(string))
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
	Content := FlowData["Content"].(string)

	// check if the the app currently runs in Windows or linux
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

	// prepare the path location of the text file
	concatinatedFilePath := txtFilePath + slash + FileName

	// open read the file available in the location
	ff, err := os.OpenFile(concatinatedFilePath, os.O_APPEND, 0666)

	if err != nil {
		ff, err = os.Create(concatinatedFilePath)
		ff, err = os.OpenFile(concatinatedFilePath, os.O_APPEND, 0666)
	}

	// insert the content into the text tile and add a new line to the text file
	_, err = ff.Write([]byte(Content))
	_, err = ff.Write([]byte("\r\n"))
	if err != nil {
		logger.Log_ACT("Error writing file", logger.Debug, FlowData["InSessionID"].(string))
		logger.Log_ACT(err.Error(), logger.Debug, FlowData["InSessionID"].(string))
		FlowData["Response"] = "There was an error writing to text file."
		activityContext.ActivityStatus = false
		activityContext.Message = "There was an error writing to text file."
	} else {
		FlowData["Response"] = "Write to text file was successful."
		logger.Log_ACT("Write to text file was successful.", logger.Debug, FlowData["InSessionID"].(string))
		activityContext.ActivityStatus = true
		activityContext.Message = "Write to text file was successful."
	}

	ff.Close()

	activityContext.ErrorState = activityError

	logger.Log_ACT("Finishing txt file write request.", logger.Debug, FlowData["InSessionID"].(string))
	// once finished return the relavent data
	return FlowData, activityContext
}
