package Get_UTC_Time

import "processengine/context"
import "time"

//import "fmt"

// invoke method on objectore to insert
func Invoke(FlowData map[string]interface{}) (flowResult map[string]interface{}, activityResult *context.ActivityContext) {

	//creating new instance of context.ActivityContext
	var activityContext = new(context.ActivityContext)

	//creating new instance of context.ActivityError
	var activityError context.ActivityError

	//setting activityError proprty values
	activityError.Encrypt = false
	activityError.ErrorString = "exception"
	activityError.Forward = false
	activityError.SeverityLevel = context.Info
	var err error

	nowTime := time.Now().UTC()
	timeInString := nowTime.Format(time.RFC3339)

	if err != nil {
		msg := "Getting UTC time : " + err.Error()
		activityContext.ActivityStatus = false
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		//		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		FlowData["response"] = msg
		FlowData["status"] = "false"
	} else {
		msg := "Getting UTC time successful!"
		activityContext.ActivityStatus = true
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		//		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		FlowData["TimeUTC"] = timeInString
	}

	return FlowData, activityContext
}
