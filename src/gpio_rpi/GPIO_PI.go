package gpio_rpi

import (
	"fmt"
	"github.com/warthog618/gpio"
	"os"
	"time"
	"sync"
	//For debugging only

)

var (
	//Should eventually create a struct for all these
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


	flowCounter1 int = 0
	flowCounter2 int = 0
	flowCounter3 int = 0
	flowCounter4 int = 0
	flowCounter5 int = 0
	flowCounter6 int = 0
	flowCounter7 int = 0
	flowCounter8 int = 0

)

func GPIO_INIT() {
	//Enable/Open GPIO
	if err := gpio.Open(); err != nil {

		fmt.Println(err)

		os.Exit(1)
	}

	// Can't use pins 2 and 3 because they have built in pull-up resistors
	// Also can't use pins 0 and 1
	//Assign pins
	pinFlowSensor1 = gpio.NewPin(14)
	pinSolenoid1 = gpio.NewPin(15)
	pinFlowSensor2 = gpio.NewPin(4)
	pinSolenoid2 = gpio.NewPin(16)
	pinFlowSensor3 = gpio.NewPin(5)
	pinSolenoid3 = gpio.NewPin(17)
	pinFlowSensor4 = gpio.NewPin(6)
	pinSolenoid4 = gpio.NewPin(18)
	pinFlowSensor5 = gpio.NewPin(7)
	pinSolenoid5 = gpio.NewPin(19)
	pinFlowSensor6 = gpio.NewPin(8)
	pinSolenoid6 = gpio.NewPin(20)
	pinFlowSensor7 = gpio.NewPin(9)
	pinSolenoid7 = gpio.NewPin(21)
	pinFlowSensor8 = gpio.NewPin(10)
	pinSolenoid8 = gpio.NewPin(22)


	//Configure Flow Sensor GPIO pin for input and PullUp
	//ALSO TRY pinFlowSensor.PullDown() if having issues!
	//Should rewrite this to use structs and a loop to set these values
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
	//Should rewrite this to use structs and a loop to set these values
	pinSolenoid1.Output()
	pinSolenoid1.High()
	pinSolenoid2.Output()
	pinSolenoid2.High()
	pinSolenoid3.Output()
	pinSolenoid3.High()
	pinSolenoid4.Output()
	pinSolenoid4.High()
	pinSolenoid5.Output()
	pinSolenoid5.High()
	pinSolenoid6.Output()
	pinSolenoid6.High()
	pinSolenoid7.Output()
	pinSolenoid7.High()
	pinSolenoid8.Output()
	pinSolenoid8.High()

}

func CloseSolenoids(tap int) {
	switch tap {
	case 0:
		pinSolenoid1.High()
		pinSolenoid2.High()
		pinSolenoid3.High()
		pinSolenoid4.High()
		pinSolenoid5.High()
		pinSolenoid6.High()
		pinSolenoid7.High()
		pinSolenoid8.High()
		fmt.Println("CLOSE ALL SOLENOIDS")
	case 1:
		pinSolenoid1.High()
	case 2:
		pinSolenoid2.High()
	case 3:
		pinSolenoid3.High()
	case 4:
		pinSolenoid4.High()
	case 5:
		pinSolenoid5.High()
	case 6:
		pinSolenoid6.High()
	case 7:
		pinSolenoid7.High()
	case 8:
		pinSolenoid8.High()
	default:
		fmt.Println("Invalid Tap # To Close!!")
	}
}

