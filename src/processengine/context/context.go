//This package will be used each of activity to share the success status of the activity and also share variables
//	within the context of the workflow
// Third party references taken from  - https://github.com/gorilla/context

package context

import "net/http"
import "sync"

// Public variables to handle the context
var (
	mutex sync.RWMutex
	data  = make(map[*http.Request]map[interface{}]interface{})
)

// Stores a value for a given key in a given request
func Set(r *http.Request, key, val interface{}) {
	mutex.Lock()
	if data[r] == nil {
		data[r] = make(map[interface{}]interface{})
	}
	data[r][key] = val
	mutex.Unlock()
}

// Returns a value stored for a given key in a given request
func Get(r *http.Request, key interface{}) interface{} {
	mutex.RLock()
	if ctx := data[r]; ctx != nil {
		value := ctx[key]
		mutex.RUnlock()
		return value
	}
	mutex.RUnlock()
	return nil
}

// GetOk returns stored value and presence state like multi-value return of map access.
func GetOk(r *http.Request, key interface{}) (interface{}, bool) {
	mutex.RLock()
	if _, ok := data[r]; ok {
		value, ok := data[r][key]
		mutex.RUnlock()
		return value, ok
	}
	mutex.RUnlock()
	return nil, false
}

// GetAll returns all stored values for the request as a map. Nil is returned for invalid requests.
func GetAll(r *http.Request) map[interface{}]interface{} {
	mutex.RLock()
	if context, ok := data[r]; ok {
		result := make(map[interface{}]interface{}, len(context))
		for k, v := range context {
			result[k] = v
		}
		mutex.RUnlock()
		return result
	}
	mutex.RUnlock()
	return nil
}

// GetAllOk returns all stored values for the request as a map and a boolean value that indicates if
// the request was registered.
func GetAllOk(r *http.Request) (map[interface{}]interface{}, bool) {
	mutex.RLock()
	context, ok := data[r]
	result := make(map[interface{}]interface{}, len(context))
	for k, v := range context {
		result[k] = v
	}
	mutex.RUnlock()
	return result, ok
}

// Delete removes a value stored for a given key in a given request.
func Delete(r *http.Request, key interface{}) {
	mutex.Lock()
	if data[r] != nil {
		delete(data[r], key)
	}
	mutex.Unlock()
}

// Clear removes all values stored for a given request.
//
// This is usually called by a handler wrapper to clean up request
// variables at the end of a request lifetime. See ClearHandler().
func Clear(r *http.Request) {
	mutex.Lock()
	clear(r)
	mutex.Unlock()
}

// clear is Clear without the lock.
func clear(r *http.Request) {
	delete(data, r)
}

// Purge removes request data stored for longer than maxAge, in seconds.
// It returns the amount of requests removed.
//
// If maxAge <= 0, all request data is removed.
//
// This is only used for sanity check: in case context cleaning was not
// properly set some request data can be kept forever, consuming an increasing
// amount of memory. In case this is detected, Purge() must be called
// periodically until the problem is fixed.
//func Purge(maxAge int) int {
//	mutex.Lock()
//	count := 0
//	if maxAge <= 0 {
//		count = len(data)
//		data = make(map[*http.Request]map[interface{}]interface{})
//	} else {
//		min := time.Now().Unix() - int64(maxAge)
//		for r := range data {
//			if datat[r] < min {
//				clear(r)
//				count++
//			}
//		}
//	}
//	mutex.Unlock()
//	return count
//}

// ClearHandler wraps an http.Handler and clears request values at the end
// of a request lifetime.
func ClearHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer Clear(r)
		h.ServeHTTP(w, r)
	})
}

//Activity error struscture used to give error information of the activity context,
//	ErrorString(string) is used to put the well formated error
//	Forward(bool) is used to mention whether the current error need to forward to next linked activites
//	SeverityLevel(Severity type) is used to define the level of error, this will be used to define the termination process of the workflow
//	Encrypt(bool) is used to mention whether the ErrorString need to encrypt or not

type ActivityError struct {
	ErrorString   string
	Forward       bool
	SeverityLevel Severity
	Encrypt       bool
}

