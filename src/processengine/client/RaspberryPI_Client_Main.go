// +build ignore

package main

//Import.Start
import "os/exec"
import "strings"
import "duov6.com/agentCore"
import "processengine/logger"
import "duov6.com/cebadapter"
import "processengine/Common" //ImportCommon
import "fmt"

import "os"
import "runtime"                //ImportFMT
import "encoding/json"          //ImportJSON
import "github.com/fatih/color" //ImportCOLOR
import "net/http"               //ImportNETHTTP
import "duov6.com/gorest"       //ImportOGREST
//import "processengine/ceb/cebadapter" //ImportCEB-Modified

import "processengine/context"                                   //ImportCONTEXT
import "processengine/client/Activities/RaspberryPI_Client_Auth" //ImportRaspberryPI_Client_Auth
import "strconv"                                                 //ImportStringConvert
import "time"                                                    //ImportTime
//Import.End

type ExecutableServiceCall string

type HelloService string

var email string

func (T SmoothFlowService) InvokeMethod(data ExecutableServiceCall) string {

	var WFContext context.WorkflowContext //InitVariableCONTEXT
	WFContext.Message = "Startng Workflow"
	WFContext.Status = true
	WFContext.ErrorCode = 1             //InitVariablesOfWFContext
	var flowData map[string]interface{} //InitVariableFLOWDATA
	jsonInput := data                   //InitVariableJSONINPUT

	InSessionID := "ignore"

	if jsonParseErr := json.Unmarshal([]byte(jsonInput), &flowData); jsonParseErr != nil {
		logger.Log_WF("The JSON input is not in correct format.", logger.Error, InSessionID)
		logger.Log_WF("Application terminated.", logger.Error, InSessionID)
		logger.Log_WF(string(jsonParseErr.Error()), logger.Error, InSessionID)
		WFContext.Message = "The JSON input is not in correct format."
		WFContext.Status = false
		WFContext.ErrorCode = 4
		var returnObj context.ReturnData
		returnObj.JSONOutput = flowData
		returnObj.WorkflowResult.Message = WFContext.Message
		returnObj.WorkflowResult.Status = WFContext.Status
		returnObj.WorkflowResult.ErrorCode = WFContext.ErrorCode
		returnDataJSON, rtDError := json.Marshal(returnObj)
		if rtDError != nil {
			logger.Log_WF("WF return data marshal error", logger.Error, InSessionID)
		}
		//fmt.Print(string(returnDataJSON))
		return string(returnDataJSON)
	} //InitVariableJsonParseInput
	InSessionID = flowData["InSessionID"].(string)

	ALLInArguments := make(map[string]interface{})
	ALLInArguments["InSessionID"] = ""
	ALLInArguments["InSecurityToken"] = ""
	ALLInArguments["InLog"] = ""
	ALLInArguments["InNamespace"] = ""

	var count int = 0
	for key, val := range ALLInArguments {
		if _, ok := flowData[key]; ok {
			logger.Log_WF("OK - "+string(key)+" - "+flowData[key].(string), logger.Debug, InSessionID)
		} else {
			if val == "" {
				logger.Log_WF("NO - "+string(key), logger.Debug, InSessionID)
				count = count + 1
			}
		}
	}
	if count > 0 {
		logger.Log_WF("All INArguments are not received", logger.Debug, InSessionID)
		logger.Log_WF("Application terminated.", logger.Debug, InSessionID)
		WFContext.Message = "All INArguments are not received"
		WFContext.Status = false
		WFContext.ErrorCode = 5
		var returnObj context.ReturnData
		returnObj.JSONOutput = flowData
		returnObj.WorkflowResult.Message = WFContext.Message
		returnObj.WorkflowResult.Status = WFContext.Status
		returnObj.WorkflowResult.ErrorCode = WFContext.ErrorCode
		returnDataJSON, rtDError := json.Marshal(returnObj)
		if rtDError != nil {
			logger.Log_WF("WF return data marshal error", logger.Error, InSessionID)
		}
		fmt.Print(string(returnDataJSON))
		return string(returnDataJSON)
	} //AllINArguments
	ALLOutArguments := make(map[string]interface{}) //AllOUTArguments
	//initiating InArgument arguments
	flowData["InSessionID"] = flowData["InSessionID"]
	flowData["InSecurityToken"] = flowData["InSecurityToken"]
	flowData["InLog"] = flowData["InLog"]
	flowData["InNamespace"] = flowData["InNamespace"]
	//INIT.INArgments
	//initiating Custom arguments

	//INIT.CustArguments
	//initiating OutArgument arguments
	ALLOutArguments["OutStatus"] = false
	//INIT.OutArguments
	//Init.Var

	logger.Log_WF("WORKFLOW STARTED!", logger.Information, InSessionID)

	//Func.Start

	if _, ok := flowData["ExecutionLevel"]; ok {
		logger.Log_WF("ExecutionLevel Value already available.", logger.Debug, InSessionID)
	} else {
		flowData["ExecutionLevel"] = "0"
		logger.Log_WF("ExecutionLevel value changed to '0'", logger.Debug, InSessionID)
	}
	if flowData["ExecutionLevel"] == "0" {
		logger.Log_WF("Invoking activity: RaspberryPI_Client_Auth", logger.Debug, InSessionID)
		logger.Log_WF("Session ID: "+InSessionID, logger.Debug, InSessionID)

		ActivityData_LgaTJSeEvlgImOnr := make(map[string]interface{})
		ActivityData_LgaTJSeEvlgImOnr["InSessionID"] = InSessionID

		LgaTJSeEvlgImOnr := make(chan bool)
		var cxOHRhyOewnPSnph = new(context.ActivityContext)
		go func() {

			//fmt.Println("Starting Activity: RaspberryPI_Client_Auth")
			kRQDgGbRyqZSJaje := time.Now()
			ActivityData_LgaTJSeEvlgImOnr, cxOHRhyOewnPSnph = RaspberryPI_Client_Auth.Invoke(ActivityData_LgaTJSeEvlgImOnr)
			LgaTJSeEvlgImOnr <- cxOHRhyOewnPSnph.ActivityStatus
			fmt.Println("Completed Activity in: ", time.Since(kRQDgGbRyqZSJaje))
			fmt.Println("")

			/* out arguments for this activity */
			if ActivityData_LgaTJSeEvlgImOnr["LoginAccess"] != "" && ActivityData_LgaTJSeEvlgImOnr["LoginAccess"] != nil {
				ALLOutArguments["OutStatus"] = ActivityData_LgaTJSeEvlgImOnr["LoginAccess"]
				flowData["OutStatus"] = ActivityData_LgaTJSeEvlgImOnr["LoginAccess"]
				logger.Log_WF("OUT - OutStatus: "+strconv.FormatBool(flowData["OutStatus"].(bool)), logger.Debug, InSessionID)
			}

			if ActivityData_LgaTJSeEvlgImOnr["Email"] != "" && ActivityData_LgaTJSeEvlgImOnr["Email"] != nil {
				ALLOutArguments["Email"] = ActivityData_LgaTJSeEvlgImOnr["Email"]
				flowData["Email"] = ActivityData_LgaTJSeEvlgImOnr["Email"]
				logger.Log_WF("OUT - Email: "+(flowData["Email"].(string)), logger.Debug, InSessionID)
			}

			if ActivityData_LgaTJSeEvlgImOnr["Password"] != "" && ActivityData_LgaTJSeEvlgImOnr["Password"] != nil {
				ALLOutArguments["Password"] = ActivityData_LgaTJSeEvlgImOnr["Password"]
				flowData["Password"] = ActivityData_LgaTJSeEvlgImOnr["Password"]
				logger.Log_WF("OUT - Password: "+(flowData["Password"].(string)), logger.Debug, InSessionID)
			}

		}()

		logger.Log_WF("Activity invoked and waiting for completion ... ", logger.Debug, InSessionID)
		<-LgaTJSeEvlgImOnr

		logger.Log_WF("Activity response was recieved.", logger.Debug, InSessionID)

		if cxOHRhyOewnPSnph.ActivityStatus == true {
			logger.Log_WF("Success - "+cxOHRhyOewnPSnph.Message, logger.Debug, InSessionID)
			WFContext.Message = cxOHRhyOewnPSnph.Message
			WFContext.Status = true
			WFContext.ErrorCode = 3
		} else {
			//logger.Log_WF("Failed - "+OfUiVrIQAzqXQMov.ErrorState.ErrorString, logger.Debug, InSessionID)
			logger.Log_WF("Failed - "+cxOHRhyOewnPSnph.Message, logger.Error, InSessionID)
			WFContext.Message = cxOHRhyOewnPSnph.ErrorState.ErrorString
			WFContext.Status = false
			WFContext.ErrorCode = 2
		}

		//initiate.View.drawboard.0
	} //INIT.executionLevels
	//initiate.View.drawboard

	logger.Log_WF("WORKFLOW COMPLETED!", logger.Error, InSessionID) //AddWFCompleted
	var returnObj context.ReturnData
	returnObj.JSONOutput = flowData
	returnObj.WFOutArguments = ALLOutArguments
	returnObj.WorkflowResult.Message = WFContext.Message
	returnObj.WorkflowResult.Status = WFContext.Status
	returnObj.WorkflowResult.ErrorCode = WFContext.ErrorCode
	returnDataJSON, rtDError := json.Marshal(returnObj)
	if rtDError != nil {
		logger.Log_WF("WF return data marshal error", logger.Information, InSessionID)
	} //ConvertWFResultToJSON
	return string(returnDataJSON) //PRINTJSONDATA
	//Func.End
}

