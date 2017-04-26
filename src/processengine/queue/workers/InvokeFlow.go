package queue

import (
	"fmt"
	"github.com/fatih/color"
	"processengine/Components"
	"processengine/logger"
	"strconv"
)

//"encoding/json"

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewInvokeWorker(id int, workerQueue chan chan Components.InvokeStruct) InvokehWorker {
	// Create, and return the worker.
	worker := InvokehWorker{
		ID:          id,
		Work:        make(chan Components.InvokeStruct),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

type InvokehWorker struct {
	ID          int
	Work        chan Components.InvokeStruct
	WorkerQueue chan chan Components.InvokeStruct
	QuitChan    chan bool
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w InvokehWorker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:

				white := color.New(color.FgWhite)
				redBackground := white.Add(color.BgRed)
				redBackground.Println("********************************************************************************")

				logger.Log_PE("Worker ID: "+strconv.Itoa(w.ID), logger.Debug, work.SessionID)
				logger.Log_PE("Method: Invoke workflow", logger.Debug, work.SessionID)
				logger.Log_PE("AppCode : "+work.AppCode, logger.Debug, work.SessionID)
				logger.Log_PE("ProcessCode: "+work.ProcessCode, logger.Debug, work.SessionID)
				logger.Log_PE("SessionID: "+work.SessionID, logger.Debug, work.SessionID)
				//logger.Log_PE("JSON string recieved : "+work, work.SessionID)
				makeExecutable := false
				result := Components.InvokeWorkflow(work, makeExecutable)
				work.ResponseMessage <- result

			case <-w.QuitChan:
				// We have been asked to stop.
				fmt.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w InvokehWorker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
