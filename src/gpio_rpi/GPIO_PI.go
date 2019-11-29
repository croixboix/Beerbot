package gpio_rpi

import (
	"fmt"
	"github.com/warthog618/gpio"
	"os"
	"time"
)

var (
	flow1Counter int = 0
	flow2Counter int = 0
	flow3Counter int = 0
	flow4Counter int = 0
	flow5Counter int = 0
	flow6Counter int = 0
	flow7Counter int = 0
	flow8Counter int = 0

	pinFlowSensor1 *gpio.Pin
	pinSolenoid1   *gpio.Pin
	pinFlowSensor2 *gpio.Pin
	pinSolenoid2   *gpio.Pin
	pinFlowSensor3 *gpio.Pin
	pinSolenoid3   *gpio.Pin
	pinFlowSensor4 *gpio.Pin
	pinSolenoid4   *gpio.Pin
	pinFlowSensor5 *gpio.Pin
	pinSolenoid5   *gpio.Pin
	pinFlowSensor6 *gpio.Pin
	pinSolenoid6   *gpio.Pin
	pinFlowSensor7 *gpio.Pin
	pinSolenoid7   *gpio.Pin
	pinFlowSensor8 *gpio.Pin
	pinSolenoid8   *gpio.Pin
)

func GPIO_INIT() {
	//Enable/Open GPIO
	if err := gpio.Open(); err != nil {

		fmt.Println(err)

		os.Exit(1)
	}

	pinFlowSensor1 = gpio.NewPin(23)
	pinSolenoid1 = gpio.NewPin(24)
	pinFlowSensor2 = gpio.NewPin(25)
	pinSolenoid2 = gpio.NewPin(8)
	pinFlowSensor3 = gpio.NewPin(7)
	pinSolenoid3 = gpio.NewPin(12)
	pinFlowSensor4 = gpio.NewPin(16)
	pinSolenoid4 = gpio.NewPin(20)
	pinFlowSensor5 = gpio.NewPin(21)
	pinSolenoid5 = gpio.NewPin(9)
	pinFlowSensor6 = gpio.NewPin(11)
	pinSolenoid6 = gpio.NewPin(5)
	pinFlowSensor7 = gpio.NewPin(6)
	pinSolenoid7 = gpio.NewPin(13)
	pinFlowSensor8 = gpio.NewPin(19)
	pinSolenoid8 = gpio.NewPin(26)

	//Configure Flow Sensor GPIO pin for input and PullUp
	//ALSO TRY pinFlowSensor.PullDown() if having issues!
	pinFlowSensor1.Input()
	pinFlowSensor1.PullUp()
	pinFlowSensor2.Input()
	pinFlowSensor2.PullUp()
	pinFlowSensor3.Input()
	pinFlowSensor3.PullUp()
	pinFlowSensor4.Input()
	pinFlowSensor4.PullUp()
	pinFlowSensor5.Input()
	pinFlowSensor5.PullUp()
	pinFlowSensor6.Input()
	pinFlowSensor6.PullUp()
	pinFlowSensor7.Input()
	pinFlowSensor7.PullUp()
	pinFlowSensor8.Input()
	pinFlowSensor8.PullUp()

	//Configure solenoid GPIO pin for output and set LOW to start
	pinSolenoid1.Output()
	pinSolenoid1.Low()
	pinSolenoid2.Output()
	pinSolenoid2.Low()
	pinSolenoid3.Output()
	pinSolenoid3.Low()
	pinSolenoid4.Output()
	pinSolenoid4.Low()
	pinSolenoid5.Output()
	pinSolenoid5.Low()
	pinSolenoid6.Output()
	pinSolenoid6.Low()
	pinSolenoid7.Output()
	pinSolenoid7.Low()
	pinSolenoid8.Output()
	pinSolenoid8.Low()
}