func (T SmoothFlowService) Invoke(data ExecutableServiceCall) {
	inputMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &inputMap)
	if err != nil {
		T.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte("Input Data Invalid : " + err.Error()))
	}

	result := T.InvokeMethod(data)

	flowOutData := make(map[string]interface{})
	err = json.Unmarshal([]byte(result), &flowOutData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON output! : " + err.Error())
	} else {
		WFArguements := flowOutData["WFOutArguments"].(map[string]interface{})

		if WFArguements["Email"] != nil {
			email = WFArguements["Email"].(string)
		}
	}
}

type SmoothFlowService struct {
	gorest.RestService `root:"/smoothflow/" consumes:"application/json" produces:"application/json"`
}

func main() {
	Splash()
	fmt.Println("Starting the SmoothFlow Client.")
	fmt.Println()
	logger.InitiateLogger()
	Common.VerifyCEBAgentConfig()

	t := SmoothFlowService{}
	t.Invoke("{\"InSessionID\":\"ignore\",\"InSecurityToken\":\"7b57e320a5b84c8a404918910edd0975\",\"InLog\":\"log\",\"InNamespace\":\"smoothflow.io\"}")

	if email != "" {
		ConnectToCEB(email)
		VerifyPieDependencies()
	} else {
		return
	}

	gorest.RegisterService(new(SmoothFlowService))

	port := "8894" //INIT.PORTNUMBER
	//Init.Settings

	portNumber := ":" + port

	http.ListenAndServe(portNumber, gorest.Handle())
}

