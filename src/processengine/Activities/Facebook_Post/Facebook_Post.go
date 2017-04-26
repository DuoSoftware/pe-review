package Facebook_Post

import "processengine/context"
import "encoding/json"
import "processengine/logger"
import httpservice "processengine/Activities/HTTP_DefaultRequest"

//import facebook "github.com/huandu/facebook"
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

	var err error
	sessionId := FlowData["InSessionID"].(string)

	if FlowData["AccessToken"] == nil || FlowData["Message"] == nil || FlowData["IsPostPage"] == nil || FlowData["IsPostProfile"] == nil {
		err = errors.New("Few Arguments are missing. Check all input arguments.")
	} else {
		accessToken := FlowData["AccessToken"].(string)
		message := FlowData["Message"].(string)
		IsPostPage := FlowData["IsPostPage"].(bool)
		IsPostProfile := FlowData["IsPostProfile"].(bool)
		var pageID string
		if IsPostPage {
			pageID = FlowData["PageID"].(string)
		} else {
			pageID = "me"
		}

		httpFlowResult := make(map[string]interface{})
		var httpActivityResult *context.ActivityContext

		httpServiceInput := make(map[string]interface{})
		httpServiceInput["InSessionID"] = FlowData["InSessionID"]
		httpServiceInput["Method"] = "POST"
		httpServiceInput["Body"] = "message=" + message

		if IsPostPage && IsPostProfile {
			httpServiceInput["URL"] = "https://graph.facebook.com/me/feed?access_token=" + accessToken
			httpFlowResult, httpActivityResult = httpservice.Invoke(httpServiceInput)

			httpServiceInput["URL"] = "https://graph.facebook.com/" + pageID + "/feed?access_token=" + accessToken
			httpFlowResult, httpActivityResult = httpservice.Invoke(httpServiceInput)

		} else if IsPostPage && !IsPostProfile {
			httpServiceInput["URL"] = "https://graph.facebook.com/" + pageID + "/feed?access_token=" + accessToken
			httpFlowResult, httpActivityResult = httpservice.Invoke(httpServiceInput)
		} else if !IsPostPage && IsPostProfile {
			httpServiceInput["URL"] = "https://graph.facebook.com/me/feed?access_token=" + accessToken
			httpFlowResult, httpActivityResult = httpservice.Invoke(httpServiceInput)
		}

		if !httpActivityResult.ActivityStatus {
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
				logger.Log_ACT("Posting to Facebook Success!", logger.Debug, sessionId)
				msg := "Request Successful!"
				FlowData["PostID"] = response["id"].(string)
				logger.Log_ACT(msg, logger.Debug, sessionId)
				activityContext.ActivityStatus = true
				activityContext.Message = msg
				activityContext.ErrorState = activityError
			}
		}

	}

	if err == nil {
		msg := "Successfully Posted to Facebook!"
		activityContext.ActivityStatus = true
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		FlowData["custMsg"] = msg
		FlowData["status"] = "true"
	} else {
		msg := "Error Posting to Facebook" + err.Error()
		activityContext.ActivityStatus = false
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		FlowData["custMsg"] = msg
		FlowData["status"] = "false"
	}

	return FlowData, activityContext
}

// Older code with Go API

/*

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

	var err error

	if FlowData["AppID"] == nil || FlowData["AppSecret"] == nil || FlowData["UserID"] == nil || FlowData["Message"] == nil {
		err = errors.New("Few Arguments are missing. Check all input arguments.")
	} else {
		appID := FlowData["AppID"].(string)
		appSecret := FlowData["AppSecret"].(string)
		userId := FlowData["UserID"].(string)
		message := FlowData["Message"].(string)

		app := facebook.New(appID, appSecret)
		accessToken := app.AppAccessToken()

		_, err = facebook.Post(userId, facebook.Params{
			"message":      message,
			"access_token": accessToken,
		})
	}

	if err == nil {
		msg := "Successfully Posted to Facebook!"
		activityContext.ActivityStatus = true
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		FlowData["custMsg"] = msg
		FlowData["status"] = "true"
	} else {
		msg := "Error Posting to Facebook" + err.Error()
		activityContext.ActivityStatus = false
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		FlowData["custMsg"] = msg
		FlowData["status"] = "false"
	}

	return FlowData, activityContext
}
*/
