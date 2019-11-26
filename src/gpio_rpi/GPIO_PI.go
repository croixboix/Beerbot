package gpio_rpi

import (
	"fmt"
	"github.com/warthog618/gpio"
	"os"
	"time"
)

var (
	flow1Counter   int = 0
	flow2Counter   int = 0

	pinFlowSensor1 *gpio.Pin
	pinSolenoid1 *gpio.Pin

	pinFlowSensor2 *gpio.Pin
	pinSolenoid2 *gpio.Pin
)

func GPIO_INIT() {
	//Enable/Open GPIO
	if err := gpio.Open(); err != nil {

		fmt.Println(err)

		os.Exit(1)
	}

	pinFlowSensor1     = gpio.NewPin(23)
	pinSolenoid1       = gpio.NewPin(24)

	pinFlowSensor2     = gpio.NewPin(25)
	pinSolenoid2       = gpio.NewPin(8)

	//Configure Flow Sensor GPIO pin for input and PullUp
	//ALSO TRY pinFlowSensor.PullDown() if having issues!
	pinFlowSensor1.Input()
	pinFlowSensor1.PullUp()
	pinFlowSensor2.Input()
	pinFlowSensor2.PullUp()
	

	//Configure solenoid GPIO pin for output and set LOW to start
	pinSolenoid1.Output()
	pinSolenoid1.Low()
	pinSolenoid2.Output()
	pinSolenoid2.Low()

}

func handleFlowEdge(pin *gpio.Pin) {
	// handle falling edge flow sensor
	flow1Counter++
	flow2Counter++
	fmt.Println(flow1Counter)
}

func Pour(size int, tap int) {	
	//Reset flow counters
	flow1Counter = 0
	flow2Counter = 0
	
	//Enable watcher
	pinFlowSensor1.Watch(gpio.EdgeFalling, handleFlowEdge)
	pinFlowSensor2.Watch(gpio.EdgeFalling, handleFlowEdge)

	//Open selected tap and meter flow
	fmt.Println("Start Pourin")
	switch tap{
	case 1:
		pinSolenoid1.Toggle()
		for flow1Counter < size {
		time.Sleep(time.Millisecond*50)
		}
	case 2:
		pinSolenoid2.Toggle()
		for flow2Counter < size {
		time.Sleep(time.Millisecond*50)
		}
	default:
		fmt.Println("Invalid Tap to open and pour!!")
	}

	//Close tap once done pouring
	switch tap{
	case 1:
		pinSolenoid1.Toggle()
	case 2:
		pinSolenoid2.Toggle()
	default:
		fmt.Println("Invalid Tap to close!!")
	}
		
	//Disable watchers
	pinFlowSensor1.Unwatch()
	pinFlowSensor2.Unwatch()
}
