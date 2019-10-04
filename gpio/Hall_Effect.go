package gpio

import (
	"fmt"
	"os"

	"github.com/stianeikeland/go-rpio"
)

var (
	pin = rpio.Pin(14)
)

/*

	return - 0 if counter has not been met
					 1 if counter has been met
*/
func pour() {
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
	for i < 468 {
		if pin.EdgeDetected() {
			i++
			fmt.Println(i)

		}
	}
	pin.Detect(rpio.NoEdge)
}
