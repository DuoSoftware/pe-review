package Get_GUID

import "processengine/context"
import "processengine/logger"
import "runtime"
import "github.com/twinj/uuid"
import "os/exec"
import "crypto/md5"
import "encoding/hex"

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

	GUID := ""

	if runtime.GOOS == "linux" {
		out, err := exec.Command("uuidgen").Output()
		h := md5.New()
		h.Write(out)
		if err == nil {
			GUID = hex.EncodeToString(h.Sum(nil))
		} else {
			h := md5.New()
			h.Write([]byte(uuid.NewV1().String()))
			GUID = hex.EncodeToString(h.Sum(nil))
		}
	} else {
		h := md5.New()
		h.Write([]byte(uuid.NewV1().String()))
		GUID = hex.EncodeToString(h.Sum(nil))
	}

	FlowData["Generated_GUID"] = GUID

	msg := "GUID Generated Successfully!"
	activityContext.ActivityStatus = true
	activityContext.Message = msg
	activityContext.ErrorState = activityError
	logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
	FlowData["custMsg"] = msg
	FlowData["status"] = "true"

	return FlowData, activityContext
}
