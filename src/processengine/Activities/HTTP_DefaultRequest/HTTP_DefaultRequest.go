package HTTP_DefaultRequest

import "processengine/context"
import "net/http"
import "bytes"
import "processengine/logger"
import "strconv"
import "time"
import "io/ioutil"
import "strings"
import "errors"

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

	// getting the details from the input arguments
	URL := FlowData["URL"].(string)
	Method := FlowData["Method"].(string)
	var Body string
	sessionId := FlowData["InSessionID"].(string)
	retryCount := 10

	Method = strings.ToUpper(Method)

	headerTokens := make(map[string]string)
	excludeResponseCodes := make([]string, 0)

	if FlowData["retryCount"] != nil {
		retryCount, _ = strconv.Atoi(FlowData["retryCount"].(string))
	}

	if FlowData["Body"] != nil {
		Body = FlowData["Body"].(string)
	}

	if FlowData["headerTokens"] != nil {
		headerTokens = FlowData["headerTokens"].(map[string]string)
	}

	if FlowData["ResponseCodeExceptions"] != nil {
		for _, value := range FlowData["ResponseCodeExceptions"].([]string) {
			excludeResponseCodes = append(excludeResponseCodes, value)
		}
	}

	err, responseMessage := SendRequest(Method, URL, Body, headerTokens, sessionId, excludeResponseCodes, retryCount)

	if err != nil {
		logger.Log_ACT("There was an error while making the HTTP request.", logger.Debug, FlowData["InSessionID"].(string))
		activityError.ErrorString = err.Error()
		activityContext.ActivityStatus = false
		activityContext.Message = "Connection to server failed! Check connectivity."
		activityContext.ErrorState = activityError
	} else {
		logger.Log_ACT("URL Request was successfull.", logger.Debug, FlowData["InSessionID"].(string))
		activityContext.ActivityStatus = true
		activityContext.Message = responseMessage
		activityContext.ErrorState = activityError
	}

	FlowData["Response"] = responseMessage

	return FlowData, activityContext
}

func SendRequest(Method string, url string, BodyString string, headerTokens map[string]string, sessionId string, excludeResponseCodes []string, retryCount int) (err error, response string) {
	var req *http.Request
	var resp *http.Response
	req, err = http.NewRequest(Method, url, bytes.NewBuffer([]byte(BodyString)))
	if err != nil {
		response = err.Error()
	} else {
		var body []byte

		if len(headerTokens) > 0 {
			for key, value := range headerTokens {
				req.Header.Set(key, value)
			}
		}

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			if retryCount > 0 && !strings.Contains(err.Error(), "permission") && (strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "EOF")) {
				logger.Log_ACT(("Error : " + err.Error() + " : Retrying Recursively Attempt : " + strconv.Itoa(retryCount)), logger.Debug, sessionId)
				err = nil
				time.Sleep(1 * time.Second)
				return SendRequest(Method, url, BodyString, headerTokens, sessionId, excludeResponseCodes, (retryCount - 1))
			} else {
				response = "Error : " + err.Error()
			}
		} else {
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				response = "Error : " + err.Error()
			}
		}

		excludeResponses := ""

		if len(excludeResponseCodes) > 0 {
			for _, value := range excludeResponseCodes {
				excludeResponses += value + " "
			}
		}

		if resp.StatusCode != 200 && !strings.Contains(excludeResponses, strconv.Itoa(resp.StatusCode)) {
			if len(body) < 4 {
				err = errors.New("Error! Empty Response Body!")
			} else {
				err = errors.New(string(body))

			}
		} else {
			response = string(body)
		}

		defer resp.Body.Close()
	}
	return
}