// Interrupt handler for flow sensor edge pulse
func handleFlowEdge(pin *gpio.Pin) {

	//Figure out which flow sensor made the interrupt call
	switch pin.Pin() {
	case pinFlowSensor1.Pin():
		flowCounter1++
		fmt.Printf("Flow Counter is %d on pinFlowSensor %d\n", flowCounter1, pinFlowSensor1.Pin())
	case pinFlowSensor2.Pin():
		flowCounter2++
		fmt.Printf("Flow Counter is %d on pinFlowSensor %d\n", flowCounter2, pinFlowSensor2.Pin())
	case pinFlowSensor3.Pin():
		flowCounter3++
		fmt.Printf("Flow Counter is %d on pinFlowSensor %d\n", flowCounter3, pinFlowSensor3.Pin())
	case pinFlowSensor4.Pin():
		flowCounter4++
		fmt.Printf("Flow Counter is %d on pinFlowSensor %d\n", flowCounter4, pinFlowSensor4.Pin())
	case pinFlowSensor5.Pin():
		flowCounter5++
		fmt.Printf("Flow Counter is %d on pinFlowSensor %d\n", flowCounter5, pinFlowSensor5.Pin())
	case pinFlowSensor6.Pin():
		flowCounter6++
		fmt.Printf("Flow Counter is %d on pinFlowSensor %d\n", flowCounter6, pinFlowSensor6.Pin())
	case pinFlowSensor7.Pin():
		flowCounter7++
		fmt.Printf("Flow Counter is %d on pinFlowSensor %d\n", flowCounter7, pinFlowSensor7.Pin())
	case pinFlowSensor8.Pin():
		flowCounter8++
		fmt.Printf("Flow Counter is %d on pinFlowSensor %d\n", flowCounter8, pinFlowSensor8.Pin())
	default:
		fmt.Println("handleFlowEdge Invalid Tap #!!")
	}

	//fmt.Printf("Go Routine %d is referencing flowCounter address %X \n", goid(), &flowCounter)
	//flowCounter++
	//fmt.Printf("Flow Counter is %d for Go Routine %d \n", flowCounter, goid())
}


//Returns true when pour is done
func Pour(size int, tap int, wg1 *sync.WaitGroup) {
	//Reset flow counter for this tap
	//var flowCounter int = 0

	defer wg1.Done()

	//Open selected tap and meter flow
	fmt.Printf("Start func Pour on tap %d of size: %d\n", tap, size)

	switch tap {
	case 1:
		//Reset this flow counter
		flowCounter1 = 0
		//Enable watcher
		pinFlowSensor1.Watch(gpio.EdgeFalling, handleFlowEdge)
		//Open solenoid
		pinSolenoid1.Low()
		//Count number of ticks from flow sensor
		for flowCounter1 < size {
			time.Sleep(time.Millisecond * 50)
		}
		//Disable flow sensor watcher
		pinFlowSensor1.Unwatch()
		//Close solenoid/tap
		pinSolenoid1.High()
	case 2:
		flowCounter2 = 0
		pinFlowSensor2.Watch(gpio.EdgeFalling, handleFlowEdge)
		pinSolenoid2.Low()
		for flowCounter2 < size {
			time.Sleep(time.Millisecond * 50)
		}
		pinFlowSensor2.Unwatch()
		pinSolenoid2.High()
	case 3:
		flowCounter3 = 0
		pinFlowSensor3.Watch(gpio.EdgeFalling, handleFlowEdge)
		pinSolenoid3.Low()
		for flowCounter3 < size {
			time.Sleep(time.Millisecond * 50)
		}
		pinFlowSensor3.Unwatch()
		pinSolenoid3.High()
	case 4:
		flowCounter4 = 0
		pinFlowSensor4.Watch(gpio.EdgeFalling, handleFlowEdge)
		pinSolenoid4.Low()
		for flowCounter4 < size {
			time.Sleep(time.Millisecond * 50)
		}
		pinFlowSensor4.Unwatch()
		pinSolenoid4.High()
	case 5:
		flowCounter5 = 0
		pinFlowSensor5.Watch(gpio.EdgeFalling, handleFlowEdge)
		pinSolenoid5.Low()
		for flowCounter5 < size {
			time.Sleep(time.Millisecond * 50)
		}
		pinFlowSensor5.Unwatch()
		pinSolenoid5.High()
	case 6:
		flowCounter6 = 0
		pinFlowSensor6.Watch(gpio.EdgeFalling, handleFlowEdge)
		pinSolenoid6.Low()
		for flowCounter6 < size {
			time.Sleep(time.Millisecond * 50)
		}
		pinFlowSensor6.Unwatch()
		pinSolenoid6.High()
	case 7:
		flowCounter7 = 0
		pinFlowSensor7.Watch(gpio.EdgeFalling, handleFlowEdge)
		pinSolenoid7.Low()
		for flowCounter7 < size {
			time.Sleep(time.Millisecond * 50)
		}
		pinFlowSensor7.Unwatch()
		pinSolenoid7.High()
	case 8:
		flowCounter8 = 0
		pinFlowSensor8.Watch(gpio.EdgeFalling, handleFlowEdge)
		pinSolenoid8.Low()
		for flowCounter8 < size {
			time.Sleep(time.Millisecond * 50)
		}
		pinFlowSensor8.Unwatch()
		pinSolenoid8.High()
	default:
		fmt.Println("Invalid Tap #!!")
	}
}
