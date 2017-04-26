package JIRA_PostCommentOnIssue

import (
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

	logger.Log_ACT("Executing JIRA_PostCommentOnIssue Activity.", logger.Debug, InSessionID)

	if FlowData["JiraHost"] != nil && FlowData["JiraAuthToken"] != nil && FlowData["JiraIssueID"] != nil && FlowData["JiraComment"] != nil {

		jiraAuthToken := FlowData["JiraAuthToken"].(string)

		if jiraAuthToken != "Unauthorized" {
			jiraHostID := FlowData["JiraHost"].(string)
			jiraHostID = strings.TrimSpace(jiraHostID)
			jiraHostID = strings.TrimSuffix(jiraHostID, "/")
			completeurl := jiraHostID + "/rest/api/2/issue/" + FlowData["JiraIssueID"].(string) + "/comment"

			httpServiceInput["InSessionID"] = InSessionID
			httpServiceInput["URL"] = completeurl
			httpServiceInput["Method"] = "POST"

			httpServiceInput["Body"] = "{\"body\":\"" + FlowData["JiraComment"].(string) + "\"}"

			headerTokens := make(map[string]string)
			headerTokens["Content-Type"] = "application/json"

			if strings.Contains(jiraAuthToken, ".") {
				headerTokens["Authorization"] = "Bearer " + jiraAuthToken
			} else {
				headerTokens["Authorization"] = "Basic " + jiraAuthToken
			}

			responseExceptions := make([]string, 1)
			responseExceptions[0] = "201"

			httpServiceInput["ResponseCodeExceptions"] = responseExceptions

			httpServiceInput["headerTokens"] = headerTokens
		} else {
			err = errors.New("Error! JIRA Authentication Failed.")
		}

	} else {
		err = errors.New("Error! Required Fields Error. Check all fields.")
	}

	if err == nil {
		_, httpActivityResult := httpservice.Invoke(httpServiceInput)
		if !httpActivityResult.ActivityStatus {
			err = errors.New(httpActivityResult.ErrorState.ErrorString)
			activityError.ErrorString = "Exception " + err.Error()
		}
	}

	if err == nil {
		msg := "Successfully posted comment on issue!"
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

	logger.Log_ACT("Finished Executing JIRA_PostCommentOnIssue Activity. Returning to Work Flow.", logger.Debug, InSessionID)

	return FlowData, activityContext
}
