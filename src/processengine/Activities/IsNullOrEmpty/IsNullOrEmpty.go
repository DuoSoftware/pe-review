package IsNullOrEmpty

import "processengine/context"
import "processengine/logger"
import "strings"
import "reflect"

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

	variableStatus := ""
	variableStatusBool := false

	if FlowData["variable"] == nil {
		variableStatus = "null"
	} else {
		if reflect.TypeOf(FlowData["variable"]).String() == "string" {
			if strings.TrimSpace(FlowData["variable"].(string)) == "" {
				variableStatus = "empty"
			} else {
				variableStatus = "valid"
				variableStatusBool = true
			}
		} else {
			variableStatus = "invalid"
		}
	}

	msg := "Successfully Checked!"
	activityContext.ActivityStatus = true
	activityContext.Message = msg
	activityContext.ErrorState = activityError
	logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
	FlowData["custMsg"] = msg
	FlowData["status"] = "true"
	FlowData["VariableStatus"] = variableStatus
	FlowData["VariableStatusInBool"] = variableStatusBool

	return FlowData, activityContext
}
