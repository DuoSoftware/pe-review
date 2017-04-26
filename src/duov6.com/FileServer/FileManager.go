package FileServer

import (
	"duov6.com/FileServer/messaging"
	"duov6.com/common"
	"duov6.com/objectstore/client"
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/toqueteos/webbrowser"
	"io"
	"io/ioutil"
	"os"
	//"time"
	//"path/filepath"
	//"strconv"
	"strings"
)

type FileManager struct {
}

type FileData struct {
	Id       string
	FileName string
	Body     string
}

func (f *FileManager) Store(request *messaging.FileRequest) messaging.FileResponse { // store disk on database

	fileResponse := messaging.FileResponse{}

	if len(request.Body) == 0 {

		//WHEN REQUEST COMES FROM A REST INTERFACE
		file, header, err := request.WebRequest.FormFile("file")

		if err != nil {
			fileResponse.IsSuccess = false
			fileResponse.Message = err.Error()
		}

		out, err := os.Create(header.Filename)
		if err != nil {
			fileResponse.IsSuccess = false
			fileResponse.Message = err.Error()
		}

		// write the content from POST to the file
		_, err = io.Copy(out, file)
		if err != nil {
			fileResponse.IsSuccess = false
			fileResponse.Message = err.Error()
		}

		file2, err2 := ioutil.ReadFile(header.Filename)

		if err2 != nil {
			fileResponse.IsSuccess = false
			fileResponse.Message = err.Error()
		}

		convertedBody := string(file2[:])
		base64Body := common.EncodeToBase64(convertedBody)

		//Create a instance of file struct
		obj := FileData{}
		obj.Id = request.Parameters["id"]
		obj.FileName = header.Filename
		obj.Body = base64Body

		var extraMap map[string]interface{}
		extraMap = make(map[string]interface{})
		extraMap["File"] = "excelFile"

		fmt.Println("Namespace : " + request.Parameters["namespace"])
		fmt.Println("Class : " + request.Parameters["class"])

		uploadContext := strings.ToLower(request.Parameters["fileContent"])

		isRawFile := false
		isIndividualData := false

		if uploadContext == "" || uploadContext == "both" || uploadContext == "raw" {
			isRawFile = true
		}
		if uploadContext == "" || uploadContext == "both" || uploadContext == "data" {
			isIndividualData = true
		}

		if isIndividualData {
			fmt.Println("Saving INDIVIDUAL DATA inside file.......... ")
			if checkIfFile(header.Filename) == "xlsx" {
				isRawFile = false
				status := SaveExcelEntries(header.Filename, request)
				if status == true {
					fmt.Println("Individual Records Saved Successfully!")
				} else {
					fmt.Println("Saving Individual Records Failed!")
				}
			}
		}

		var returnParams []map[string]interface{}
		returnParams = make([]map[string]interface{}, 1)
		if isRawFile {
			fmt.Println("Saving the RAW file.......... ")
			returnParams = client.GoExtra(request.Parameters["securityToken"], request.Parameters["namespace"], request.Parameters["class"], extraMap).StoreObject().WithKeyField("Id").AndStoreOne(obj).FileOk()
			if len(returnParams) > 0 {
				fmt.Fprintf(request.WebResponse, returnParams[0]["ID"].(string))
			} else {
				fmt.Fprintf(request.WebResponse, "FAILED!")
			}
		} else {
			fmt.Fprintf(request.WebResponse, header.Filename)
		}

		//close the files
		err = out.Close()
		err = file.Close()

		if err != nil {
			fileResponse.IsSuccess = false
			fileResponse.Message = err.Error()
		}

		//remove the temporary stored file from the disk
		err2 = os.Remove(header.Filename)

		if err2 != nil {
			fileResponse.IsSuccess = false
			fileResponse.Message = err2.Error()
		}

		if err == nil && err2 == nil {
			fileResponse.IsSuccess = true
			fileResponse.Message = "Storing file successfully completed"
		} else {
			fileResponse.IsSuccess = false
			fileResponse.Message = "Storing file was unsuccessfull!" + "\n" + err.Error() + "\n" + err2.Error()
		}

	} else {

		//WHEN REQUEST COMES FROM A NON REST INTERFACE
		convertedBody := string(request.Body[:])
		base64Body := common.EncodeToBase64(convertedBody)

		//store file in the DB as a single file
		obj := FileData{}
		obj.Id = request.Parameters["id"]
		obj.FileName = request.FileName
		obj.Body = base64Body

		client.Go(request.Parameters["securityToken"], request.Parameters["namespace"], request.Parameters["class"]).StoreObject().WithKeyField("Id").AndStoreOne(obj).FileOk()

		fileResponse.IsSuccess = true
		fileResponse.Message = "Storing file successfully completed"

	}

	return fileResponse
}

