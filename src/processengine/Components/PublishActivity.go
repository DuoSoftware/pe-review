package Components

import (
	"os"
	// "strconv"  // commenting cuz the activity backup is now commented
	// "time"     // commenting cuz the activity backup is now commented
	"io/ioutil"
	"processengine/context"
	"processengine/logger"
	"runtime"
	"strings"
)

func PublishActivityFile(activityStruct ActivityStruct, activityName, sessionId string) *context.FlowResult {

	if len(strings.TrimSpace(activityName)) == 0 {
		var flowResult = new(context.FlowResult)
		flowResult.Status = false
		flowResult.Message = "No Flow Name specified!"
		flowResult.SessionID = sessionId
		return flowResult
	} //else continue....

	var flowResult = new(context.FlowResult)
	flowResult.FlowName = activityName
	flowResult.SessionID = sessionId
	flowResult.Status = true
	flowResult.Message = "Starting activity publish method."

	logger.Log_PE("~~Initiating Activity Publish", logger.Information, sessionId)
	logger.Log_PE("Activity Name: "+activityName, logger.Debug, sessionId)
	logger.Log_PE("Session ID:"+sessionId, logger.Debug, sessionId)

	/*red := color.New(color.FgWhite)
	ErrorColorscheme := red.Add(color.BgRed)



	decoder := json.NewDecoder(request.Body)
	var activityStruct ActivityStruct
	decodeError := decoder.Decode(&activityStruct)
	if decodeError != nil {
		fmt.Println("There was an error Decoding the jsonData sent to activity publish method.")
		ErrorColorscheme.Println(decodeError.Error())
		flowResult.Message = flowResult.Message + decodeError.Error() + " -> "
		flowResult.Status = false
		fmt.Println("")
		return flowResult
		}*/

	var slash = ""
	if runtime.GOOS == "windows" {
		slash = "\\"
	} else {
		slash = "/"
	}

	// check if the Activity folder exists or not, create one if not available
	pwd, _ := os.Getwd()
	publishPathRoot := pwd + slash + "src" + slash + "processengine" + slash + "Activities"
	//publishPathRoot := ".." + slash + "Activities"
	_, publishRooterr := os.Stat(publishPathRoot)
	if publishRooterr != nil {
		// create folder in the given path and permissions
		os.Mkdir(publishPathRoot, 0777)
	} else {
		msg := "Activity folder already created."
		logger.Log_PE(msg, logger.Debug, sessionId)
		flowResult.Message = msg
	}

	// check if the publishing activity folder exists, if not create one.

	activityRootPath := publishPathRoot + slash + activityName
	_, activityRooterr := os.Stat(activityRootPath)
	if activityRooterr != nil {
		// create folder in the given path and permissions
		os.Mkdir(activityRootPath, 0777)
	} else {
		logger.Log_PE("Activity Root folder creation failed for some reason.", logger.Debug, sessionId)
	}

	// check if the file already exist. if so rename it and make a backup
	archiveFlow := IsExists(activityRootPath + slash + activityName + ".go")
	if archiveFlow == true {
		/*newName := strconv.Itoa(time.Now().Year()) + "_" + time.Now().Month().String() + "_" + strconv.Itoa(time.Now().Day()) + "_" + strconv.Itoa(time.Now().Hour()) + "_" + strconv.Itoa(time.Now().Minute()) + "_" + strconv.Itoa(time.Now().Second())
		logger.Log_PE("Archiving existing activity file", logger.Debug, sessionId)
		flowResult.Message = "Archiving existing activity file."
		err := os.Rename(activityRootPath+slash+activityName+".go", activityRootPath+slash+activityName+"_"+newName+".go")
		if err != nil {
			flowResult.Message = err.Error()
			logger.Log_PE(err.Error(), logger.Debug, sessionId)
			flowResult.Status = false
			return flowResult
		}*/
	} else {
		msg := "New Activity file was created with the name: " + activityName + ".go"
		logger.Log_PE(msg, logger.Information, sessionId)
		flowResult.Message = msg
	}

	// save the go gile in the given path
	ioutil.WriteFile(activityRootPath+slash+activityName+".go", []byte(activityStruct.GoCode), 0777)

	return flowResult
}
