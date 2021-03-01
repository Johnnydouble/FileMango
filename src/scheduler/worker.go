package scheduler

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

func createWorker(msg chan message, programPath string) (chan message, error) {
	var cmd *exec.Cmd

	cmd = exec.Command(programPath)
	in, _ := cmd.StdinPipe()
	out, _ := cmd.StdoutPipe()
	outChan := readIntoChannel(out)
	if cmd.Start() != nil { //start module
		//if modules fails to start properly
		return outChan, errors.New("Module Failed to Start" + programPath)
	}

	//communicate with worker
	go func() {
		for {
			msgObj := <-msg
			switch msgObj.Input.Type {
			case analyze, resume, suspend:
				in.Write(marshalMessage(msgObj))
				in.Write([]byte("\n"))
			case queryUsage:
				switch msgObj.Input.Data {
				case "cpu":
					pid := cmd.Process.Pid
					*msgObj.Header.recipient <- generateInternalResponse(queryUsage, "cpu", fmt.Sprint(getProcessCpu(int(pid))))
				case "mem":
					pid := cmd.Process.Pid
					*msgObj.Header.recipient <- generateInternalResponse(queryUsage, "mem", fmt.Sprint(getProcessMem(int(pid))))
				}
			case noop:
				//yup

			default:
				if msgObj.Input.Type == stop {
					in.Write(marshalMessage(msgObj))
					out.Close()
				}
			}
		}
	}()
	return outChan, nil
}

func generateInternalResponse(msgType messageType, args string, val string) message {
	return message{
		Input: input{
			Type: msgType,
			Data: args,
		},
		Output: output{
			Pairs: []pair{{
				Key:   args,
				Value: val,
			}},
		},
	}
}

func marshalMessage(msg message) []byte {
	message, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return message
}

func readIntoChannel(rc io.ReadCloser) chan message {
	out := make(chan message)
	go func() {
		reader := bufio.NewScanner(rc)
		for reader.Scan() {
			out <- parseMessage(reader.Text())
		}
		close(out)
		fmt.Println("SUCCESS")
	}()
	return out
}

func parseMessage(in string) message {
	var target message
	if json.Unmarshal([]byte(in), &target) != nil {
		fmt.Println("Error: Failure to unmarshal JSON from external module.")
		return message{header{}, input{noop, in}, output{}} //convert invalid json to noop message
	}
	return target
}

//Message Data Structures

type message struct {
	Header header
	Input  input
	Output output
}

type input struct {
	Type messageType
	Data string
}

type output struct {
	Pairs []pair
}

type pair struct {
	Key   string
	Value string
}

type header struct {
	recipient *chan message
}

type messageType int

const (
	noop messageType = iota
	analyze
	suspend
	resume
	stop
	//internal types
	queryUsage
	modErr
)
