package Email_AuthenticatedSMTP

import "processengine/context"
import "processengine/logger"
import "net/smtp"

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

	// to read more about email please read the following http://golang.org/pkg/net/smtp/

	Username := FlowData["Username"]
	Password := FlowData["Password"]
	Host := FlowData["Host"]
	Port := FlowData["Port"]
	Recipient := FlowData["Recipient"]
	Sender := FlowData["Sender"]
	MessageBody := FlowData["MessageBody"]

	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		Username.(string),
		Password.(string),
		Host.(string),
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		Host.(string)+":"+Port.(string),
		auth,
		Sender.(string),
		[]string{Recipient.(string)},
		[]byte(MessageBody.(string)),
	)
	if err != nil {
		FlowData["custMsg"] = "Ooops, There was an error!"
		logger.Log_ACT("Ooops, There was an error: "+err.Error(), logger.Debug, FlowData["InSessionID"].(string))
	} else {
		FlowData["custMsg"] = "Email successfully sent!"
		logger.Log_ACT("Email successfully sent!", logger.Debug, FlowData["InSessionID"].(string))
	}

	//setting activityContext property values
	activityContext.ActivityStatus = true
	activityContext.Message = FlowData["custMsg"].(string)
	activityContext.ErrorState = activityError

	return FlowData, activityContext
}
