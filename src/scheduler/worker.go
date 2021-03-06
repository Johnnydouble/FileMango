package scheduler

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func createWorker(msg chan message, programCmd string) (chan message, error) {
	var cmd *exec.Cmd

	programPath, args := processCommand(programCmd)

	cmd = exec.Command(programPath, args...)
	in, _ := cmd.StdinPipe()
	out, _ := cmd.StdoutPipe()
	outChan := readIntoChannel(out)
	if cmd.Start() != nil { //start module
		//if modules fails to start properly
		return outChan, errors.New("Module Failed to Start" + programCmd)
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
		fmt.Println("Module exited successfully")
	}()
	return out
}

func parseMessage(in string) message {
	var target message
	if json.Unmarshal([]byte(in), &target) != nil {
		fmt.Println("Error: Failure to unmarshal JSON from external module.")
		return message{header{}, input{noop, in, ""}, output{}} //convert invalid json to noop message
	}
	return target
}

func processCommand(fullCommand string) (string, []string) {
	cmdParts := strings.Split(fullCommand, " ")

	programPath := ""
	var args []string

	if len(cmdParts) > 0 {
		programPath = cmdParts[0]
		for i := 1; i < len(cmdParts); i++ {
			args = append(args, cmdParts[i])
		}
	}
	return programPath, args
}

//Message Data Structures

type message struct {
	Header header
	Input  input
	Output output
}

type input struct {
	Type    messageType
	Data    string
	ModPath string
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
