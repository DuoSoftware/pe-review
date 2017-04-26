package main

import (
	//"encoding/json"
	"fmt"
	"objectstore"
)

func main() {

	// n := objectstore.GetAll{}

	// var parameters map[string]interface{}
	// parameters = make(map[string]interface{})
	// parameters["securityToken"] = "securityToken"
	// parameters["log"] = "log"
	// parameters["namespace"] = "com.duosoftware.com"
	// parameters["class"] = "duodigin_dashboardddd"

	// ss := n.Invoke(parameters)
	// fmt.Println(string(ss.SharedContext))
	// fmt.Println(string(ss.ResultMessage))

	//....................................................................

	// n := objectstore.GetByKey{}

	// var parameters map[string]interface{}
	// parameters = make(map[string]interface{})
	// parameters["securityToken"] = "securityToken"
	// parameters["log"] = "log"
	// parameters["namespace"] = "com.duosoftware.com"
	// parameters["class"] = "Inventory"
	// parameters["key"] = "4"

	// ss := n.Invoke(parameters)
	// fmt.Println(string(ss.SharedContext))
	// fmt.Println(string(ss.ResultMessage))

	//.........................................................

	n := objectstore.Insert{} //Insert, Update, Delete

	var parameters map[string]interface{}
	parameters = make(map[string]interface{})
	parameters["securityToken"] = "securityToken"
	parameters["log"] = "log"
	parameters["namespace"] = "com.duosoftware.com"
	parameters["class"] = "Inventory"
	parameters["JSON"] = "{\"Object\":{\"Id\":\"700\", \"Name\":\"Prasad\"}, \"Parameters\":{\"KeyProperty\":\"Id\"}}"

	ss := n.Invoke(parameters)
	fmt.Println(string(ss.ResultMessage))
}
