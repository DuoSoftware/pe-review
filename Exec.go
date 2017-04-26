// +build ignore

package main

//Import.Start
import "processengine/logger"
import "processengine/Common"   //ImportCommon
import "fmt"                    //ImportFMT
import "encoding/json"          //ImportJSON
import "github.com/fatih/color" //ImportCOLOR
import "net/http"               //ImportNETHTTP
import "duov6.com/gorest"       //ImportOGREST
import "duov6.com/cebadapter"   //ImportCEB-Modified

//Import.End

type ExecutableServiceCall string

type HelloService string

func (T SmoothFlowService) InvokeMethod(data ExecutableServiceCall, Headers map[string]string) string {

	//Init.Var

	logger.Log_WF("WORKFLOW STARTED!", logger.Information, flowData["InSessionID"].(string))

	//Func.Start

	//initiate.View.drawboard

	//Func.End
}

func (T SmoothFlowService) Invoke(data ExecutableServiceCall) {

	/*securityToken := T.Context.Request().Header.Get("Securitytoken")
	authenticateStruct := make(map[string]interface{})
	if securityToken == "" {
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Security Token Empty"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	if !Common.AuthenticateSecurityToken(securityToken) {
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Invalid Security Token"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}*/

	white := color.New(color.FgWhite)
	redBackground := white.Add(color.BgRed)
	redBackground.Println("********************************************************************************")

	// validating headers
	ALLHeaders := make(map[string]string)
	//INIT.Headers

	var count int = 0
	for key, val := range ALLHeaders {
		if _, ok := ALLHeaders[key]; ok {
			if val != "" {
				fmt.Println("OK - " + string(key) + " - " + ALLHeaders[key])
			} else {
				fmt.Println("NO - " + string(key))
				count = count + 1
			}
		}
	}
	if count > 0 {
		T.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte(`{"Reason":"All headers are not present.","Status":false}`))
		return
	}

	//Rest Request prints starts here..
	logColor := color.New(color.FgYellow)
	inputMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &inputMap)
	if err != nil {
		T.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte(`{"Reason":"Input Data Invalid : ` + err.Error() + `","Status":false}`))
		return
	}
	logColor.Println("In Arguments : {")
	for key, value := range inputMap {
		logColor.Print("		" + key + " : ")
		logColor.Println(value)
	}
	logColor.Println("		}")

	//Rest Request print ends here
	fmt.Println("")
	fmt.Println("Starting the Workflow.")

	result := T.InvokeMethod(data, ALLHeaders)

	fmt.Println("Completed the Workflow.")
	fmt.Println("")

	//DSF-188 Correction Starts here....
	//jsonData, _ := json.Marshal(result)
	//It goes as a string since Marshalling a string variable which contains escape chars.
	//Just casting to byte is the correction.
	jsonData := []byte(result)
	//Correction ends here
	T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(jsonData)
}

func (T SmoothFlowService) Hello() (obj HelloService) {
	obj = `{"Message":"Hello from the otherside!"}`
	return
	//T.ResponseBuilder().SetResponseCode(200).WriteAndOveride([]byte("Hello from the otherside!!!!"))
}

func (T SmoothFlowService) ToggleLogs() string {
	msg := logger.ToggleLogs()
	return msg
}

type SmoothFlowService struct {
	gorest.RestService `root:"/smoothflow/" consumes:"application/json" produces:"application/json"`
	invoke             gorest.EndPoint `method:"POST" path:"/Invoke/" postdata:"ExecutableServiceCall"`
	hello              gorest.EndPoint `method:"GET" path:"/Hello/" output:"HelloService"`
	toggleLogs         gorest.EndPoint `method:"GET" path:"/Logs/" output:"string"`
	getConfig          gorest.EndPoint `method:"GET" path:"/GetConfig/" output:"string"`
	getConfigElement   gorest.EndPoint `method:"GET" path:"/GetConfig/{Key:string}" output:"string"`
	setConfigElement   gorest.EndPoint `method:"GET" path:"/SetConfig/{Key:string}/{Value:string}" output:"string"`
}

func (T SmoothFlowService) GetConfig() string {
	byteValue, _ := json.Marshal(Common.GetConfig())
	T.ResponseBuilder().SetResponseCode(200)
	return string(byteValue)
}

func (T SmoothFlowService) GetConfigElement(Key string) string {
	byteValue, _ := json.Marshal(Common.GetConfigValue(Key))
	T.ResponseBuilder().SetResponseCode(200)
	return string(byteValue)
}

func (T SmoothFlowService) SetConfigElement(Key string, Value string) string {
	config := Common.GetConfig()
	config[Key] = Value
	Common.SaveConfig(config)
	byteValue, _ := json.Marshal(config)
	T.ResponseBuilder().SetResponseCode(200)
	return string(byteValue)
}

func ConnectToCEB() {
	forever := make(chan bool)
	cebadapter.Attach("wf_"+"//INIT.FlowNameCEB", func(s bool) {
		forever <- false
		color.Yellow("Successfully registered in CEB")
	})

	<-forever
}

func main() {
	Common.VerifyCEBAgentConfig()
	logger.InitiateLogger()
	err := Common.CheckConnectionToCEB()
	if err == nil {
		ConnectToCEB()
	} else {
		fmt.Println("Running as a Offline Workflow. CEB service offline.")
	}
	gorest.RegisterService(new(SmoothFlowService))

	//Init.Settings

	portNumber := ":" + port

	fmt.Println()
	logColor := color.New(color.FgYellow)
	logColor.Println("Flow Name : //INIT.FlowName")
	logColor.Println("Port : " + port)
	logColor.Println("Build Platform : //INIT.FlowBuildPlatform")
	logColor.Println("REST Method : POST; URL : http://localhost" + portNumber + "/smoothflow/Invoke/")
	logColor.Println("REST Method :  GET; URL : http://localhost" + portNumber + "/smoothflow/Hello/")
	fmt.Println()
	logColor.Println("//INIT.PrintInArguments")
	logColor.Println("//INIT.PrintOutArguments")
	logColor.Println("//INIT.PrintWfBodySample")
	fmt.Println()

	fmt.Println("Service is running...")

	http.ListenAndServe(portNumber, gorest.Handle())
}
