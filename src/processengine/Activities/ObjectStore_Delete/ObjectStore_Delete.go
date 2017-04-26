package ObjectStore_Delete

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

	// initiating the object creating the parameters for the object store call
	IsManyObjects := false
	DataObjectString := FlowData["ID"].(string)
	if string(DataObjectString[0]) == "[" {
		IsManyObjects = true
	}

	var insertObj = new(Components.InsertTemplate)

	if IsManyObjects {
		var multipleObjects []string
		err := json.Unmarshal([]byte(DataObjectString), &multipleObjects)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			deleteMultipleMap := make([]map[string]interface{}, len(multipleObjects))
			for _, singleID := range multipleObjects {
				singleEntry := make(map[string]interface{})
				singleEntry[FlowData["KeyProperty"].(string)] = singleID
				deleteMultipleMap = append(deleteMultipleMap, singleEntry)
			}
			insertObj.Objects = deleteMultipleMap
		}
	} else {
		singleObject := make(map[string]interface{})
		singleObject[FlowData["KeyProperty"].(string)] = FlowData["ID"].(string)
		insertObj.Object = singleObject
	}
	//insertObj.Object = `{"` + FlowData["KeyProperty"].(string) + `":"` + FlowData["ID"].(string) + `"}`
	insertObj.Parameters = parameterObj

	// converting the struct into a JSON stucture
	convertedObj, err := json.Marshal(insertObj)
	if err != nil {
		logger.Log_ACT(err.Error(), logger.Debug, FlowData["InSessionID"].(string))
	}
	logger.Log_ACT(string(convertedObj), logger.Debug, FlowData["InSessionID"].(string))

	// getting and assigning the values from the Map to the internal variables
	securityToken := FlowData["securityToken"]
	log := FlowData["log"]
	namespace := FlowData["namespace"]
	class := FlowData["class"]

	// create the instance from objectstore
	n := objectstore.Delete{} //Insert, Update, Delete

	// assigning values to the method calling map objects
	var parameters map[string]interface{}
	parameters = make(map[string]interface{})
	parameters["securityToken"] = securityToken
	parameters["log"] = log
	parameters["namespace"] = namespace
	parameters["class"] = class
	parameters["JSON"] = string(convertedObj)

	// invoking the method Delete on object stores
	ss := n.Invoke(parameters)
	logger.Log_ACT(string(ss.Message), logger.Debug, FlowData["InSessionID"].(string))

	//setting activityContext property values
	activityContext.ActivityStatus = true
	activityContext.Message = "The data deleted successfully"
	activityContext.ErrorState = activityError

	FlowData["custMsg"] = "The data deleted successfully"

	return FlowData, activityContext
}
