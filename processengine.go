package main

import (
	"duov6.com/cebadapter"
	"duov6.com/gorest"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"os"
	"os/exec"
	"processengine/Common"
	"processengine/Components"
	"processengine/context"
	"processengine/logger"
	"processengine/plugins/jira"
	"processengine/queue/workers"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type ProcessEngineService struct {
	gorest.RestService   `root:"/processengine/" consumes:"application/json" produces:"application/json"`
	buildFlow            gorest.EndPoint `method:"POST" path:"/BuildFlow/{flowName:string}/{sessionId:string}" postdata:"JsonFlow"`
	invokeFlow           gorest.EndPoint `method:"POST" path:"/InvokeFlow/" postdata:"InvokeStruct"`
	invokeHibernatedFlow gorest.EndPoint `method:"POST" path:"/InvokeHibernatedFlow/{appCode:string}/{processCode:string}/{sessionId:string}" postdata:"InvokeStruct"`
	publishActivity      gorest.EndPoint `method:"POST" path:"/PublishActivity/{activityName:string}/{sessionId:string}" postdata:"ActivityStruct"`
	removeActivity       gorest.EndPoint `method:"POST" path:"/RemoveActivity/{activityName:string}/{sessionId:string}" postdata:"ActivityStruct"`
	getVersionInfo       gorest.EndPoint `method:"GET" path:"/" output:"VersionResponse"`
	getEnv               gorest.EndPoint `method:"GET" path:"/env/" output:"bool"`
	getSessionDetails    gorest.EndPoint `method:"GET" path:"/GetSessionDetails/{sessionID:string}/{sessionType:string}" output:"SessionTranDetails"`
	testWorkflow         gorest.EndPoint `method:"POST" path:"/TestWorkflow/{flowName:string}/{sessionId:string}" postdata:"TestWorkflowInvoke"`
	testActivity         gorest.EndPoint `method:"POST" path:"/TestActivity/{flowName:string}/{sessionId:string}" postdata:"TestWorkflowInvoke"`
	installPackages      gorest.EndPoint `method:"POST" path:"/InstallPackages/{sessionId:string}" postdata:"PackageInstaller"`
	downloadExecutable   gorest.EndPoint `method:"POST" path:"/DownloadExecutable/{flowName:string}/{sessionId:string}" postdata:"JsonFlow"`
	publishToDocker      gorest.EndPoint `method:"POST" path:"/PublishToDocker/" postdata:"DockerDeployement"`
	removeDocker         gorest.EndPoint `method:"POST" path:"/RemoveDocker/" postdata:"DockerDeployement"`
	toggleLogs           gorest.EndPoint `method:"GET" path:"/Logs" output:"string"`
	toggleLogstash       gorest.EndPoint `method:"GET" path:"/Logstash" output:"string"`

	//Jira
	jiraTrigger     gorest.EndPoint `method:"POST" path:"/JiraTrigger/" postdata:"JiraTriggerRequest"`
	authorizeJira   gorest.EndPoint `method:"GET" path:"/AuthorizeJira/{jiraDomain:string}" output:"string"`
	unauthorizeJira gorest.EndPoint `method:"GET" path:"/UnauthorizeJira/{jiraDomain:string}" output:"string"`
}

// gorest documentation and sample code
// https://code.google.com/p/gorest/

//list of work queues for different types of processes
var WFPublishQueue = make(chan Components.JsonFlow, 100)
var WFInvokeQueue = make(chan Components.InvokeStruct, 100)

//list of worker queues
var WFPublishWorkerQueue chan chan Components.JsonFlow
var InvokeWorkerQueue chan chan Components.InvokeStruct

//Jira Methods

func (T ProcessEngineService) AuthorizeJira(jiraDomain string) string {

	studioCrowdToken := T.Context.Request().Header.Get("studio.crowd.tokenkey")
	jSession := T.Context.Request().Header.Get("JSESSIONID")
	xsrfToken := T.Context.Request().Header.Get("atlassian.xsrf.token")
	sessionToken := T.Context.Request().Header.Get("cloud.session.token")

	fmt.Println("Jira Tokens Recieved : ")
	fmt.Println("studio.crowd.tokenkey : " + studioCrowdToken)
	fmt.Println("JSESSIONID : " + jSession)
	fmt.Println("atlassian.xsrf.token : " + xsrfToken)
	fmt.Println("cloud.session.token : " + sessionToken)
	fmt.Println()

	if jSession == "" || studioCrowdToken == "" {
		res := jira.JiraResponse{}
		res.IsSuccess = false
		res.Message = "Jira Auth Tokens Missing."
		byteData, _ := json.Marshal(res)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteData).AddHeader("Content-Type", "application/json")
	}

	jiraTokens := make(map[string]string)
	jiraTokens["studio.crowd.tokenkey"] = studioCrowdToken
	jiraTokens["JSESSIONID"] = jSession
	jiraTokens["atlassian.xsrf.token"] = xsrfToken
	jiraTokens["cloud.session.token"] = sessionToken

	result := jira.JiraAuthorize(jiraDomain, jiraTokens)

	jsonData, _ := json.Marshal(result)

	if result.IsSuccess {
		T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
	} else {
		if strings.Contains(string(jsonData), "Server Error") {
			T.ResponseBuilder().SetResponseCode(500).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
		} else {
			T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
		}
	}

	return ""
}

