package Components

import (
	// "go/parser"
	// "go/token"
	"crypto/rand"
	"encoding/json"
	linq "github.com/ahmetb/go-linq"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	//import "time"
	//import "strconv"
	"processengine/context"
	//import "net/http"
	"fmt"
	"processengine/logger"
	"runtime"
)

func InitializeGoFlow(flowStruct JsonFlow, flowName, sessionId string, makeExecutable bool) *context.FlowResult {

	// convert the JSON object to struct
	//flowStruct := GetJsonData(jsonInput)
	var flowResult = new(context.FlowResult)
	flowResult.FlowName = flowName
	flowResult.SessionID = sessionId
	flowResult.Status = true
	flowResult.Message = "Building flow has initialized"

	var SupportedDataTypes []string
	SupportedDataTypes = append(SupportedDataTypes, "string", "int", "float", "boolean", "array(string)")

	var wfcontext = new(context.WorkflowContext)

	//logger.Log_PE("********************************************************************************")
	logger.Log_PE("~~ Initiating Workflow Build", logger.Information, sessionId)

	/*
		red := color.New(color.FgWhite)
		ErrorColorscheme := red.Add(color.BgRed)

		decoder := json.NewDecoder(request.Body)
		var flowStruct JsonFlow
		decodeError := decoder.Decode(&flowStruct)
		if decodeError != nil {
						logger.Log_PE("There was an error Decoding the jsonData sent to process.")
						ErrorColorscheme.Println(decodeError.Error())
						flowResult.Message = flowResult.Message + decodeError.Error() + " -> "
						flowResult.Status = false
						logger.Log_PE("")
						return flowResult
		}
		convertedStruct, _ := json.Marshal(flowStruct)
		logger.Log_PE("Json string recieved : " + string(convertedStruct))

		logger.Log_PE("")*/
	//logger.Log_PE(flowStruct)

	var slash = ""
	if runtime.GOOS == "windows" {
		slash = "\\"
	} else {
		slash = "/"
	}

	logger.Log_PE("Initializing flow: "+flowName, logger.Debug, sessionId)
	logger.Log_PE("Session ID: "+sessionId, logger.Debug, sessionId)

	pwd, _ := os.Getwd()
	buildPath := ""

	if makeExecutable == false {
		buildPath = pwd + slash + "BuiltFlows"
	} else {
		buildPath = pwd + slash + "BuiltExecutables"
	}

	// check for existing workflow file, if exist to rename it as a backup
	checkPath := buildPath + slash + flowName + slash + flowName + ".go"
	logger.Log_PE("File location path: "+checkPath, logger.Debug, sessionId)
	archiveFlow := IsExists(checkPath)
	if archiveFlow == true {
		//newName := strconv.Itoa(time.Now().Year()) + "_" + time.Now().Month().String() + "_" + strconv.Itoa(time.Now().Day()) + "_" + strconv.Itoa(time.Now().Hour()) + "_" + strconv.Itoa(time.Now().Minute()) + "_" + strconv.Itoa(time.Now().Second())
		//logger.Log_PE("Archiving existing flow file")
		logger.Log_PE("Removing existing flow file", logger.Debug, sessionId)
		flowResult.Message = "Removing existing flow file"
		err := os.Remove(buildPath + slash + flowName + slash + flowName + ".go")
		if err != nil {
			logger.Log_PE(err.Error(), logger.Error, sessionId)
			flowResult.Message = err.Error()
			flowResult.Status = false
			return flowResult
		}
		logger.Log_PE("File exists. File removed.", logger.Debug, sessionId)
	} else {
		logger.Log_PE("File did not exist and continuing.", logger.Debug, sessionId)
		// this is a comment
	}
	// arguments will be initiated only once
	argumentsInitiated := false
	var nodeDataObj nodedata

	//initialize argument variables
	logger.HighLight("-- Task: Initializing Workflow Arguments", sessionId)

	ConstructGoFlow(sessionId, flowName, "addimport", "encoding/json", "ImportJSON", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
	//ConstructGoFlow(sessionId, flowName, "addimport", "github.com/fatih/color", "ImportCOLOR", "initiate.View.drawboard", nodeData.(nodedata), flowResult, flowStruct,makeExecutable)
	ConstructGoFlow(sessionId, flowName, "addimport", "processengine/context", "ImportCONTEXT", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
	ConstructGoFlow(sessionId, flowName, "init.var", "var WFContext context.WorkflowContext", "InitVariableCONTEXT", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)

	ConstructGoFlow(sessionId, flowName, "init.var", `WFContext.Message = "Startng Workflow";WFContext.Status = true;WFContext.ErrorCode = 1`, "InitVariablesOfWFContext", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)

	//ConstructGoFlow(sessionId, flowName, "init.var", "activityTrace := make(map[string]interface{})", "InitVariableActivityTrace", "initiate.View.drawboard", nodeData.(nodedata), flowResult, flowStruct,makeExecutable)
	ConstructGoFlow(sessionId, flowName, "init.var", "var flowData map[string]interface{}", "InitVariableFLOWDATA", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
	if makeExecutable == false {
		ConstructGoFlow(sessionId, flowName, "addimport", "fmt", "ImportFMT", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
		ConstructGoFlow(sessionId, flowName, "addimport", "os", "ImportOS", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
		ConstructGoFlow(sessionId, flowName, "init.var", "jsonInput := os.Args[1]", "InitVariableJSONINPUT", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
		ConstructGoFlow(sessionId, flowName, "init.var", `if jsonParseErr := json.Unmarshal([]byte(jsonInput), &flowData); jsonParseErr != nil { logger.Log_WF("The JSON input is not in correct format.", logger.Error, flowData["InSessionID"].(string)); logger.Log_WF("Application terminated.", logger.Debug, flowData["InSessionID"].(string)); logger.Log_WF(string(jsonParseErr.Error()), logger.Error, flowData["InSessionID"].(string));WFContext.Message = "The JSON input is not in correct format.";WFContext.Status = false;WFContext.ErrorCode = 4;var returnObj context.ReturnData;returnObj.JSONOutput = flowData;returnObj.WorkflowResult.Message = WFContext.Message;returnObj.WorkflowResult.Status = WFContext.Status;returnObj.WorkflowResult.ErrorCode = WFContext.ErrorCode;returnDataJSON, rtDError := json.Marshal(returnObj);if rtDError != nil {logger.Log_WF("WF return data marshal error", logger.Error, flowData["InSessionID"].(string));};fmt.Print(string(returnDataJSON));return;}`, "InitVariableJsonParseInput", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)

		/*ConstructGoFlow(sessionId, flowName, "init.headers", getHeaderInformation(flowStruct, makeExecutable), "INIT.Headers", "INIT.Headers", nodeDataObj, flowResult, flowStruct, makeExecutable)
		ConstructGoFlow(sessionId, flowName, "init.var", "headerInput := os.Args[2]", "InitHeaderJSONINPUT", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
		ConstructGoFlow(sessionId, flowName, "init.var", `if jsonParseErr := json.Unmarshal([]byte(headerInput), &Headers); jsonParseErr != nil { logger.Log_WF("Required headers are not available", logger.Debug, flowData["InSessionID"].(string)); logger.Log_WF("Application terminated.", logger.Debug, flowData["InSessionID"].(string)); logger.Log_WF(string(jsonParseErr.Error()), logger.Debug, flowData["InSessionID"].(string));WFContext.Message = "Required headers are not available.";WFContext.Status = false;WFContext.ErrorCode = 4;var returnObj context.ReturnData;returnObj.JSONOutput = flowData;returnObj.WorkflowResult.Message = WFContext.Message;returnObj.WorkflowResult.Status = WFContext.Status;returnObj.WorkflowResult.ErrorCode = WFContext.ErrorCode;returnDataJSON, rtDError := json.Marshal(returnObj);if rtDError != nil {logger.Log_WF("WF return data marshal error", logger.Debug, flowData["InSessionID"].(string));};fmt.Print(string(returnDataJSON));return;}`, "InitHeaderJsonParseInput", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)*/
	} else {
		ConstructGoFlow(sessionId, flowName, "addimport", "fmt", "ImportFMT", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
		ConstructGoFlow(sessionId, flowName, "init.var", "jsonInput := data", "InitVariableJSONINPUT", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
		ConstructGoFlow(sessionId, flowName, "init.var", `if jsonParseErr := json.Unmarshal([]byte(jsonInput), &flowData); jsonParseErr != nil { logger.Log_WF("The JSON input is not in correct format.", logger.Debug, flowData["InSessionID"].(string)); logger.Log_WF("Application terminated.", logger.Debug, flowData["InSessionID"].(string)); logger.Log_WF(string(jsonParseErr.Error()), logger.Debug, flowData["InSessionID"].(string));WFContext.Message = "The JSON input is not in correct format.";WFContext.Status = false;WFContext.ErrorCode = 4;var returnObj context.ReturnData;returnObj.JSONOutput = flowData;returnObj.WorkflowResult.Message = WFContext.Message;returnObj.WorkflowResult.Status = WFContext.Status;returnObj.WorkflowResult.ErrorCode = WFContext.ErrorCode;returnDataJSON, rtDError := json.Marshal(returnObj);if rtDError != nil {logger.Log_WF("WF return data marshal error", logger.Debug, flowData["InSessionID"].(string));};fmt.Print(string(returnDataJSON));return string(returnDataJSON);}`, "InitVariableJsonParseInput", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)

		ConstructGoFlow(sessionId, flowName, "dlExec", `port := "`+flowStruct.Port+`"`, "INIT.PORTNUMBER", "", nodeDataObj, flowResult, flowStruct, makeExecutable)
		ConstructGoFlow(sessionId, flowName, "printInfoNameForCEB", flowName, "INIT.FLOWNAMECEB", "", nodeDataObj, flowResult, flowStruct, makeExecutable)
		//Added by Jay
		ConstructGoFlow(sessionId, flowName, "printInfoName", flowName, "INIT.FLOWNAME", "", nodeDataObj, flowResult, flowStruct, makeExecutable)
		ConstructGoFlow(sessionId, flowName, "printInfoBuild", flowStruct.OSCode+"-"+flowStruct.SysArch, "INIT.FLOWBUILDPLATFORM", "INIT.FlowBuildPlatform", nodeDataObj, flowResult, flowStruct, makeExecutable)
		ConstructGoFlow(sessionId, flowName, "printInArugements", getInArgumentPrintDetails(flowStruct), "INIT.PRINTINARGUMENTS", "INIT.PrintInArguments", nodeDataObj, flowResult, flowStruct, makeExecutable)
		ConstructGoFlow(sessionId, flowName, "printOutArugements", getOutArgumentPrintDetails(flowStruct), "INIT.PRINTOUTARGUMENTS", "INIT.PrintOutArguments", nodeDataObj, flowResult, flowStruct, makeExecutable)
		ConstructGoFlow(sessionId, flowName, "printWfBodySample", getInJSONPrintDetails(flowStruct), "INIT.PRINTWFBODYSAMPLE", "INIT.PrintWfBodySample", nodeDataObj, flowResult, flowStruct, makeExecutable)

		// initiating headers
		ConstructGoFlow(sessionId, flowName, "init.headers", getHeaderInformation(flowStruct, makeExecutable), "INIT.Headers", "INIT.Headers", nodeDataObj, flowResult, flowStruct, makeExecutable)
	}

	ConstructGoFlow(sessionId, flowName, "dlExec", getConfigContent(flowStruct, makeExecutable), "INIT.CONFIG", "", nodeDataObj, flowResult, flowStruct, makeExecutable)

	ConstructGoFlow(sessionId, flowName, "init.var", getAllINArgumentList(flowStruct, makeExecutable), "AllINArguments", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)

	ConstructGoFlow(sessionId, flowName, "init.var", "ALLOutArguments := make(map[string]interface{});", "AllOUTArguments", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
	// initiate inarguments
	ConstructGoFlow(sessionId, flowName, "init.arguments", "InArgument", "INIT.INArgments", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
	// initiate custom variables
	ConstructGoFlow(sessionId, flowName, "init.arguments", "Custom", "INIT.CustArguments", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)

	ConstructGoFlow(sessionId, flowName, "init.arguments", "OutArgument", "INIT.OutArguments", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)

	//ConstructGoFlow(sessionId, flowName, "addimport", "processengine/Components", "ImportComponents", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct,makeExecutable)
	//ConstructGoFlow(sessionId, flowName, "init.executionLevels", getBasicExecutionalFunctionData(), "INIT.executionLevels", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct,makeExecutable)
	//GG
	//ConstructGoFlow(sessionId, flowName, "init.executionLevels", getExecutionalData(flowStruct, "drawboard", sessionId), "INIT.executionLevels", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)

	// iterate through []views and create flow for each view separately
	for _, view := range flowStruct.Views {

		ConstructGoFlow(sessionId, flowName, "init.executionLevels", getExecutionalData(flowStruct, view, sessionId), "INIT.executionLevels", "initiate.View."+view, nodeDataObj, flowResult, flowStruct, makeExecutable)

		logger.HighLight("Processing View: "+view, sessionId)

		//this line searchs for the starting node of the view
		startNode, _, _ := linq.From(flowStruct.Nodes).Where(
			func(in linq.T) (bool, error) {
				return in.(nodedata).LibraryID == "0" && in.(nodedata).ParentView == view, nil
			}).First()

		viewNodes, _ := linq.From(flowStruct.Nodes).Where(
			func(in linq.T) (bool, error) {
				return in.(nodedata).ParentView == view, nil
			}).Results()

		stopFlow := false
		fmt.Println("Start node:", startNode)
		fmt.Println("Other nodes:", viewNodes)

		executionlevel := "0"
		previous_executionlevel := "0"
		targetNodeId := ""

		if len(viewNodes) == 0 {
			stopFlow = true
		} else {
			if startNode == nil {
				logger.Log_PE("Start node was not found in view "+view, logger.Debug, sessionId)
				flowResult.Message = "Start node was not found in view " + view
				flowResult.Status = false
				return flowResult
			} else {
				logger.HighLight("View: initiate.View."+view+"."+executionlevel, sessionId)
				targetNodeId = startNode.(nodedata).SchemaID
			}
		}
		/*if startNode == nil {
		logger.Log_PE("Start node was not found in view "+view, sessionId)
		flowResult.Message = "Start node was not found in view " + view
		flowResult.Status = false
		return flowResult
		}*/

		//logger.Log_PE("Start node " + startNode.(nodedata).SchemaID)

		// check if any hibernate values are present.. if so print comment lines for them through a loop

		/*hibernateVariables, _ := linq.From(flowStruct.Nodes).Where(func(in linq.T) (bool, error) {
						return in.(nodedata).LibraryID == "3" && in.(nodedata).ParentView == view, nil
		}).Results()
		convertedNodes, _ := json.Marshal(hibernateVariables)
		logger.Log_PE("Available Nodes for the view: "+view, sessionId)
		logger.Log_PE(string(convertedNodes), sessionId)

		execString := ``

		// add hibernate value 0 object to the list
		//execString = execString + "\n\n" + `if (flowData["ExecutionLevel"] == "0") {logger.Log_WF("Execution executionlevel: 0", logger.Debug, flowData["InSessionID"].(string));` + "\n" + `//initiate.View.` + view + `.0`

		//tot := cap(hibernateVariables)
		var count int = 0
		for _, node := range hibernateVariables {
						count = count + 1
						executionLevel := node.(nodedata).Variables[0].Value
						execString = execString + "\n\n" + `if (flowData["ExecutionLevel"] == "` + executionLevel + `") {logger.Log_WF("Execution executionlevel: ` + executionLevel + `", logger.Debug, flowData["InSessionID"].(string));` + "\n" + `//initiate.View.` + view + `.` + executionLevel + `
										}`
						// if its the last round, it will add the delete function
						//if tot == count {
						//  execString = execString + "\n\n" + `if flowData["ExecutionLevel"] == "` + executionLevel + `" {logger.Log_WF("Deleting session details.", logger.Debug, flowData["InSessionID"].(string));Components.DeleteSession(flowData);}`
						//  }
		}

		// print the following
		ConstructGoFlow(sessionId, flowName, "init.executionLevels", execString, "INIT.hibernateValues", "initiate.View."+view, nodeDataObj, flowResult, flowStruct,makeExecutable)*/

		for stopFlow == false {

			targetNode, _, _ := linq.From(flowStruct.Connections).Where(
				func(in linq.T) (bool, error) { return in.(connection).SourceId == targetNodeId, nil }).First()

			targetNodeId = targetNode.(connection).TargetId

			nodeData, _, _ := linq.From(flowStruct.Nodes).Where(
				func(in linq.T) (bool, error) { return in.(nodedata).SchemaID == targetNodeId, nil }).First()

			targetLibraryId := nodeData.(nodedata).LibraryID
			//logger.Log_PE("Sibling node " + targetNode.(connection).TargetId + "   " + nodeData.(nodedata).Name)
			logger.Log_PE("Sibling node "+targetNode.(connection).TargetId+"  "+nodeData.(nodedata).Name, logger.Debug, sessionId)
			flowResult.Message = "Sibling node " + targetNode.(connection).TargetId + "  " + nodeData.(nodedata).Name

			for argumentsInitiated == false {

				argumentsInitiated = true
			}

			//ConstructGoFlow(sessionId, flowName, "initiateView", view, "initiate.View."+view, nodeData.(nodedata), flowResult, flowStruct,makeExecutable)
			if targetLibraryId == "2" {
				// adding if control

				logger.HighLight("-- Task: Add IF condition - "+nodeData.(nodedata).DisplayName, sessionId)

				ifCondition, _, _ := linq.From(flowStruct.Ifconditions).Where(
					func(in linq.T) (bool, error) { return in.(ifcondition).ID == targetNodeId, nil }).First()

				convertedIfcondition, _ := json.Marshal(ifCondition)
				logger.Log_PE("If Condition picked: "+string(convertedIfcondition), logger.Debug, sessionId)

				ifNode, _, _ := linq.From(flowStruct.Nodes).Where(
					func(in linq.T) (bool, error) { return in.(nodedata).SchemaID == ifCondition.(ifcondition).ID, nil }).First()
				logger.Log_PE("If Node picked", logger.Debug, sessionId)

				ifVariable, _, _ := linq.From(ifNode.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "ValueOne", nil }).First()
				logger.Log_PE("ValueOne picked: "+string(ifVariable.(variable).Value), logger.Debug, sessionId)

				ifCheck, _, _ := linq.From(ifNode.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "Condition", nil }).First()
				logger.Log_PE("Condition picked: "+string(ifCheck.(variable).Value), logger.Debug, sessionId)

				ifValue, _, _ := linq.From(ifNode.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "ValueTwo", nil }).First()
				logger.Log_PE("ValueTwo picked: "+string(ifValue.(variable).Value), logger.Debug, sessionId)

				logger.Log_PE("True node "+ifCondition.(ifcondition).True+"   False node "+ifCondition.(ifcondition).False, logger.Debug, sessionId)
				logger.Log_PE("If variable - "+ifVariable.(variable).Value+"   If condition - "+ifCheck.(variable).Value+"   If Value - "+ifValue.(variable).Value, logger.Debug, sessionId)

				if ifVariable.(variable).DataType != "string" {
					ConstructGoFlow(sessionId, flowName, "addimport", "strconv", "ImportStringConvert", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				}

				ConstructGoFlow(sessionId, flowName, "if", addIfCondition(ifVariable.(variable), ifCheck.(variable), ifValue.(variable), ifCondition.(ifcondition)), "AddIfConfition_"+ifNode.(nodedata).SchemaID, getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)

				if flowResult.Status == false {
					break
				}

				//targetNodeId = ifCondition.(ifcondition).True

			} else if targetLibraryId == "1" {

				stopFlow = true

			} else if targetLibraryId == "4" {
				// message control
				printMessage, _, _ := linq.From(nodeData.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "MessageBody", nil }).First()

				logger.HighLight("-- Task: Add print - "+printMessage.(variable).Value, sessionId)

				if checkIfConversionAvailable(nodeData.(nodedata)) == true {
					ConstructGoFlow(sessionId, flowName, "addimport", "strconv", "ImportStringConvert", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				}

				ConstructGoFlow(sessionId, flowName, "addimport", "processengine/Common", "ImportCommon", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				ConstructGoFlow(sessionId, flowName, "print", getMessageContent(printMessage.(variable)), "", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				if flowResult.Status == false {
					break
				}

			} else if targetLibraryId == "5" {
				// foreach control
				logger.HighLight("-- Task: Add Foreach control", sessionId)

				ForeachStatement, _, _ := linq.From(flowStruct.Forloops).Where(
					func(in linq.T) (bool, error) { return in.(foreachStuct).ID == targetNodeId, nil }).First()

				CollectionName, _, _ := linq.From(nodeData.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "CollectionName", nil }).First()

				viewUUID := nodeData.(nodedata).OtherData.ForeachUUID

				ConstructGoFlow(sessionId, flowName, "foreach", getForeachContent(viewUUID, nodeData.(nodedata), ForeachStatement.(foreachStuct), flowStruct), "ForeachCtrl"+CollectionName.(variable).Value, getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				if flowResult.Status == false {
					break
				}

			} else if targetLibraryId == "3" {
				// add hibernate control
				logger.HighLight("-- Task: Add Hibernate control - "+nodeData.(nodedata).DisplayName, sessionId)
				// if view != "drawboard" {
				// 	previous_executionlevel = executionlevel
				// 	executionlevel = nodeData.(nodedata).Variables[0].Value
				// 	ConstructGoFlow(sessionId, flowName, "addimport", "processengine/Components", "ImportComponents", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
				// 	ConstructGoFlow(sessionId, flowName, "init.executionLevels", getExecutionalFooter(executionlevel, flowName+".exe", view), "INIT.executionFooterfor."+executionlevel, getViewName(view, previous_executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				// 	logger.Log_PE("Early PEL: "+previous_executionlevel, logger.Debug, sessionId)
				// 	logger.Log_PE("Early EL: "+executionlevel, logger.Debug, sessionId)
				// 	//executionlevel = nodeData.(nodedata).Variables[0].Value
				// 	logger.Log_PE("Current PEL: "+previous_executionlevel, logger.Debug, sessionId)
				// 	logger.Log_PE("Current EL: "+executionlevel, logger.Debug, sessionId)
				// } else {
				previous_executionlevel = executionlevel
				ConstructGoFlow(sessionId, flowName, "addimport", "processengine/Components", "ImportComponents", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
				ConstructGoFlow(sessionId, flowName, "init.executionLevels", getExecutionalFooter(previous_executionlevel, flowName+".exe", view), "INIT.executionFooterfor."+previous_executionlevel, getViewName(view, previous_executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				logger.Log_PE("Early PEL: "+previous_executionlevel, logger.Debug, sessionId)
				logger.Log_PE("Early EL: "+executionlevel, logger.Debug, sessionId)
				executionlevel = nodeData.(nodedata).Variables[0].Value
				logger.Log_PE("Current PEL: "+previous_executionlevel, logger.Debug, sessionId)
				logger.Log_PE("Current EL: "+executionlevel, logger.Debug, sessionId)
				//	}

				if flowResult.Status == false {
					break
				}

			} else if targetLibraryId == "7" {
				// adding assign control
				fromControl, _, _ := linq.From(nodeData.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "Capture", nil }).First()

				toControl, _, _ := linq.From(nodeData.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "AssignTo", nil }).First()

				logger.HighLight("-- Task: Add Assign control - from "+fromControl.(variable).Value+" to "+toControl.(variable).Value, sessionId)
				ConstructGoFlow(sessionId, flowName, "assign", getAssignmentControlDetails(fromControl.(variable), toControl.(variable)), "AssignCtrl"+fromControl.(variable).Value+toControl.(variable).Value, getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				if flowResult.Status == false {
					break
				}
			} else if targetLibraryId == "8" {
				// adding switch statement
				logger.HighLight("-- Task: Add Switch control", sessionId)

				SwitchStatement, _, _ := linq.From(flowStruct.Switches).Where(
					func(in linq.T) (bool, error) { return in.(switchStuct).ID == targetNodeId, nil }).First()

				VariableObj, _, _ := linq.From(nodeData.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "Variable", nil }).First()

				DataTypeObj, _, _ := linq.From(nodeData.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "DataType", nil }).First()

				viewUUID := nodeData.(nodedata).OtherData.SwitchUUID

				if DataTypeObj.(variable).Value != "string" {
					ConstructGoFlow(sessionId, flowName, "addimport", "strconv", "ImportStringConvert", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				}

				ConstructGoFlow(sessionId, flowName, "switch", getSwitchConditionContent(viewUUID, nodeData.(nodedata), SwitchStatement.(switchStuct), flowStruct), "SwitchCtrl"+VariableObj.(variable).Value, getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				if flowResult.Status == false {
					break
				}

			} else if targetLibraryId == "9" {
				// adding case
				logger.HighLight("-- Task: Add Case control", sessionId)

				CaseNode, _, _ := linq.From(flowStruct.Nodes).Where(
					func(in linq.T) (bool, error) { return in.(nodedata).SchemaID == targetNodeId, nil }).First()

				CaseValue, _, _ := linq.From(nodeData.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "CaseValue", nil }).First()

				cstate := CaseNode.(nodedata).OtherData.CaseUUID
				ConstructGoFlow(sessionId, flowName, "case", getCaseContent(cstate, sessionId, flowStruct), "CaseCtrl"+CaseValue.(variable).Value, getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				if flowResult.Status == false {
					break
				}

			} else if targetLibraryId == "10" {
				// adding default case
				logger.HighLight("-- Task: Add Default Case control", sessionId)

				CaseNode, _, _ := linq.From(flowStruct.Nodes).Where(
					func(in linq.T) (bool, error) { return in.(nodedata).SchemaID == targetNodeId, nil }).First()

				cstate := CaseNode.(nodedata).OtherData.DefaultUUID
				ConstructGoFlow(sessionId, flowName, "case", getDefaultCaseContent(cstate, sessionId, flowStruct), "DefaultCaseCtrl", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				if flowResult.Status == false {
					break
				}

			} else if targetLibraryId == "12" {
				// adding default case
				logger.HighLight("-- Task: Termination control", sessionId)

				termincationNode, _, _ := linq.From(flowStruct.Nodes).Where(
					func(in linq.T) (bool, error) { return in.(nodedata).SchemaID == targetNodeId, nil }).First()

				reason, _, _ := linq.From(termincationNode.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "Reason", nil }).First()

				status, _, _ := linq.From(termincationNode.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "Status", nil }).First()

				ConstructGoFlow(sessionId, flowName, "terminate", getFlowTerminationContent("terminate", reason.(variable).Value, status.(variable).Value, "6"), "TerminateControl", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				if makeExecutable == false {
					ConstructGoFlow(sessionId, flowName, "terminate", `fmt.Print(string(returnDataJSON));return;`, "PRINTJSONDATA", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				} else {
					ConstructGoFlow(sessionId, flowName, "terminate", `return string(returnDataJSON)`, "PRINTJSONDATA", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				}
				if flowResult.Status == false {
					break
				}
				stopFlow = true

			} else if targetLibraryId == "13" {
				// adding switch statement
				logger.HighLight("-- Task: Add While loop control", sessionId)

				fmt.Println("While Loops: ", flowStruct.WhileLoops)
				fmt.Println("targetNode id: ", targetNodeId)

				WhileStatement, _, _ := linq.From(flowStruct.WhileLoops).Where(
					func(in linq.T) (bool, error) { return in.(whileStuct).ID == targetNodeId, nil }).First()

				fmt.Println("While Statement: ", WhileStatement)

				VariableOne, _, _ := linq.From(nodeData.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "ValueOne", nil }).First()

				fmt.Println("ValueOne: ", VariableOne)

				/*	ConditionObj, _, _ := linq.From(nodeData.(nodedata).Variables).Where(
						func(in linq.T) (bool, error) { return in.(variable).Key == "Condition", nil }).First()

					VariableTwo, _, _ := linq.From(nodeData.(nodedata).Variables).Where(
						func(in linq.T) (bool, error) { return in.(variable).Key == "VariableTwo", nil }).First()

						viewUUID := nodeData.(nodedata).OtherData.WhileUUID*/

				if VariableOne.(variable).Value != "string" {
					ConstructGoFlow(sessionId, flowName, "addimport", "strconv", "ImportStringConvert", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				}

				ConstructGoFlow(sessionId, flowName, "while", getWhileLoopContent(nodeData.(nodedata), WhileStatement.(whileStuct), flowStruct), "WhileCtrl"+WhileStatement.(whileStuct).WhileState, getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				if flowResult.Status == false {
					break
				}

			} else if targetLibraryId == "14" {
				// New Collection control
				logger.HighLight("-- Task: Add New Collection control", sessionId)

				Collection, _, _ := linq.From(nodeData.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "CollectionName", nil }).First()

				ConstructGoFlow(sessionId, flowName, "newcollection", getCollectionData(Collection.(variable)), "NewCollectionCtrl", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				if flowResult.Status == false {
					break
				}

			} else if targetLibraryId == "15" {
				// Add to Collection control
				logger.HighLight("-- Task: Add 'Add to Collection' control", sessionId)

				Collection, _, _ := linq.From(nodeData.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "CollectionName", nil }).First()

				CollValue, _, _ := linq.From(nodeData.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "Value", nil }).First()

				ConstructGoFlow(sessionId, flowName, "newcollection", getAddtoCollectionContent(Collection.(variable), CollValue.(variable)), "AddToCollectionCtrl", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				if flowResult.Status == false {
					break
				}

			} else if targetLibraryId == "16" {
				// Add to Collection control
				logger.HighLight("-- Task: Add 'Remove from Collection' control", sessionId)

				Collection, _, _ := linq.From(nodeData.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "CollectionName", nil }).First()

				CollValue, _, _ := linq.From(nodeData.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "Value", nil }).First()

				ConstructGoFlow(sessionId, flowName, "newcollection", getRemoveFromCollectionContent(Collection.(variable), CollValue.(variable)), "AddRemoveFromCollectionCtrl", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				if flowResult.Status == false {
					break
				}

			} else if targetLibraryId == "17" {
				// Add to Collection control
				logger.HighLight("-- Task: Add 'Remove Collection' control", sessionId)

				Collection, _, _ := linq.From(nodeData.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "CollectionName", nil }).First()

				ConstructGoFlow(sessionId, flowName, "newcollection", Collection.(variable).Value+"=nil", "RemovingCollectionCtrl", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				if flowResult.Status == false {
					break
				}

			} else if targetLibraryId == "18" {
				// Add Calculation control
				logger.HighLight("-- Task: Add 'Calculation' control", sessionId)

				Expression, _, _ := linq.From(nodeData.(nodedata).Variables).Where(
					func(in linq.T) (bool, error) { return in.(variable).Key == "Expression", nil }).First()

				ConstructGoFlow(sessionId, flowName, "calculation", Expression.(variable).Value, "CalculationCtrl", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				if flowResult.Status == false {
					break
				}

			} else {

				logger.HighLight("-- Task: Add activity - "+nodeData.(nodedata).Name, sessionId)
				ConstructGoFlow(sessionId, flowName, "addimport", "processengine/context", "ImportCONTEXT", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				ConstructGoFlow(sessionId, flowName, "addimport", "processengine/Activities/"+nodeData.(nodedata).Name, "Import"+nodeData.(nodedata).Name, getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable) //
				if checkIfConversionAvailable(nodeData.(nodedata)) == true {
					ConstructGoFlow(sessionId, flowName, "addimport", "strconv", "ImportStringConvert", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				}
				if makeExecutable == true {
					ConstructGoFlow(sessionId, flowName, "addimport", "time", "ImportTime", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				}
				ConstructGoFlow(sessionId, flowName, "init.activity", nodeData.(nodedata).Name, "", getViewName(view, executionlevel), nodeData.(nodedata), flowResult, flowStruct, makeExecutable)
				//logger.Log_PE(flowResult.Status)
				//logger.Log_PE("flow message" + flowResult.Message)
				/*if flowResult.Status == false {
				break
				}*/
			}
		}
	}

	logger.HighLight("-- Task: Finishing workflow", sessionId)
	logger.Log_PE("Flow has met with stop node", logger.Debug, sessionId)
	flowResult.Message = "Flow has met with stop node"

	//nodeDataObj := new(nodedata)
	ConstructGoFlow(sessionId, flowName, "addimport", "fmt", "ImportFMT", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
	//ConstructGoFlow(sessionId, flowName, "addimport", "Components", "ImportComponents", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct,makeExecutable)
	//ConstructGoFlow(sessionId, flowName, "init.end", `jsonOut, MarshalError := json.Marshal(flowData); if MarshalError != nil {logger.Log_WF("JSONOut marshal error", logger.Debug, flowData["InSessionID"].(string));WFContext.Message = "marshal error";WFContext.Status = false;WFContext.ErrorCode = 4}`, "MarshalJSONOutput", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct,makeExecutable)
	//ConstructGoFlow(sessionId, flowName, "init.end", `wfOut, wfmarshalError :=json.Marshal(ALLOutArguments);if wfmarshalError != nil {logger.Log_WF("WF result marshal error", logger.Debug, flowData["InSessionID"].(string));WFContext.Message = "WF result marshal error";WFContext.Status = false;WFContext.ErrorCode = 4}`, "MarshalWFOutput", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct,makeExecutable)
	//ConstructGoFlow(sessionId, flowName, "init.end", "WFContext.ActivityTrace = activityTrace", "ASSIGNINGTractTOActivityOBJ", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct,makeExecutable)
	//ConstructGoFlow(sessionId, flowName, "init.end", `logger.Log_WF("", logger.Debug, flowData["InSessionID"].(string));logger.Log_WF("JSONOutput: "+string(jsonOut), logger.Debug, flowData["InSessionID"].(string))`, "AddOutputJson", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct,makeExecutable)
	//ConstructGoFlow(sessionId, flowName, "init.end", `logger.Log_WF("", logger.Debug, flowData["InSessionID"].(string));logger.Log_WF("WFOutArguments: "+string(wfOut), logger.Debug, flowData["InSessionID"].(string))`, "AddWFOutputJson", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct,makeExecutable)
	ConstructGoFlow(sessionId, flowName, "init.end", `logger.Log_WF("WORKFLOW COMPLETED!", logger.Debug, flowData["InSessionID"].(string))`, "AddWFCompleted", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
	ConstructGoFlow(sessionId, flowName, "init.end", getFlowTerminationContent("flowend", "", "", ""), "ConvertWFResultToJSON", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
	//returnObj.WorkflowTrace = Common.GetWFLog(flowData["InSessionID"].(string));returnObj.WorkflowResult.ActivityTrace = activityTrace;
	if makeExecutable == false {
		ConstructGoFlow(sessionId, flowName, "init.end", `fmt.Print(string(returnDataJSON));return;`, "PRINTJSONDATA", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
	} else {
		ConstructGoFlow(sessionId, flowName, "init.end", `return string(returnDataJSON)`, "PRINTJSONDATA", "initiate.View.drawboard", nodeDataObj, flowResult, flowStruct, makeExecutable)
	}

	isSuccess, exeLocation := buildOrRunGoFlow(flowStruct, flowName, sessionId, makeExecutable, flowResult)
	if isSuccess == true {
		msg := "Processing completed. Workflow built successfully"
		logger.Log_PE(msg, logger.Debug, sessionId)

		flowResult.Message = msg
		wfcontext.Message = exeLocation
		wfcontext.Status = true
	} else {
		msg := "Processing completed with an error building the workflow"
		logger.Log_PE(msg, logger.Debug, sessionId)

		flowResult.Message = msg
		wfcontext.Message = exeLocation
		wfcontext.Status = false
	}

	//flowResult.ReturnData.WorkflowTrace = Common.GetPELog(sessionId)
	flowResult.ReturnData.WorkflowResult.ErrorCode = wfcontext.ErrorCode
	flowResult.ReturnData.WorkflowResult.Status = wfcontext.Status
	flowResult.ReturnData.WorkflowResult.Message = wfcontext.Message
	flowResult.Status = wfcontext.Status

	return flowResult
}

func ConstructGoFlow(sessionId, flowName, functionName, content, validateKey, view string, nodeData nodedata, flowResult *context.FlowResult, flowStruct JsonFlow, makeExecutable bool) {

	logger.Log_PE("~ Constructing flow", logger.Debug, sessionId)

	var slash = ""
	if runtime.GOOS == "windows" {
		slash = "\\"
	} else {
		slash = "/"
	}

	pwd, _ := os.Getwd()
	buildPath := ""
	initFile := ""
	//activityPath := ".." + slash + "Activities"

	if makeExecutable == false {
		initFile = "Init.go"
		buildPath = pwd + slash + "BuiltFlows"
	} else {
		initFile = "Exec.go"
		buildPath = pwd + slash + "BuiltExecutables"
	}

	checkPath := buildPath + slash + flowName + slash + flowName + ".go"

	var srcFile []byte
	existFlowFile := IsExists(checkPath)
	//logger.Log_PE("Checking location: " + buildPath + slash + flowName + ".go")
	if existFlowFile == true {
		//logger.Log_PE("WF file already avaiable. Using it to proceed.")
		goFile, init_error := ioutil.ReadFile(buildPath + slash + flowName + slash + flowName + ".go")
		if init_error != nil {
			println(init_error.Error())
			flowResult.Message = init_error.Error()
			flowResult.Status = false
			return
		} else {
			srcFile = goFile
			//logger.Log_PE("Modifying existing file " + flowName + ".go")
			flowResult.Message = "Modifying existing file " + flowName + ".go"
		}
	} else {
		//logger.Log_PE("WF file created using Init.go")
		goFile, init_error := ioutil.ReadFile(pwd + slash + initFile)
		if init_error != nil {
			println(init_error.Error())
			flowResult.Message = init_error.Error()
			flowResult.Status = false
			return
		} else {
			srcFile = goFile
			//logger.Log_PE("Building new file " + flowName + ".go with existing " + initFile + " initial file")
			flowResult.Message = "Building new file " + flowName + ".go with existing " + initFile + " initial file"
		}
	}

	src := string(srcFile)
	replaceWith := "//Func.Body"

	if functionName == "initiateView" {
		logger.Log_PE("Initiate view: "+content, logger.Debug, sessionId)
		content = `//` + validateKey + `
												//` + view
		replaceWith = "//" + view
	} else if functionName == "init.arguments" {
		logger.Log_PE("Initiate arguments: "+content, logger.Debug, sessionId)

		addVariable := contentAlreadyAvailable(src, validateKey)

		if addVariable != true {
			arguments, _ := linq.From(flowStruct.Arguments).Where(func(in linq.T) (bool, error) {
				return in.(variable).Category == content, nil
			}).Results()
			convertedArguments, _ := json.Marshal(arguments)
			logger.Log_PE(string(convertedArguments), logger.Debug, sessionId)
			// validate if all INArguments are available in flowData
			// write it when it is published
			concatinatedString := `// initiating ` + content + ` arguments
			`
			if content == "InArgument" {
				for _, argument := range arguments {
					// if the arguent is hardcorded
					if argument.(variable).Type == "hardcoded" {
						if argument.(variable).DataType == "string" {
							concatinatedString = concatinatedString + `if _, ok := flowData["` + argument.(variable).Key + `"]; ok {
									flowData["` + argument.(variable).Key + `"] = flowData["` + argument.(variable).Key + `"]
								} else {
									flowData["` + argument.(variable).Key + `"] = "` + argument.(variable).Value + `"
								}
								`
							//concatinatedString = concatinatedString + `flowData["` + argument.(variable).Key + `"] = "` + argument.(variable).Value + `";`
						} else if argument.(variable).DataType == "int" {
							concatinatedString = concatinatedString + `if _, ok := flowData["` + argument.(variable).Key + `"]; ok {
									flowData["` + argument.(variable).Key + `"] = flowData["` + argument.(variable).Key + `"]
								} else {
									flowData["` + argument.(variable).Key + `"] = int64(` + argument.(variable).Value + `);
								}
								`
							//concatinatedString = concatinatedString + `flowData["` + argument.(variable).Key + `"] = int64(` + argument.(variable).Value + `);`
						} else if argument.(variable).DataType == "float" {
							concatinatedString = concatinatedString + `if _, ok := flowData["` + argument.(variable).Key + `"]; ok {
									flowData["` + argument.(variable).Key + `"] = flowData["` + argument.(variable).Key + `"]
								} else {
									flowData["` + argument.(variable).Key + `"] = float64(` + argument.(variable).Value + `);
								}
								`
							//concatinatedString = concatinatedString + `flowData["` + argument.(variable).Key + `"] = float64(` + argument.(variable).Value + `);`
						} else if argument.(variable).DataType == "boolean" {
							concatinatedString = concatinatedString + `if _, ok := flowData["` + argument.(variable).Key + `"]; ok {
									flowData["` + argument.(variable).Key + `"] = flowData["` + argument.(variable).Key + `"]
								} else {
									flowData["` + argument.(variable).Key + `"] = bool(` + argument.(variable).Value + `);
								}
								`
							//concatinatedString = concatinatedString + `flowData["` + argument.(variable).Key + `"] = bool(` + argument.(variable).Value + `);`
						}
					}
					if argument.(variable).Type == "dynamic" {
						if argument.(variable).DataType == "string" {
							concatinatedString = concatinatedString + `flowData["` + argument.(variable).Key + `"] = flowData["` + argument.(variable).Key + `"].(string);`
						} else if argument.(variable).DataType == "int" {
							concatinatedString = concatinatedString + `val` + argument.(variable).Key + `, _ := strconv.ParseInt(flowData["` + argument.(variable).Key + `"].(string), 10, 64);
							flowData["` + argument.(variable).Key + `"] = val` + argument.(variable).Key + `;`
						} else if argument.(variable).DataType == "float" {
							concatinatedString = concatinatedString + `val` + argument.(variable).Key + `, _ := strconv.ParseFloat(flowData["` + argument.(variable).Key + `"].(string), 64)
							flowData["` + argument.(variable).Key + `"] = val` + argument.(variable).Key + `;`
						} else if argument.(variable).DataType == "boolean" {
							concatinatedString = concatinatedString + `val` + argument.(variable).Key + `, _ := strconv.ParseBool(flowData["` + argument.(variable).Key + `"].(string))
							flowData["` + argument.(variable).Key + `"] = val` + argument.(variable).Key + `;`
						}
						/*if argument.(variable).Value != "" {
							concatinatedString = concatinatedString + `flowData["` + argument.(variable).Key + `"] = flowData["` + argument.(variable).Key + `"];`
						} else {
							concatinatedString = concatinatedString + `flowData["` + argument.(variable).Key + `"] = flowData["` + argument.(variable).Key + `"];`
						}*/
					}
				}
				concatinatedString = concatinatedString + getArrayDataInitiation(arguments)
			} else if content == "Custom" {
				for _, argument := range arguments {
					// if the arguent is hardcorded
					if argument.(variable).Type == "hardcoded" {
						if argument.(variable).DataType == "string" {
							concatinatedString = concatinatedString + `flowData["` + argument.(variable).Key + `"] = "` + argument.(variable).Value + `";`
						} else if argument.(variable).DataType == "int" {
							concatinatedString = concatinatedString + `flowData["` + argument.(variable).Key + `"] = int64(` + argument.(variable).Value + `);`
						} else if argument.(variable).DataType == "float" {
							concatinatedString = concatinatedString + `flowData["` + argument.(variable).Key + `"] = float64(` + argument.(variable).Value + `);`
						} else if argument.(variable).DataType == "boolean" {
							concatinatedString = concatinatedString + `flowData["` + argument.(variable).Key + `"] = bool(` + argument.(variable).Value + `);`
						}
					}
					if argument.(variable).Type == "dynamic" {
						if argument.(variable).Value != "" {
							concatinatedString = concatinatedString + `flowData["` + argument.(variable).Key + `"] = flowData["` + argument.(variable).Value + `"];`
						} else {
							if argument.(variable).DataType == "string" {
								concatinatedString = concatinatedString + `flowData["` + argument.(variable).Key + `"] = "";`
							} else if argument.(variable).DataType == "int" {
								concatinatedString = concatinatedString + `flowData["` + argument.(variable).Key + `"] = 0;`
							} else if argument.(variable).DataType == "float" {
								concatinatedString = concatinatedString + `flowData["` + argument.(variable).Key + `"] = 0.0;`
							} else if argument.(variable).DataType == "boolean" {
								concatinatedString = concatinatedString + `flowData["` + argument.(variable).Key + `"] = false;`
							}
						}
					}
				}
				concatinatedString = concatinatedString + getArrayDataInitiation(arguments)
			} else if content == "OutArgument" {
				for _, argument := range arguments {
					// if the arguent is hardcorded
					if argument.(variable).Type == "hardcoded" {
						concatinatedString = concatinatedString + `ALLOutArguments["` + argument.(variable).Key + `"] = "` + argument.(variable).Value + `";`
					}
					if argument.(variable).Type == "dynamic" {
						if argument.(variable).Value != "" {
							concatinatedString = concatinatedString + `ALLOutArguments["` + argument.(variable).Key + `"] = flowData["` + argument.(variable).Value + `"];`
						} else {
							if argument.(variable).DataType == "string" {
								concatinatedString = concatinatedString + `ALLOutArguments["` + argument.(variable).Key + `"] = "";`
							} else if argument.(variable).DataType == "int" {
								concatinatedString = concatinatedString + `ALLOutArguments["` + argument.(variable).Key + `"] = 0;`
							} else if argument.(variable).DataType == "float" {
								concatinatedString = concatinatedString + `ALLOutArguments["` + argument.(variable).Key + `"] = 0.0;`
							} else if argument.(variable).DataType == "boolean" {
								concatinatedString = concatinatedString + `ALLOutArguments["` + argument.(variable).Key + `"] = false;`
							}
						}

					}
				}
				concatinatedString = concatinatedString + getArrayDataInitiation(arguments)
			}
			content = concatinatedString + "//" + validateKey + `
												//Init.Var`
			replaceWith = "//Init.Var"
		} else {
			logger.Log_PE("Arguments already added", logger.Debug, sessionId)
			replaceWith = ""
		}
	} else if functionName == "dlExec" {
		logger.Log_PE("Setting: "+content, logger.Debug, sessionId)
		content = content + " //" + validateKey + "\n" + "//Init.Settings"
		replaceWith = "//Init.Settings"
		logger.Log_PE("Execution variables added", logger.Debug, sessionId)
	} else if functionName == "init.executionLevels" {
		logger.Log_PE("Adding Execution variables: "+content, logger.Debug, sessionId)
		content = content + " //" + validateKey + "\n" + "//" + view
		replaceWith = "//" + view
		logger.Log_PE("Execution variables added", logger.Debug, sessionId)

	} else if functionName == "init.var" {
		logger.Log_PE("Adding new Variable: "+content, logger.Debug, sessionId)
		// validate and check if the variable is already available in the code
		addVariable := contentAlreadyAvailable(src, validateKey)
		// add the import or else skip the function
		if addVariable != true {
			logger.Log_PE("New variable added", logger.Debug, sessionId)
			flowResult.Message = "Adding new variable"
			content = content + " " + "//" + validateKey + "\n" + "//Init.Var"
			replaceWith = "//Init.Var"
		} else {
			logger.Log_PE("Variable already added", logger.Debug, sessionId)
			replaceWith = ""
		}

	} else if functionName == "printInfoName" { //Added by Jay

		tokens := strings.Split(content, "_")
		if len(tokens) > 2 {
			content = ""
			for x := 1; x < len(tokens); x++ {
				content += tokens[x] + " "
			}
		} else {
			if len(tokens) == 2 {
				content = tokens[1]
			}
		}

		logger.Log_PE("Adding new Print: "+content, logger.Debug, sessionId)

		replaceWith = "//INIT.FlowName"
		// removed a code because it was unnecessary
	} else if functionName == "printInfoNameForCEB" { //Added by Jay
		logger.Log_PE("Adding new Print: "+content, logger.Debug, sessionId)
		replaceWith = "//INIT.FlowNameCEB"
		// removed a code because it was unnecessary
	} else if functionName == "printInfoBuild" { //Added by Jay
		logger.Log_PE("Adding new Print: "+content, logger.Debug, sessionId)
		// validate import, see if it already exist
		addFlowPort := contentAlreadyAvailable(src, validateKey)
		//addImport := false
		// add the import or else skip the function
		if addFlowPort != true {
			logger.Log_PE("New print added", logger.Debug, sessionId)
			//flowResult.Message = "Adding new imports"
			content += "\") //INIT.FLOWBUILDPLATFORM"
			replaceWith = "//INIT.FlowBuildPlatform\")"
		} else {
			logger.Log_PE("Print already added", logger.Debug, sessionId)
			replaceWith = ""
		}
	} else if functionName == "printInArugements" { //Added by Jay
		logger.Log_PE("Adding new Print: "+content, logger.Debug, sessionId)
		// validate import, see if it already exist
		addFlowPort := contentAlreadyAvailable(src, validateKey)
		//addImport := false
		// add the import or else skip the function
		if addFlowPort != true {
			logger.Log_PE("New print added", logger.Debug, sessionId)
			//flowResult.Message = "Adding new imports"
			content += "\") //INIT.PRINTINARGUMENTS"
			replaceWith = "//INIT.PrintInArguments\")"
		} else {
			logger.Log_PE("Print already added", logger.Debug, sessionId)
			replaceWith = ""
		}
	} else if functionName == "printOutArugements" { //Added by Jay
		logger.Log_PE("Adding new Print: "+content, logger.Debug, sessionId)
		// validate import, see if it already exist
		addFlowPort := contentAlreadyAvailable(src, validateKey)
		//addImport := false
		// add the import or else skip the function
		if addFlowPort != true {
			logger.Log_PE("New print added", logger.Debug, sessionId)
			//flowResult.Message = "Adding new imports"
			content += "\") //INIT.PRINTOUTARGUMENTS"
			replaceWith = "//INIT.PrintOutArguments\")"
		} else {
			logger.Log_PE("Print already added", logger.Debug, sessionId)
			replaceWith = ""
		}
	} else if functionName == "printWfBodySample" { //Added by Jay
		logger.Log_PE("Adding new Print: "+content, logger.Debug, sessionId)
		// validate import, see if it already exist
		addFlowPort := contentAlreadyAvailable(src, validateKey)
		//addImport := false
		// add the import or else skip the function
		if addFlowPort != true {
			logger.Log_PE("New print added", logger.Debug, sessionId)
			//flowResult.Message = "Adding new imports"
			content += "\") //INIT.PRINTWFBODYSAMPLE"
			replaceWith = "//INIT.PrintWfBodySample\")"
		} else {
			logger.Log_PE("Print already added", logger.Debug, sessionId)
			replaceWith = ""
		}
	} else if functionName == "addimport" {
		logger.Log_PE("Adding new Import: "+content, logger.Debug, sessionId)
		// validate import, see if it already exist
		addImport := contentAlreadyAvailable(src, validateKey)
		// add the import or else skip the function
		if addImport != true {
			logger.Log_PE("New import added", logger.Debug, sessionId)
			flowResult.Message = "Adding new imports"
			content = "import \"" + content + "\" //" + validateKey + "\n" + "//Import.End"
			replaceWith = "//Import.End"
		} else {
			logger.Log_PE("Import already added", logger.Debug, sessionId)
			replaceWith = ""
		}
	} else if functionName == "init.end" {
		logger.Log_PE("Adding end content: "+content, logger.Debug, sessionId)
		flowResult.Message = "Adding new end variable"
		content = content + " //" + validateKey + "\n" + "//Func.End"
		replaceWith = "//Func.End"
	} else if functionName == "init.headers" {
		logger.Log_PE("Adding header information: "+content, logger.Debug, sessionId)
		content = content + " //" + validateKey + "\n"
		replaceWith = "//INIT.Headers"
	} else if functionName == "if" {
		logger.Log_PE("Adding if content: "+content, logger.Debug, sessionId)
		content = content + " //" + validateKey + "\n" + "//" + view
		replaceWith = "//" + view
	} else if functionName == "print" {
		logger.Log_PE("Adding new print activity: "+content, logger.Debug, sessionId)
		flowResult.Message = "Adding print activity"
		content = content + `
					//` + view
		replaceWith = "//" + view
	} else if functionName == "assign" {
		logger.Log_PE("Adding assign content: "+content, logger.Debug, sessionId)
		content = content + " //" + validateKey + "\n" + "//" + view
		replaceWith = "//" + view
	} else if functionName == "switch" {
		logger.Log_PE("Adding switch statement", logger.Debug, sessionId)
		content = content + " //" + validateKey + "\n" + "//" + view
		replaceWith = "//" + view // switch
	} else if functionName == "case" {
		logger.Log_PE("Adding case statement", logger.Debug, sessionId)
		content = content + " //" + validateKey + "\n" + "//" + view
		replaceWith = "//" + view // switch
	} else if functionName == "terminate" {
		logger.Log_PE("Adding termination control", logger.Debug, sessionId)
		content = content + " //" + validateKey + "\n" + "//" + view
		replaceWith = "//" + view // switch
	} else if functionName == "while" {
		logger.Log_PE("Adding while control", logger.Debug, sessionId)
		content = content + " //" + validateKey + "\n" + "//" + view
		replaceWith = "//" + view // switch
	} else if functionName == "newcollection" {
		logger.Log_PE("Adding New Collection variable", logger.Debug, sessionId)
		content = content + " //" + validateKey + "\n" + "//" + view
		replaceWith = "//" + view // collection
	} else if functionName == "calculation" {
		logger.Log_PE("Adding New Calculation control", logger.Debug, sessionId)
		content = content + " //" + validateKey + "\n" + "//" + view
		replaceWith = "//" + view // collection
	} else if functionName == "foreach" {
		logger.Log_PE("Adding Foreach variable", logger.Debug, sessionId)
		content = content + " //" + validateKey + "\n" + "//" + view
		replaceWith = "//" + view //
		//foreach
	} else if functionName == "init.activity" {
		logger.Log_PE("Adding activity: "+content, logger.Debug, sessionId)
		//fset := token.NewFileSet()
		filename := content
		//activityFile, _ := parser.ParseFile(fset, activityPath+slash+content+slash+content+".go", nil, 0)

		flowResult.Message = "Adding custom activity"

		invokeStatusPropertyName := randStr(16, "alpha")
		activityResultPropertyName := randStr(16, "alpha")
		timeVariable := randStr(16, "alpha")

		var variableName = "ActivityData_" + invokeStatusPropertyName
		activityINArguments := getINArgumentString(nodeData, variableName, flowStruct, makeExecutable)
		activityOUTArguments := getOutArgumentString(nodeData, variableName)

		// if the built is an executable, add timing to be printed on the workflow with time
		var timeNow = ""
		var timeSince = ""
		if makeExecutable == true {
			timeNow = `fmt.Println("Starting Activity: ` + filename + `");` + timeVariable + `:= time.Now()`
			timeSince = `fmt.Println("Completed Activity in: ",time.Since(` + timeVariable + `));fmt.Println("");`
		}

		content = `
			logger.Log_WF("", logger.Debug, flowData["InSessionID"].(string))
			logger.Log_WF("Invoking activity: ` + filename + `", logger.Debug, flowData["InSessionID"].(string))
			logger.Log_WF("Session ID: "+flowData["InSessionID"].(string), logger.Debug, flowData["InSessionID"].(string))

			` + activityINArguments + `

			` + invokeStatusPropertyName + ` := make(chan bool)
			var ` + activityResultPropertyName + ` = new(context.ActivityContext)
			go func() {

				` + timeNow + `
				` + variableName + `, ` + activityResultPropertyName + ` = ` + content + `.Invoke(` + variableName + `)
				` + invokeStatusPropertyName + ` <- ` + activityResultPropertyName + `.ActivityStatus
				` + timeSince + `

				` + activityOUTArguments + `
				}()

				logger.Log_WF("Activity invoked and waiting for completion ... ", logger.Debug, flowData["InSessionID"].(string))
				<-` + invokeStatusPropertyName + `

				logger.Log_WF("Activity response was recieved.", logger.Debug, flowData["InSessionID"].(string))

				if(` + activityResultPropertyName + `.ActivityStatus == true){
					logger.Log_WF("Success - "+` + activityResultPropertyName + `.Message, logger.Debug, flowData["InSessionID"].(string))
					WFContext.Message = ` + activityResultPropertyName + `.Message
					WFContext.Status = true
					WFContext.ErrorCode = 3
				}else{
																																																												//logger.Log_WF("Failed - "+OfUiVrIQAzqXQMov.ErrorState.ErrorString, logger.Debug, flowData["InSessionID"].(string))
					logger.Log_WF("Failed - "+` + activityResultPropertyName + `.Message, logger.Debug, flowData["InSessionID"].(string))
					WFContext.Message = ` + activityResultPropertyName + `.ErrorState.ErrorString
					WFContext.Status = false
					WFContext.ErrorCode = 2
				}

																																//` + view
		replaceWith = "//" + view
		logger.Log_PE("New activity added", logger.Debug, sessionId)
		// the following was removed from the above code because to make the logs be retrived later saparately
		//` + activityFile.Name.String() + `_` + activityResultPropertyName + ` := Common.GetActLog(flowData["InSessionID"].(string))
		//  activityTrace["` + activityFile.Name.String() + `"] = ` + activityFile.Name.String() + `_` + activityResultPropertyName + `
	}

	if replaceWith != "" {
		src = strings.Replace(src, replaceWith, content, 1)
	}
	flowResult.Message = "Constructing go flow"
	contructGoFile(src, flowName, flowResult, makeExecutable)
}

func contructGoFile(goCode, flowName string, flowResult *context.FlowResult, makeExecutable bool) {

	//logger.Log_PE("")
	//logger.Log_PE("~ Contructing Go file")
	var slash = ""
	if runtime.GOOS == "windows" {
		slash = "\\"
	} else {
		slash = "/"
	}

	pwd, _ := os.Getwd()
	buildPath := ""

	if makeExecutable == false {
		buildPath = pwd + slash + "BuiltFlows"
	} else {
		buildPath = pwd + slash + "BuiltExecutables"
	}

	// check if BuiltFlows or Executable folder exists or not, create one if not available
	_, err := os.Stat(buildPath)
	if err != nil {
		// create folder in the given path and permissions
		os.Mkdir(buildPath, 0777)
	}

	WFFolder := buildPath + slash + flowName
	_, wferr := os.Stat(WFFolder)
	if wferr != nil {
		// create folder in the given path and permissions
		os.Mkdir(WFFolder, 0777)
	}

	// save the go gile in the given path
	ioutil.WriteFile(WFFolder+slash+flowName+".go", []byte(goCode), 0777)

	//logger.Log_PE("GoFlow saved: " + buildPath + slash + flowName + ".go")
	//flowResult.Message = flowResult.Message + "Go file saved and directory " + buildPath + " -> "
}

func buildOrRunGoFlow(flowStruct JsonFlow, flowName, sessionId string, makeExecutable bool, flowResult *context.FlowResult) (bool, string) {

	// dettermine the slash type of the running OS
	var slash = ""
	if runtime.GOOS == "windows" {
		slash = "\\"
	} else {
		slash = "/"
	}

	pwd, _ := os.Getwd()
	buildPath := ""

	fmt.Println("MAKE EXECUTABLE : ")
	fmt.Println(makeExecutable)

	logger.Log_PE("~ Building flow", logger.Information, sessionId)
	logger.Log_PE("Runtime Platform: Environment - "+string(runtime.GOOS), logger.Debug, sessionId)

	if makeExecutable == false {
		//set the folder path for if the normal build is happened to the cloud environment
		buildPath = pwd + slash + "BuiltFlows"
		flowResult.FlowName = flowName
		flowResult.Status = true
		cmd, err := exec.Command("go", "build", buildPath+slash+flowName+slash+flowName+".go").CombinedOutput()
		fmt.Print(string(cmd))
		if err == nil {
			logger.Log_PE("No error building flow", logger.Debug, sessionId)
			logger.Log_PE(string(cmd), logger.Debug, sessionId)
			return true, ""
		} else {
			logger.Log_PE("Error building flow", logger.Error, sessionId)
			logger.Log_PE(string(err.Error()), logger.Error, sessionId)
			return false, ""
		}
	} else {
		// set the folder path if the building process is for the standalone system type
		buildPath = pwd + slash + "BuiltExecutables"
		//Port := flowStruct.Port
		OSCode := flowStruct.OSCode
		SysArch := flowStruct.SysArch
		filefolder := buildPath + slash + flowName + slash
		fullfilePath := buildPath + slash + flowName + slash + flowName + ".go"

		logger.Log_PE("OSCode: "+OSCode, logger.Debug, sessionId)
		logger.Log_PE("SysArch: "+SysArch, logger.Debug, sessionId)
		logger.Log_PE("FileFolder: "+filefolder, logger.Debug, sessionId)
		logger.Log_PE("FilePath:"+fullfilePath, logger.Debug, sessionId)
		//set GOARCH=amd64
		//set GOOS=linux

		architecture := SysArch
		if SysArch == "86zip" || SysArch == "386zip" {
			architecture = "386"
		} else if SysArch == "amd64zip" {
			architecture = "amd64"
		}

		if runtime.GOOS == "windows" {
			// this execution happens if smoothflow is hosted on a window server like IIS or something
			logger.Log_PE("Starting executable process", logger.Information, sessionId)
			out, err := exec.Command("generateExecutable.bat", OSCode, SysArch, filefolder, flowName+".go", "CMD", "start").CombinedOutput()

			if err != nil {
				logger.Log_PE("Error building executable", logger.Error, sessionId)
				logger.Log_PE(string(err.Error()), logger.Debug, sessionId)
				return false, ""
			} else {
				logger.Log_PE("Success building executable", logger.Error, sessionId)
				logger.Log_PE(string(out), logger.Debug, sessionId)
				// if the generated on is windows, add .exe at last
				fileLocation := ""
				if OSCode == "windows" {
					//fileLocation = filefolder + flowName + ".exe"
					if strings.Contains(flowStruct.SysArch, "zip") {
						//tar -czf dd.tar.gz --directory="/var/www/html/engine/BuiltExecutables/ziketmaileme101comqwerty" ziketmaileme101comqwerty.exe
						cmd := "C:\\Program Files\\7-Zip\\7z.exe a -r " + filefolder + flowName + "zip -w " + filefolder + "\" " + flowName + ".exe -mem=AES256"
						_, _ = exec.Command("cmd", "/C", cmd).Output()
						fileLocation = filefolder + flowName + ".zip"
					} else {
						fileLocation = filefolder + flowName + ".exe"
					}
				} else {
					//fileLocation = filefolder + flowName
					if strings.Contains(flowStruct.SysArch, "zip") {
						cmd := "C:\\Program Files\\7-Zip\\7z.exe a -r " + filefolder + flowName + "zip -w " + filefolder + "\" " + flowName + " -mem=AES256"
						_, _ = exec.Command("cmd", "/C", cmd).Output()
						fileLocation = filefolder + flowName + ".zip"
					} else {
						fileLocation = filefolder + flowName
					}
				}
				// according to the windows server the file location should be changed. refer the linux implementation
				return true, fileLocation
			}

		} else {
			// this execution happens if smoothflow is hosted on a linux server
			logger.Log_PE("Starting executable process", logger.Information, sessionId)

			out, err := exec.Command(pwd+slash+"generateExecutable.sh", OSCode, architecture, filefolder, flowName+".go", "sh").CombinedOutput()
			if err != nil {
				logger.Log_PE("Error building executable", logger.Error, sessionId)
				logger.Log_PE(string(err.Error()), logger.Error, sessionId)
				return false, ""
			} else {
				logger.Log_PE("Success building executable", logger.Debug, sessionId)
				logger.Log_PE(string(out), logger.Information, sessionId)

				// if the generated on is windows, add .exe at last
				fileLocation := ""
				if OSCode == "windows" {
					if strings.Contains(flowStruct.SysArch, "zip") {
						//tar -czf dd.tar.gz --directory="/var/www/html/engine/BuiltExecutables/ziketmaileme101comqwerty" ziketmaileme101comqwerty.exe
						cmd := "tar -czf " + filefolder + flowName + ".tar.gz --directory=\"" + filefolder + "\" " + flowName + ".exe"
						_, _ = exec.Command("sh", "-c", cmd).Output()
						fileLocation = filefolder + flowName + ".tar.gz"
					} else {
						fileLocation = filefolder + flowName + ".exe"
					}

				} else {
					if strings.Contains(flowStruct.SysArch, "zip") {
						cmd := "tar -czf " + filefolder + flowName + ".tar.gz --directory=\"" + filefolder + "\" " + flowName
						_, _ = exec.Command("sh", "-c", cmd).Output()
						fileLocation = filefolder + flowName + ".tar.gz"
					} else {
						fileLocation = filefolder + flowName
					}
					//_, _ = exec.Command("sh", "-c", ("tar -czf " + flowName + ".tar.gz " + flowName)).Output()
					//fileLocation = filefolder + flowName
				}
				// make the file location relavent to a linux server
				fileLocation = strings.Replace(fileLocation, "/var/www/html", "", -1)
				return true, fileLocation
			}
		}
	}

}

func getMessageContent(msg variable) string {
	var template = ""
	var stringvariable = ""
	var convertionDetails = ""
	switch msg.DataType {
	case "arrayItem(string)":
		{
			stringvariable = msg.Value
		}
	default:
		{
			if msg.Type == "hardcoded" {
				stringvariable = `"` + msg.Value + `"`
			}
			if msg.Type == "dynamic" {
				if msg.DataType == "string" {
					stringvariable = `flowData["` + msg.Value + `"].(string)`
				}
				if msg.DataType == "int" {
					stringvariable = `strconv.FormatInt(flowData["` + msg.Value + `"].(int64),10)`
				}
				if msg.DataType == "float" {
					stringvariable = `strconv.FormatFloat(flowData["` + msg.Value + `"].(float64), 'f', 0, 64)`
				}
			}
			if msg.Type == "custom" {
				stringvariable, convertionDetails = getCustomValueListString(msg)
			}
		}
	}
	template = convertionDetails
	template = template + `logger.Log_WF(` + stringvariable + `, logger.Debug, flowData["InSessionID"].(string))`

	return template
}

func getCustomValueListString(obj variable) (string, string) {
	var convertionDetails = ""
	var returnString = ""
	var count int
	// loop through the set of valuelist available to generate msg output stirng
	for _, listItem := range obj.ValueList {
		var trimmedValue = strings.TrimSpace(listItem.Value)
		ItemID := randStr(5, "alpha")
		switch listItem.DataType {
		case "string":
			if listItem.Type == "hardcoded" {
				returnString = returnString + `"` + trimmedValue + `"`
			}
			if listItem.Type == "dynamic" {
				returnString = returnString + `flowData["` + trimmedValue + `"].(string)`
			}
		case "int":
			if listItem.Type == "hardcoded" {
				convertionDetails = convertionDetails + ItemID + `, _ := strconv.Atoi("` + trimmedValue + `");`
				returnString = returnString + "string(" + ItemID + ")"
			}
			if listItem.Type == "dynamic" {
				convertionDetails = convertionDetails + ItemID + `, _ := strconv.Atoi(flowData["` + trimmedValue + `"].(string));`
				//returnString = returnString + "string(" + ItemID + ")"
				returnString = returnString + "strconv.Itoa(" + ItemID + ")"
			}
		case "float":
			if listItem.Type == "hardcoded" {
				convertionDetails = convertionDetails + ItemID + `, _ := strconv.ParseFloat("` + trimmedValue + `");`
				returnString = returnString + "string(" + ItemID + ")"
			}
			if listItem.Type == "dynamic" {
				convertionDetails = convertionDetails + ItemID + `, _ := strconv.ParseFloat(strconv.FormatFloat(flowData["` + trimmedValue + `"].(float64), 'f', 0, 64), 64);`
				returnString = returnString + "strconv.FormatFloat(" + ItemID + ", 'f', 0, 64)"
			}
		}
		if count < len(obj.ValueList)-1 {
			returnString = returnString + ` +" " +`
		}
		count++
	}
	return returnString, convertionDetails
}

func getArrayDataInitiation(arguments []linq.T) string {
	var variableInit = `
	`
	for _, argumentVal := range arguments {
		// check for variables which are arrays or collections
		var argument = argumentVal.(variable)
		switch argument.DataType {
		case "array(string)":
			variableInit = variableInit + `var ` + argument.Key + ` []string`
		}
	}
	return variableInit
}

func getFlowTerminationContent(flag, Messsage, Status, ErrorCode string) string {
	msg := ""
	status := ""
	errorCode := ""
	if flag == "terminate" {
		msg = `"` + Messsage + `"`
		status = Status
		errorCode = ErrorCode
	} else if flag == "flowend" {
		msg = "WFContext.Message"
		status = "WFContext.Status"
		errorCode = "WFContext.ErrorCode"
	}
	return `var returnObj context.ReturnData;returnObj.JSONOutput = flowData;returnObj.WFOutArguments = ALLOutArguments;returnObj.WorkflowResult.Message = ` + msg + `;returnObj.WorkflowResult.Status = ` + status + `;returnObj.WorkflowResult.ErrorCode = ` + errorCode + `;returnDataJSON, rtDError := json.Marshal(returnObj);if rtDError != nil {logger.Log_WF("WF return data marshal error", logger.Debug, flowData["InSessionID"].(string));};`
}

func getViewName(view, executionlevel string) string {
	return "initiate.View." + view + "." + executionlevel
}

func getCollectionData(variableObject variable) string {
	var stringValue = ""
	var datatype = ""
	switch variableObject.DataType {
	case "int":
		{
			datatype = "int"
		}
	case "string":
		{
			datatype = "string"
		}
	case "float":
		{
			datatype = "float32"
		}
	}
	//var intarray []int
	stringValue = "var " + variableObject.Value + " []" + datatype
	return stringValue
}

func getRemoveFromCollectionContent(colVariable, valVariable variable) string {
	var itemtoRemove = ""
	switch valVariable.DataType {
	case "string":
		{
			itemtoRemove = `"` + valVariable.Value + `"`
		}
	default:
		{
			itemtoRemove = valVariable.Value
		}
	}

	var indexName = "index" + colVariable.Value
	var returnString = `
		var ` + indexName + ` = 0;
		for index, value := range ` + colVariable.Value + ` {
			if(value == ` + itemtoRemove + `){
				` + indexName + ` = index
			}
		}
		` + colVariable.Value + ` = append(` + colVariable.Value + `[:` + indexName + `],` + colVariable.Value + `[` + indexName + `+1:]...);
	`
	return returnString
}

func getAddtoCollectionContent(colVariable, valVariable variable) string {
	var returnString = ""
	//intarray = append(intarray,1,2,3,77,87,788,96)
	switch valVariable.DataType {
	case "string":
		{
			if colVariable.Type == "hardcoded" {
				returnString = colVariable.Value + " = append(" + colVariable.Value + ",`" + valVariable.Value + "`);"
			}
			if colVariable.Type == "dynamic" {
				returnString = colVariable.Value + ` = append(` + colVariable.Value + ` ,"` + valVariable.Value + `");
				flowData["` + colVariable.Value + `"] = ` + colVariable.Value
			}
		}
	default:
		{
			if colVariable.Type == "hardcoded" {
				returnString = colVariable.Value + " = append(" + colVariable.Value + "," + valVariable.Value + ");"
			}
			if colVariable.Type == "dynamic" {
				returnString = colVariable.Value + ` = append(` + colVariable.Value + ` ,` + valVariable.Value + `);
				flowData["` + colVariable.Value + `"] = ` + colVariable.Value
			}
		}
	}
	return returnString
}

/*func getBasicExecutionalFunctionData() string {
return `if _, ok := flowData["ExecutionLevel"]; ok {logger.Log_WF("ExecutionLevel Value already available.", logger.Debug, flowData["InSessionID"].(string));} else {flowData["ExecutionLevel"] = "0";logger.Log_WF("ExecutionLevel value changed to '0'", logger.Debug, flowData["InSessionID"].(string));}; if (flowData["ExecutionLevel"] == "0") {` + "\n" + `//initiate.View.drawboard.0` + "\n\n" + `}`
}*/

func getCaseContent(CaseState, sessionId string, flowstruct JsonFlow) string {

	returnString := `case Case_` + CaseState + `:{`
	// check if Hibernate controls are present
	/*hibernateNodes, _ := linq.From(flowstruct.Nodes).Where(
	func(in linq.T) (bool, error) {
					return in.(nodedata).LibraryID == "3" && in.(nodedata).ParentView == CaseState, nil
					}).Results()*/

	//logger.Log_PE("Hibernated Nodes: "+hibernateNodes.([]nodedata).length,logger.Debug, sessionId)
	returnString = returnString + `//initiate.View.` + CaseState + `.0
					`
	/*if hibernateNodes != nil {
	// if present add multiple view lines
	for _, node := range hibernateNodes {
					Levelobj, _, _ := linq.From(node.(nodedata).Variables).Where(
									func(in linq.T) (bool, error) { return in.(variable).Key == "Level", nil }).First()
					returnString = returnString + `//initiate.View.` + CaseState + `.` + Levelobj.(variable).Value + `
																	`
	}
	}*/

	returnString = returnString + `}
	`

	return returnString
}

func getDefaultCaseContent(CaseState, sessionId string, flowstruct JsonFlow) string {
	returnString := `default :{`
	// check if Hibernate controls are present
	/*hibernateNodes, _ := linq.From(flowstruct.Nodes).Where(
	func(in linq.T) (bool, error) {
					return in.(nodedata).LibraryID == "3" && in.(nodedata).ParentView == CaseState, nil
					}).Results()*/

	//logger.Log_PE("Hibernated Nodes: "+hibernateNodes.([]nodedata).length,logger.Debug, sessionId)
	returnString = returnString + `//initiate.View.` + CaseState + `.0
					`
	/*  if hibernateNodes != nil {
					// if present add multiple view lines
					for _, node := range hibernateNodes {
									Levelobj, _, _ := linq.From(node.(nodedata).Variables).Where(
													func(in linq.T) (bool, error) { return in.(variable).Key == "Level", nil }).First()
									returnString = returnString + `//initiate.View.` + CaseState + `.` + Levelobj.(variable).Value + `
																									`
					}
	}
	*/
	returnString = returnString + `}
	`

	return returnString
}

func getForeachContent(view string, nodeData nodedata, foreachStuct foreachStuct, flowStruct JsonFlow) string {
	CollectionName, _, _ := linq.From(nodeData.Variables).Where(
		func(in linq.T) (bool, error) { return in.(variable).Key == "CollectionName", nil }).First()

	ItemName, _, _ := linq.From(nodeData.Variables).Where(
		func(in linq.T) (bool, error) { return in.(variable).Key == "ItemName", nil }).First()

	var returnString = `
	for _, ` + ItemName.(variable).Value + ` := range ` + CollectionName.(variable).Value + ` {
		//initiate.View.` + foreachStuct.ForloopState + `.0
	}
	`

	return returnString
}

func getSwitchConditionContent(view string, nodeData nodedata, switchStruct switchStuct, flowStruct JsonFlow) string {
	// if the variable is hardcoded and dynamic, it should be checked here before assigning to the switch statement
	returnstring := ``

	VariableObj, _, _ := linq.From(nodeData.Variables).Where(
		func(in linq.T) (bool, error) { return in.(variable).Key == "Variable", nil }).First()

	DataTypeObj, _, _ := linq.From(nodeData.Variables).Where(
		func(in linq.T) (bool, error) { return in.(variable).Key == "DataType", nil }).First()

	variableString := getConvertedMainVariableString(DataTypeObj.(variable).Value, VariableObj.(variable).Type, VariableObj.(variable).Value)
	listOfCases := getListofCases(DataTypeObj.(variable).Value, view, flowStruct)
	/*
		MainVariable := "valueFrom Outside"
		Case_2814de := "aaa"
		Case_47370d := "bbb"
	*/

	returnstring = variableString + listOfCases + `
		switch MainVariable {
//initiate.View.` + switchStruct.SwitchState + `.0
			}`

	fmt.Println("Switch content: ")
	fmt.Println(returnstring)

	fmt.Println("Variable String: ")
	fmt.Println(variableString)

	fmt.Println("Cases: ")
	fmt.Println(listOfCases)

	return returnstring
}

func getWhileLoopContent(nodeData nodedata, whileObj whileStuct, flowStruct JsonFlow) string {
	// if the variable is hardcoded and dynamic, it should be checked here before assigning to the switch statement
	returnstring := ``

	ValueOneObj, _, _ := linq.From(nodeData.Variables).Where(
		func(in linq.T) (bool, error) { return in.(variable).Key == "ValueOne", nil }).First()

	ConditionObj, _, _ := linq.From(nodeData.Variables).Where(
		func(in linq.T) (bool, error) { return in.(variable).Key == "Condition", nil }).First()

	ValueTwoObj, _, _ := linq.From(nodeData.Variables).Where(
		func(in linq.T) (bool, error) { return in.(variable).Key == "ValueTwo", nil }).First()

	//variableString := getConvertedMainVariableString(DataTypeObj.(variable).Value, VariableObj.(variable).Type, VariableObj.(variable).Value)
	returnstring = getConvertedWhileVariableString(ValueOneObj.(variable).Value, ValueOneObj.(variable).Type, ValueOneObj.(variable).DataType, ConditionObj.(variable).Value, ValueTwoObj.(variable).Value, ValueTwoObj.(variable).Type, ValueTwoObj.(variable).DataType, whileObj.WhileState)

	fmt.Println("While content: ")
	fmt.Println(returnstring)

	return returnstring
}

func getListofCases(Datatype, view string, flowStruct JsonFlow) string {
	caseControls, _ := linq.From(flowStruct.Nodes).Where(func(in linq.T) (bool, error) {
		return in.(nodedata).LibraryID == "9" && in.(nodedata).ParentView == view, nil
	}).Results()

	variableString := ``

	for _, caseObj := range caseControls {
		value := caseObj.(nodedata).Variables[0].Value
		inputType := caseObj.(nodedata).Variables[0].Type
		viewID := caseObj.(nodedata).OtherData.CaseUUID

		switch Datatype {
		case "string":
			if inputType == "hardcoded" {
				variableString = variableString + `Case_` + viewID + ` := "` + value + `";`
			}
			if inputType == "dynamic" {
				variableString = variableString + `Case_` + viewID + ` := flowData["` + value + `"].(string);`
			}
		case "int":
			if inputType == "hardcoded" {
				variableString = variableString + `Case_` + viewID + ` := strconv.Atoi("` + value + `");`
			}
			if inputType == "dynamic" {
				variableString = variableString + `Case_` + viewID + ` := strconv.Atoi(flowData["` + value + `"].(string));`
			}
		case "float":
			if inputType == "hardcoded" {
				variableString = variableString + `Case_` + viewID + ` := strconv.ParseFloat("` + value + `");`
			}
			if inputType == "dynamic" {
				variableString = variableString + `Case_` + viewID + ` := strconv.ParseFloat(flowData["` + value + `"].(string));`
			}
		default:
			variableString = ""
		}
	}
	return variableString
}

func getConvertedMainVariableString(datatype, inputType, value string) string {
	variableString := ``
	switch datatype {
	case "string":
		if inputType == "hardcoded" {
			variableString = `MainVariable := "` + value + `";`
		}
		if inputType == "dynamic" {
			variableString = `MainVariable := flowData["` + value + `"].(string);`
		}
	case "int":
		if inputType == "hardcoded" {
			variableString = `MainVariable := strconv.Atoi("` + value + `");`
		}
		if inputType == "dynamic" {
			variableString = `MainVariable := strconv.Atoi(flowData["` + value + `"].(string));`
		}
	case "float":
		if inputType == "hardcoded" {
			variableString = `MainVariable := strconv.ParseFloat("` + value + `");`
		}
		if inputType == "dynamic" {
			variableString = `MainVariable := strconv.ParseFloat(flowData["` + value + `"].(string));`
		}
	default:
		variableString = ""
	}
	return variableString
}

func getConvertedWhileVariableString(valueOneValue, valueOneType, valueOneDatatype, Condition, valueTwoValue, valueTwoType, valueTwoDatatype, View string) string {
	variableString := ``
	fmt.Println("DataType: " + valueOneDatatype)
	fmt.Println("ValueOne Value: " + valueOneValue)
	fmt.Println("ValueOne type: " + valueOneType)
	fmt.Println("ValueOne datatype: " + valueOneDatatype)
	fmt.Println("Conditioin: " + Condition)
	fmt.Println("ValueTwo Value: " + valueTwoValue)
	fmt.Println("ValueTwo type: " + valueTwoType)
	fmt.Println("DataType: " + valueTwoDatatype)
	fmt.Println("View: " + View)

	switch valueOneDatatype {
	case "string":
		fmt.Println("Entered Strings")
		//valueone
		if strings.EqualFold("hardcoded", valueOneType) {
			variableString = variableString + `ValueOne := "` + valueOneValue + `";`
		}
		if strings.EqualFold("dynamic", valueOneType) {
			variableString = variableString + `ValueOne := flowData["` + valueOneValue + `"].(string);`
		}
		// valuetwo
		if strings.EqualFold("hardcoded", valueTwoType) {
			variableString = variableString + `ValueTwo := "` + valueTwoValue + `";`
		}
		if strings.EqualFold("dynamic", valueTwoType) {
			variableString = variableString + `ValueTwo := flowData["` + valueTwoValue + `"].(string);`
		}
		fmt.Println("Exiting string")
		fmt.Println(variableString)
	case "int":
		fmt.Println("Entered Int")
		// value one
		if strings.EqualFold("hardcoded", valueOneType) {
			variableString = variableString + `ValueOne, _ := strconv.Atoi("` + valueOneValue + `");`
		}
		if strings.EqualFold("dynamic", valueOneType) {
			variableString = variableString + `ValueOne, _ := strconv.Atoi(flowData["` + valueOneValue + `"].(string));`
		}
		// value two
		if strings.EqualFold("hardcoded", valueTwoType) {
			variableString = variableString + `ValueTwo, _ := strconv.Atoi("` + valueTwoValue + `");`
		}
		if strings.EqualFold("dynamic", valueTwoType) {
			variableString = variableString + `ValueTwo, _ := strconv.Atoi(flowData["` + valueTwoValue + `"].(string));`
		}
		fmt.Println("Exiting int")
		fmt.Println(variableString)
	case "float":
		fmt.Println("Entered Float")
		// value one
		if strings.EqualFold("hardcoded", valueOneType) {
			variableString = variableString + `ValueOne, _ := strconv.ParseFloat("` + valueOneValue + `");`
		}
		if strings.EqualFold("dynamic", valueOneType) {
			variableString = variableString + `ValueOne, _ := strconv.ParseFloat(flowData["` + valueOneValue + `"].(string));`
		}
		// value two
		if strings.EqualFold("hardcoded", valueTwoType) {
			variableString = variableString + `ValueTwo, _ := strconv.ParseFloat("` + valueTwoValue + `");`
		}
		if strings.EqualFold("dynamic", valueTwoType) {
			variableString = variableString + `ValueTwo, _ := strconv.ParseFloat(flowData["` + valueTwoValue + `"].(string));`
		}
		fmt.Println("Exiting float")
		fmt.Println(variableString)
	default:
		variableString = ""
		fmt.Println("Exiting default")
		fmt.Println(variableString)
	}
	variableString = variableString + `
			for ValueOne ` + Condition + ` ValueTwo {
																//initiate.View.` + View + `.0
				ValueOne = ValueOne + 1
			}
			`
	return variableString
}

func getExecutionalFooter(executionLevel, WFName, view string) string {

	execString := `
			status_hibernateLevel_` + executionLevel + ` := Components.HibernateWorkflow(flowData,"` + executionLevel + `","` + WFName + `"); if status_hibernateLevel_` + executionLevel + ` == true {logger.Log_WF("Hibernation level ` + executionLevel + ` hibernated succesfully!", logger.Debug, flowData["InSessionID"].(string));` + "\n" + `//initiate.View.` + view + `.` + executionLevel + "\n" +
		`} else {logger.Log_WF("Hibernation level ` + executionLevel + ` failed to hibernated succesfully.", logger.Debug, flowData["InSessionID"].(string));}
			`
	return execString
}

func getAssignmentControlDetails(fromControl, toControl variable) string {

	var template = ""
	var stringvariable = ""
	var convertionDetails = ""
	to := toControl.Value
	from := fromControl.Value
	switch fromControl.ValueType {
	case "arrayItem(string)":
		{
			stringvariable = `"` + fromControl.Value + `"`
		}
	default:
		{
			if fromControl.Type == "hardcoded" {
				stringvariable = `flowData["` + to + `"] = flowData["` + from + `"]`
			}
			if fromControl.Type == "dynamic" {
				if fromControl.ValueType == "string" {
					stringvariable = `flowData["` + to + `"] = flowData["` + from + `"].(string)`
				}
				if fromControl.ValueType == "int" {
					stringvariable = `flowData["` + to + `"] = strconv.FormatInt(flowData["` + from + `"].(int64),10)`
				}
				if fromControl.ValueType == "float" {
					stringvariable = `flowData["` + to + `"] = strconv.FormatFloat(flowData["` + from + `"].(float64), 'f', 0, 64)`
				}
			}
			if fromControl.Type == "custom" {
				stringvariable, convertionDetails = getCustomValueListString(fromControl)
			}
		}
	}
	template = convertionDetails + stringvariable
	return template
}

func getExecutionalData(flowStruct JsonFlow, view, sessionId string) string {

	hibernateVariables, _ := linq.From(flowStruct.Nodes).Where(func(in linq.T) (bool, error) {
		return in.(nodedata).LibraryID == "3" && in.(nodedata).ParentView == view, nil
	}).Results()
	convertedNodes, _ := json.Marshal(hibernateVariables)
	logger.Log_PE("Available Nodes for the view: "+view, logger.Debug, sessionId)
	logger.Log_PE(string(convertedNodes), logger.Debug, sessionId)
	execString := ""
	if view == "drawboard" {
		execString = `if _, ok := flowData["ExecutionLevel"]; ok {logger.Log_WF("ExecutionLevel Value already available.", logger.Debug, flowData["InSessionID"].(string));} else {flowData["ExecutionLevel"] = "0";logger.Log_WF("ExecutionLevel value changed to '0'", logger.Debug, flowData["InSessionID"].(string));}; if (flowData["ExecutionLevel"] == "0") {` + "\n" + `//initiate.View.drawboard.0
			}`
	}

	tot := cap(hibernateVariables)
	var count int = 0
	for _, node := range hibernateVariables {
		count = count + 1
		// argument.(variable).Type
		executionLevel := node.(nodedata).Variables[0].Value
		execString = execString + "\n\n" + `// save the objects upto hear in objectstore` + "\n" + `// stop the above process and continue with the following.` + "\n" + `if (flowData["ExecutionLevel"] == "` + executionLevel + `") {logger.Log_WF("Execution executionlevel: ` + executionLevel + `", logger.Debug, flowData["InSessionID"].(string));` + "\n" + `//initiate.View.` + view + `.` + executionLevel + `
				}`
		// if its the last round, it will add the delete function
		if tot == count {
			execString = execString + "\n\n" + `if flowData["ExecutionLevel"] == "` + executionLevel + `" {logger.Log_WF("Deleting session details.", logger.Debug, flowData["InSessionID"].(string));Components.DeleteSession(flowData);}`
		}
	}
	return execString
}

func getConfigContent(flowStruct JsonFlow, makeExecutable bool) string {
	var concatinatedString string = `defaultSettings := make(map[string]string);`
	var str string = ""
	for _, variable := range flowStruct.Arguments {
		if strings.EqualFold("Configuration", variable.Category) {
			concatinatedString = concatinatedString + `defaultSettings["` + variable.Key + `"] = "` + variable.Value + `";`
		}
	}
	if makeExecutable == false {
		str = `settings := Common.GetConfig();if len(settings) > 0 {for key, value := range defaultSettings {_, ok := settings[key];if ok == false {settings[key] = value;}else{settings[key] = value;}};Common.SaveConfig(settings);} else {Common.SaveConfig(defaultSettings);}`
	} else {
		str = `settings := Common.GetConfig();if len(settings) > 0 {for key, value := range defaultSettings {_, ok := settings[key];if ok == false {settings[key] = value;}};Common.SaveConfig(settings);} else {Common.SaveConfig(defaultSettings);}`
	}
	concatinatedString = concatinatedString + str
	return concatinatedString
}

func getHeaderInformation(flowStruct JsonFlow, makeExecutable bool) string {
	var concatinatedString = ""
	arguments, _ := linq.From(flowStruct.Arguments).Where(func(in linq.T) (bool, error) {
		return in.(variable).Category == "Header", nil
	}).Results()

	if makeExecutable {
		if len(arguments) > 0 {
			for _, argument := range arguments {
				if argument.(variable).Type == "hardcoded" {
					concatinatedString = concatinatedString + `ALLHeaders["` + argument.(variable).Key + `"] = "` + argument.(variable).Value + `";`
				}
				if argument.(variable).Type == "dynamic" {
					if argument.(variable).Value != "" {
						concatinatedString = concatinatedString + `ALLHeaders["` + argument.(variable).Key + `"] = T.Context.Request().Header.Get("` + argument.(variable).Key + `");`
					} else {
						concatinatedString = concatinatedString + `ALLHeaders["` + argument.(variable).Key + `"] = T.Context.Request().Header.Get("` + argument.(variable).Key + `");`
					}
				}
			}
		}

	} else {
		// if not a executable what will happen
	}

	return concatinatedString
}

func getAllINArgumentList(flowStruct JsonFlow, makeExecutable bool) string {
	var concatinatedString string = `ALLInArguments := make(map[string]interface{});`
	var str string = ""
	for _, variable := range flowStruct.Arguments {
		if strings.EqualFold("InArgument", variable.Category) {
			if strings.EqualFold("dynamic", variable.Type) {
				concatinatedString = concatinatedString + `ALLInArguments["` + variable.Key + `"] = flowData["` + variable.Key + `"];`
			}
			if strings.EqualFold("hardcoded", variable.Type) {
				concatinatedString = concatinatedString + `ALLInArguments["` + variable.Key + `"] = "` + variable.Value + `";`
			}
		}
	}
	if makeExecutable == false {
		str = `logger.Log_WF("", logger.Debug, flowData["InSessionID"].(string));var count int = 0;for key, val := range ALLInArguments {if _, ok := flowData[key]; ok {logger.Log_WF("OK - " + string(key) + " - " + flowData[key].(string), logger.Debug, flowData["InSessionID"].(string));} else {if val == "" {logger.Log_WF("NO - "+ string(key), logger.Debug, flowData["InSessionID"].(string));count = count + 1;}}};if count > 0 {logger.Log_WF("", logger.Debug, flowData["InSessionID"].(string));logger.Log_WF("All INArguments are not received", logger.Debug, flowData["InSessionID"].(string));logger.Log_WF("Application terminated.", logger.Debug, flowData["InSessionID"].(string));WFContext.Message = "All INArguments are not received";WFContext.Status = false;WFContext.ErrorCode = 5;var returnObj context.ReturnData;returnObj.JSONOutput = flowData;returnObj.WorkflowResult.Message = WFContext.Message;returnObj.WorkflowResult.Status = WFContext.Status;returnObj.WorkflowResult.ErrorCode = WFContext.ErrorCode;returnDataJSON, rtDError := json.Marshal(returnObj);if rtDError != nil {logger.Log_WF("WF return data marshal error", logger.Debug, flowData["InSessionID"].(string));};fmt.Print(string(returnDataJSON));return;}`
	} else {
		str = `logger.Log_WF("", logger.Debug, flowData["InSessionID"].(string));var count int = 0;for key, val := range ALLInArguments {if _, ok := flowData[key]; ok {logger.Log_WF("OK - " + string(key) + " - " + flowData[key].(string), logger.Debug, flowData["InSessionID"].(string));} else {if val == "" {logger.Log_WF("NO - "+ string(key), logger.Debug, flowData["InSessionID"].(string));count = count + 1;}}};if count > 0 {logger.Log_WF("", logger.Debug, flowData["InSessionID"].(string));logger.Log_WF("All INArguments are not received", logger.Debug, flowData["InSessionID"].(string));logger.Log_WF("Application terminated.", logger.Debug, flowData["InSessionID"].(string));WFContext.Message = "All INArguments are not received";WFContext.Status = false;WFContext.ErrorCode = 5;var returnObj context.ReturnData;returnObj.JSONOutput = flowData;returnObj.WorkflowResult.Message = WFContext.Message;returnObj.WorkflowResult.Status = WFContext.Status;returnObj.WorkflowResult.ErrorCode = WFContext.ErrorCode;returnDataJSON, rtDError := json.Marshal(returnObj);if rtDError != nil {logger.Log_WF("WF return data marshal error", logger.Debug, flowData["InSessionID"].(string));};fmt.Print(string(returnDataJSON));return string(returnDataJSON);}`
	}

	concatinatedString = concatinatedString + str
	return concatinatedString
}

func getInArgumentPrintDetails(flowStruct JsonFlow) (result string) {
	result = "In Arguments : \\n"
	for _, argument := range flowStruct.Arguments {
		if strings.EqualFold(argument.Category, "InArgument") {
			result += "\\t " + argument.Key + " <" + argument.DataType + "> \\n"
		}
	}
	return
}

func getInJSONPrintDetails(flowStruct JsonFlow) (result string) {
	result = "Work Flow Invoke Body : \\n \\t { \\n"

	argumentArray := make([]variable, 0)

	for _, args := range flowStruct.Arguments {
		if strings.EqualFold(args.Category, "InArgument") {
			argumentArray = append(argumentArray, args)
		}
	}

	for x := 0; x < len(argumentArray); x++ {
		argument := argumentArray[x]
		if strings.EqualFold(argument.Category, "InArgument") && x < (len(argumentArray)-1) {
			result += " \\t \\t " + "\\\"" + argument.Key + "\\\"" + " : " + "\\\" " + "\\\", \\n"
		} else {
			result += " \\t \\t " + "\\\"" + argument.Key + "\\\"" + " : " + "\\\" " + "\\\" \\n"
		}
	}

	result += " \\t }"

	return
}

func getOutArgumentPrintDetails(flowStruct JsonFlow) (result string) {
	result = "Out Arguments : \\n"
	for _, argument := range flowStruct.Arguments {
		if strings.EqualFold(argument.Category, "OutArgument") {
			result += "\\t " + argument.Key + " <" + argument.DataType + "> \\n"
		}
	}
	return
}

func checkIfConversionAvailable(module nodedata) bool {

	var convertionAvailable = false
	for _, varObj := range module.Variables {
		fmt.Println(varObj)
		if varObj.ConvertTo != "" {
			if !strings.EqualFold("string", varObj.ConvertTo) {
				convertionAvailable = true
				fmt.Println(convertionAvailable)
			}
		}
		if varObj.Type == "custom" {
			for _, item := range varObj.ValueList {
				if !strings.EqualFold("string", item.DataType) {
					convertionAvailable = true
					fmt.Println(convertionAvailable)
				}
			}
		}
		if varObj.DataType != "string" {
			convertionAvailable = true
			fmt.Println(convertionAvailable)
		}
	}
	return convertionAvailable
}

func getINArgumentString(nodeData nodedata, variableName string, flowStruct JsonFlow, makeExecutable bool) string {
	var concatinatedString string = variableName + ` := make(map[string]interface{});`
	for _, variable := range nodeData.Variables {
		fmt.Println("inargument - ", variable)
		if variable.Category == "InArgument" {
			fmt.Println("-" + variable.Category)
			fmt.Println("-" + variable.ConvertTo)
			fmt.Println("-" + variable.Type)
			fmt.Println("-" + variable.ValueType)
			// in here should check if the ConverTo is having a value or not, if not it can continue
			if variable.ConvertTo == "" {
				// if the arguent is hardcorded
				if variable.Type == "hardcoded" {
					concatinatedString = concatinatedString + variableName + `["` + variable.Key + `"] = "` + variable.Value + `";`
				}
				// if the argument is directly taken from flowData
				if variable.Type == "dynamic" {
					concatinatedString = concatinatedString + variableName + `["` + variable.Key + `"] = flowData["` + variable.Value + `"];`
				}
				if variable.Type == "custom" {
					stringvariable, convertionDetails := getCustomValueListString(variable)
					concatinatedString = concatinatedString + convertionDetails
					concatinatedString = concatinatedString + variableName + `["` + variable.Key + `"] = ` + stringvariable + `;`
				}
			} else {
				switch variable.ConvertTo {
				case "string":
					{
						if variable.Type == "hardcoded" {
							concatinatedString = concatinatedString + variableName + `["` + variable.Key + `"] = "` + variable.Value + `";`
						}
						if variable.Type == "dynamic" {
							if variable.Value != "" {
								switch variable.ValueType {
								case "string":
									{
										concatinatedString = concatinatedString + variableName + `["` + variable.Key + `"] = flowData["` + variable.Value + `"].(string);`
									}
								case "int":
									{
										concatinatedString = concatinatedString + variableName + `["` + variable.Key + `"] = strconv.FormatInt(flowData["` + variable.Value + `"].(int64),10) ;`
									}
								case "float":
									{
										concatinatedString = concatinatedString + variableName + `["` + variable.Key + `"] = strconv.FormatFloat(flowData["` + variable.Value + `"].(float64), 'f', 0, 64);`
									}
								case "boolean":
									{
										concatinatedString = concatinatedString + variableName + `["` + variable.Key + `"] = strconv.FormatBool(flowData["` + variable.Value + `"].(bool));`
									}
								}
							} else {
								concatinatedString = concatinatedString + variableName + `["` + variable.Key + `"] = "";`
							}
						}
						if variable.Type == "custom" {
							stringvariable, convertionDetails := getCustomValueListString(variable)
							concatinatedString = concatinatedString + convertionDetails
							concatinatedString = concatinatedString + variableName + `["` + variable.Key + `"] = ` + stringvariable + `;`
						}
					}
				case "int":
					{
						if variable.Type == "hardcoded" {
							concatinatedString = concatinatedString + variableName + variable.Key + `,` + variableName + variable.Key + `Err := strconv.Atoi("` + variable.Value + `");
									if ` + variableName + variable.Key + `Err != nil {
										logger.Log_WF("InArgument - ` + variable.Key + ` type convertion failed.", logger.Debug, flowData["InSessionID"].(string))}
										` + variableName + `["` + variable.Key + `"] = ` + variableName + variable.Key + `
										`
						}
						if variable.Type == "dynamic" {
							if variable.Value != "" {
								concatinatedString = concatinatedString + variableName + variable.Key + `,` + variableName + variable.Key + `Err := strconv.Atoi(strconv.FormatInt(flowData["` + variable.Value + `"].(int64),10));
											if ` + variableName + variable.Key + `Err != nil {
												logger.Log_WF("InArgument - ` + variable.Key + ` type convertion failed.", logger.Debug, flowData["InSessionID"].(string))}
												` + variableName + `["` + variable.Key + `"] = ` + variableName + variable.Key + `
												`
							} else {
								concatinatedString = concatinatedString + variableName + `["` + variable.Key + `"] = 0;`
							}
						}
					}
				case "float":
					{
						if variable.Type == "hardcoded" {
							concatinatedString = concatinatedString + variableName + variable.Key + `,` + variableName + variable.Key + `Err := strconv.ParseFloat("` + variable.Value + `", 64)
											if ` + variableName + variable.Key + `Err != nil {
												logger.Log_WF("InArgument - ` + variable.Key + ` type convertion failed.", logger.Debug, flowData["InSessionID"].(string))}
												` + variableName + `["` + variable.Key + `"] = ` + variableName + variable.Key + `
												`
						}
						if variable.Type == "dynamic" {
							if variable.Value != "" {

								switch variable.ValueType {
								case "string":
									{
										concatinatedString = concatinatedString + variableName + variable.Key + `,` + variableName + variable.Key + `Err := strconv.ParseFloat(flowData["` + variable.Value + `"].(string), 64);
													if ` + variableName + variable.Key + `Err != nil {
														logger.Log_WF("InArgument - ` + variable.Key + ` type convertion failed.", logger.Debug, flowData["InSessionID"].(string))}
														` + variableName + `["` + variable.Key + `"] = ` + variableName + variable.Key + `
														`
									}
								case "int":
									{
										concatinatedString = concatinatedString + variableName + variable.Key + `,` + variableName + variable.Key + `Err := strconv.ParseFloat(strconv.FormatInt(flowData["` + variable.Value + `"].(int64),10), 64);
													if ` + variableName + variable.Key + `Err != nil {
														logger.Log_WF("InArgument - ` + variable.Key + ` type convertion failed.", logger.Debug, flowData["InSessionID"].(string))}
														` + variableName + `["` + variable.Key + `"] = ` + variableName + variable.Key + `
														`
									}
								case "float":
									{
										concatinatedString = concatinatedString + variableName + variable.Key + `,` + variableName + variable.Key + `Err := strconv.ParseFloat(strconv.FormatFloat(flowData["` + variable.Value + `"].(float64), 'f', 0, 64), 64);
													if ` + variableName + variable.Key + `Err != nil {
														logger.Log_WF("InArgument - ` + variable.Key + ` type convertion failed.", logger.Debug, flowData["InSessionID"].(string))}
														` + variableName + `["` + variable.Key + `"] = ` + variableName + variable.Key + `
														`
									}
								}
							} else {
								concatinatedString = concatinatedString + variableName + `["` + variable.Key + `"] = "";`
							}
						}
					}
				case "boolean":
					{
						if variable.Type == "hardcoded" {
							concatinatedString = concatinatedString + variableName + variable.Key + `,` + variableName + variable.Key + `Err := strconv.ParseBool("` + variable.Value + `")
											if ` + variableName + variable.Key + `Err != nil {
												logger.Log_WF("InArgument - ` + variable.Key + ` type convertion failed.", logger.Debug, flowData["InSessionID"].(string))}
												` + variableName + `["` + variable.Key + `"] = ` + variableName + variable.Key + `
												`
						}
						if variable.Type == "dynamic" {
							if variable.Value != "" {
								concatinatedString = concatinatedString + variableName + variable.Key + `,` + variableName + variable.Key + `Err := strconv.ParseBool(strconv.FormatBool(flowData["` + variable.Value + `].(bool)))
													if ` + variableName + variable.Key + `Err != nil {
														logger.Log_WF("InArgument - ` + variable.Key + ` type convertion failed.", logger.Debug, flowData["InSessionID"].(string))}
														` + variableName + `["` + variable.Key + `"] = ` + variableName + variable.Key + `
														`
							} else {
								concatinatedString = concatinatedString + variableName + `["` + variable.Key + `"] = false;`
							}
						}
					}
				}
			}
		}
	}
	if makeExecutable {
		concatinatedString = concatinatedString + `
	//adding headers into the activity
	`
		arguments, _ := linq.From(flowStruct.Arguments).Where(func(in linq.T) (bool, error) {
			return in.(variable).Category == "Header", nil
		}).Results()

		if len(arguments) > 0 {
			for _, obj := range arguments {
				concatinatedString = concatinatedString + variableName + `["` + obj.(variable).Key + `"] = Headers["` + obj.(variable).Key + `"];`
			}
		}
	}
	return concatinatedString
}

func getOutArgumentString(nodeData nodedata, variableName string) string {
	var concatinatedString string = `/* out arguments for this activity */`
	str := ""
	for _, variable := range nodeData.Variables {
		if variable.Category == "OutArgument" {
			if variable.Value != "" {
				if variable.Type == "dynamic" {
					if variable.DataType == "string" {
						str = `if ` + variableName + `["` + variable.Key + `"] != "" && ` + variableName + `["` + variable.Key + `"] != nil {/*ALLOutArguments["` + variable.Value + `"] = ` + variableName + `["` + variable.Key + `"]*/;flowData["` + variable.Value + `"] = ` + variableName + `["` + variable.Key + `"];
						logger.Log_WF("OUT - ` + variable.Value + `: "+flowData["` + variable.Value + `"].(string), logger.Debug, flowData["InSessionID"].(string))
					}
											`
					}
					if variable.DataType == "int" {
						str = `if ` + variableName + `["` + variable.Key + `"] != "" && ` + variableName + `["` + variable.Key + `"] != nil {/*ALLOutArguments["` + variable.Value + `"] = ` + variableName + `["` + variable.Key + `"]*/;flowData["` + variable.Value + `"] = ` + variableName + `["` + variable.Key + `"];
						logger.Log_WF("OUT - ` + variable.Value + `: "+strconv.FormatInt(flowData["` + variable.Value + `"].(int64),10), logger.Debug, flowData["InSessionID"].(string))
					}
											`
					}
					if variable.DataType == "float" {
						str = `if ` + variableName + `["` + variable.Key + `"] != "" && ` + variableName + `["` + variable.Key + `"] != nil {/*ALLOutArguments["` + variable.Value + `"] = ` + variableName + `["` + variable.Key + `"]*/;flowData["` + variable.Value + `"] = ` + variableName + `["` + variable.Key + `"];
						logger.Log_WF("OUT - ` + variable.Value + `: "+strconv.FormatFloat(flowData["` + variable.Value + `"].(float64), 'f', 0, 64), logger.Debug, flowData["InSessionID"].(string))
					}
											`
					}
					if variable.DataType == "boolean" {
						str = `if ` + variableName + `["` + variable.Key + `"] != "" && ` + variableName + `["` + variable.Key + `"] != nil {/*ALLOutArguments["` + variable.Value + `"] = ` + variableName + `["` + variable.Key + `"]*/;flowData["` + variable.Value + `"] = ` + variableName + `["` + variable.Key + `"];
						logger.Log_WF("OUT - ` + variable.Value + `: "+ strconv.FormatBool(flowData["` + variable.Value + `"].(bool)), logger.Debug, flowData["InSessionID"].(string))
					}
											`
					}
					/* before removing adding the out argument to the Outarguments map */
					/*if variable.DataType == "boolean" {
						str = `if ` + variableName + `["` + variable.Key + `"] != "" && ` + variableName + `["` + variable.Key + `"] != nil {ALLOutArguments["` + variable.Value + `"] = ` + variableName + `["` + variable.Key + `"];flowData["` + variable.Value + `"] = ` + variableName + `["` + variable.Key + `"];
						logger.Log_WF("OUT - ` + variable.Value + `: "+strconv.FormatFloat(ALLOutArguments["` + variable.Value + `"].(float64), 'f', 0, 64), logger.Debug, flowData["InSessionID"].(string))
					}
											`
					}*/
					concatinatedString = concatinatedString + str
				}
			}
		}
	}

	return concatinatedString
}

func contentAlreadyAvailable(content, validatekey string) bool {
	if strings.Contains(content, validatekey) {
		return true
	} else {
		return false
	}
}

func randStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func IsExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func addIfCondition(valOne, check, valTwo variable, condition ifcondition) string {
	//setting the data type of the if condition
	dataType := valOne.DataType
	valOneType := valOne.Type
	valTwoType := valTwo.Type
	dataConvertionVal1 := ""
	dataConvertionVal2 := ""
	valueOne := valOne.Value
	conditionParameter := check.Value
	valueTwo := valTwo.Value
	VariableOneName := "Val1_" + condition.ID
	VariableOneError := VariableOneName + "_Error"
	VariableTwoName := "Val2_" + condition.ID
	VariableTwoError := VariableTwoName + "_Error"

	switch dataType {
	case "string":
		if valOneType == "hardcoded" {
			dataConvertionVal1 = `
									` + VariableOneName + ` := "` + valueOne + `"
									`
		} else if valOneType == "dynamic" {
			dataConvertionVal1 = `
										` + VariableOneName + ` := flowData["` + valueOne + `"].(string)
										`
		}

		if valTwoType == "hardcoded" {
			dataConvertionVal2 = `
										` + VariableTwoName + ` := "` + valueTwo + `"
										`
		} else if valTwoType == "dynamic" {
			dataConvertionVal2 = `
											` + VariableTwoName + ` := flowData["` + valueTwo + `"].(string)
											`
		}

	case "int":

		if valOneType == "hardcoded" {
			dataConvertionVal1 = `
											` + VariableOneName + `, ` + VariableOneError + ` := strconv.Atoi("` + valueOne + `")
											if ` + VariableOneError + ` != nil {
												logger.Log_WF("Value 1 type convertion failed.", logger.Debug, flowData["InSessionID"].(string))
											}
											`
			//strconv.Atoi("` + valueOne + `")
		} else if valOneType == "dynamic" {
			dataConvertionVal1 = `
												` + VariableOneName + `, ` + VariableOneError + ` := strconv.Atoi(strconv.FormatInt(flowData["` + valueOne + `"].(int64), 10))
												if ` + VariableOneError + ` != nil {
													logger.Log_WF("Value 1 type convertion failed.", logger.Debug, flowData["InSessionID"].(string))
												}
												`
		}

		if valTwoType == "hardcoded" {
			dataConvertionVal2 = `
												` + VariableTwoName + `, ` + VariableTwoError + ` := strconv.Atoi("` + valueTwo + `")
												if ` + VariableTwoError + ` != nil {
													logger.Log_WF("Value 2 type convertion failed.", logger.Debug, flowData["InSessionID"].(string))
													}`
		} else if valTwoType == "dynamic" {
			dataConvertionVal2 = `
														` + VariableTwoName + `, ` + VariableTwoError + ` := strconv.Atoi(strconv.FormatInt(flowData["` + valueOne + `"].(int64), 10))
														if ` + VariableTwoError + ` != nil {
															logger.Log_WF("Value 2 type convertion failed.", logger.Debug, flowData["InSessionID"].(string))
															}`
		}

	case "float":
		if valOneType == "hardcoded" {
			dataConvertionVal1 = `
															` + VariableOneName + `, ` + VariableOneError + ` := strconv.ParseFloat("` + valueOne + `", 64)
															if ` + VariableOneError + ` != nil {
																logger.Log_WF("Value 1 type convertion failed.", logger.Debug, flowData["InSessionID"].(string))
															}
															`
		} else if valOneType == "dynamic" {
			dataConvertionVal1 = `
																` + VariableOneName + `, ` + VariableOneError + ` := strconv.ParseFloat(strconv.FormatFloat(flowData["` + valueOne + `"].(float64), 'f', 0, 64), 64)
																if ` + VariableOneError + ` != nil {
																	logger.Log_WF("Value 1 type convertion failed.", logger.Debug, flowData["InSessionID"].(string))
																}
																`
		}

		if valTwoType == "hardcoded" {
			dataConvertionVal2 = `
																` + VariableTwoName + `, ` + VariableTwoError + ` := strconv.ParseFloat("` + valueTwo + `", 64)
																if ` + VariableTwoError + ` != nil {
																	logger.Log_WF("Value 2 type convertion failed.", logger.Debug, flowData["InSessionID"].(string))
																	}`
		} else if valTwoType == "dynamic" {
			dataConvertionVal2 = `
																		` + VariableTwoName + `, ` + VariableTwoError + ` := strconv.ParseFloat(strconv.FormatFloat(flowData["` + valueTwo + `"].(float64), 'f', 0, 64), 64)
																		if ` + VariableTwoError + ` != nil {
																			logger.Log_WF("Value 2 type convertion failed.", logger.Debug, flowData["InSessionID"].(string))
																			}`
		}

	case "boolean":

		if valOneType == "hardcoded" {
			dataConvertionVal1 = `
											` + VariableOneName + `, ` + VariableOneError + ` := strconv.ParseBool("` + valueOne + `")
											if ` + VariableOneError + ` != nil {
												logger.Log_WF("Value 1 type convertion failed.", logger.Debug, flowData["InSessionID"].(string))
											}
											`
		} else if valOneType == "dynamic" {
			dataConvertionVal1 = `
												` + VariableOneName + `, ` + VariableOneError + ` := strconv.ParseBool(strconv.FormatBool(flowData["` + valueOne + `"].(bool)))
												if ` + VariableOneError + ` != nil {
													logger.Log_WF("Value 1 type convertion failed.", logger.Debug, flowData["InSessionID"].(string))
												}
												`
		}

		if valTwoType == "hardcoded" {
			dataConvertionVal2 = `
												` + VariableTwoName + `, ` + VariableTwoError + ` := strconv.ParseBool("` + valueTwo + `")
												if ` + VariableTwoError + ` != nil {
													logger.Log_WF("Value 2 type convertion failed.", logger.Debug, flowData["InSessionID"].(string))
													}`
		} else if valTwoType == "dynamic" {
			dataConvertionVal2 = `
														` + VariableTwoName + `, ` + VariableTwoError + ` := strconv.ParseBool(strconv.FormatBool(flowData["` + valueTwo + `"].(bool)))
														if ` + VariableTwoError + ` != nil {
															logger.Log_WF("Value 2 type convertion failed.", logger.Debug, flowData["InSessionID"].(string))
															}`
		}

	default:
		dataConvertionVal1 = ""
		dataConvertionVal2 = ""
	}

	ifCode := `
																	` + dataConvertionVal1 + `
																	` + dataConvertionVal2 + `
																	if(` + VariableOneName + conditionParameter + VariableTwoName + `){logger.Log_WF("", logger.Debug, flowData["InSessionID"].(string));logger.Log_WF("Entering True side of: ` + condition.ID + `", logger.Debug, flowData["InSessionID"].(string));logger.Log_WF("Condition: ` + valueOne + ` ` + conditionParameter + ` ` + valueTwo + `", logger.Debug, flowData["InSessionID"].(string));logger.Log_WF("", logger.Debug, flowData["InSessionID"].(string));
//initiate.View.` + condition.True + `.0
																	}else{logger.Log_WF("", logger.Debug, flowData["InSessionID"].(string));logger.Log_WF("Entering False side of: ` + condition.ID + `", logger.Debug, flowData["InSessionID"].(string));logger.Log_WF("Condition: ` + valueOne + ` ` + conditionParameter + ` ` + valueTwo + `", logger.Debug, flowData["InSessionID"].(string));logger.Log_WF("", logger.Debug, flowData["InSessionID"].(string));
//initiate.View.` + condition.False + `.0
																}
																`
	return ifCode
}
