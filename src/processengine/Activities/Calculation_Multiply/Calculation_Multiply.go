package Calculation_Multiply

import "processengine/context"

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

	ValueOne := FlowData["ValueOne"].(float64)
	ValueTwo := FlowData["ValueTwo"].(float64)

	result := ValueOne * ValueTwo

	if err != nil {
		msg := "Getting UTC time : " + err.Error()
		activityContext.ActivityStatus = false
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		//		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		FlowData["Result"] = msg
		FlowData["status"] = "false"
	} else {
		msg := "Multiply successfull."
		activityContext.ActivityStatus = true
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		//		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		FlowData["Result"] = result
	}

	return FlowData, activityContext
}
