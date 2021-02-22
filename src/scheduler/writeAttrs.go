package scheduler

import (
	"fmt"
	"github.com/pkg/xattr"
)

func writeAttributes(msg message) {
	for _, pair := range msg.Output.Pairs {
		_ = xattr.Set(msg.Input.Data, pair.Key, []byte(pair.Value))
	}
	fmt.Print("WROTE: ")
	fmt.Print(msg)
	fmt.Println(" TO ATTRIBUTES")
}
