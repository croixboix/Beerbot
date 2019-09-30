package main

import (
	"fmt"
	"os"

	"github.com/stianeikeland/go-rpio"
)

var (
	pin = rpio.Pin(14)
)

func main() {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer rpio.Close()

	pin.Input()
	pin.PullUp()
	pin.Detect(rpio.FallEdge)

	fmt.Println("Start Pourin")
	i := 0
	for {
		if pin.EdgeDetected() {
			i++
			fmt.Println(i)

		}
	}
	pin.Detect(rpio.NoEdge)
}
