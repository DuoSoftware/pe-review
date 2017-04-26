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
func NewWFPublishWorker(id int, workerQueue chan chan Components.JsonFlow) WFPublishWorker {
	// Create, and return the worker.
	worker := WFPublishWorker{
		ID:          id,
		Work:        make(chan Components.JsonFlow),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

type WFPublishWorker struct {
	ID          int
	Work        chan Components.JsonFlow
	WorkerQueue chan chan Components.JsonFlow
	QuitChan    chan bool
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w WFPublishWorker) Start() {
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
				logger.Log_PE("Method: Build workflow", logger.Debug, work.SessionID)
				logger.Log_PE("Flow name recieved : "+work.FlowName, logger.Debug, work.SessionID)
				logger.Log_PE("Session ID recieved : "+work.SessionID, logger.Debug, work.SessionID)
				//logger.Log_PE("JSON string recieved : "+work, logger.Debug, work.SessionID)

				result := Components.InitializeGoFlow(work, work.FlowName, work.SessionID, false)
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
func (w WFPublishWorker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
