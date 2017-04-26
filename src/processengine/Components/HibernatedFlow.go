package Components

import (
	"encoding/json"
	"processengine/logger"
	"processengine/objectstore"
	"time"
)

// method used to call a hibernated workflow
func HibernateWorkflow(FlowData map[string]interface{}, ExecutionLevel, WFName string) (Status bool) {
	status := false

	// initiating the struct with relavent data to do the processing
	var wfInstance HibernatedWF
	wfInstance.FlowData = FlowData
	wfInstance.DateTime = time.Now().String()
	wfInstance.SessionID = FlowData["InSessionID"].(string)
	wfInstance.WFName = WFName
	wfInstance.ExecutionLevel = FlowData["ExecutionLevel"].(string)
	sessionID := FlowData["InSessionID"].(string)

	logger.Log_WF("Starting Hibernated WF object.", logger.Information, sessionID)

	// convert the struct to JSON
	wfconverted, convertionerr := json.Marshal(wfInstance)
	if convertionerr != nil {
		logger.Log_WF("Hibernated WF object convertion failed.", logger.Error, sessionID)
	}

	// prepare the JSON string for the process
	JSON := `{"Object":` + string(wfconverted) + `,"Parameters":{"KeyProperty":"SessionID"}}`
	logger.Log_WF("NEW Object: "+JSON, logger.Debug, sessionID)

	// get and assign data from MAP to internal variables
	securityToken := FlowData["InSecurityToken"]
	log := FlowData["InLog"]
	namespace := FlowData["InNamespace"]
	class := "hibernated_workflows"

	// if the execution level is 0, the hibernated data will be inserted to the flow.. else it should be updated.
	// create the instance from objectstore
	n := objectstore.Insert{} //Insert, Update, Delete

	// prepare data related for the process
	var parameters map[string]interface{}
	parameters = make(map[string]interface{})
	parameters["securityToken"] = securityToken
	parameters["log"] = log
	parameters["namespace"] = namespace
	parameters["class"] = class
	parameters["JSON"] = string(JSON)

	// performe the process
	ss := n.Invoke(parameters)
	logger.Log_WF(string(ss.Message), logger.Debug, sessionID)
	//logger.Log_WF(ss.ActivityStatus)

	status = ss.ActivityStatus

	// return the status once complted
	return status
}

// once a hibernated flow is completed its stored data on objectstore coule be removed with this method
func DeleteSession(FlowData map[string]interface{}) {
	sessionID := FlowData["InSessionID"].(string)
	// prepare the JSON object for the request
	JSON := `{"Object":{"SessionID":"` + sessionID + `"},"Parameters":{"KeyProperty":"SessionID"}}`
	logger.Log_WF("NEW Object: "+JSON, logger.Debug, sessionID)

	// get and assign data from MAP to internal variables
	securityToken := FlowData["InSecurityToken"].(string)
	log := FlowData["InLog"].(string)
	namespace := FlowData["InNamespace"].(string)
	class := "hibernated_workflows"

	// create the instance from objectstore
	n := objectstore.Delete{} //Insert, Update, Delete

	// prepare MAP
	var parameters map[string]interface{}
	parameters = make(map[string]interface{})
	parameters["securityToken"] = securityToken
	parameters["log"] = log
	parameters["namespace"] = namespace
	parameters["class"] = class
	parameters["JSON"] = JSON

	// invoking the method
	ss := n.Invoke(parameters)
	logger.Log_WF(string(ss.Message), logger.Debug, sessionID)
}
