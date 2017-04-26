package RaspberryPie_Camera_CapturePhoto

import "processengine/context"
import "errors"
import "duov6.com/common"
import "os/exec"
import "io/ioutil"
import "processengine/logger"
import "processengine/Common"
import "time"
import "encoding/json"
import "strings"
import "os"

type InsParameters struct {
	KeyProperty string `json:"KeyProperty"`
}
type InsertTemplate struct {
	Object     map[string]interface{}   `json:"Object"`
	Objects    []map[string]interface{} `json:"Objects"`
	Parameters InsParameters
}

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

	pythonFileName := common.GetGUID() + ".py"

	if Common.VerifyGPIOCapability() && Common.VerifyDependencies() {

		_, errr := os.Stat("MyImages")
		if errr != nil {
			os.Mkdir("MyImages", 0777)
		}

		imageName := "Image_" + time.Now().UTC().Format("2006-01-02T15:04:05") + ".jpg"
		fileContent := "import picamera\ncamera = picamera.PiCamera()\ncamera.capture('MyImages/" + imageName + "')"
		_ = ioutil.WriteFile(pythonFileName, []byte(fileContent), 0666)

		_, err = exec.Command("sh", "-c", ("python " + pythonFileName)).Output()

		_ = os.Remove(pythonFileName)

		if FlowData["ObjectStoreURL"] != nil && FlowData["Tenant"] != nil {

			file2, err2 := ioutil.ReadFile("MyImages/" + imageName)
			if err2 == nil {
				base64Body := common.EncodeToBase64(string(file2))

				obj := make(map[string]interface{})
				obj["Id"] = imageName
				obj["FileName"] = imageName
				obj["Body"] = base64Body

				url := FlowData["ObjectStoreURL"].(string)
				url = strings.TrimSpace(url)
				url = strings.Replace(url, "http://", "", -1)
				url = strings.Replace(url, "https://", "", -1)
				url = strings.Replace(url, "/", "", -1)
				url = strings.Replace(url, ":3000", "", -1)
				url = "http://" + url + ":3000/" + FlowData["Tenant"].(string) + "/MyImages?securityToken=ignore"

				paramObject := InsParameters{}
				paramObject.KeyProperty = "Id"

				insertObject := InsertTemplate{}
				insertObject.Object = obj
				insertObject.Parameters = paramObject

				byteArray, _ := json.Marshal(insertObject)

				common.HTTP_POST(url, nil, byteArray, false)
			}
		}

	} else {
		err = errors.New("GPIO Dependencies not met. Check for Operating system and Architecture.")
	}

	if err != nil {
		//setting activityContext property values
		activityContext.ActivityStatus = false
		activityContext.Message = "Rasperry Pi Capture Image Failed : " + err.Error()
		logger.Log_ACT(activityContext.Message, logger.Debug, FlowData["InSessionID"].(string))
		activityContext.ErrorState = activityError
	} else {
		//setting activityContext property values
		activityContext.ActivityStatus = true
		activityContext.Message = "Rasperry Pi Capture Image Completed Successfully!"
		logger.Log_ACT(activityContext.Message, logger.Debug, FlowData["InSessionID"].(string))
	}

	return FlowData, activityContext
}
