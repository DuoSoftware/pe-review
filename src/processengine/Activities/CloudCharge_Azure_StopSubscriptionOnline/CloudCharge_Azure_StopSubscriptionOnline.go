package CloudCharge_Azure_StopSubscriptionOnline

import (
	"encoding/json"
	"errors"
	"fmt"
	httpservice "processengine/Activities/HTTP_DefaultRequest"
	"processengine/context"
	"processengine/logger"
)

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

	sessionId := FlowData["InSessionID"].(string)

	// Start input variable
	logger.Log_ACT("CloudCharge_Azure_StopSubscriptionOnline", logger.Information, sessionId)

	var err error

	httpServiceInput := make(map[string]interface{})

	fmt.Println(FlowData)

	if FlowData["SubscriptionKey"] != nil && FlowData["email"] != nil && FlowData["planCode"] != nil && FlowData["changeType"] != nil {

		securityToken := FlowData["SubscriptionKey"].(string)

		object := make(map[string]interface{})
		object["email"] = FlowData["email"].(string)
		object["oldplanCode"] = FlowData["planCode"].(string)
		object["changeType"] = FlowData["changeType"].(string)

		objectInBytes, _ := json.Marshal(object)

		completeurl := "https://cloudcharge.azure-api.net/Subscription-Online/stopSubscription"

		httpServiceInput["InSessionID"] = FlowData["InSessionID"]
		httpServiceInput["URL"] = completeurl
		httpServiceInput["Method"] = "POST"
		httpServiceInput["Body"] = string(objectInBytes)

		headerTokens := make(map[string]string)
		headerTokens["Ocp-Apim-Subscription-Key"] = securityToken
		httpServiceInput["headerTokens"] = headerTokens

	} else {
		err = errors.New("Error! Required Fields Error. Check all fields.")
	}

	// End input variable
	if err != nil {
		FlowData["status"] = "false"
		logger.Log_ACT(err.Error(), logger.Error, sessionId)
		activityError.ErrorString = err.Error()
		activityContext.ActivityStatus = false
		activityContext.Message = err.Error()
		activityContext.ErrorState = activityError
	} else {
		httpFlowResult, httpActivityResult := httpservice.Invoke(httpServiceInput)

		if !httpActivityResult.ActivityStatus {
			FlowData["status"] = "false"
			logger.Log_ACT(("Request Error! : " + httpActivityResult.ErrorState.ErrorString), logger.Error, sessionId)
			activityError.ErrorString = "Exception " + httpActivityResult.ErrorState.ErrorString
			activityContext.ActivityStatus = false
			activityContext.Message = "Respond Err"
			activityContext.ErrorState = activityError
		} else {
			response := make(map[string]interface{})
			fmt.Println("--------------- Stop Subscription Response --------------------")
			fmt.Println(httpFlowResult["Response"])
			err = json.Unmarshal([]byte(httpFlowResult["Response"].(string)), &response)
			if err != nil {
				FlowData["status"] = "false"
				logger.Log_ACT(err.Error(), logger.Error, sessionId)
				activityError.ErrorString = err.Error()
				activityContext.ActivityStatus = false
				activityContext.Message = err.Error()
				activityContext.ErrorState = activityError
			} else {
				if FlowData["Response"] != nil {
					delete(FlowData, "Response")
				}
				data := response["data"].(map[string]interface{})
				logger.Log_ACT("Request Successful!", logger.Information, sessionId)
				msg := data["message"].(string)
				FlowData["status"] = "true"
				logger.Log_ACT(msg, logger.Debug, sessionId)
				activityContext.ActivityStatus = true
				activityContext.Message = msg
				activityContext.ErrorState = activityError
			}
		}
	}

	return FlowData, activityContext
}
