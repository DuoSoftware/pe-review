package ObjectStore_GetByKey

import "processengine/context"
import "processengine/objectstore"

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

	// getting the details from the input arguments
	securityToken := FlowData["securityToken"]
	log := FlowData["log"]
	namespace := FlowData["namespace"]
	class := FlowData["class"]
	key := FlowData["key"]

	// create the instance from objectstore
	n := objectstore.GetAll{}

	var parameters map[string]interface{}
	parameters = make(map[string]interface{})
	parameters["securityToken"] = securityToken
	parameters["log"] = log
	parameters["namespace"] = namespace
	parameters["class"] = class
	parameters["key"] = key

	ss := n.Invoke(parameters)
	//fmt.Println(string(ss.ResultMessage))

	//setting activityContext property values
	activityContext.ActivityStatus = true
	activityContext.Message = "New fields are added to the map"
	activityContext.ErrorState = activityError

	FlowData["OutData"] = string(ss.SharedContext)
	FlowData["custMsg"] = "All the data retrived successfully"

	return FlowData, activityContext
}