func (f *FileManager) Remove(request *messaging.FileRequest) messaging.FileResponse { // remove file from disk and database
	fileResponse := messaging.FileResponse{}
	var saveServerPath string = request.RootSavePath
	file, err := ioutil.ReadFile(saveServerPath + request.FilePath + request.FileName)

	if len(file) > 0 {
		err = os.Remove(saveServerPath + request.FilePath + request.FileName)
	}

	if err == nil {
		fileResponse.IsSuccess = true
		fileResponse.Message = "Deletion of file successfully completed"
	} else {
		fileResponse.IsSuccess = true
		fileResponse.Message = "Deletion of file Aborted"
	}

	obj := FileData{}
	obj.Id = request.Parameters["id"]
	obj.FileName = request.FileName

	client.Go(request.Parameters["securityToken"], request.Parameters["namespace"], request.Parameters["class"]).StoreObjectWithOperation("delete").WithKeyField("Id").AndStoreOne(obj).Ok()
	fileResponse.IsSuccess = true
	fileResponse.Message = "Deletion of file successfully completed"

	return fileResponse
}

func (f *FileManager) Download(request *messaging.FileRequest) messaging.FileResponse { // save the file to ftp and download via ftp on browser
	fileResponse := messaging.FileResponse{}

	if len(request.Body) == 0 {

	} else {
		var saveServerPath string = request.RootSavePath
		var accessServerPath string = request.RootGetPath

		file := FileData{}
		json.Unmarshal(request.Body, &file)

		temp := common.DecodeFromBase64(file.Body)
		ioutil.WriteFile((saveServerPath + request.FilePath + file.FileName), []byte(temp), 0666)
		err := webbrowser.Open(accessServerPath + request.FilePath + file.FileName)
		if err != nil {
			fileResponse.IsSuccess = false
			fileResponse.Message = "Downloading Failed!" + err.Error()
		} else {
			fileResponse.IsSuccess = true
			fileResponse.Message = "Downloading file successfully completed"
		}
	}

	return fileResponse
}

