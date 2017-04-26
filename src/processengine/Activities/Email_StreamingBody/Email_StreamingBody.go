package Email_StreamingBody

import "processengine/context"
import "bytes"
import "processengine/logger"
import "log"
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

	Host := FlowData["Host"]
	Port := FlowData["Port"]
	Recipient := FlowData["Recipient"]
	Sender := FlowData["Sender"]
	MessageBody := FlowData["MessageBody"]

	// Set up authentication information.
	// Connect to the remote SMTP server.
	c, err := smtp.Dial(Host.(string) + ":" + Port.(string))
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	// Set the sender and recipient.
	c.Mail(Sender.(string))
	c.Rcpt(Recipient.(string))
	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	defer wc.Close()
	buf := bytes.NewBufferString(MessageBody.(string))
	_, err = buf.WriteTo(wc)
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
