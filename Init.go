// +build ignore

package main

//Import.Start
import "processengine/logger"
import "processengine/Common" //ImportCommon

//Import.End

func main() {
	//_ = Common.GetResourceClass()
	//Init.Settings

	//INIT.Headers

	//Init.Var

	logger.Log_WF("WORKFLOW STARTED!", logger.Information, flowData["InSessionID"].(string))

	//Func.Start

	//initiate.View.drawboard

	//Func.End
}
