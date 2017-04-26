package PayPal_Get_Access_Token

import "processengine/context"
import "processengine/logger"
import "github.com/logpacker/PayPal-Go-SDK"
import "errors"

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

	var clientID string
	var secretID string
	var accessToken *paypalsdk.TokenResponse
	var err error

	InSessionID := FlowData["InSessionID"].(string)

	if FlowData["clientID"] != nil {
		clientID = FlowData["clientID"].(string)

		if FlowData["secretID"] != nil {
			secretID = FlowData["secretID"].(string)
		} else {
			err = errors.New("secretID not defined!")
		}

	} else {
		err = errors.New("clientID not defined!")
	}

	if err == nil {
		c, err := paypalsdk.NewClient(clientID, secretID, paypalsdk.APIBaseSandBox)
		if err == nil {
			accessToken, err = c.GetAccessToken()
		}
	}

	if err != nil {
		msg := err.Error()
		activityContext.ActivityStatus = false
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, InSessionID)
		FlowData["custMsg"] = msg
		FlowData["status"] = "false"
	} else {
		msg := "Successfully Retrieved Access Token"
		activityContext.ActivityStatus = true
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, InSessionID)
		FlowData["custMsg"] = msg
		FlowData["status"] = "true"
		FlowData["RefreshToken"] = accessToken.RefreshToken
		FlowData["Token"] = accessToken.Token
	}

	return FlowData, activityContext
}
