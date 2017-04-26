package CloudCharge_Azure_SearchProfile

import (
	"encoding/json"
	"errors"
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
	logger.Log_ACT("CloudCharge_Azure_SearchProfile", logger.Information, sessionId)

	var err error

	httpServiceInput := make(map[string]interface{})

	if FlowData["skip"] != nil && FlowData["take"] != nil && FlowData["order"] != nil && FlowData["column"] != nil && FlowData["keyword"] != nil && FlowData["SubscriptionKey"] != nil {

		skip := FlowData["skip"].(string)
		take := FlowData["take"].(string)
		order := FlowData["order"].(string)
		column := FlowData["column"].(string)
		keyword := FlowData["keyword"].(string)

		completeurl := "https://cloudcharge.azure-api.net/Profile/getProfile/?skip=" + skip + "&take=" + take + "&order=" + order + "&column=" + column + "&keyword=" + keyword

		httpServiceInput["InSessionID"] = FlowData["InSessionID"]
		httpServiceInput["URL"] = completeurl
		httpServiceInput["Method"] = "GET"
		httpServiceInput["Body"] = ""

		securityToken := FlowData["SubscriptionKey"].(string)
		headerTokens := make(map[string]string)
		headerTokens["Ocp-Apim-Subscription-Key"] = securityToken
		httpServiceInput["headerTokens"] = headerTokens

	} else {
		err = errors.New("Error! Required Fields Error. Check all fields.")
	}

	// End input variable
	if err != nil {
		FlowData["status"] = false
		logger.Log_ACT(err.Error(), logger.Error, sessionId)
		activityError.ErrorString = err.Error()
		activityContext.ActivityStatus = false
		activityContext.Message = err.Error()
		activityContext.ErrorState = activityError
	} else {
		httpFlowResult, httpActivityResult := httpservice.Invoke(httpServiceInput)

		if !httpActivityResult.ActivityStatus {
			FlowData["status"] = false
			logger.Log_ACT(("Request Error! : " + httpActivityResult.ErrorState.ErrorString), logger.Error, sessionId)
			activityError.ErrorString = "Exception " + httpActivityResult.ErrorState.ErrorString
			activityContext.ActivityStatus = false
			activityContext.Message = "Respond Err"
			activityContext.ErrorState = activityError
		} else {
			var response []map[string]interface{}
			err = json.Unmarshal([]byte(httpFlowResult["Response"].(string)), &response)
			if err != nil {
				FlowData["status"] = false
				logger.Log_ACT(err.Error(), logger.Error, sessionId)
				activityError.ErrorString = err.Error()
				activityContext.ActivityStatus = false
				activityContext.Message = err.Error()
				activityContext.ErrorState = activityError
			} else {
				logger.Log_ACT("Request Successful!", logger.Information, sessionId)
				msg := "Request Successful!"
				FlowData["status"] = true
				FlowData["Response"] = response[0]
				logger.Log_ACT(msg, logger.Debug, sessionId)
				activityContext.ActivityStatus = true
				activityContext.Message = msg
				activityContext.ErrorState = activityError
			}
		}
	}

	return FlowData, activityContext
}