func (T ProcessEngineService) UnauthorizeJira(jiraDomain string) string {

	studioCrowdToken := T.Context.Request().Header.Get("studio.crowd.tokenkey")
	jSession := T.Context.Request().Header.Get("JSESSIONID")
	xsrfToken := T.Context.Request().Header.Get("atlassian.xsrf.token")
	sessionToken := T.Context.Request().Header.Get("cloud.session.token")

	fmt.Println("Jira Tokens Recieved : ")
	fmt.Println("studio.crowd.tokenkey : " + studioCrowdToken)
	fmt.Println("JSESSIONID : " + jSession)
	fmt.Println("atlassian.xsrf.token : " + xsrfToken)
	fmt.Println("cloud.session.token : " + sessionToken)
	fmt.Println()

	if jSession == "" || studioCrowdToken == "" {
		res := jira.JiraResponse{}
		res.IsSuccess = false
		res.Message = "Jira Auth Tokens Missing."
		byteData, _ := json.Marshal(res)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteData).AddHeader("Content-Type", "application/json")
	}

	jiraTokens := make(map[string]string)
	jiraTokens["studio.crowd.tokenkey"] = studioCrowdToken
	jiraTokens["JSESSIONID"] = jSession
	jiraTokens["atlassian.xsrf.token"] = xsrfToken
	jiraTokens["cloud.session.token"] = sessionToken

	result := jira.JiraUnAuthorize(jiraDomain, jiraTokens)

	jsonData, _ := json.Marshal(result)

	if result.IsSuccess {
		T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
	} else {
		if strings.Contains(string(jsonData), "Server Error") {
			T.ResponseBuilder().SetResponseCode(500).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
		} else {
			T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
		}
	}

	return ""
}

