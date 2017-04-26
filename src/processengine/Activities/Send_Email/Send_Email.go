package Send_Email

import "processengine/context"

//import httpservice "processengine/Activities/HTTP_DefaultRequest"
import "processengine/logger"
import "net/http"
import "io/ioutil"
import "strings"

//import "fmt"
import "encoding/json"
import "bytes"
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

	isError := false
	errorMessage := ""

	emailSendMap := make(map[string]interface{})
	emailSendMap["type"] = "newemail"

	if FlowData["ToAddress"].(string) == "" {
		isError = true
		errorMessage = "No Reciever Address Defined!"
	} else {
		emailSendMap["to"] = FlowData["ToAddress"].(string)
	}

	var subject string

	if FlowData["subject"] != nil {
		subject = FlowData["subject"].(string)
	}

	if FlowData["Subject"] != nil {
		subject = FlowData["Subject"].(string)
	}

	if subject == "" {
		emailSendMap["subject"] = "New Email from SmoothFlow!"
	} else {
		emailSendMap["subject"] = subject
	}

	if FlowData["FromAddress"].(string) == "" {
		isError = true
		errorMessage = "No Sender Address Defined!"
	} else {
		emailSendMap["from"] = FlowData["FromAddress"].(string)
	}

	if FlowData["Body"].(string) == "" {
		isError = true
		errorMessage = "No Message Body Defined!"
	} else {
		emailSendMap["body"] = FlowData["Body"].(string)
	}

	url := ""

	if FlowData["URL"].(string) == "" {
		isError = true
		errorMessage = "No CEB URL Defined!"
	} else {
		url = FlowData["URL"].(string)
	}

	securityToken := "ignore"

	if FlowData["securityToken"].(string) != "" {
		securityToken = FlowData["securityToken"].(string)
	}

	if isError {
		msg := "Email Sending Failed : " + errorMessage
		activityContext.ActivityStatus = false
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		FlowData["custMsg"] = msg
		FlowData["status"] = "false"
	} else {
		url = strings.TrimSuffix(url, "/")
		url += ":3500/command/notification"
		headers := make(map[string]string)
		headers["securityToken"] = securityToken
		headers["Content-Type"] = "application/json"

		dataBody, err := json.Marshal(emailSendMap)

		err = HTTP_POST(url, headers, dataBody)

		if err != nil {
			msg := "Email Sending Failed : " + err.Error()
			activityContext.ActivityStatus = false
			activityContext.Message = msg
			activityContext.ErrorState = activityError
			logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
			FlowData["response"] = msg
			FlowData["status"] = "false"
		} else {
			msg := "Email Successfully Sent!"
			activityContext.ActivityStatus = true
			activityContext.Message = msg
			activityContext.ErrorState = activityError
			logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
			FlowData["response"] = msg
			FlowData["status"] = "true"
		}

		// if err != nil {
		// 	msg := "Email Sending Failed : " + err.Error()
		// 	activityContext.ActivityStatus = false
		// 	activityContext.Message = msg
		// 	activityContext.ErrorState = activityError
		// 	logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		// 	FlowData["response"] = msg
		// 	FlowData["status"] = "false"
		// } else {

		// 	httpServiceInput := make(map[string]interface{})
		// 	httpServiceInput["InSessionID"] = FlowData["InSessionID"]
		// 	httpServiceInput["URL"] = url
		// 	httpServiceInput["Method"] = "POST"

		// 	headerTokens := make(map[string]string)
		// 	headerTokens["securityToken"] = FlowData["securityToken"].(string)

		// 	httpServiceInput["Body"] = string(dataBody)
		// 	httpServiceInput["headerTokens"] = headerTokens

		// 	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
		// 	fmt.Println(string(dataBody))
		// 	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
		// 	gg, httpActivityResult := httpservice.Invoke(httpServiceInput)

		// 	if !httpActivityResult.ActivityStatus {
		// 		msg := "Email Sending Failed : " + httpActivityResult.ErrorState.ErrorString
		// 		activityContext.ActivityStatus = false
		// 		activityContext.Message = msg
		// 		activityContext.ErrorState = activityError
		// 		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		// 		FlowData["response"] = msg
		// 		FlowData["status"] = "false"
		// 	} else {
		// 		fmt.Println(gg["Response"])

		// 		msg := "Email Successfully Sent!"
		// 		activityContext.ActivityStatus = true
		// 		activityContext.Message = msg
		// 		activityContext.ErrorState = activityError
		// 		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		// 		FlowData["response"] = msg
		// 		FlowData["status"] = "true"
		// 	}
		// }

	}

	return FlowData, activityContext
}

func HTTP_POST(url string, headers map[string]string, JSON_DATA []byte) (err error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(JSON_DATA))

	for headerName, headerValue := range headers {
		req.Header.Set(headerName, headerValue)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.New("Connection Failed!")
	} else {
		var body []byte
		body, _ = ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			err = errors.New(string(body))
		}
		defer resp.Body.Close()
	}
	return
}
