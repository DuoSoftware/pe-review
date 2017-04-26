package BigTable_Delete

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
	JSONObject := FlowData["JSONObject"].(string)

	dtype := FlowData["type"].(string)
	private_key_id := FlowData["private_key_id"].(string)
	private_key := FlowData["private_key"].(string)
	client_email := FlowData["client_email"].(string)
	client_id := FlowData["client_id"].(string)
	auth_uri := FlowData["auth_uri"].(string)
	token_uri := FlowData["token_uri"].(string)
	auth_provider_x509_cert_url := FlowData["auth_provider_x509_cert_url"].(string)
	client_x509_cert_url := FlowData["client_x509_cert_url"].(string)

	//make byte array
	byteVal := []byte(JSONObject)
	//make interface
	object := make(map[string]interface{})
	//unmarshall
	_ = json.Unmarshal(byteVal, &object)

	settings := make(map[string]interface{})
	settings["DB_Type"] = "GoogleBigTable"
	settings["type"] = dtype
	settings["private_key_id"] = private_key_id
	settings["private_key"] = private_key
	settings["client_email"] = client_email
	settings["client_id"] = client_id
	settings["auth_uri"] = auth_uri
	settings["token_uri"] = token_uri
	settings["auth_provider_x509_cert_url"] = auth_provider_x509_cert_url
	settings["client_x509_cert_url"] = client_x509_cert_url

	err := client.GoSmoothFlow(securityToken, namespace, class, settings).DeleteObject().WithKeyField(keyProperty).AndDeleteOne(object).Ok()

	if err == nil {
		msg := "Successfully Deleted!"
		activityContext.ActivityStatus = true
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		FlowData["custMsg"] = msg
		FlowData["status"] = "true"
	} else {
		msg := "Error when deleting"
		activityContext.ActivityStatus = false
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		FlowData["custMsg"] = msg
		FlowData["status"] = "false"
	}

	return FlowData, activityContext
}
