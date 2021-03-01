package scheduler

import (
	"FileMango/src/db"
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
	db.DequeueFile(msg.Input.Data)
	fmt.Printf("wrote: %q to %q attributes and dequeued\n", msg.Output.Pairs, msg.Input.Data)
}
