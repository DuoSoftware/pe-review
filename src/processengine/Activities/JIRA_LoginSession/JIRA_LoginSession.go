package JIRA_LoginSession

import (
	"duov6.com/common"
	"errors"
	httpservice "processengine/Activities/HTTP_DefaultRequest"
	"processengine/context"
	"processengine/logger"
	"strings"
)

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
	var token string
	InSessionID := FlowData["InSessionID"].(string)
	httpServiceInput := make(map[string]interface{})

	logger.Log_ACT("Executing JIRA_LoginSession Activity.", logger.Debug, InSessionID)

	if FlowData["JiraHost"] != nil && FlowData["Username"] != nil && FlowData["Password"] != nil {

		username := FlowData["Username"].(string)
		password := FlowData["Password"].(string)
		jiraHostID := FlowData["JiraHost"].(string)
		jiraHostID = strings.TrimSpace(jiraHostID)
		jiraHostID = strings.TrimSuffix(jiraHostID, "/")
		completeurl := jiraHostID + "/rest/api/2/myself"

		token = common.EncodeToBase64(username + ":" + password)

		httpServiceInput["InSessionID"] = InSessionID
		httpServiceInput["URL"] = completeurl
		httpServiceInput["Method"] = "GET"
		httpServiceInput["Body"] = ""

		headerTokens := make(map[string]string)
		headerTokens["Content-Type"] = "application/json"
		headerTokens["Authorization"] = "Basic " + token
		httpServiceInput["headerTokens"] = headerTokens
	} else {
		err = errors.New("Error! Required Fields Error. Check all fields.")
	}

	if err == nil {
		_, httpActivityResult := httpservice.Invoke(httpServiceInput)
		if !httpActivityResult.ActivityStatus {
			err = errors.New(httpActivityResult.ErrorState.ErrorString)
			activityError.ErrorString = "Exception " + err.Error()
		}
		delete(FlowData, "Username")
		delete(FlowData, "Password")
	}

	FlowData["JiraAuthMethod"] = "LoginSession"

	if err == nil {
		msg := "Successfully Logged in via Session HTTP method!"
		activityContext.ActivityStatus = true
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, InSessionID)
		FlowData["JiraAuthToken"] = token
	} else {
		msg := "Error : " + err.Error()
		activityContext.ActivityStatus = false
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, InSessionID)
		FlowData["JiraAuthToken"] = "Unauthorized"
	}

	logger.Log_ACT("Finished Executing JIRA_LoginSession Activity. Returning to Work Flow.", logger.Debug, InSessionID)

	return FlowData, activityContext
}
