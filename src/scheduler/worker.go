package scheduler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

func createWorker(msg chan message, programPath string) <-chan string {
	var cmd *exec.Cmd
	done := make(chan bool, 2)

	cmd = exec.Command(programPath)
	in, _ := cmd.StdinPipe()
	out, _ := cmd.StdoutPipe()
	outChan := readIntoChannel(out, done)
	cmd.Start()

	//communicate with worker
	go func() {
		for {
			msgObj := <-msg
			switch msgObj.Type {
			case analyze, resume, suspend:
				in.Write(marshalMessage(msgObj))
				in.Write([]byte("\n"))
			default:
				if msgObj.Type == stop {
					in.Write(marshalMessage(msgObj))
					done <- true

					out.Close()
					fmt.Println("RECIEVED")
				}
			}
		}
	}()
	return outChan
}

type message struct {
	Type messageType
	Data string
}

type messageType int

const (
	analyze messageType = iota
	suspend
	resume
	stop
)

func marshalMessage(msg message) []byte {
	message, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return message
}

func readIntoChannel(rc io.ReadCloser, done <-chan bool) chan string {
	out := make(chan string)
	go func() {
		reader := bufio.NewScanner(rc)
		for {
			if !reader.Scan() {
				fmt.Println("done")
				return
			}
			fmt.Println("line read...")
			select {
			case out <- reader.Text():
			case <-done:
				fmt.Println("done")
				return
			}
		}
	}()
	return out
}
