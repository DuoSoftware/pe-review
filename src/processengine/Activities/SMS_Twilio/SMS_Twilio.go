package SMS_Twilio

import "processengine/context"
import "github.com/subosito/twilio"
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

	AccountSid := FlowData["AccountSid"].(string)
	AuthToken := FlowData["AuthToken"].(string)
	From := FlowData["From"].(string)
	To := FlowData["To"].(string)
	TextMessage := FlowData["TextMessage"].(string)

	// Initialize twilio client
	c := twilio.NewClient(AccountSid, AuthToken, nil)

	// You can set custom Client, eg: you're using `appengine/urlfetch` on Google's appengine
	// a := appengine.NewContext(r) // r is a *http.Request
	// f := urlfetch.Client(a)
	// c := twilio.NewClient(AccountSid, AuthToken, f)

	// Send Message
	params := twilio.MessageParams{
		Body: TextMessage,
	}
	response, _, err := c.Messages.Send(From, To, params)
	/*// You can also using lower level function: Create
	s, response, err = c.Messages.Create(url.Values{
		"From": {From},
		"To":   {To},
		"Body": {TextMessage},
		})*/
	if err != nil {
		activityContext.ActivityStatus = false
		FlowData["custMsg"] = "Ooops, There was an error!"
		FlowData["Response"] = err.Error()
		logger.Log_ACT("Ooops, There was an error: ", logger.Debug, FlowData["InSessionID"].(string))
	} else {
		activityContext.ActivityStatus = true
		FlowData["custMsg"] = "SMS successfully sent!"
		FlowData["Response"] = response.Body
		logger.Log_ACT("SMS successfully sent!", logger.Debug, FlowData["InSessionID"].(string))
	}

	//setting activityContext property values
	activityContext.Message = FlowData["custMsg"].(string)
	activityContext.ErrorState = activityError

	return FlowData, activityContext
}
