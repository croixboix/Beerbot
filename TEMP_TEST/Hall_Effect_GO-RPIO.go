package gpio

import (
	"fmt"
	//"os"
	"github.com/stianeikeland/go-rpio"
	"time"
)

var (
	pinRead = rpio.Pin(23)
)

func Pour() {

	/*if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}*/

	//defer rpio.Close()
	//pinRead.Input()
	//pinRead.PullUp()

	//Enable edge detection
	pinRead.Detect(rpio.FallEdge)

	//TEST DELAY, remove eventually!
	time.Sleep(time.Second / 10)

	fmt.Println("Start Pourin")
	i := 0
	for i < 234 {
		if pinRead.EdgeDetected() {
			i++
			fmt.Println(i)
		}
	}

	//Disable edge detection
	pinRead.Detect(rpio.NoEdge)
	time.Sleep(time.Second / 10)
}
