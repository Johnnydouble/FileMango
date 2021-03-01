package scheduler

import (
	"fmt"
	"github.com/pkg/xattr"
	"log"
)

func writeAttributes(msg message) {
	for _, pair := range msg.Output.Pairs {
		err := xattr.Set(msg.Input.Data, "user."+pair.Key, []byte(pair.Value))
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Print("WROTE: ")
	fmt.Print(msg)
	fmt.Println(" TO ATTRIBUTES")
	//todo: delete file_queue.txt entry
}
