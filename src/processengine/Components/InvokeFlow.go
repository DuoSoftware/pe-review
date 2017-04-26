package Components

import (
	"encoding/json"
	linq "github.com/ahmetb/go-linq"
	"os"
	"os/exec"
	"processengine/Common"
	"processengine/context"
	"processengine/logger"
	"processengine/objectstore"
	"runtime"
)

func InvokeWorkflow(invokeFlowData InvokeStruct, IsHibernated bool) *context.FlowResult {

	//geting the initial session id from
	sessionId := invokeFlowData.SessionID

	var flowResult = new(context.FlowResult)
	flowResult.FlowName = invokeFlowData.ProcessCode
	flowResult.SessionID = sessionId
	flowResult.Status = true
	flowResult.Message = "Invoking workflow"

	logger.Log_PE("~~ Initiating Workflow invoke", logger.Information, sessionId)
	/*red := color.New(color.FgWhite)
	ErrorColorscheme := red.Add(color.BgRed)

	//logger.Log_PE("********************************************************************************")


	decoder := json.NewDecoder(request.Body)
	var invokeFlowData InvokeStruct
	decodeError := decoder.Decode(&invokeFlowData)
	if decodeError != nil {
		logger.Log_PE("There was an error Decoding the jsonData sent to activity invoke method.")
		ErrorColorscheme.Println(decodeError.Error())
		flowResult.Message = flowResult.Message + decodeError.Error() + " -> "
		flowResult.Status = false
		logger.Log_PE("")
		return flowResult
		}*/

	var slash = ""
	if runtime.GOOS == "windows" {
		slash = "\\"
	} else {
		slash = "/"
	}

	// get all the process codes "for now" and get the process record relavent to the process code
	logger.Log_PE("~~ Getting Process mapping details", logger.Debug, sessionId)
	logger.Log_PE("Received Namespace: "+invokeFlowData.Namespace, logger.Debug, sessionId)
	logger.Log_PE("Received SecurityToken: "+invokeFlowData.SecurityToken, logger.Debug, sessionId)

	getmapping := objectstore.GetByQuery{}

	var getAll_parameters map[string]interface{}
	getAll_parameters = make(map[string]interface{})
	getAll_parameters["securityToken"] = invokeFlowData.SecurityToken
	getAll_parameters["log"] = invokeFlowData.Log
	getAll_parameters["namespace"] = invokeFlowData.Namespace
	getAll_parameters["class"] = "process_mapping"
	getAll_parameters["query"] = "select * from process_mapping where ProcessCode = '" + invokeFlowData.ProcessCode + "'"

	logger.Log_PE("Query: "+getAll_parameters["query"].(string), logger.Debug, sessionId)

	var mappingList []ProcessMapping
	all_mappings := getmapping.Invoke(getAll_parameters)

	hasMap := false

	logger.Log_PE("Received Data:", logger.Debug, sessionId)
	logger.Log_PE(all_mappings.SharedContext, logger.Debug, sessionId)

	err := json.Unmarshal([]byte(string(all_mappings.SharedContext)), &mappingList)
	if err != nil {
		msg := "There was an error converting objectStore result to JsonObject. Error: "
		logger.Log_PE(msg, logger.Error, sessionId)
		logger.Log_PE(string(err.Error()), logger.Error, sessionId)
		flowResult.Message = msg
		hasMap = false
		return flowResult
	} else {
		logger.Log_PE("Process mapping details received.", logger.Debug, sessionId)
		logger.Log_PE(string(all_mappings.SharedContext), logger.Debug, sessionId)
		hasMap = true
	}
	numMaps := "Number of Mappings recieving: " + string(cap(mappingList))
	logger.Log_PE(numMaps, logger.Debug, sessionId)
	flowResult.Message = numMaps
	logger.Log_PE("Filtering mapping details.", logger.Debug, sessionId)
	// filter the mapping record from the list
	mapping, hasMap, _ := linq.From(mappingList).Where(
		func(mapp linq.T) (bool, error) {
			return mapp.(ProcessMapping).ProcessCode == invokeFlowData.ProcessCode, nil
		}).First()

	if hasMap == false {
		msg := "Mapping not found on the Process Mapping list."
		logger.Log_PE(msg, logger.Error, sessionId)
		flowResult.Message = msg
		flowResult.Status = false
		logger.Log_PE("", logger.Error, sessionId)
		return flowResult
	} else {
		logger.Log_PE("Workflow ID: "+string(mapping.(ProcessMapping).WorkflowID), logger.Debug, sessionId)
		logger.Log_PE("Workflow Name: "+string(mapping.(ProcessMapping).Name), logger.Debug, sessionId)
		logger.Log_PE("Exe Name: "+string(mapping.(ProcessMapping).Name)+".exe", logger.Debug, sessionId)
	}

	workflowname := string(mapping.(ProcessMapping).Name)

	// get the tagged WF id for the Process code and sending it to be retrived
	/*logger.Log_PE("", logger.Debug, sessionId)
	logger.Log_PE("~~ Get WF details for the ProcessCode:"+string(mapping.(ProcessMapping).ProcessCode), logger.Debug, sessionId)

	getWFrecord := objectstore.GetByKey{}

	var getWF_parameters map[string]interface{}
	getWF_parameters = make(map[string]interface{})
	getWF_parameters["securityToken"] = invokeFlowData.SecurityToken
	getWF_parameters["log"] = invokeFlowData.Log
	getWF_parameters["namespace"] = invokeFlowData.Namespace
	getWF_parameters["class"] = "process_flows"
	getWF_parameters["key"] = mapping.(ProcessMapping).WorkflowID

	var workflowObject nodedata
	WF := getWFrecord.Invoke(getWF_parameters)

	WFerr := json.Unmarshal([]byte(string(WF.SharedContext)), &workflowObject)
	//camelCasedWFName := ""
	if WFerr != nil {
		msg := "There was an error converting objectStore result to JsonObject. Error: "
		logger.Log_PE(msg, logger.Debug, sessionId)
		logger.Log_PE(string(err.Error()), logger.Debug, sessionId)

		flowResult.Message = msg
		flowResult.Status = false
		logger.Log_PE("", logger.Debug, sessionId)
	} else {
		//camelCasedWFName = Common.MakeFirstLowerCase(Common.CamelCase(workflowObject.Name))
		logger.Log_PE("Workflow Name: "+string(workflowObject.Name), logger.Debug, sessionId)
		logger.Log_PE("Workflow Display name: "+string(workflowObject.DisplayName), logger.Debug, sessionId)
		logger.Log_PE("Exe Name: "+workflowObject.Name+".exe", logger.Debug, sessionId)
		}*/

	// the required keys will be added to the flowData
	WFParameters := ""
	TempParameters := ""
	var tempData map[string]interface{}
	if jsonParseErr2 := json.Unmarshal([]byte(invokeFlowData.JSONData), &tempData); jsonParseErr2 != nil {
		logger.Log_PE("The JSON input is not in correct format.", logger.Debug, sessionId)
	}
	tempData["InSessionID"] = invokeFlowData.SessionID
	tempData["InSecurityToken"] = invokeFlowData.SecurityToken
	tempData["InLog"] = invokeFlowData.Log
	tempData["InNamespace"] = invokeFlowData.Namespace
	temConv, _ := json.Marshal(tempData)
	TempParameters = string(temConv)

	// if the call is hibernated, it will function differently
	if IsHibernated == true {
		// get stored objects for the session
		getHData := objectstore.GetByKey{}

		var getHD_parameters map[string]interface{}
		getHD_parameters = make(map[string]interface{})
		getHD_parameters["securityToken"] = invokeFlowData.SecurityToken
		getHD_parameters["log"] = invokeFlowData.Log
		getHD_parameters["namespace"] = invokeFlowData.Namespace
		getHD_parameters["class"] = "hibernated_workflows"
		getHD_parameters["key"] = sessionId

		var HWFData HibernatedWF
		HWFD := getHData.Invoke(getHD_parameters)

		// get the hibernated data from object store
		if HWFDerr := json.Unmarshal([]byte(string(HWFD.SharedContext)), &HWFData); HWFDerr != nil {
			logger.Log_PE("Failed to get details of HibernatedWF from objectstore.", logger.Error, sessionId)
		}
		// the details send for the Hibernated data will be added to the saved data from this point onwards and the next Executional level will be updated
		var flowData map[string]interface{}
		if jsonParseErr := json.Unmarshal([]byte(TempParameters), &flowData); jsonParseErr != nil {
			logger.Log_PE("The JSON input is not in correct format.", logger.Error, sessionId)
		} else {
			// updating the past Executional level with the next level to be executed
			HWFData.FlowData["ExecutionLevel"] = flowData["ExecutionLevel"]
			// the new fields are been added to the hibernated object
			for key, val := range flowData {
				if _, ok := HWFData.FlowData[key]; ok {
					//logger.Log_PE("OK - " + string(key) + " - " + flowData[key].(string))
				} else {
					logger.Log_PE("New Key,Value: "+string(key)+" - "+flowData[key].(string), logger.Debug, sessionId)
					HWFData.FlowData[key] = val
				}
			}
		}
		// convert the updated data object to json string
		convertedStruct, _ := json.Marshal(HWFData.FlowData)
		logger.Log_PE("Updated Arguments : "+string(convertedStruct), logger.Debug, sessionId)
		WFParameters = string(convertedStruct)
	} else {
		WFParameters = TempParameters
	}

	// caling workflow with the specific name
	pwd, _ := os.Getwd()
	isExist := false
	if runtime.GOOS == "windows" {
		isExist = IsExists(pwd + slash + workflowname + ".exe")
	} else {
		isExist = IsExists(pwd + slash + workflowname)
	}

	if isExist == false {
		logger.Log_PE("Existing "+workflowname+".exe has not found to execute", logger.Error, sessionId)
		flowResult.Status = false
		flowResult.Message = "Existing " + workflowname + ".exe has not found to execute."

	} else {
		if runtime.GOOS == "windows" {
			// if execute is running in windows
			logger.Log_PE("Runtime Platform: Environment - "+string(runtime.GOOS), logger.Debug, sessionId)
			flowResult.Message = "In windows environment"
			out, err := exec.Command(pwd+slash+workflowname+".exe", WFParameters, "CMD", "start").Output()

			if err != nil {
				flowResult.Status = false
				flowResult.Message = err.Error()
				logger.Log_PE("Error on execute.", logger.Error, sessionId)
			}

			logger.Log_PE("Starting Unmarshal the response.", logger.Information, sessionId)

			res := context.ReturnData{}
			jerr := json.Unmarshal([]byte(string(out)), &res)
			if jerr != nil {
				logger.Log_PE("Error occured on Unmarshal.", logger.Error, sessionId)
			} else {
				flowResult.ReturnData = res
			}

			flowResult.Status = true
			flowResult.Message = "Workflow has completed executing with an output."
			logger.Log_PE(workflowname+".exe"+" has executed with an output.", logger.Debug, sessionId)
			logger.Log_PE(string(out), logger.Debug, sessionId)
			logger.Log_PE("Please check the WorkFlowLog for more details.", logger.Debug, sessionId)
			logger.Log_PE("Workflow Invoke completed!", logger.Debug, sessionId)

			var logfile string
			logfile = Common.GetPELog(sessionId)
			flag := false
			if flag == true {
				logger.Log_PE(logfile, logger.Debug, sessionId)
			}

		} else {
			// if execute is running in linux
			logger.Log_PE("Runtime Platform: Environment - "+string(runtime.GOOS), logger.Debug, sessionId)
			flowResult.Message = "In linux environment -> "
			out, err := exec.Command(pwd+slash+workflowname, WFParameters, "sh").Output()

			if err != nil {
				flowResult.Status = false
				flowResult.Message = err.Error()
				logger.Log_PE("Error on execute.", logger.Error, sessionId)
			}

			logger.Log_PE("Starting Unmarshal the response.", logger.Error, sessionId)

			res := context.ReturnData{}
			jerr := json.Unmarshal([]byte(string(out)), &res)
			if jerr != nil {
				logger.Log_PE("Error occured on Unmarshal.", logger.Error, sessionId)
			} else {
				flowResult.ReturnData = res
			}

			flowResult.Status = true
			flowResult.Message = "Workflow has completed executing with an output."
			logger.Log_PE(workflowname+" has executed with an output", logger.Debug, sessionId)
			logger.Log_PE(string(out), logger.Debug, sessionId)
			logger.Log_PE("Please check the WorkFlowLog for more details.", logger.Debug, sessionId)
			logger.Log_PE("Workflow Invoke completed!", logger.Information, sessionId)

			var logfile string
			logfile = Common.GetPELog(sessionId)
			flag := false
			if flag == true {
				logger.Log_PE(logfile, logger.Debug, sessionId)
			}
		}
	}

	return flowResult
}