func ConnectToCEB(email string) {
	forever := make(chan bool)
	cebadapter.Attach(("client_" + email), func(s bool) {

		fmt.Println("Successfully connected.")
		agentCore.GetInstance().Client.OnCommand("InstallWorkflow", InstallWorkflow)
		forever <- false
		fmt.Println("Successfully registered in CEB")
	})
	<-forever
}

func InstallWorkflow(from string, name string, data map[string]interface{}, resources map[string]interface{}) {
	fmt.Println("Deploying Standalone Worlflows!")

	if len(data) > 0 {
		for _, value := range data["urls"].(map[string]interface{}) {

			if strings.ToLower(runtime.GOOS) == "linux" {
				URL := value.(string)
				URL = strings.TrimSpace(URL)
				tokens := strings.Split(URL, "/")
				workFlowName := tokens[(len(tokens) - 1)]

				if strings.Contains(workFlowName, ".exe") {
					color.Red("Not a linux executable file. Please download a compatible Workflow version.")
					return
				}

				fmt.Println("Deploying Workflow : " + workFlowName)

				wgetCommand := "wget " + URL
				nohupCommand := "nohup ./" + workFlowName + " > " + workFlowName + ".log 2>&1 &"

				_, _ = exec.Command("sh", "-c", ("rm " + workFlowName)).Output()
				_, _ = exec.Command("sh", "-c", wgetCommand).Output()
				_, _ = exec.Command("sh", "-c", "chmod 777 *").Output()
				_, _ = exec.Command("sh", "-c", nohupCommand).Output()
			} else if strings.ToLower(runtime.GOOS) == "windows" {

				URL := value.(string)
				URL = strings.TrimSpace(URL)
				tokens := strings.Split(URL, "/")
				workFlowName := tokens[(len(tokens) - 1)]

				if !strings.Contains(workFlowName, ".exe") {
					color.Red("Not a Windows EXE file. Please download a compatible Workflow version.")
					return
				}

				fmt.Println("Deploying Workflow : " + workFlowName)

				currentWorkingDirectory, _ := os.Getwd()

				currentWorkingDirectory = strings.Replace(currentWorkingDirectory, "/", "\\", -1)
				currentWorkingDirectory += "\\"

				wgetCommand := "bitsadmin /transfer myDownloadJob /download /priority normal " + URL + " " + currentWorkingDirectory + workFlowName

				fmt.Println("Downloaded File Saved at : " + currentWorkingDirectory + workFlowName)

				nohupCommand := "start " + currentWorkingDirectory + workFlowName

				_, _ = exec.Command("cmd", "/C", ("del " + workFlowName)).Output()
				_, _ = exec.Command("cmd", "/C", wgetCommand).Output()
				_, _ = exec.Command("cmd", "/C", nohupCommand).Output()
			} else {
				color.Red("Unsupported Operating system : Failed executing downloaded executable!")
			}
		}

		fmt.Println("Completed Deploying!")
		fmt.Println()
	} else {
		fmt.Println("No Workflows to be Deployed!")
	}

}

