package JIRA_CreateIssue

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

	logger.Log_ACT("Executing JIRA_CreateIssue Activity.", logger.Debug, InSessionID)

	if FlowData["JiraHost"] != nil && FlowData["JiraAuthToken"] != nil && FlowData["JiraProjectID"] != nil && FlowData["JiraProjectKey"] != nil && FlowData["JiraIssueSummary"] != nil && FlowData["JiraIssueType"] != nil {

		jiraAuthToken := FlowData["JiraAuthToken"].(string)

		if jiraAuthToken != "Unauthorized" {
			jiraHostID := FlowData["JiraHost"].(string)
			jiraHostID = strings.TrimSpace(jiraHostID)
			jiraHostID = strings.TrimSuffix(jiraHostID, "/")
			completeurl := jiraHostID + "/rest/api/2/issue"

			httpServiceInput["InSessionID"] = InSessionID
			httpServiceInput["URL"] = completeurl
			httpServiceInput["Method"] = "POST"
			httpServiceInput["Body"] = GetCreateIssueRequestBody(FlowData)

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
		httpFlowResult, httpActivityResult := httpservice.Invoke(httpServiceInput)
		if !httpActivityResult.ActivityStatus {
			err = errors.New(httpActivityResult.ErrorState.ErrorString)
			activityError.ErrorString = "Exception " + err.Error()
		} else {
			var data map[string]interface{}
			_ = json.Unmarshal([]byte(httpFlowResult["Response"].(string)), &data)
			FlowData["JiraCreateIssueResponse"] = data
		}
	}

	if err == nil {
		msg := "Successfully created the Issue!"
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

	logger.Log_ACT("Finished Executing JIRA_CreateIssue Activity. Returning to Work Flow.", logger.Debug, InSessionID)

	return FlowData, activityContext
}

func GetCreateIssueRequestBody(FlowData map[string]interface{}) (JsonString string) {
	issueMap := make(map[string]interface{})
	fieldsMap := make(map[string]interface{})
	//Map assigning logic goes here

	var projectID string
	var projectKey string
	var summary string
	var issueType string
	var assignee string
	var reporter string
	var priority string
	var labels []string
	var originalEstimate string
	var remainingEstimate string
	var description string
	var duedate string
	var versions []string
	var fixVersions []string

	projectID = FlowData["JiraProjectID"].(string)
	projectKey = FlowData["JiraProjectKey"].(string)
	summary = FlowData["JiraIssueSummary"].(string)
	issueType = FlowData["JiraIssueType"].(string)

	projectMap := make(map[string]interface{})
	projectMap["id"] = projectID
	projectMap["key"] = projectKey
	fieldsMap["project"] = projectMap

	fieldsMap["summary"] = summary

	issueTypeMap := make(map[string]interface{})
	issueTypeMap["name"] = issueType
	fieldsMap["issuetype"] = issueTypeMap

	if FlowData["JiraIssueAssignee"] != nil {
		assignee = FlowData["JiraIssueAssignee"].(string)

		assigneeMap := make(map[string]interface{})
		assigneeMap["name"] = assignee
		fieldsMap["assignee"] = assigneeMap
	}

	if FlowData["JiraIssueReporter"] != nil {
		reporter = FlowData["JiraIssueReporter"].(string)

		reporterMap := make(map[string]interface{})
		reporterMap["name"] = reporter
		fieldsMap["reporter"] = reporterMap
	}

	if FlowData["JiraIssuePriority"] != nil {
		priority = FlowData["JiraIssuePriority"].(string)

		priorityMap := make(map[string]interface{})
		priorityMap["name"] = priority
		fieldsMap["priority"] = priorityMap
	}

	if FlowData["JiraIssueLabels"] != nil {
		labels = FlowData["JiraIssueLabels"].([]string)
		fieldsMap["labels"] = labels
	}

	timetrackingMap := make(map[string]interface{})

	if FlowData["JiraIssueOriginalEstimate"] != nil {
		originalEstimate = FlowData["JiraIssueOriginalEstimate"].(string)
		timetrackingMap["originalEstimate"] = originalEstimate
	}

	if FlowData["JiraIssueRemainingEstimate"] != nil {
		remainingEstimate = FlowData["JiraIssueRemainingEstimate"].(string)
		timetrackingMap["remainingEstimate"] = remainingEstimate
	}

	if len(timetrackingMap) > 0 {
		fieldsMap["timetracking"] = timetrackingMap
	}

	if FlowData["JiraIssueDescription"] != nil {
		description = FlowData["JiraIssueDescription"].(string)
		fieldsMap["description"] = description
	}

	if FlowData["JiraIssueDueDate"] != nil {
		duedate = FlowData["JiraIssueDueDate"].(string)
		fieldsMap["duedate"] = duedate
	}

	if FlowData["JiraIssueVersions"] != nil {
		versions = FlowData["JiraIssueVersions"].([]string)

		var versionMap []map[string]interface{}
		for _, value := range versions {
			tempMap := make(map[string]interface{})
			tempMap["name"] = value
			versionMap = append(versionMap, tempMap)
		}
		fieldsMap["versions"] = versionMap
	}

	if FlowData["JiraIssueFixVersions"] != nil {
		fixVersions = FlowData["JiraIssueFixVersions"].([]string)

		var versionMap []map[string]interface{}
		for _, value := range fixVersions {
			tempMap := make(map[string]interface{})
			tempMap["name"] = value
			versionMap = append(versionMap, tempMap)
		}
		fieldsMap["fixVersions"] = versionMap
	}

	issueMap["fields"] = fieldsMap

	byteBody, _ := json.Marshal(issueMap)
	JsonString = string(byteBody)

	return
}
