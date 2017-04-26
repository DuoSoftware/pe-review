package RaspberryPI_Client_Auth

import (
	"duov6.com/common"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	httpService "processengine/Activities/HTTP_DefaultRequest"
	"processengine/Common"
	"processengine/context"
	"strings"
)

var InSessionId string
var host string

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

	InSessionId = FlowData["InSessionID"].(string)
	var email string
	var password string
	var deviceID string
	_ = deviceID

	host = "auth.smoothflow.io:3048"

	content, err1 := ioutil.ReadFile("agent.config")
	if err1 == nil {
		config := make(map[string]interface{})
		_ = json.Unmarshal(content, &config)
		host = config["authUrl"].(string)
	}

	var err error

	fmt.Println("Welcome to Duo SmoothFlow Automated Client!")
	fmt.Println()

	email, password, err = EvaluateSettingsAndClearOnDemand()

	if err == nil {
		msg := "Successfully Completed Authentication!"
		logColor := color.New(color.FgGreen)
		logColor.Println(msg)
		activityContext.ActivityStatus = true
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		Common.LogACT(msg, InSessionId)
		FlowData["LoginAccess"] = true
		FlowData["Email"] = email
		FlowData["Password"] = password
		if strings.Contains(email, "@") {
			FlowData["UserType"] = "User"
		} else {
			FlowData["UserType"] = "Token"
		}
	} else {
		msg := "Error Authentication : " + err.Error()
		logColor := color.New(color.FgRed)
		logColor.Println(msg)
		activityContext.ActivityStatus = false
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		Common.LogACT(msg, InSessionId)
		FlowData["LoginAccess"] = false
	}

	return FlowData, activityContext
}

func GetSettings() (object map[string]interface{}) {
	content, err := ioutil.ReadFile("settings.config")
	object = make(map[string]interface{})
	if err == nil {
		_ = json.Unmarshal(content, &object)
	}
	return
}

func SaveSettings(object map[string]interface{}) {
	byteArray, _ := json.Marshal(object)
	_ = ioutil.WriteFile("settings.config", byteArray, 0666)
}

func Authenticate(host, email, password, InSessionId string) (err error) {
	//Authenticate
	restAuthData := make(map[string]interface{})

	restAuthData["URL"] = "http://" + host + "/Login/" + email + "/" + password + "/smoothflow.io"
	restAuthData["Method"] = "GET"
	restAuthData["Body"] = ""
	restAuthData["InSessionID"] = InSessionId

	httpFlowResult, httpActivityResult := httpService.Invoke(restAuthData)

	if !httpActivityResult.ActivityStatus {
		//Auth Error
		authStruct := make(map[string]interface{})
		json.Unmarshal([]byte(httpActivityResult.ErrorState.ErrorString), &authStruct)
		Common.LogACT(("Authentication Error : " + authStruct["Message"].(string)), InSessionId)
		err = errors.New("Authentication Error : " + authStruct["Message"].(string))
		//Clear settings
		settings := make(map[string]interface{})
		SaveSettings(settings)
	} else {
		//Auth Successfull
		settings := make(map[string]interface{})
		settings["Email"] = email
		settings["Password"] = password
		SaveSettings(settings)
		fmt.Println()
		authStruct := make(map[string]interface{})
		fmt.Println("This Device is Registerd to : " + email)
		json.Unmarshal([]byte(httpFlowResult["Response"].(string)), &authStruct)
		Common.LogACT(("Authentication Successful! Your Security Token : " + authStruct["SecurityToken"].(string)), InSessionId)
	}

	return
}

func EvaluateSettingsAndClearOnDemand() (email, password string, err error) {
	settings := GetSettings()

	var yesOrNo string
	if len(settings) > 0 {
		color.Yellow("Existing Login Information found. Proceed with Existing Login Information?")
		fmt.Println("Press Y to continue. Press N to login as a new user.")

		_, _ = fmt.Scanln(&yesOrNo)
		fmt.Println()
		yesOrNo = strings.ToLower(yesOrNo)

		if yesOrNo != "y" && yesOrNo != "yes" && yesOrNo != "n" && yesOrNo != "no" {
			color.Red("Invalid Response!")
			os.Exit(0)
		}

		if strings.Contains(yesOrNo, "y") {
			fmt.Println("Logging in with existing information...")
			//authenticate
			if settings["DeviceID"] != nil {
				email = settings["DeviceID"].(string)
				fmt.Println("This Device is Registered as : " + email)
				//how to authenticate?
			} else {
				email = settings["Email"].(string)
				password = settings["Password"].(string)
				err = Authenticate(host, email, password, InSessionId)
			}
		} else {
			settings := make(map[string]interface{})
			SaveSettings(settings)
			color.Yellow("Old Login info deleted... Please provide Login Option.")
			email, password = EvaluateOptionsAndGetAuthData()
			if strings.Contains(email, "@") {
				err = Authenticate(host, email, password, InSessionId)
			} else {
				fmt.Println("This Device is Registered as : " + email)
				//how to authenticate?
			}
		}
	} else {
		//Get login credentials from user
		email, password = EvaluateOptionsAndGetAuthData()
		if strings.Contains(email, "@") {
			err = Authenticate(host, email, password, InSessionId)
		} else {
			fmt.Println("This Device is Registered as : " + email)
			//how to authenticate?
		}
	}
	return
}

func EvaluateOptionsAndGetAuthData() (email, password string) {
	var yesOrNo string
	fmt.Println("Press 1 to Login as SmoothFlow User.")
	fmt.Println("Press 2 to Login as SmoothFlow Token.")
	_, _ = fmt.Scanln(&yesOrNo)

	if yesOrNo != "1" && yesOrNo != "2" {
		color.Red("Invalid Response!")
		os.Exit(0)
	}

	if strings.Contains(yesOrNo, "1") {
		//Get User prompt username and password
		fmt.Println()
		fmt.Println("Please Enter Your SmoothFlow Account Email and press ENTER")
		_, _ = fmt.Scanln(&email)

		if !strings.Contains(email, "@") {
			color.Red("Invalid Email!")
			os.Exit(0)
		}

		fmt.Println()
		fmt.Println("Please Enter Your SmoothFlow Account Password and press ENTER")
		_, _ = fmt.Scanln(&password)
		fmt.Println()
	} else {
		//get GUID for device
		email = GetSerialNumber()
		settings := make(map[string]interface{})
		settings["DeviceID"] = email
		SaveSettings(settings)
	}
	return
}

func GetSerialNumber() (serial string) {
	serial = strings.ToUpper(common.GetGUID())
	serial = strings.Replace(serial, "-", "", -1)
	serial = serial[:20]
	serial = serial[0:4] + "-" + serial[15:19]
	return
}

// test commit
