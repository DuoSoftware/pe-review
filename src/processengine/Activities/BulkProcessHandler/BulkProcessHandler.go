package BulkProcessHandler

import "processengine/context"
import "processengine/logger"
import "github.com/tealeg/xlsx"
import "errors"
import duoCommon "duov6.com/common"
import "io/ioutil"
import "encoding/csv"
import "os"
import "strings"
import "reflect"

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

	var err error
	var filePath string
	var requiredFields []string
	var toBeDeletedFields []string
	var numberOfRows int
	var skipRows int
	skipRows = 0
	fileType := "excel" //default type

	exceldata := make([]map[string]interface{}, 0)

	if FlowData["Take"] != nil {
		numberOfRows = FlowData["Take"].(int)
	} else {
		// 0 means output all rows. dont limit
		numberOfRows = 0
	}

	if FlowData["FileType"] != nil {
		fileType = FlowData["FileType"].(string)
	}

	if FlowData["Skip"] != nil {
		skipRows = FlowData["Skip"].(int)
	}

	if FlowData["RequiredFields"] != nil {
		tokens := strings.Split(FlowData["RequiredFields"].(string), ",")
		requiredFields = make([]string, len(tokens))

		for x := 0; x < len(tokens); x++ {
			value := tokens[x]
			value = strings.TrimSpace(value)
			requiredFields[x] = value
		}
	}

	if FlowData["FilePath"] == nil && FlowData["FileContent"] == nil {
		err = errors.New("No CSV or Excel File is found in Input Parameters.")
	} else {
		if FlowData["FilePath"] == nil {
			var fileName string
			if fileType == "excel" {
				fileName = duoCommon.GetGUID() + ".xlsx"
			} else {
				fileName = duoCommon.GetGUID() + ".csv"
			}
			if reflect.TypeOf(FlowData["FileContent"]).String() == "string" {
				ioutil.WriteFile(fileName, []byte(FlowData["FileContent"].(string)), 0666)
				filePath = fileName
			} else {
				ioutil.WriteFile(fileName, FlowData["FileContent"].([]byte), 0666)
				filePath = fileName
			}
		} else {
			filePath = FlowData["FilePath"].(string)
		}

		if filePath != "" {
			if strings.Contains(filePath, ".xlsx") {
				exceldata, toBeDeletedFields = ReadExcelFile(filePath)
			} else {
				exceldata, toBeDeletedFields = ReadCsvFile(filePath)
			}
		}
	}

	if len(requiredFields) != 0 {
		//Remove required fields from ToBeDeletedFields Array

		for x := 0; x < len(toBeDeletedFields); x++ {
			for y := 0; y < len(requiredFields); y++ {
				if toBeDeletedFields[x] == requiredFields[y] {
					toBeDeletedFields = append(toBeDeletedFields[:x], toBeDeletedFields[x+1:]...)
				}
			}
		}

		//Delete fields from map

		for index, _ := range exceldata {
			for _, value := range toBeDeletedFields {
				delete(exceldata[index], value)
			}
		}
	}

	if err == nil {
		msg := "Successfully Read File!"
		activityContext.ActivityStatus = true
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		FlowData["custMsg"] = msg
		FlowData["status"] = "true"
		FlowData["Data"] = exceldata[skipRows : skipRows+numberOfRows]
	} else {
		msg := "Error Reading File : " + err.Error()
		activityContext.ActivityStatus = false
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		FlowData["custMsg"] = msg
		FlowData["status"] = "false"
	}

	return FlowData, activityContext
}

func ReadExcelFile(filePath string) (exceldata []map[string]interface{}, toBeDeletedFields []string) {
	colunmcount := 0
	rowcount := 0
	var colunName []string

	//file read
	xlFile, err := xlsx.OpenFile(filePath)

	if err == nil {
		for _, sheet := range xlFile.Sheets {
			rowcount = (sheet.MaxRow - 1)
			colunmcount = sheet.MaxCol
			colunName = make([]string, colunmcount)
			toBeDeletedFields = make([]string, colunmcount)
			for _, row := range sheet.Rows {
				for j, cel := range row.Cells {
					colunName[j] = cel.String()
					toBeDeletedFields[j] = cel.String()
				}
				break
			}

			exceldata = make(([]map[string]interface{}), rowcount)
			if err == nil {
				for _, sheet := range xlFile.Sheets {
					for rownumber, row := range sheet.Rows {
						currentRow := make(map[string]interface{})
						if rownumber != 0 {
							exceldata[rownumber-1] = currentRow
							for cellnumber, cell := range row.Cells {
								if cellnumber == 0 {
									exceldata[rownumber-1][colunName[cellnumber]] = cell.String()
								} else if cell.Type() == 0 {
									exceldata[rownumber-1][colunName[cellnumber]] = cell.String()
								} else if cell.Type() == 2 {
									dd, _ := cell.Float()
									exceldata[rownumber-1][colunName[cellnumber]] = float64(dd)
								} else if cell.Type() == 3 {
									exceldata[rownumber-1][colunName[cellnumber]] = cell.Bool()
								} else {
									exceldata[rownumber-1][colunName[cellnumber]] = cell.String()
								}
							}
						}
					}
				}
			}
		}

	}
	return
}

func ReadCsvFile(path string) (exceldata []map[string]interface{}, toBeDeletedFields []string) {
	csvfile, err := os.Open(path)

	if err != nil {
		return
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = -1

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		return
	}

	fieldNames := make([]string, len(rawCSVdata[0]))

	exceldata = make([]map[string]interface{}, (len(rawCSVdata) - 1))

	for x := 0; x < len(rawCSVdata[0]); x++ {
		fieldNames[x] = rawCSVdata[0][x]
	}

	toBeDeletedFields = fieldNames

	for x := 1; x < len(rawCSVdata); x++ {
		currentRow := make(map[string]interface{})
		exceldata[x-1] = currentRow
		for cellnumber := 0; cellnumber < len(rawCSVdata[x]); cellnumber++ {
			exceldata[x-1][fieldNames[cellnumber]] = rawCSVdata[x][cellnumber]
		}
	}
	return
}
