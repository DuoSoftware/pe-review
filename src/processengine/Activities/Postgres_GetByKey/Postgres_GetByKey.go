package Postgres_GetByKey

import "processengine/context"
import "processengine/logger"
import "duov6.com/objectstore/client"

//import "encoding/json"

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

	securityToken := FlowData["securityToken"].(string)
	KeyProperty := FlowData["KeyProperty"].(string)
	//log := FlowData["log"].(string)
	namespace := FlowData["namespace"].(string)
	class := FlowData["class"].(string)
	Username := FlowData["Username"].(string)
	Password := FlowData["Password"].(string)
	Url := FlowData["Url"].(string)
	Port := FlowData["Port"].(string)

	//make byte array
	//byteVal := []byte(JSONObject)
	//make interface
	//object := make(map[string]interface{})
	//unmarshall
	//_ = json.Unmarshal(byteVal, &object)

	settings := make(map[string]interface{})
	settings["DB_Type"] = "POSTGRES"
	settings["Username"] = Username
	settings["Password"] = Password
	settings["Url"] = Url
	settings["Port"] = Port

	data, err := client.GoSmoothFlow(securityToken, namespace, class, settings).GetOne().ByUniqueKey(KeyProperty).Ok()

	if err == "" {
		msg := "Successfully retrieved data!"
		activityContext.ActivityStatus = true
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		FlowData["custMsg"] = msg
		FlowData["status"] = "true"
		FlowData["Data"] = string(data)
	} else {
		msg := "Error when retriving data"
		activityContext.ActivityStatus = false
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		FlowData["custMsg"] = msg
		FlowData["status"] = "false"
		FlowData["Data"] = ""
	}

	return FlowData, activityContext
}
