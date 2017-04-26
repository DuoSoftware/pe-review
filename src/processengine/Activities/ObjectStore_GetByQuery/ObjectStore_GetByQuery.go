package ObjectStore_GetByQuery

import "processengine/context"
import "processengine/objectstore"
import "processengine/logger"

func Invoke(FlowData map[string]interface{}) (flowResult map[string]interface{}, activityResult *context.ActivityContext) {

	logger.Log_ACT("Starting Get by Query request.", logger.Debug, FlowData["InSessionID"].(string))
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
	query := FlowData["query"]

	// create the instance from objectstore
	n := objectstore.GetByQuery{}

	var parameters map[string]interface{}
	parameters = make(map[string]interface{})
	parameters["securityToken"] = securityToken
	parameters["log"] = log
	parameters["namespace"] = namespace
	parameters["class"] = class
	parameters["query"] = query

	ss := n.Invoke(parameters)
	//fmt.Println(string(ss.ResultMessage))

	//setting activityContext property values
	activityContext.ActivityStatus = true
	activityContext.Message = "New fields are added to the map"
	activityContext.ErrorState = activityError

	logger.Log_ACT("Received Data:", logger.Debug, FlowData["InSessionID"].(string))
	logger.Log_ACT(string(ss.SharedContext), logger.Debug, FlowData["InSessionID"].(string))

	FlowData["OutData"] = string(ss.SharedContext)
	FlowData["custMsg"] = "All the data retrived successfully"

	return FlowData, activityContext
}
