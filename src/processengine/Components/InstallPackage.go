package Components

import (
	"encoding/json"
	"os"
	"os/exec"
	"processengine/context"
	"processengine/logger"
	"runtime"
)

// method used to test the workflow with dummy data
func InstallPackages(PackagesContent PackageInstaller, sessionId string) (flowResult context.InstallerResponse) {

	var slash = ""
	if runtime.GOOS == "windows" {
		slash = "\\"
	} else {
		slash = "/"
	}

	pwd, _ := os.Getwd()

	// initiating default variables
	flowResult.Status = true
	flowResult.Message = "Starting package install method"
	// starting the process
	logger.Log_PE("~~ Starting Package Install method", logger.Information, sessionId)

	// checking the code tempororaily
	tempContent, _ := json.Marshal(PackagesContent)
	logger.Log_PE("Received package data: ", logger.Debug, sessionId)
	logger.Log_PE(string(tempContent), logger.Debug, sessionId)

	// iterate through each array item and install the packages onto the server
	logger.Log_PE("Runtime Platform: Environment - "+string(runtime.GOOS), logger.Debug, sessionId)

	// the loop runs trough the list of packages to install eachone in every iterate
	for _, packagename := range PackagesContent.Content {

		// declaring package detail object
		packageDetail := context.PackageDetail{}
		packageDetail.PackageName = packagename
		packageDetail.Status = true
		packageDetail.Message = "Installing package"

		if runtime.GOOS == "windows" {
			// if execute is running in windows
			out, err := exec.Command("installPackage.bat", packagename, "CMD", "start").Output()

			// check if there is an error or not
			if err != nil {
				msg := "Error on installing package: " + packagename
				flowResult.Status = false
				flowResult.Message = msg
				packageDetail.Message = string(out)
				logger.Log_PE(msg, logger.Error, sessionId)
				logger.Log_PE(err.Error(), logger.Error, sessionId)
			} else {
				msg := "Installation Success"
				flowResult.Status = true
				flowResult.Message = msg
				packageDetail.Message = msg
				logger.Log_PE(msg, logger.Debug, sessionId)
			}
			flowResult.PackageDetails = append(flowResult.PackageDetails, packageDetail)

		} else {
			// if execute is running in linux
			//logger.Log_PE("Runtime Platform: Environment - "+string(runtime.GOOS), logger.Debug, sessionId)
			flowResult.Message = "In linux environment -> "
			out, err := exec.Command(pwd+slash+"installPackage.sh", packagename, "sh").Output()

			// check if there is an error or not
			if err != nil {
				msg := "Error on installing package: " + packagename
				flowResult.Status = false
				flowResult.Message = msg
				packageDetail.Message = string(out)
				logger.Log_PE(msg, logger.Error, sessionId)
				logger.Log_PE(err.Error(), logger.Error, sessionId)
			} else {
				msg := "Installation Success"
				flowResult.Status = true
				flowResult.Message = msg
				packageDetail.Message = msg
				logger.Log_PE(msg, logger.Debug, sessionId)
			}
			flowResult.PackageDetails = append(flowResult.PackageDetails, packageDetail)
		}
	}
	// when the process is complete it passes the response back to the front
	return flowResult
}
