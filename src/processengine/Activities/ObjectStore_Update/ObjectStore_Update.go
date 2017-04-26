package ObjectStore_Update

import "processengine/context"
import "processengine/Components"
import "processengine/logger"
import "processengine/objectstore"
import "encoding/json"

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

	//create the insert object
	var parameterObj Components.InsParameters
	parameterObj.KeyProperty = FlowData["KeyProperty"].(string)

	// creating the object structure for the parameters for objectstore insert
	IsManyObjects := false
	DataObjectString := FlowData["JSONObject"].(string)
	if string(DataObjectString[0]) == "[" {
		IsManyObjects = true
	}

	var insertObj = new(Components.InsertTemplate)

	if IsManyObjects {
		var multipleObjects []map[string]interface{}
		err := json.Unmarshal([]byte(DataObjectString), &multipleObjects)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			insertObj.Objects = multipleObjects
		}
	} else {
		var singleObject map[string]interface{}
		err := json.Unmarshal([]byte(DataObjectString), &singleObject)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			insertObj.Object = singleObject
		}
	}

	insertObj.Parameters = parameterObj

	convertedObj, err := json.Marshal(insertObj)
	if err != nil {
		logger.Log_ACT(err.Error(), logger.Debug, FlowData["InSessionID"].(string))
	}
	logger.Log_ACT(string(convertedObj), logger.Debug, FlowData["InSessionID"].(string))

	securityToken := FlowData["securityToken"]
	log := FlowData["log"]
	namespace := FlowData["namespace"]
	class := FlowData["class"]

	// create the instance from objectstore
	n := objectstore.Update{} //Insert, Update, Delete

	var parameters map[string]interface{}
	parameters = make(map[string]interface{})
	parameters["securityToken"] = securityToken
	parameters["log"] = log
	parameters["namespace"] = namespace
	parameters["class"] = class
	parameters["JSON"] = string(convertedObj)

	ss := n.Invoke(parameters)
	logger.Log_ACT(string(ss.Message), logger.Debug, FlowData["InSessionID"].(string))

	//setting activityContext property values
	activityContext.ActivityStatus = true
	activityContext.Message = "New fields are added to the map"
	activityContext.ErrorState = activityError

	FlowData["custMsg"] = "The data updated successfully"

	return FlowData, activityContext
}