func Splash() {
	logColor := color.New(color.FgGreen)
	fmt.Println()
	fmt.Println()
	logColor.Println(" _____                       _   _    ______ _               ")
	logColor.Println("/  ___|                     | | | |   |  ___| |              ")
	logColor.Println("\\ `--. _ __ ___   ___   ___ | |_| |__ | |_  | | _____      __")
	logColor.Println(" `--. \\ '_ ` _ \\ / _ \\ / _ \\| __| '_ \\|  _| | |/ _ \\ \\ /\\ / /")
	logColor.Println("/\\__/ / | | | | | (_) | (_) | |_| | | | |   | | (_) \\ V  V / ")
	logColor.Println("\\____/|_| |_| |_|\\___/ \\___/ \\__|_| |_\\_|   |_|\\___/ \\_/\\_/  ")
	logColor.Println()
	fmt.Println()
	logColor.Println("Welcome to smoothflow.io")
	fmt.Println()
}

func VerifyPieDependencies() {
	color.Yellow("Verifying status of GPIO compatibility with your Operating System........")
	if strings.ToLower(runtime.GOOS) == "linux" && strings.ToLower(runtime.GOARCH) == "arm" {
		fmt.Println("This might take a few minutes depending on internet connection speed.")
		Common.VerifyDependencies()
	} else {
		color.Yellow("Warning : Operating system not supported for GPIO usage. Activities with GPIO usage won't be executed!")
	}
}
