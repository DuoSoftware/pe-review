package Email_GoMail

import "processengine/context"
import "processengine/logger"
import "gopkg.in/gomail.v2-unstable"

// method can be used to send Emails using the Go implemented mail service
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

	// get and assign values from MAP to internal variables
	Username := FlowData["Username"]
	Password := FlowData["Password"]
	Host := FlowData["Host"]
	Port := FlowData["Port"]
	Recipient := FlowData["Recipient"]
	Sender := FlowData["Sender"]
	SenderName := FlowData["SenderName"]
	Subject := FlowData["Subject"]
	MessageBody := FlowData["MessageBody"]

	// perpare the data structure required for the request
	m := gomail.NewMessage()
	m.SetHeader("From", Sender.(string), SenderName.(string))
	m.SetHeader("To", Recipient.(string))
	m.SetHeader("Subject", Subject.(string))
	m.SetBody("text/html", MessageBody.(string))

	// make the call
	d := gomail.NewPlainDialer(Host.(string), Port.(int), Username.(string), Password.(string))

	// Send the email to Bob, Cora and Dan.
	err := d.DialAndSend(m)
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

	// once completed return the set of data
	return FlowData, activityContext
}
