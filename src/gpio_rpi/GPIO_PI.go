package gpio_rpi

import (
	"fmt"
	"github.com/warthog618/gpio"
	"os"
	"time"
)

var (
	pinFlowSensor1     = gpio.NewPin(23)
	pinSolenoid1       = gpio.NewPin(24)
	flow1Counter   int = 0
)

func GPIO_INIT() {

	//Enable/Open GPIO
	if err := gpio.Open(); err != nil {

		fmt.Println(err)

		os.Exit(1)

	}

	//time.Sleep(time.Second / 10)

	//Configure Flow Sensor GPIO pin for input and PullUp
	pinFlowSensor1.Input()
	pinFlowSensor1.PullUp()
	//ALSO TRY pinFlowSensor.PullDown() if having issues!

	//Configure solenoid GPIO pin for output and set LOW to start
	pinSolenoid1.Output()
	pinSolenoid1.Low()

}

func Flow1CounterIncrement(*Pin) {
	// handle change in pin value
	flow1Counter++
	fmt.Println(flow1Counter)
}

func Pour() {
	//Enable edge detection
	//pinFlowSensor1.Detect(rpio.FallEdge)

	//Reset flow counter
	flow1Counter = 0

	//Enable wacher
	pinFlowSensor1.Watch(gpio.EdgeFalling, Flow1CounterIncrement())

	//TEST DELAY, remove eventually!
	//time.Sleep(time.Second / 10)

	fmt.Println("Start Pourin")
	for flow1Counter < 234 {
		time.Sleep()
	}

	//Disable edge detection
	//pinFlowSensor1.Detect(rpio.NoEdge)
	//time.Sleep(time.Second / 10)

	//Disable watcher
	pinFlowSensor1.Unwatch()
}

func Toggle() {

	//time.Sleep(time.Second / 10)

	pinSolenoid1.Toggle()

	//time.Sleep(time.Second / 10)

}
