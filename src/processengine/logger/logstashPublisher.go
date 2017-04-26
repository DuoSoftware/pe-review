package logger

import (
	"duov6.com/common"
	"encoding/json"
	"strings"
	"time"
)

func PublishToLogstash(url, sessionID, logMessage string, logType int, logCategory int) {
	domain := GetDomainBySessionID(sessionID)
	logNames := GetLogNamesByCategory(logCategory, sessionID)
	otherData := make(map[string]interface{})

	for _, logName := range logNames {
		logName = strings.TrimSuffix(logName, ".log")
		if strings.Contains(logName, "_") {
			//request one to actual user domain
			PostToLogStash(url, domain, logName, sessionID, logMessage, GetLogTypeByID(logType), GetCatTypeByID(logCategory), otherData)
		} else {
			//request two to common log
			if len(logNames) > 1 {
				PostToLogStash(url, logName, strings.TrimSuffix(logNames[1], ".log"), sessionID, logMessage, GetLogTypeByID(logType), GetCatTypeByID(logCategory), otherData)
			}
		}
	}

}

func PostToLogStash(url, domain, class, sessionID, logMessage, logType, category string, otherData interface{}) {
	payload := make(map[string]interface{})
	domain = strings.ToLower(domain)
	nowTime := time.Now()
	payload["Domain"] = domain
	payload["Class"] = class
	payload["SessionID"] = sessionID
	// payload["TimeStamp"] = nowTime.Format("2006-01-02T15:04:05")
	payload["TimeStamp"] = nowTime.Format("2006-01-02T15:04:05.000Z")
	payload["Type"] = logType
	payload["Message"] = logMessage
	payload["Category"] = category

	byteData, _ := json.Marshal(payload)
	go common.HTTP_POST(url, nil, byteData, false)
}