func (T ProcessEngineService) JiraTrigger(wfDetails jira.JiraTriggerRequest) {
	color.Cyan("Jira Trigger Executed.")

	ipTokens := strings.Split(T.Context.Request().RemoteAddr, ":")
	originIP := ipTokens[0]
	fmt.Println("Origin IP : " + originIP)

	domainFromTrigger := jira.GetJiraDomain(wfDetails)
	domainIP := ""
	fmt.Println("Jira Site Domain : " + domainFromTrigger)

	//Get ip for domain
	var err error
	var byteData []byte

	if runtime.GOOS == "windows" {
		byteData, err = exec.Command("cmd", "/C", "nslookup "+domainFromTrigger).Output()
	} else {
		byteData, err = exec.Command("sh", "-c", "nslookup "+domainFromTrigger).Output()
	}

	if err != nil {
		fmt.Println(err.Error())
	} else {
		outTokens := strings.SplitAfter(string(byteData), "\n")
		possibleTokens := make([]string, 3)

		tokenIndex := 0
		for x := 0; x < len(outTokens); x++ {
			if strings.Contains(outTokens[x], "Address:") {
				possibleTokens[tokenIndex] = outTokens[x]
				tokenIndex++
			}
		}

		if len(possibleTokens) > 1 {
			domainIP = strings.TrimPrefix(possibleTokens[1], "Address: ")
		} else {
			domainIP = strings.TrimPrefix(possibleTokens[0], "Address: ")
		}

		domainIP = strings.TrimSpace(domainIP)

	}

	fmt.Println("Domain IP : " + domainIP)

	if domainIP != originIP {
		authRes := jira.JiraResponse{}
		authRes.IsSuccess = false
		authRes.Message = "Unauthorized Request Origin! Request was expected from " + domainFromTrigger + "(" + domainIP + "). But request was made from : " + originIP
		jsonDataa, _ := json.Marshal(authRes)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(jsonDataa).AddHeader("Content-Type", "application/json")
		return
	}

	result := jira.InvokeJiraTrigger(wfDetails)

	// conver the result back to JSON to return to the request origin location
	jsonData, _ := json.Marshal(result)

	if result.IsSuccess {
		T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
	} else {
		T.ResponseBuilder().SetResponseCode(500).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
	}
}