//Predefined severity levels (You must be carefull when putting these values, in future workflow termination process will be defined according to this level)
//	Info will define the ErrorString of the ActivityError structure is just an information message
//	Warning will define the ErrorString of the ActivityError structure is a warning and need user action
//	Error will define the ErrorString of the ActivityError structure is a error and need user action
//	Critical will define the ErrorString of the ActivityError structure is critical and need user action

type Severity int

const (
	Info Severity = 1 + iota
	Warning
	Error
	Critical
)

// Activity contet structure, All activities written in go should use this structure and modify the the values accordingly
//	 ActivityStatus(bool) is used to get or set the status of the specific activity
//	 ResultMessage(string) is used to get or set the result message of the specific activity
//	 CustomerMessage(string) is used to get or set the message which is giving out from the workflow to other sources
//	 ErrorState(ActivityError type) is used to get or set the specific activity error information
//	 ErrorState([]byte) json object is used to get or set all shared variables, properties within the context of the workflow,
//	 proprty values can be added or modified according to the activity logic but it is not recommended to delete any proprty
//	 until the workflow get finished, those information will be use for report and analytics data in future

type ActivityContext struct {
	ActivityStatus bool
	Message        string
	ErrorState     ActivityError
	SharedContext  []byte
}

type WorkflowContext struct {
	Message       string                 `json:"Message"`
	Status        bool                   `json:"Status"`
	ErrorCode     int                    `json:"ErrorCode"`
	WorkflowTrace string                 `json:"WorkflowTrace"`
	ActivityTrace map[string]interface{} `json:"ActivityTrace"`
}

type FlowResult struct {
	Message    string     `json:"Message"`
	Status     bool       `json:"Status"`
	FlowName   string     `json:"FlowName"`
	SessionID  string     `json:"SessionID"`
	ReturnData ReturnData `json:"ReturnData"`
}

/*type ReturnData struct {
JSONOutput     map[string]interface{} `json:"JSONOutput"`
WFOutArguments map[string]interface{} `json:"WFOutArguments"`
WorkflowResult WorkflowContext        `json:"WorkflowResult"`
WorkflowTrace  string                 `json:"WorkflowTrace"`
}*/

type ReturnData struct {
	JSONOutput     map[string]interface{} `json:"JSONOutput"`
	WFOutArguments map[string]interface{} `json:"WFOutArguments"`
	WorkflowResult struct {
		ErrorCode    int    `json:"ErrorCode"`
		ErrorDetails string `json:"ErrorDetails"`
		Message      string `json:"Message"`
		Status       bool   `json:"Status"`
	} `json:"WorkflowResult"`
}

/*type ReturnData struct {
	JSONOutput     map[string]interface{} `json:"JSONOutput"`
	WFOutArguments map[string]interface{} `json:"WFOutArguments"`
	WorkflowResult struct {
		ActivityTrace map[string]interface{} `json:"ActivityTrace"`
		ErrorCode     int                    `json:"ErrorCode"`
		Message       string                 `json:"Message"`
		Status        bool                   `json:"Status"`
	} `json:"WorkflowResult"`
	WorkflowTrace string `json:"WorkflowTrace"`
}
*/

type TestWorkflowResponse struct {
	Status       bool       `json:"Status"`
	Message      string     `json:"Message"`
	ErrorDetails string     `json:"ErrorDetails"`
	ErrorCode    int        `json:"ErrorCode"`
	ResponseData ReturnData `json:"ResponseData"`
}

type InstallerResponse struct {
	Status         bool            `json:"Status"`
	Message        string          `json:"Message"`
	PackageDetails []PackageDetail `json:"PackageDetails"`
	ErrorCode      int             `json:"ErrorCode"`
}

type PackageDetail struct {
	PackageName string `json:"PackageName"`
	Status      bool   `json:"Status"`
	Message     string `json:"Message"`
}

type ObjectStoreResponse struct {
	ErrorCode    string                   `json:"errorCode"`
	ErrorLog     []interface{}            `json:"errorLog"`
	ErrorMessage string                   `json:"errorMessage"`
	Result       []map[string]interface{} `json:"result"`
	Success      bool                     `json:"success"`
}