//huehuehue
func SaveExcelEntries(excelFileName string, request *messaging.FileRequest) bool {
	fmt.Println("Inserting Records to Database....")
	rowcount := 0
	colunmcount := 0
	var exceldata []map[string]interface{}
	var colunName []string

	//file read
	xlFile, error := xlsx.OpenFile(excelFileName)
	if error == nil {
		for _, sheet := range xlFile.Sheets {
			rowcount = (sheet.MaxRow - 1)
			colunmcount = sheet.MaxCol
			colunName = make([]string, colunmcount)
			for _, row := range sheet.Rows {
				for j, cel := range row.Cells {
					colunName[j] = cel.String()
				}
				break
			}
			exceldata = make(([]map[string]interface{}), rowcount)
			if error == nil {
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

			Id := colunName[0]
			var extraMap map[string]interface{}
			extraMap = make(map[string]interface{})
			extraMap["File"] = "exceldata"
			fmt.Println("Namespace : " + request.Parameters["namespace"])
			fmt.Println("Keyfield : " + Id)
			fmt.Println("filename : " + getExcelFileName(excelFileName))

			client.GoExtra(request.Parameters["securityToken"], request.Parameters["namespace"], getExcelFileName(excelFileName), extraMap).StoreObject().WithKeyField(Id).AndStoreMapInterface(exceldata).Ok()
			return true
		}

	}
	return false
}

//Original - working huge overhead
/*func SaveExcelEntries(excelFileName string, request *messaging.FileRequest) bool {
	fmt.Println("Inserting Records to Database....")
	rowcount := 0
	colunmcount := 0
	var exceldata []map[string]interface{}
	var colunName []string

	//file read
	xlFile, error := xlsx.OpenFile(excelFileName)
	if error == nil {
		for _, sheet := range xlFile.Sheets {
			rowcount = (sheet.MaxRow - 1)
			colunmcount = sheet.MaxCol
			colunName = make([]string, colunmcount)
			for _, row := range sheet.Rows {
				for j, cel := range row.Cells {
					colunName[j] = cel.String()
				}
				break
			}
			exceldata = make(([]map[string]interface{}), rowcount)
			if error == nil {
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

			Id := colunName[0]
			var extraMap map[string]interface{}
			extraMap = make(map[string]interface{})
			extraMap["File"] = "exceldata"
			fmt.Println("Namespace : " + request.Parameters["namespace"])
			fmt.Println("Keyfield : " + Id)
			fmt.Println("filename : " + getExcelFileName(excelFileName))

			noOfElementsPerSet, _ := strconv.Atoi(request.Parameters["BlockSize"])
			fmt.Println("---------------")
			fmt.Println(noOfElementsPerSet)
			fmt.Println("---------------")
			noOfSets := (len(exceldata) / noOfElementsPerSet)
			remainderFromSets := 0
			remainderFromSets = (len(exceldata) - (noOfSets * noOfElementsPerSet))

			startIndex := 0
			stopIndex := noOfElementsPerSet

			for x := 0; x < noOfSets; x++ {
				client.GoExtra(request.Parameters["securityToken"], request.Parameters["namespace"], getExcelFileName(excelFileName), extraMap).StoreObject().WithKeyField(Id).AndStoreMapInterface(exceldata[startIndex:stopIndex]).Ok()
				startIndex += noOfElementsPerSet
				stopIndex += noOfElementsPerSet
				if noOfElementsPerSet < 500 {
					time.Sleep(2 * time.Second)
				}
			}

			if remainderFromSets > 0 {
				start := len(exceldata) - remainderFromSets
				client.GoExtra(request.Parameters["securityToken"], request.Parameters["namespace"], getExcelFileName(excelFileName), extraMap).StoreObject().WithKeyField(Id).AndStoreMapInterface(exceldata[start:len(exceldata)]).Ok()
				if noOfElementsPerSet < 500 {
					time.Sleep(2 * time.Second)
				}
			}
			return true
		}

	}
	return false
}*/

//working new
/*func SaveExcelEntries(excelFileName string, request *messaging.FileRequest) bool {
	fmt.Println("Inserting Records to Database....")
	rowcount := 0
	colunmcount := 0
	var exceldata []map[string]interface{}
	var colunName []string

	blockSizeRecords, _ := strconv.Atoi(request.Parameters["BlockSize"])

	//file read
	xlFile, error := xlsx.OpenFile(excelFileName)
	fmt.Println("File Opened")
	if error == nil {
		for _, sheet := range xlFile.Sheets {
			rowcount = (sheet.MaxRow - 1)
			colunmcount = sheet.MaxCol
			colunName = make([]string, colunmcount)
			for _, row := range sheet.Rows {
				for j, cel := range row.Cells {
					colunName[j] = cel.String()
				}
				break
			}
			fmt.Println(rowcount)
			wholeRowIndex := 1
			//exceldata = make(([]map[string]interface{}), blockSizeRecords-1)
			//exceldata = make(([]map[string]interface{}), rowcount)
			if error == nil {
				for _, sheet := range xlFile.Sheets {
					rowIndex := 1
					for rownumber, row := range sheet.Rows {
						currentRow := make(map[string]interface{})
						if rownumber != 0 {
							//exceldata[rownumber-1] = currentRow
							for cellnumber, cell := range row.Cells {
								if cellnumber == 0 {
									//exceldata[rownumber-1][colunName[cellnumber]] = cell.String()
									currentRow[colunName[cellnumber]] = cell.String()
								} else if cell.Type() == 0 {
									currentRow[colunName[cellnumber]] = cell.String()
									//exceldata[rownumber-1][colunName[cellnumber]] = cell.String()
								} else if cell.Type() == 2 {
									dd, _ := cell.Float()
									//exceldata[rownumber-1][colunName[cellnumber]] = float64(dd)
									currentRow[colunName[cellnumber]] = float64(dd)
								} else if cell.Type() == 3 {
									//exceldata[rownumber-1][colunName[cellnumber]] = cell.Bool()
									currentRow[colunName[cellnumber]] = cell.Bool()
								} else {
									//exceldata[rownumber-1][colunName[cellnumber]] = cell.String()
									currentRow[colunName[cellnumber]] = cell.String()
								}
							}
							exceldata = append(exceldata, currentRow)
							if rowIndex == blockSizeRecords || wholeRowIndex == rowcount {
								fmt.Println(wholeRowIndex)
								//fmt.Println(exceldata)
								Id := colunName[0]
								var extraMap map[string]interface{}
								extraMap = make(map[string]interface{})
								extraMap["File"] = "exceldata"
								client.GoExtra(request.Parameters["securityToken"], request.Parameters["namespace"], getExcelFileName(excelFileName), extraMap).StoreObject().WithKeyField(Id).AndStoreMapInterface(exceldata).Ok()
								rowIndex = 1
								exceldata = nil
							} else {
								rowIndex++
							}
							wholeRowIndex++
						}

					}
				}
			}
			return true
		}

	}
	return false
}
*/
//test code
/*func SaveExcelEntries(excelFileName string, request *messaging.FileRequest) bool {
	fmt.Println("Inserting Records to Database....")

	//split and save massive xcel files to smaller files
	splitExcelData(excelFileName)
	//get filenames in directory
	releventFiles := getReleventExcelSplits(excelFileName)
	rowcount := 0
	colunmcount := 0
	var exceldata []map[string]interface{}
	var colunName []string

	blockSizeRecords, _ := strconv.Atoi(request.Parameters["BlockSize"])

	//file read
	xlFile, error := xlsx.OpenFile(excelFileName)
	fmt.Println("File Opened")
	if error == nil {
		for _, sheet := range xlFile.Sheets {
			rowcount = (sheet.MaxRow - 1)
			colunmcount = sheet.MaxCol
			colunName = make([]string, colunmcount)
			for _, row := range sheet.Rows {
				for j, cel := range row.Cells {
					colunName[j] = cel.String()
				}
				break
			}
			fmt.Println(rowcount)
			wholeRowIndex := 1
			//exceldata = make(([]map[string]interface{}), blockSizeRecords-1)
			//exceldata = make(([]map[string]interface{}), rowcount)
			if error == nil {
				for _, sheet := range xlFile.Sheets {
					rowIndex := 1
					for rownumber, row := range sheet.Rows {
						currentRow := make(map[string]interface{})
						if rownumber != 0 {
							//exceldata[rownumber-1] = currentRow
							for cellnumber, cell := range row.Cells {
								if cellnumber == 0 {
									//exceldata[rownumber-1][colunName[cellnumber]] = cell.String()
									currentRow[colunName[cellnumber]] = cell.String()
								} else if cell.Type() == 0 {
									currentRow[colunName[cellnumber]] = cell.String()
									//exceldata[rownumber-1][colunName[cellnumber]] = cell.String()
								} else if cell.Type() == 2 {
									dd, _ := cell.Float()
									//exceldata[rownumber-1][colunName[cellnumber]] = float64(dd)
									currentRow[colunName[cellnumber]] = float64(dd)
								} else if cell.Type() == 3 {
									//exceldata[rownumber-1][colunName[cellnumber]] = cell.Bool()
									currentRow[colunName[cellnumber]] = cell.Bool()
								} else {
									//exceldata[rownumber-1][colunName[cellnumber]] = cell.String()
									currentRow[colunName[cellnumber]] = cell.String()
								}
							}
							exceldata = append(exceldata, currentRow)
							if rowIndex == blockSizeRecords || wholeRowIndex == rowcount {
								fmt.Println(wholeRowIndex)
								//fmt.Println(exceldata)
								Id := colunName[0]
								var extraMap map[string]interface{}
								extraMap = make(map[string]interface{})
								extraMap["File"] = "exceldata"
								client.GoExtra(request.Parameters["securityToken"], request.Parameters["namespace"], getExcelFileName(excelFileName), extraMap).StoreObject().WithKeyField(Id).AndStoreMapInterface(exceldata).Ok()
								rowIndex = 1
								exceldata = nil
							} else {
								rowIndex++
							}
							wholeRowIndex++
						}

					}
				}
			}
			return true
		}

	}
	return false
}*/

func checkIfFile(params string) (fileType string) {
	var tempArray []string
	tempArray = strings.Split(params, ".")
	if len(tempArray) > 1 {
		fileType = tempArray[len(tempArray)-1]
	} else {
		fileType = "NAF"
	}
	return
}

func getExcelFileName(path string) (fileName string) {
	subsets := strings.Split(path, "\\")
	subfilenames := strings.Split(subsets[len(subsets)-1], ".")
	fileName = subfilenames[0]
	return
}

/*func createExcelFile(fileName string, columns []string, data []map[string]interface{}) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	file = xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}

	row = sheet.AddRow()
	for x := 0; x < len(columns); x++ {
		cell = row.AddCell()
		cell.Value = columns[x]
	}

	for x := 0; x < len(data); x++ {
		row = sheet.AddRow()
		for y := 0; y < len(columns); y++ {
			cell = row.AddCell()
			cell.Value = data[x][columns[y]].(string)
		}
	}

	err = file.Save(fileName)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func splitExcelData(fileName string) {
	rowcount := 0
	colunmcount := 0
	var exceldata []map[string]interface{}
	var colunName []string
	blockSizeRecords := 50000
	blockindex := 0

	xlFile, err := xlsx.OpenFile(fileName)

	if err == nil {
		for _, sheet := range xlFile.Sheets {
			rowcount = (sheet.MaxRow - 1)
			colunmcount = sheet.MaxCol
			colunName = make([]string, colunmcount)
			for _, row := range sheet.Rows {
				for j, cel := range row.Cells {
					colunName[j] = cel.String()
				}
				break
			}
			wholeRowIndex := 1
			if err == nil {
				for _, sheet := range xlFile.Sheets {
					rowIndex := 1
					for rownumber, row := range sheet.Rows {
						currentRow := make(map[string]interface{})
						if rownumber != 0 {
							for cellnumber, cell := range row.Cells {
								currentRow[colunName[cellnumber]] = cell.String()
							}
							exceldata = append(exceldata, currentRow)
							if rowIndex == blockSizeRecords || wholeRowIndex == rowcount {
								createExcelFile((fileName + strconv.Itoa(blockindex) + ".xlsx"), colunName, exceldata)
								rowIndex = 1
								exceldata = nil
								blockindex++
							} else {
								rowIndex++
							}
							wholeRowIndex++
						}

					}
				}
			}
		}

	}
}

func getReleventExcelSplits(fileName string) []string {
	fileName = strings.TrimSpace(fileName)
	tokens := strings.Split(fileName, ".")
	files1, _ := filepath.Glob("*" + tokens[0])
	return files1
}
*/