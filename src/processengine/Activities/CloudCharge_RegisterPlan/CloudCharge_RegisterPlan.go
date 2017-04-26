package CloudCharge_RegisterPlan

import (
	"encoding/json"
	"errors"
	httpservice "processengine/Activities/HTTP_DefaultRequest"
	"processengine/context"
	"processengine/logger"
	"strings"
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
	logger.Log_ACT("CloudCharge_RegisterPlan", logger.Debug, sessionId)

	var err error
	inputJson := make(map[string]interface{})
	var inputByteBody []byte
	var inputBodyInString string
	var host string
	httpServiceInput := make(map[string]interface{})

	if FlowData["Host"] != nil && FlowData["uom"] != nil && FlowData["code"] != nil && FlowData["product_name"] != nil && FlowData["category"] != nil && FlowData["price_of_unit"] != nil && FlowData["status"] != nil && FlowData["type"] != nil && FlowData["interval"] != nil && FlowData["securityToken"] != nil && FlowData["currency"] != nil {
		host = FlowData["Host"].(string)
		host = strings.TrimSpace(host)
		host = strings.TrimSuffix(host, "/")
		completeurl := host + "/duosoftware.paymentgateway.service/stripe/registerPlan"

		inputJson["uom"] = FlowData["uom"]
		inputJson["code"] = FlowData["code"]
		inputJson["product_name"] = FlowData["product_name"]
		inputJson["category"] = FlowData["category"]
		inputJson["price_of_unit"] = FlowData["price_of_unit"]
		inputJson["status"] = FlowData["status"]
		inputJson["type"] = FlowData["type"]
		subscriptionDetails := make(map[string]interface{})
		subscriptionDetails["interval"] = FlowData["interval"]
		subscriptionDetails["currency"] = FlowData["currency"]
		inputJson["subscriptionDetails"] = subscriptionDetails

		inputByteBody, _ = json.Marshal(inputJson)
		inputBodyInString = string(inputByteBody)

		httpServiceInput["InSessionID"] = FlowData["InSessionID"]
		httpServiceInput["URL"] = completeurl
		httpServiceInput["Method"] = "POST"
		httpServiceInput["Body"] = inputBodyInString

		securityToken := FlowData["securityToken"].(string)
		headerTokens := make(map[string]string)
		headerTokens["securityToken"] = securityToken
		httpServiceInput["headerTokens"] = headerTokens

	} else {
		err = errors.New("Error! Required Fields Error. Check all fields.")
	}

	// End input variable
	if err != nil {
		FlowData["status"] = false
		logger.Log_ACT(err.Error(), logger.Debug, sessionId)
		activityError.ErrorString = err.Error()
		activityContext.ActivityStatus = false
		activityContext.Message = err.Error()
		activityContext.ErrorState = activityError
	} else {

		httpFlowResult, httpActivityResult := httpservice.Invoke(httpServiceInput)

		if !httpActivityResult.ActivityStatus {
			FlowData["status"] = false
			logger.Log_ACT(("Request Error! : " + httpActivityResult.ErrorState.ErrorString), logger.Debug, sessionId)
			activityError.ErrorString = "Exception " + httpActivityResult.ErrorState.ErrorString
			activityContext.ActivityStatus = false
			activityContext.Message = "Respond Err"
			activityContext.ErrorState = activityError
		} else {

			response := make(map[string]interface{})
			err = json.Unmarshal([]byte(httpFlowResult["Response"].(string)), &response)
			if err != nil {
				FlowData["status"] = false
				logger.Log_ACT(err.Error(), logger.Debug, sessionId)
				activityError.ErrorString = err.Error()
				activityContext.ActivityStatus = false
				activityContext.Message = err.Error()
				activityContext.ErrorState = activityError
			} else {
				logger.Log_ACT("Request Successful!", logger.Debug, sessionId)
				msg := "Request Successful!"
				FlowData["status"] = response["status"].(bool)
				logger.Log_ACT(msg, logger.Debug, sessionId)
				activityContext.ActivityStatus = true
				activityContext.Message = msg
				activityContext.ErrorState = activityError
			}
		}
	}

	return FlowData, activityContext
}
