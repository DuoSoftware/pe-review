package jira

import (
	"duov6.com/common"
	"encoding/json"
	"fmt"
	"strings"
)

type JiraTriggerRequest struct {
	Id           int                    `json:"id"`
	Issue        map[string]interface{} `json:"issue"`
	User         map[string]interface{} `json:"user"`
	Changelog    map[string]interface{} `json:"changelog"`
	Comment      map[string]interface{} `json:"comment"`
	Timestamp    string                 `json:"timestamp"`
	WebhookEvent string                 `json:"webhookEvent"`
}

func InvokeJiraTrigger(trigger JiraTriggerRequest) (response JiraResponse) {
	jiraCon := GetAddonLink(GetJiraDomain(trigger))
	fmt.Println("")
	if jiraCon.Email == "" {
		response.IsSuccess = false
		response.Message = "No smoothflow account is authorized for this Jira Domain. Please contact your Jira Administrator."
	} else {
		triggerType := strings.Replace(trigger.WebhookEvent, "jira:", "", -1)

		fmt.Println("Tigger Event : " + triggerType)

		postBody := GetInvokeBody(trigger)

		isAllOkay := true

		for _, triggerUrl := range jiraCon.Triggers[triggerType] {
			err, _ := common.HTTP_POST(triggerUrl, nil, postBody, false)
			if err != nil {
				fmt.Println(err.Error())
				isAllOkay = false
			}
		}

		if isAllOkay {
			response.IsSuccess = true
			response.Message = "All trigger URLs have been invoked"
		} else {
			response.IsSuccess = false
			response.Message = "Error Invoking all or some of the trigger URLs"
		}
	}

	return
}

func GetJiraDomain(trigger JiraTriggerRequest) (domain string) {
	domain = trigger.User["self"].(string)
	domain = strings.Replace(domain, "https://", "", -1)
	domain = strings.Replace(domain, "http://", "", -1)
	tokens := strings.Split(domain, "/")
	domain = tokens[0]
	return
}

func GetInvokeBody(trigger JiraTriggerRequest) (body []byte) {
	bodyMap := make(map[string]interface{})
	bodyMap["InSessionID"] = CreateSessionID(trigger)
	bodyMap["InSecurityToken"] = "7b57e320a5b84c8a404918910edd0975"
	bodyMap["InLog"] = "log"
	bodyMap["InNamespace"] = GetJiraDomain(trigger)
	bodyMap["JiraTriggerID"] = trigger.Id
	bodyMap["JiraTriggerTimestamp"] = trigger.Timestamp
	bodyMap["JiraTriggerIssue"] = trigger.Issue
	bodyMap["JiraTriggerUser"] = trigger.User
	bodyMap["JiraTriggerChangelog"] = trigger.Changelog
	bodyMap["JiraTriggerComment"] = trigger.Comment
	bodyMap["JiraTriggerWebHookEvent"] = trigger.WebhookEvent
	body, _ = json.Marshal(bodyMap)
	return
}

func CreateSessionID(trigger JiraTriggerRequest) (session string) {
	domain := GetJiraDomain(trigger)
	session = common.EncodeToBase64(domain + "-" + common.RandomString(6))
	return
}

/*
{
	"id": 2,
 	"timestamp": "2009-09-09T00:08:36.796-0500",
	"issue": {
		"expand":"renderedFields,names,schema,transitions,operations,editmeta,changelog",
		"id":"99291",
		"self":"https://jira.atlassian.com/rest/api/2/issue/99291",
		"key":"JRA-20002",
		"fields":{
			"summary":"I feel the need for speed",
			"created":"2009-12-16T23:46:10.612-0600",
			"description":"Make the issue nav load 10x faster",
			"labels":["UI", "dialogue", "move"],
			"priority": "Minor"
		}
	},
	"user": {
		"self":"https://jira.atlassian.com/rest/api/2/user?username=brollins",
		"name":"prasad",
		"emailAddress":"prasad@duosoftware.com",
		"avatarUrls":{
			"16x16":"https://jira.atlassian.com/secure/useravatar?size=small&avatarId=10605",
			"48x48":"https://jira.atlassian.com/secure/useravatar?avatarId=10605"
		},
		"displayName":"Bryan Rollins [Atlassian]",
		"active" : "true"
	},
  	"changelog": {
        "items": [
            {
                "toString": "A new summary.",
                "to": null,
                "fromString": "What is going on here?????",
                "from": null,
                "fieldtype": "jira",
                "field": "summary"
            },
            {
                "toString": "New Feature",
                "to": "2",
                "fromString": "Improvement",
                "from": "4",
                "fieldtype": "jira",
                "field": "issuetype"
            }
        ],
		"id": 10124
	},
	"comment" : {
		"self":"https://jira.atlassian.com/rest/api/2/issue/10148/comment/252789",
		"id":"252789",
		"author":{
			"self":"https://jira.atlassian.com/rest/api/2/user?username=brollins",
			"name":"brollins",
			"emailAddress":"bryansemail@atlassian.com",
			"avatarUrls":{
				"16x16":"https://jira.atlassian.com/secure/useravatar?size=small&avatarId=10605",
				"48x48":"https://jira.atlassian.com/secure/useravatar?avatarId=10605"
			},
			"displayName":"Bryan Rollins [Atlassian]",
			"active":true
		},
		"body":"Just in time for AtlasCamp!",
		"updateAuthor":{
			"self":"https://jira.atlassian.com/rest/api/2/user?username=brollins",
			"name":"brollins",
			"emailAddress":"brollins@atlassian.com",
			"avatarUrls":{
				"16x16":"https://jira.atlassian.com/secure/useravatar?size=small&avatarId=10605",
				"48x48":"https://jira.atlassian.com/secure/useravatar?avatarId=10605"
			},
			"displayName":"Bryan Rollins [Atlassian]",
			"active":true
		},
		"created":"2011-06-07T10:31:26.805-0500",
		"updated":"2011-06-07T10:31:26.805-0500"
	},
    "webhookEvent": "jira:issue_updated"
}
*/
