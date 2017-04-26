//PACKAGEID

import (
"processengine/logger"
"processengine/context"
)

// this is the main method
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

	logger.Log_ACT("This is a Tistuslabs message generated from an acitivity", FlowData["InSessionID"].(string))

	// getting the details from the input argument

	activityContext.ErrorState = activityError
	return FlowData, activityContext
}