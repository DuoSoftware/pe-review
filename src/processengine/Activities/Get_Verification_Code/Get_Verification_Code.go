package Get_Verification_Code

import "processengine/context"
import "processengine/logger"
import "math/rand"
import "time"

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

	count := FlowData["VerifyCodeLength"].(int)

	if count <= 1 {
		count = 4
	}

	verifyCode := RandomString(count)
	FlowData["VerifyCode"] = verifyCode

	msg := "Verification Code Generated Successfully!"
	activityContext.ActivityStatus = true
	activityContext.Message = msg
	activityContext.ErrorState = activityError
	logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
	FlowData["custMsg"] = msg
	FlowData["status"] = "true"

	return FlowData, activityContext
}

func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
