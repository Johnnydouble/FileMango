package scheduler

import "fmt"

func writeAttributes(msg message) {
	fmt.Print("WROTE: ")
	fmt.Print(msg)
	fmt.Println(" TO ATTRIBUTES")
}