// method is used to generate a workflow file when the request is made from the Prcess Designer application
func (T ProcessEngineService) BuildFlow(FlowData Components.JsonFlow, flowName, sessionId string) {

	securityToken := T.Context.Request().Header.Get("Securitytoken")
	authenticateStruct := make(map[string]interface{})
	if securityToken == "" {
		logger.Log("SecurityToken not found executing BuildFlow. Domain : "+Common.GetDomainBySessionID(sessionId), logger.Error, sessionId, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Security Token Empty"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	authStatus, authSession := Common.AuthenticateSecurityToken(securityToken)
	_ = authSession
	if !authStatus {
		logger.Log("Invalid SecurityToken executing BuildFlow. Domain : "+Common.GetDomainBySessionID(sessionId), logger.Error, sessionId, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Invalid Security Token"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	// adding necessary data for the work request object
	FlowData.FlowName = flowName
	FlowData.SessionID = sessionId
	FlowData.ResponseMessage = make(chan *context.FlowResult)

	// Push the work onto the queue.
	WFPublishQueue <- FlowData

	for {
		select {
		case resp := <-FlowData.ResponseMessage:
			jsonData, _ := json.Marshal(resp)
			T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
			fmt.Println(string(jsonData))
			return
		}
	}
}

// used to invoke a existing workflow and give back the result from it
func (T ProcessEngineService) InvokeFlow(invokeFlowData Components.InvokeStruct) {

	securityToken := T.Context.Request().Header.Get("Securitytoken")
	authenticateStruct := make(map[string]interface{})
	if securityToken == "" {
		logger.Log("SecurityToken not found executing InvokeFlow. Domain : "+Common.GetDomainBySessionID(invokeFlowData.SessionID), logger.Error, invokeFlowData.SessionID, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Security Token Empty"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	authStatus, authSession := Common.AuthenticateSecurityToken(securityToken)
	_ = authSession
	if !authStatus {
		logger.Log("Invalid SecurityToken executing InvokeFlow. Domain : "+Common.GetDomainBySessionID(invokeFlowData.SessionID), logger.Error, invokeFlowData.SessionID, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Invalid Security Token"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	invokeFlowData.ResponseMessage = make(chan *context.FlowResult)

	// Push the work onto the queue.
	WFInvokeQueue <- invokeFlowData

	for {
		select {
		case resp := <-invokeFlowData.ResponseMessage:
			jsonData, _ := json.Marshal(resp)
			T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
			fmt.Println(string(jsonData))
			return
		}
	}

	/*	white := color.New(color.FgWhite)
		redBackground := white.Add(color.BgRed)
		redBackground.Println("********************************************************************************")
		fmt.Println("")
		fmt.Println("Method: Invoke workflow")
		fmt.Println("")
		fmt.Println("App code recieved : " + appCode)
		fmt.Println("ProcessCode recieved : " + processCode)
		fmt.Println("Session ID recieved : " + sessionId)
		fmt.Println("")
		result := Components.InvokeWorkflow(invokeFlowData, appCode, processCode, sessionId, false)

		jsonData, _ := json.Marshal(result)
		T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(jsonData)*/
}

// used when a new activity is created from DevStudio. This method will be used to publish it
func (T ProcessEngineService) PublishActivity(activityStruct Components.ActivityStruct, activityName, sessionId string) {

	securityToken := T.Context.Request().Header.Get("Securitytoken")
	authenticateStruct := make(map[string]interface{})
	if securityToken == "" {
		logger.Log("SecurityToken not found executing PublishActivity. Domain : "+Common.GetDomainBySessionID(sessionId), logger.Error, sessionId, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Security Token Empty"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	authStatus, authSession := Common.AuthenticateSecurityToken(securityToken)
	_ = authSession
	if !authStatus {
		logger.Log("Invalid SecurityToken executing PublishActivity. Domain : "+Common.GetDomainBySessionID(sessionId), logger.Error, sessionId, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Invalid Security Token"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	white := color.New(color.FgWhite)
	redBackground := white.Add(color.BgRed)
	redBackground.Println("********************************************************************************")
	fmt.Println("")
	fmt.Println("Method: Publish activity")
	fmt.Println("Activity Name : " + activityName)
	fmt.Println("Session ID recieved : " + sessionId)
	fmt.Println("")

	result := Components.PublishActivityFile(activityStruct, activityName, sessionId)

	jsonData, _ := json.Marshal(result)
	T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
}

//used when an already existing activity is removed from the processengine
func (T ProcessEngineService) RemoveActivity(activityStruct Components.ActivityStruct, activityName, sessionId string) {

	securityToken := T.Context.Request().Header.Get("Securitytoken")
	authenticateStruct := make(map[string]interface{})
	if securityToken == "" {
		logger.Log("SecurityToken not found executing RemoveActivity. Domain : "+Common.GetDomainBySessionID(sessionId), logger.Error, sessionId, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Security Token Empty"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	authStatus, authSession := Common.AuthenticateSecurityToken(securityToken)
	_ = authSession
	if !authStatus {
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Invalid Security Token"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	white := color.New(color.FgWhite)
	redBackground := white.Add(color.BgRed)
	redBackground.Println("********************************************************************************")
	fmt.Println("")
	fmt.Println("Method: Remove activity")
	fmt.Println("Activity Name : " + activityName)
	fmt.Println("Session ID recieved : " + sessionId)
	fmt.Println("")

	result := Components.RemoveActivityFile(activityStruct, activityName, sessionId)

	jsonData, _ := json.Marshal(result)
	T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
}

// InvokeHibernatedFlow method is used to invoke hibernated/
func (T ProcessEngineService) InvokeHibernatedFlow(invokeFlowData Components.InvokeStruct, appCode, processCode, sessionId string) {

	securityToken := T.Context.Request().Header.Get("Securitytoken")
	authenticateStruct := make(map[string]interface{})
	if securityToken == "" {
		logger.Log("SecurityToken not found executing InvokeHibernatedFlow. Domain : "+Common.GetDomainBySessionID(sessionId), logger.Error, sessionId, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Security Token Empty"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	authStatus, authSession := Common.AuthenticateSecurityToken(securityToken)
	_ = authSession
	if !authStatus {
		logger.Log("Invalid SecurityToken executing InvokeHibernatedFlow. Domain : "+Common.GetDomainBySessionID(sessionId), logger.Error, sessionId, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Invalid Security Token"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	white := color.New(color.FgWhite)
	redBackground := white.Add(color.BgRed)
	redBackground.Println("********************************************************************************")
	fmt.Println("")
	fmt.Println("Method: Invoke Hibernated workflow")
	fmt.Println("")
	fmt.Println("App code recieved : " + appCode)
	fmt.Println("ProcessCode recieved : " + processCode)
	fmt.Println("Session ID recieved : " + sessionId)
	fmt.Println("")
	result := Components.InvokeWorkflow(invokeFlowData, true)

	jsonData, _ := json.Marshal(result)
	T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
}

// Get version information method is used to get all the details about the version of process engine
func (T ProcessEngineService) GetVersionInfo() (versionResp Components.VersionResponse) {
	white := color.New(color.FgWhite)
	redBackground := white.Add(color.BgRed)
	redBackground.Println("********************************************************************************")
	fmt.Println("")
	fmt.Println("Method: Get Version information")
	fmt.Println("")
	versionResp = Components.GetVersionDetails("Y29tLnNtb290aGZsb3cuaW8tMDAwMA")

	//T.ResponseBuilder().SetResponseCode(200).Overide(true)
	return
}

func (T ProcessEngineService) ToggleLogs() string {
	msg := logger.ToggleLogs()
	return msg
}

func (T ProcessEngineService) ToggleLogstash() string {
	msg := logger.ToggleLogstash()
	return msg
}

func (T ProcessEngineService) GetEnv() bool {
	path := os.Getenv("PATH")
	gopath := os.Getenv("GOPATH")
	results := make(map[string]string)
	results["PATH"] = path
	results["GOPATH"] = gopath
	bytes, _ := json.Marshal(results)
	T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(bytes)
	return true
}

// Get session details method is used to get all the details about a particular session id.
// session types PE , WF, ACT, ALL
func (T ProcessEngineService) GetSessionDetails(sessionID, sessionType string) (sessionDetail Components.SessionTranDetails) {

	securityToken := T.Context.Request().Header.Get("Securitytoken")
	authenticateStruct := make(map[string]interface{})
	if securityToken == "" {
		logger.Log("SecurityToken not found executing GetSessionDetails. Domain : "+Common.GetDomainBySessionID(sessionID), logger.Error, sessionID, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Security Token Empty"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	authStatus, authSession := Common.AuthenticateSecurityToken(securityToken)
	_ = authSession
	if !authStatus {
		logger.Log("Invalid SecurityToken executing GetSessionDetails. Domain : "+Common.GetDomainBySessionID(sessionID), logger.Error, sessionID, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Invalid Security Token"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	white := color.New(color.FgWhite)
	redBackground := white.Add(color.BgRed)
	redBackground.Println("********************************************************************************")
	fmt.Println("")
	fmt.Println("Method: Get Session Details")
	fmt.Println("")
	sessionDetail = Components.GetSessionDetails(sessionID, sessionType)
	//T.ResponseBuilder().SetResponseCode(200).Overide(true)
	return
}

// this method is used to test the build or opened workflow from the app itself
func (T ProcessEngineService) TestWorkflow(TestData Components.TestWorkflowInvoke, flowName, sessionId string) {

	securityToken := T.Context.Request().Header.Get("Securitytoken")
	authenticateStruct := make(map[string]interface{})
	if securityToken == "" {
		logger.Log("SecurityToken not found executing TestWorkflow. Domain : "+Common.GetDomainBySessionID(sessionId), logger.Error, sessionId, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Security Token Empty"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	authStatus, authSession := Common.AuthenticateSecurityToken(securityToken)
	_ = authSession
	if !authStatus {
		logger.Log("Invalid SecurityToken executing TestWorkflow. Domain : "+Common.GetDomainBySessionID(sessionId), logger.Error, sessionId, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Invalid Security Token"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	white := color.New(color.FgWhite)
	redBackground := white.Add(color.BgRed)
	redBackground.Println("********************************************************************************")
	fmt.Println("")
	fmt.Println("Method: Test workflow")
	fmt.Println("")
	fmt.Println("Flow name recieved : " + flowName)
	fmt.Println("Session ID recieved : " + sessionId)
	fmt.Println("JSON string recieved : ", TestData)
	fmt.Println("")

	result := Components.TestWorkflow(TestData, flowName, sessionId, true)

	jsonData, _ := json.Marshal(result)
	T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
}

// this method is used to test the build or opened activity from the app itself
func (T ProcessEngineService) TestActivity(TestData Components.TestWorkflowInvoke, flowName, sessionId string) {

	securityToken := T.Context.Request().Header.Get("Securitytoken")
	authenticateStruct := make(map[string]interface{})
	if securityToken == "" {
		logger.Log("SecurityToken not found executing TestActivity. Domain : "+Common.GetDomainBySessionID(sessionId), logger.Error, sessionId, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Security Token Empty"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	authStatus, authSession := Common.AuthenticateSecurityToken(securityToken)
	_ = authSession
	if !authStatus {
		logger.Log("Invalid SecurityToken executing TestActivity. Domain : "+Common.GetDomainBySessionID(sessionId), logger.Error, sessionId, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Invalid Security Token"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	white := color.New(color.FgWhite)
	redBackground := white.Add(color.BgRed)
	redBackground.Println("********************************************************************************")
	fmt.Println("")
	fmt.Println("Method: Test Activity")
	fmt.Println("")
	fmt.Println("Activity name recieved : " + flowName)
	fmt.Println("Session ID recieved : " + sessionId)
	fmt.Println("JSON string recieved : ", TestData)
	fmt.Println("")

	result := Components.TestWorkflow(TestData, flowName, sessionId, false)

	jsonData, _ := json.Marshal(result)
	T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
}

// this method is used to install new thrid party packages to the backend
func (T ProcessEngineService) InstallPackages(installList Components.PackageInstaller, sessionId string) {

	securityToken := T.Context.Request().Header.Get("Securitytoken")
	authenticateStruct := make(map[string]interface{})
	if securityToken == "" {
		logger.Log("SecurityToken not found executing Install Packages. Domain : "+Common.GetDomainBySessionID(sessionId), logger.Error, sessionId, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Security Token Empty"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	authStatus, authSession := Common.AuthenticateSecurityToken(securityToken)
	_ = authSession
	if !authStatus {
		logger.Log("Invalid SecurityToken executing Install Packages. Domain : "+Common.GetDomainBySessionID(sessionId), logger.Error, sessionId, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Invalid Security Token"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	white := color.New(color.FgWhite)
	redBackground := white.Add(color.BgRed)
	redBackground.Println("********************************************************************************")
	fmt.Println("")
	fmt.Println("Method: Install Package")
	fmt.Println("")
	fmt.Println("Session ID recieved : " + sessionId)
	fmt.Println("JSON string recieved : ", installList)
	fmt.Println("")

	result := Components.InstallPackages(installList, sessionId)

	jsonData, _ := json.Marshal(result)
	T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
}

// this method is used to create an executable version of the workflow designed and download it directly from smoothflow server
func (T ProcessEngineService) DownloadExecutable(FlowData Components.JsonFlow, flowName, sessionId string) {

	securityToken := T.Context.Request().Header.Get("Securitytoken")
	authenticateStruct := make(map[string]interface{})
	if securityToken == "" {
		logger.Log("SecurityToken not found executing DownloadExecutable. Domain : "+Common.GetDomainBySessionID(sessionId), logger.Error, sessionId, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Security Token Empty"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	authStatus, authSession := Common.AuthenticateSecurityToken(securityToken)
	_ = authSession
	if !authStatus {
		logger.Log("Invalid SecurityToken executing DownloadExecutable. Domain : "+Common.GetDomainBySessionID(sessionId), logger.Error, sessionId, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Invalid Security Token"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	// displaying the method details on the console
	white := color.New(color.FgWhite)
	redBackground := white.Add(color.BgRed)
	redBackground.Println("********************************************************************************")

	fmt.Println("")
	fmt.Println("Method: Download Executable")
	fmt.Println("")
	fmt.Println("Flow name recieved : " + flowName)
	fmt.Println("Session ID recieved : " + sessionId)
	fmt.Println("")

	//result := Components.DownloadExecutable(FlowData, sessionId)
	makeExecutable := true
	result := Components.InitializeGoFlow(FlowData, flowName, sessionId, makeExecutable)

	// conver the result back to JSON to return to the request origin location
	jsonData, _ := json.Marshal(result)
	T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
}

// this function is used to publish workflows to a docker
func (T ProcessEngineService) PublishToDocker(wfDetails Components.DockerDeployement) {

	securityToken := T.Context.Request().Header.Get("Securitytoken")
	authenticateStruct := make(map[string]interface{})
	if securityToken == "" {
		logger.Log("SecurityToken not found executing PublishToDocker. Domain : "+Common.GetDomainBySessionID(wfDetails.SessionID), logger.Error, wfDetails.SessionID, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Security Token Empty"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	authStatus, authSession := Common.AuthenticateSecurityToken(securityToken)
	_ = authSession
	if !authStatus {
		logger.Log("Invalid SecurityToken executing PublishToDocker. Domain : "+Common.GetDomainBySessionID(wfDetails.SessionID), logger.Error, wfDetails.SessionID, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Invalid Security Token"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	// displaying the method details on the console
	white := color.New(color.FgWhite)
	redBackground := white.Add(color.BgRed)
	redBackground.Println("********************************************************************************")

	fmt.Println("")
	fmt.Println("Method: Publish to Docker")
	fmt.Println("")
	fmt.Println("Flow name recieved : " + wfDetails.WFName)
	fmt.Println("Session ID recieved : " + wfDetails.SessionID)
	fmt.Println("")

	if wfDetails.RAM == "" {
		//set default ram size
		wfDetails.RAM = "300m"
	} else {
		wfDetails.RAM += "m"
	}

	cpuCount := float64(runtime.NumCPU())
	if wfDetails.CPU == "" {
		//set default CPU size
		wfDetails.CPU = strconv.FormatFloat(cpuCount, 'f', 6, 64)
	} else {
		cpuRatio, _ := strconv.ParseFloat(wfDetails.CPU, 64)
		calculatedCpuLimit := (cpuCount / (100 / cpuRatio))
		limitInString := strconv.FormatFloat(calculatedCpuLimit, 'f', 6, 64)
		wfDetails.CPU = limitInString[:3]
	}

	// publish the exe to a docker
	result := Components.PublishToDocker(wfDetails)

	// conver the result back to JSON to return to the request origin location
	jsonData, _ := json.Marshal(result)
	T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
}

// this function is used to publish workflows to a docker
func (T ProcessEngineService) RemoveDocker(wfDetails Components.DockerDeployement) {

	securityToken := T.Context.Request().Header.Get("Securitytoken")
	authenticateStruct := make(map[string]interface{})
	if securityToken == "" {
		logger.Log("SecurityToken not found executing RemoveDocker. Domain : "+Common.GetDomainBySessionID(wfDetails.SessionID), logger.Error, wfDetails.SessionID, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Security Token Empty"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	authStatus, authSession := Common.AuthenticateSecurityToken(securityToken)
	_ = authSession
	if !authStatus {
		logger.Log("Invalid SecurityToken executing RemoveDocker. Domain : "+Common.GetDomainBySessionID(wfDetails.SessionID), logger.Error, wfDetails.SessionID, logger.ProcessEngine)
		authenticateStruct["Status"] = false
		authenticateStruct["Message"] = "Invalid Security Token"
		byteArray, _ := json.Marshal(authenticateStruct)
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride(byteArray)
		return
	}

	// displaying the method details on the console
	white := color.New(color.FgWhite)
	redBackground := white.Add(color.BgRed)
	redBackground.Println("********************************************************************************")

	fmt.Println("")
	fmt.Println("Method: Remove Docker")
	fmt.Println("")
	fmt.Println("Docker name recieved : " + wfDetails.WFName)
	fmt.Println("Session ID recieved : " + wfDetails.SessionID)
	fmt.Println("")

	// remove the docker
	result := Components.RemoveDocker(wfDetails)

	// conver the result back to JSON to return to the request origin location
	jsonData, _ := json.Marshal(result)
	T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(jsonData).AddHeader("Content-Type", "application/json")
}

func main() {

	//Verify agent.config is present. If not create one with settings for LIVE server
	Common.VerifyCEBAgentConfig()
	runtime.GOMAXPROCS(runtime.NumCPU())

	logger.InitiateLogger()

	Components.ProcessEngineStartTime = time.Now()

	cebadapter.Attach("ProcessEngine", func(s bool) {
		cebadapter.GetLatestGlobalConfig("StoreConfig", func(data []interface{}) {
			logger.TermLog("Store Configuration Successfully Loaded...", 10)

			agent := cebadapter.GetAgent()

			agent.Client.OnEvent("globalConfigChanged.StoreConfig", func(from string, name string, data map[string]interface{}, resources map[string]interface{}) {
				cebadapter.GetLatestGlobalConfig("StoreConfig", func(data []interface{}) {
					logger.TermLog("Store Configuration Successfully Updated...", 10)
				})
			})
		})
		logger.TermLog("Successfully registered in CEB", 10)
	})

	white := color.New(color.FgWhite)
	redBackground := white.Add(color.BgRed)

	redBackground.Println("********************************************************************************")
	color.Cyan("")
	color.Cyan("       ______                              _____            _            ")
	color.Cyan("       | ___ \\                            |  ___|          (_)           ")
	color.Cyan("       | |_/ / __ ___   ___ ___  ___ ___  | |__ _ __   __ _ _ _ __   ___ ")
	color.Cyan("       |  __/ '__/ _ \\ / __/ _ \\/ __/ __| |  __| '_ \\ / _` | | '_ \\ / _ \\")
	color.Cyan("       | |  | | | (_) | (_|  __/\\__ \\__ \\ | |__| | | | (_| | | | | |  __/")
	color.Cyan("       \\_|  |_|  \\___/ \\___\\___||___/___/ \\____/_| |_|\\__, |_|_| |_|\\___|")
	color.Cyan("                                                       __/ |             ")
	color.Cyan("                                                      |___/              ")
	color.Cyan("")
	fmt.Println("Process Engine service has started on port 8093")
	fmt.Println("Methods:")
	fmt.Println("\t /TestBuildActivity/{activityName}/{sessionId}")
	fmt.Println("\t /TestRunActivity/{activityName}/{sessionId}")
	fmt.Println("")

	go runRestFul()
	go runQueues()

	forever := make(chan bool)
	<-forever
}

func runRestFul() {
	gorest.RegisterService(new(ProcessEngineService))

	err := http.ListenAndServe(":8093", gorest.Handle())
	if err != nil {
		logger.TermLog(err.Error(), logger.Error)
		return
	}

}

func runQueues() {
	numberOfWorkers := 4
	WFPublishWorkerQueue = make(chan chan Components.JsonFlow, numberOfWorkers)
	InvokeWorkerQueue = make(chan chan Components.InvokeStruct, numberOfWorkers)

	// Now, create all of our workers for WFPublish method
	for i := 0; i < numberOfWorkers; i++ {
		fmt.Println("Starting WFPublisher worker: ", i+1)
		worker := queue.NewWFPublishWorker(i+1, WFPublishWorkerQueue)
		worker.Start()
	}
	for i := 0; i < numberOfWorkers; i++ {
		fmt.Println("Starting Invoke worker: ", i+1)
		worker := queue.NewInvokeWorker(i+1, InvokeWorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WFPublishQueue:
				fmt.Println("")
				fmt.Println("Workflow publish requeust received")
				go func() {
					worker := <-WFPublishWorkerQueue
					fmt.Println("Dispatching now...")
					worker <- work
				}()
			}
		}
	}()
	go func() {
		for {
			select {
			case work := <-WFInvokeQueue:
				fmt.Println("")
				fmt.Println("Workflow invoke requeust received")
				go func() {
					worker := <-InvokeWorkerQueue
					fmt.Println("Dispatching now...")
					worker <- work
				}()
			}
		}
	}()
}
