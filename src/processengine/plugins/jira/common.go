package jira

import (
	"errors"
	"io/ioutil"
	"net/http"
)

type JiraUser struct {
	Self             string                 `json:"self"`
	Key              string                 `json:"key"`
	Name             string                 `json:"name"`
	EmailAddress     string                 `json:"emailAddress"`
	AvatarUrls       map[string]interface{} `json:"avatarUrls"`
	DisplayName      string                 `json:"displayName"`
	Active           bool                   `json:"active"`
	TimeZone         string                 `json:"timeZone"`
	Locale           string                 `json:"locale"`
	Groups           map[string]interface{} `json:"groups"`
	ApplicationRoles map[string]interface{} `json:"applicationRoles"`
	Expand           string                 `json:"expand"`
}

type JiraResponse struct {
	IsSuccess bool
	Message   string
	Data      interface{}
}

//------------------ Accesser methods --------------------

func Jira_HTTP_GET(url string, tokens map[string]string) (err error, body []byte) {
	req, err := http.NewRequest("GET", url, nil)

	for key, value := range tokens {
		if value != "" {
			c := &http.Cookie{Name: key, Value: value, HttpOnly: false}
			req.AddCookie(c)
		}
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		err = errors.New("Connection Failed!")
	} else {
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			err = errors.New(string(body))
		}
	}
	return
}
