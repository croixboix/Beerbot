package gpio

import (

	//"fmt"

	//"os"

	"time"

	"github.com/stianeikeland/go-rpio"
)

var (
	pin = rpio.Pin(24)
)

func Taggle() {

	/*if err := rpio.Open(); err != nil {

		fmt.Println(err)

		os.Exit(1)

	}*/

	//defer rpio.Close()

	//pin.Output()

	//time.Sleep(time.Second / 5)

	//pin.Low()

	time.Sleep(time.Second / 5)

	pin.Toggle()

	time.Sleep(time.Second / 5)

}
