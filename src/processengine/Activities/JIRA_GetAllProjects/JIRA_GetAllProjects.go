package JIRA_GetAllProjects

import (
	"encoding/json"
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
	InSessionID := FlowData["InSessionID"].(string)
	httpServiceInput := make(map[string]interface{})

	logger.Log_ACT("Executing JIRA_GetAllProjects Activity.", logger.Debug, InSessionID)

	if FlowData["JiraHost"] != nil && FlowData["JiraAuthToken"] != nil {

		jiraAuthToken := FlowData["JiraAuthToken"].(string)

		if jiraAuthToken != "Unauthorized" {
			jiraHostID := FlowData["JiraHost"].(string)
			jiraHostID = strings.TrimSpace(jiraHostID)
			jiraHostID = strings.TrimSuffix(jiraHostID, "/")
			completeurl := jiraHostID + "/rest/api/2/project"

			httpServiceInput["InSessionID"] = InSessionID
			httpServiceInput["URL"] = completeurl
			httpServiceInput["Method"] = "GET"
			httpServiceInput["Body"] = ""

			headerTokens := make(map[string]string)
			headerTokens["Content-Type"] = "application/json"

			if strings.Contains(jiraAuthToken, ".") {
				headerTokens["Authorization"] = "Bearer " + jiraAuthToken
			} else {
				headerTokens["Authorization"] = "Basic " + jiraAuthToken
			}

			httpServiceInput["headerTokens"] = headerTokens
		} else {
			err = errors.New("Error! JIRA Authentication Failed.")
		}

	} else {
		err = errors.New("Error! Required Fields Error. Check all fields.")
	}

	if err == nil {
		httpFlowResult, httpActivityResult := httpservice.Invoke(httpServiceInput)
		if !httpActivityResult.ActivityStatus {
			err = errors.New(httpActivityResult.ErrorState.ErrorString)
			activityError.ErrorString = "Exception " + err.Error()
		} else {
			var data []map[string]interface{}
			_ = json.Unmarshal([]byte(httpFlowResult["Response"].(string)), &data)
			FlowData["JiraProjects"] = data
		}
	}

	if err == nil {
		msg := "Successfully Retrieved All Projects for user!"
		activityContext.ActivityStatus = true
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, InSessionID)
	} else {
		msg := "Error : " + err.Error()
		activityContext.ActivityStatus = false
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, InSessionID)
	}

	logger.Log_ACT("Finished Executing JIRA_GetAllProjects Activity. Returning to Work Flow.", logger.Debug, InSessionID)

	return FlowData, activityContext
}
