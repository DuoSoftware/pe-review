package SMS_GoNexmo

import "processengine/context"
import "github.com/njern/gonexmo"
import "strconv"
import "time"
import "processengine/logger"

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

	// to read more about email please read the following https://godoc.org/gopkg.in/gomail.v2-unstable

	APIKey := FlowData["APIKey"]
	APIPassword := FlowData["APIPassword"]
	From := FlowData["From"]
	To := FlowData["To"]
	TextMessage := FlowData["TextMessage"]

	nexmoClient, _ := nexmo.NewClientFromAPI(APIKey.(string), APIPassword.(string))

	// Test if it works by retrieving your account balance
	balance, _ := nexmoClient.Account.GetBalance()

	// Send an SMS
	// See https://docs.nexmo.com/index.php/sms-api/send-message for details.
	message := &nexmo.SMSMessage{
		From:            From.(string),
		To:              To.(string),
		Type:            nexmo.Text,
		Text:            TextMessage.(string),
		ClientReference: "Sent from DuoWorld " + strconv.FormatInt(time.Now().Unix(), 10),
		Class:           nexmo.Standard,
	}

	_, err := nexmoClient.SMS.Send(message)
	if err != nil {
		FlowData["custMsg"] = "Ooops, There was an error!"
		logger.Log_ACT("Ooops, There was an error: ", logger.Debug, FlowData["InSessionID"].(string))
	} else {
		FlowData["custMsg"] = "SMS successfully sent!"
		logger.Log_ACT("SMS successfully sent!", logger.Debug, FlowData["InSessionID"].(string))
	}
	s := strconv.FormatFloat(balance, 'E', -1, 64)
	FlowData["Balance"] = s

	//setting activityContext property values
	activityContext.ActivityStatus = true
	activityContext.Message = FlowData["custMsg"].(string)
	activityContext.ErrorState = activityError

	return FlowData, activityContext
}
