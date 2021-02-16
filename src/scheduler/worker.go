package scheduler

import (
	"encoding/json"
	"fmt"
)

var testMsg = message{msgType: analyze, msgData: "/home/melkor/Pictures/aur-plugin.png"}

func runJob(msg chan<- message) {
	marshalMessage(testMsg) //test Marshal function
}

type message struct {
	msgType messageType
	msgData string
}

type messageType int

const (
	analyze = iota
	suspend
	resume
	stop
)

func marshalMessage(msg message) string {
	message, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(string(message))
	return string(message)
}