func handleFlowEdge(pin *gpio.Pin) {
	// handle falling edge flow sensor
	flow1Counter++
	flow2Counter++
	flow3Counter++
	flow4Counter++
	flow5Counter++
	flow6Counter++
	flow7Counter++
	flow8Counter++
	fmt.Println(flow1Counter)
	fmt.Println(flow2Counter)
	fmt.Println(flow3Counter)
	fmt.Println(flow4Counter)
	fmt.Println(flow5Counter)
	fmt.Println(flow6Counter)
	fmt.Println(flow7Counter)
	fmt.Println(flow8Counter)
}

func Pour(size int, tap int) {
	//Reset flow counters
	flow1Counter = 0
	flow2Counter = 0
	flow3Counter = 0
	flow4Counter = 0
	flow5Counter = 0
	flow6Counter = 0
	flow7Counter = 0
	flow8Counter = 0

	//Open selected tap and meter flow
	fmt.Println("Start Pourin")
	switch tap {
	case 1:
		//Enable watcher
		pinFlowSensor1.Watch(gpio.EdgeFalling, handleFlowEdge)
		//Open solenoid
		pinSolenoid1.Toggle()
		//Count number of ticks from flow sensor
		for flow1Counter < size {
			time.Sleep(time.Millisecond * 50)
		}
		//Disable flow sensor watcher
		pinFlowSensor1.Unwatch()
	case 2:
		pinFlowSensor2.Watch(gpio.EdgeFalling, handleFlowEdge)
		pinSolenoid2.Toggle()
		for flow2Counter < size {
			time.Sleep(time.Millisecond * 50)
		}
		pinFlowSensor2.Unwatch()
	case 3:
		pinFlowSensor3.Watch(gpio.EdgeFalling, handleFlowEdge)
		pinSolenoid3.Toggle()
		for flow2Counter < size {
			time.Sleep(time.Millisecond * 50)
		}
		pinFlowSensor3.Unwatch()
	case 4:
		pinFlowSensor4.Watch(gpio.EdgeFalling, handleFlowEdge)
		pinSolenoid4.Toggle()
		for flow2Counter < size {
			time.Sleep(time.Millisecond * 50)
		}
		pinFlowSensor4.Unwatch()
	case 5:
		pinFlowSensor5.Watch(gpio.EdgeFalling, handleFlowEdge)
		pinSolenoid5.Toggle()
		for flow2Counter < size {
			time.Sleep(time.Millisecond * 50)
		}
		pinFlowSensor5.Unwatch()
	case 6:
		pinFlowSensor6.Watch(gpio.EdgeFalling, handleFlowEdge)
		pinSolenoid6.Toggle()
		for flow2Counter < size {
			time.Sleep(time.Millisecond * 50)
		}
		pinFlowSensor6.Unwatch()
	case 7:
		pinFlowSensor7.Watch(gpio.EdgeFalling, handleFlowEdge)
		pinSolenoid7.Toggle()
		for flow2Counter < size {
			time.Sleep(time.Millisecond * 50)
		}
		pinFlowSensor7.Unwatch()
	case 8:
		pinFlowSensor8.Watch(gpio.EdgeFalling, handleFlowEdge)
		pinSolenoid8.Toggle()
		for flow2Counter < size {
			time.Sleep(time.Millisecond * 50)
		}
		pinFlowSensor8.Unwatch()
	default:
		fmt.Println("Invalid Tap # attempted to open!!")
	}

	//Close tap once done pouring
	switch tap {
	case 1:
		pinSolenoid1.Toggle()
	case 2:
		pinSolenoid2.Toggle()
	case 3:
		pinSolenoid3.Toggle()
	case 4:
		pinSolenoid4.Toggle()
	case 5:
		pinSolenoid5.Toggle()
	case 6:
		pinSolenoid6.Toggle()
	case 7:
		pinSolenoid7.Toggle()
	case 8:
		pinSolenoid8.Toggle()
	default:
		fmt.Println("Invalid Tap # attempted to close!!")
	}
}
