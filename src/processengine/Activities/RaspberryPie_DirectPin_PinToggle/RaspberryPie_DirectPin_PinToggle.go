package RaspberryPie_DirectPin_PinToggle

import "processengine/context"
import "github.com/stianeikeland/go-rpio"
import "errors"
import "processengine/logger"
import "processengine/Common"
import "strconv"

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

	var err error
	if Common.VerifyGPIOCapability() && Common.VerifyDependencies() {
		pinNumber := 0

		if FlowData["PinNumber"] != nil {
			pinNumber, err = strconv.Atoi(FlowData["PinNumber"].(string))
			if err == nil {
				err = rpio.Open()
				if err == nil {
					pin := rpio.Pin(pinNumber)
					pin.Output()
					pin.Toggle()
					rpio.Close()
				}
			}
		} else {
			err = errors.New("Pin Number not found!")
		}
	} else {
		err = errors.New("GPIO Dependencies not met. Check for Operating system and Architecture.")
	}
	if err != nil {
		//setting activityContext property values
		activityContext.ActivityStatus = false
		activityContext.Message = "Rasperry Pi Pin Toggle Failed : " + err.Error()
		activityContext.ErrorState = activityError
		FlowData["Response"] = "Rasperry Pi Pin Toggle Failed : " + err.Error()
		logger.Log_ACT(activityContext.Message, logger.Debug, FlowData["InSessionID"].(string))
	} else {
		//setting activityContext property values
		activityContext.ActivityStatus = true
		activityContext.Message = "Rasperry Pi Pin Toggle completed Successfully!"
		FlowData["Response"] = "Rasperry Pi Pin Toggle completed Successfully!"
		logger.Log_ACT(activityContext.Message, logger.Debug, FlowData["InSessionID"].(string))
	}

	return FlowData, activityContext
}
