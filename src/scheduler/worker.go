package scheduler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

func createWorker(msg chan message, programPath string) chan message {
	var cmd *exec.Cmd

	cmd = exec.Command(programPath)
	in, _ := cmd.StdinPipe()
	out, _ := cmd.StdoutPipe()
	outChan := readIntoChannel(out)
	cmd.Start()

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
					outChan <- generateQueryUsageResponse("cpu", fmt.Sprint(getProcessCpu(int(pid))))

				case "mem":
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
	return outChan
}

func generateQueryUsageResponse(t string, val string) message {
	return message{
		Input: input{
			Type: queryUsage,
			Data: t,
		},
		Output: output{
			Pairs: []pair{{
				Key:   t,
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
		return message{input{noop, in}, output{}} //convert invalid json to noop message
	}
	return target
}

//Message Data Structures

type message struct {
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

type messageType int

const (
	noop messageType = iota
	analyze
	suspend
	resume
	stop
	queryUsage
)
