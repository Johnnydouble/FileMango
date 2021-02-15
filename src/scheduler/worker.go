package scheduler

func runJob(msg chan<- Message) {

}

type Message struct {
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
