package scheduler

import (
	"FileMango/src/config"
	"FileMango/src/watch"
	"bufio"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"os"
)

var queuePath string = "./res/file_queue.txt"

//a function that runs analysis on files that are in a given queue file
func RunAnalysis() {
	as := analysisSystem{}

	watcher, _ := fsnotify.NewWatcher()
	watcher.Add(queuePath)
	queueFile, _ := os.Open(queuePath)
	scanner := bufio.NewScanner(io.Reader(queueFile))

	//add new jobs
	go func() {
		select {
		//if the file queue is edited, check for any files that have not already begun being processed.
		case event := <-watcher.Events:
			switch event.Op {
			case fsnotify.Write:
				as.addNewJobs(scanner)
			}

		case err := <-watcher.Errors:
			fmt.Println("FILE QUEUE WATCH ERROR: ", err)
		}
	}()

	//fan in program output to single channel
	//todo: im like 98% sure that this doesn't work
	//todo: replace with reflect based solution found here: https://stackoverflow.com/questions/19992334/how-to-listen-to-n-channels-dynamic-select-statement
	aggregate := make(chan chan message)
	go func() {
		for _, job := range as.jobs {
			go func() {
				aggregate <- job.output
			}()
		}
	}()

	go func() {
		select {
		case msg := <-aggregate:
			msginner := <-msg
			writeAttributes(msginner)
		}
	}()
}

//todo: this is bad, i dont like this, jobs should go into a jobQueue before they go here, this need only exist to track what jobs are already being processed
type analysisSystem struct {
	jobs []job
}

type job struct {
	output chan message
	input  <-chan string //this should really probably be a message tbh...
}

func (as analysisSystem) addNewJobs(qF *bufio.Scanner) {
	FileTypes := config.GetFileTypes()
	for qF.Scan() {
		//todo: might want to handle the potential errors...
		file, _ := os.Open(qF.Text())
		fileType, _ := watch.GetFileContentType(file)
		for _, supportedType := range FileTypes {
			if fileType == supportedType.Type {
				for _, module := range supportedType.ModulePaths {
					as.jobs = append(as.jobs, initJob(module, file.Name()))
				}
			}
		}
	}
}

func initJob(modulePath string, filePath string) job {
	in := make(chan message)
	out := createWorker(in, modulePath)
	in <- message{
		Input:  input{Type: analyze, Data: filePath},
		Output: output{},
	}
	return job{in, out}
}
