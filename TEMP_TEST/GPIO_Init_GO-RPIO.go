package gpio

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os"
	"time"
)

var (
	pinFlowSensor = rpio.Pin(23)
	pinSolenoid   = rpio.Pin(24)
)

func GPIO_INIT() {

	//Enable/Open GPIO
	if err := rpio.Open(); err != nil {

		fmt.Println(err)

		os.Exit(1)

	}

	time.Sleep(time.Second / 5)

	//Configure Flow Sensor GPIO pin for input and PullUp
	pinFlowSensor.Input()
	pinFlowSensor.PullUp()
	//ALSO TRY pinFlowSensor.PullDown() if having issues!

	//Configure solenoid GPIO pin for output and set LOW to start
	pinSolenoid.Output()
	pinSolenoid.Low()

}
