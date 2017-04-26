package Hive_InsertMultiple

import "processengine/context"
import "processengine/logger"
import "duov6.com/objectstore/client"
import "encoding/json"

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

	keyProperty := FlowData["KeyProperty"].(string)
	securityToken := FlowData["securityToken"].(string)
	//log := FlowData["log"].(string)
	namespace := FlowData["namespace"].(string)
	class := FlowData["class"].(string)
	JSONObjects := FlowData["JSONObjects"].(string)

	Host := FlowData["Host"].(string)
	Port := FlowData["Port"].(string)

	//make byte array
	byteVal := []byte(JSONObjects)
	//make interface
	var object []map[string]interface{}
	//unmarshall
	_ = json.Unmarshal(byteVal, &object)

	settings := make(map[string]interface{})
	settings["DB_Type"] = "HIVE"
	settings["Host"] = Host
	settings["Port"] = Port

	objectInterface := make([]interface{}, len(object))
	for i := range object {
		objectInterface[i] = object[i]
	}

	err := client.GoSmoothFlow(securityToken, namespace, class, settings).StoreObject().WithKeyField(keyProperty).AndStoreMany(objectInterface).Ok()

	if err == nil {
		msg := "Successfully Inserted multiple objects!"
		activityContext.ActivityStatus = true
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		FlowData["custMsg"] = msg
		FlowData["status"] = "true"
	} else {
		msg := "Error when inserting multiple objects"
		activityContext.ActivityStatus = false
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		FlowData["custMsg"] = msg
		FlowData["status"] = "false"
	}

	return FlowData, activityContext
}
