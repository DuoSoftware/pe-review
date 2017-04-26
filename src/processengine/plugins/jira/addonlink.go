package jira

import (
	"duov6.com/common"
	"encoding/json"
	"processengine/Common"
)

type JiraAddonLink struct {
	Email      string
	JiraDomain string //PK
	Triggers   map[string][]string
}

func GetAddonLink(JiraDomain string) (jiraCon JiraAddonLink) {
	jiraCon = JiraAddonLink{}
	url := "http://" + Common.OBJECTSTORE_URL + "/com.smoothflow.io/jiraconnection/" + JiraDomain + "?securityToken=ignore"
	err, bodyBytes := common.HTTP_GET(url, nil, false)
	if err == nil && len(bodyBytes) > 4 {
		_ = json.Unmarshal(bodyBytes, &jiraCon)
	}
	return
}

func SetAddonLink(jiraCon JiraAddonLink) (err error) {
	url := "http://" + Common.OBJECTSTORE_URL + "/com.smoothflow.io/jiraconnection?securityToken=ignore"
	bytesVal, _ := json.Marshal(jiraCon)
	payload := `{"Object":` + string(bytesVal) + `, "Parameters":{"KeyProperty":"` + jiraCon.JiraDomain + `"}}`
	err, _ = common.HTTP_POST(url, nil, []byte(payload), false)
	return
}
