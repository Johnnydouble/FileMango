package scheduler

import (
	"FileMango/src/config"
	"FileMango/src/db"
	"FileMango/src/watch"
	"fmt"
	"io"
	"os"
	"time"
)

var queuePath = "./res/file_queue.txt"
var pool = activePool{jobs: []job{}, poolSize: 2} //initialize pool with size of 2
var queue = fileQueue{}
var metaOutput = make(chan chan message)

//a function that runs analysis on files that are in a given queue file
func RunAnalysis() {
	go func() {
		for {
			path := <-db.Ambassador.Path
			queue.createJob(path)
		}
	}()

	//do an initial check for jobs
	addInitialJobs()
	go handleOutput(metaOutput)
	go poolManager() //periodically calls managePool
}

func poolManager() {
	go func() {
		for {
			managePool()
			time.Sleep(100 * time.Millisecond) //todo: this might be able to be increased if todos in managePool() are resolved
		}
	}()
}

//creates new active jobs when there is more space in the pool's size than there are actual job objects
//called whenever a job may be ready to enter or leave the pool
func managePool() {
	// figure out what the size of the pool should be
	pool.poolSize = determinePoolSize(pool.poolSize)
	for len(pool.jobs) < pool.poolSize {
		if len(queue.preJobs) == 0 {
			return
		}
		//get first jobInfo in fileQueue
		ji := queue.preJobs[0]
		//delete it from the preJobs slice
		queue.preJobs = disorderlyRemove(queue.preJobs, 0)
		//construct a new job using initJob and add it to activePool
		newJob, err := initJob(ji)
		if err != nil { //todo: might want to handle this a bit better in the future by iterating until initJob is successful
			return
		}

		pool.jobs = append(pool.jobs, newJob)
		//push the output channel to a meta channel that is being fanned into the writing system
		metaOutput <- newJob.output
	}
}

func handleOutput(metaOut chan chan message) {
	messages := fanIn(metaOut)

	//err printer
	printModuleErr := func(txt string, msg message) {
		out := fmt.Sprintln(txt, "There is likely an issue with the module.\nOFFENDING MESSAGE: ", msg)
		_, _ = io.WriteString(os.Stdout, out)
	}

	for msg := range messages {
		switch msg.Input.Type {
		case noop:
			printModuleErr("noop type returned from external module.", msg)
		case queryUsage:
			printModuleErr("queryUsage message type is handled internally.", msg)
		case suspend, resume, stop:
			//do nothing for now
		case analyze:
			writeAttributes(msg)
			managePool() //add new jobs if space exists in the pool
		case modErr:
			fmt.Println(msg)
			//todo: check for potential race condition in the future here, messages might need some sort of UUID so that there can be a dialogue between modules and the system about closing.
			managePool()
		default:
			printModuleErr("Unrecognized message type.", msg)
		}
	}
}

//listens on a channel of channels for a new channel and produces a channel of messages
func fanIn(metaChan chan chan message) chan message {
	out := make(chan message)

	//todo: check if this needs to be in a goroutine
	go func() {
		for msgChan := range metaChan {
			msgChan := msgChan
			go func() {
				for msg := range msgChan {
					out <- msg
				}
			}()
		}
	}()
	return out
}

//a list of active jobs and a size that comprise a pool of workers
type activePool struct {
	jobs     []job
	poolSize int
}

//an in-program representation of the file_queue
type fileQueue struct {
	preJobs     []jobInfo
	pathHistory []string
}

//adds jobs that have not already been added / are not currently processing to the fileQueue
func (q *fileQueue) createJob(path string) {
	if !stringSliceContains(q.pathHistory, path) {
		file, _ := os.Open(path)
		q.pathHistory = append(q.pathHistory, path) //add file to history to prevent duplicates
		q.analyzeFile(file)
		file.Close()
	}
}

//export createJob
func AddJob(path string) {
	queue.createJob(path)
}

func addInitialJobs() {
	add := func(key []byte) error {
		AddJob(string(key))
		return nil
	}
	err := db.FoldQueue(add)
	if err != nil {
	}
}

//create jobs with corresponding providers and add them to the specified fileQueue
func (q *fileQueue) analyzeFile(file *os.File) {
	FileTypes := config.GetFileTypes()
	//todo: might want to handle the potential errors...
	fileType, _ := watch.GetFileContentType(file)
	for _, supportedType := range FileTypes {
		if fileType == supportedType.Type {
			for _, module := range supportedType.ModulePaths {
				q.preJobs = append(q.preJobs, jobInfo{module, file.Name()})
			}
		}
	}
}

//information about a given job; used when to constructing a new active job
type jobInfo struct {
	modulePath string
	filePath   string
}

//represents a task currently being handled by an external module
type job struct {
	input  chan message
	output chan message
	info   jobInfo
}

//creates a job given a jobInfo and begins analysis
func initJob(info jobInfo) (job, error) {
	in := make(chan message)
	out, err := createWorker(in, info.modulePath)
	if err != nil {
		close(in) //note that out is already closed
		return job{}, err
	}
	in <- message{
		Input:  input{Type: analyze, Data: info.filePath},
		Output: output{},
	}
	return job{in, out, info}, nil
}

//returns if a string slice contains an element
func stringSliceContains(slice []string, elem string) bool {
	for _, currentElem := range slice {
		if currentElem == elem {
			return true
		}
	}
	return false
}

//removes an element from a slice without preserving indices
func disorderlyRemove(s []jobInfo, i int) []jobInfo {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
