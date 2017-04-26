package Json_Converter

import "processengine/context"
import "processengine/logger"
import "encoding/json"
import "errors"
import "strings"
import "fmt"

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

	sessionID := FlowData["InSessionID"].(string)

	var err error
	var JSON string

	if FlowData["MapData"] != nil && FlowData["ListData"] != nil && FlowData["ValueData"] != nil {
		err = errors.New("Error : No values, lists or map inputs recieved.")
	} else {
		if FlowData["ValueData"] != nil {
			JSON, err = ConvertValueDataToJsonString(FlowData["ValueData"].(string))
		} else if FlowData["ListData"] != nil {
			JSON, err = ConvertListDataToJsonString(FlowData["ListData"])
		} else if FlowData["MapData"] != nil {
			JSON, err = ConvertMapDataToJsonString(FlowData["MapData"])
		}
	}

	fmt.Println("--------------------------------------------")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(JSON)
	}
	fmt.Println("--------------------------------------------")
	//When a set of values are given as input, will convert it to a JSON list
	//When a list variable is given, will convert to the corresponding JSON list
	//When a Dictionary is given, should convert it into the corresponding JSON string.

	if err == nil {
		msg := "Successfully Converted!"
		activityContext.ActivityStatus = true
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, sessionID)
		FlowData["custMsg"] = msg
		FlowData["status"] = "true"
	} else {
		msg := "Error when Converting : " + err.Error()
		activityContext.ActivityStatus = false
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, sessionID)
		FlowData["custMsg"] = msg
		FlowData["status"] = "false"
	}

	return FlowData, activityContext
}

func ConvertValueDataToJsonString(value string) (jsonString string, err error) {
	var byteArray []byte
	value = strings.Replace(value, "{", "", -1)
	value = strings.Replace(value, "}", "", -1)
	value = strings.Replace(value, "[", "", -1)
	value = strings.Replace(value, "]", "", -1)
	value = strings.Replace(value, "(", "", -1)
	value = strings.Replace(value, ")", "", -1)
	tokens := strings.Split(value, ",")

	for x := 0; x < len(tokens); x++ {
		tokens[x] = strings.TrimSpace(tokens[x])
	}

	jsonMap := make(map[string]interface{})
	jsonMap["Value"] = tokens

	byteArray, err = json.Marshal(jsonMap)
	jsonString = string(byteArray)
	return
}

func ConvertListDataToJsonString(input interface{}) (jsonString string, err error) {
	var byteArray []byte
	jsonMap := make(map[string]interface{})
	jsonMap["Value"] = input
	byteArray, err = json.Marshal(jsonMap)
	jsonString = string(byteArray)
	return
}

func ConvertMapDataToJsonString(value interface{}) (jsonString string, err error) {
	var byteArray []byte
	byteArray, err = json.Marshal(value)
	jsonString = string(byteArray)
	return
}
