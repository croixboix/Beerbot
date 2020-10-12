package gpio_rpi

import (
	"fmt"
	"github.com/warthog618/gpio"
	"os"
	"sync"
	"time"

	//For debugging only
	"runtime"
	"strconv"
	"strings"
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

	//Assign pins
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

// Interrupt handler for flow sensor edge pulse
func handleFlowEdge(pin *gpio.Pin) {

	//Figure out which flow sensor made the interrupt call
	switch pin.Pin() {
	case pinFlowSensor1.Pin():
		flowCounter1++
		fmt.Printf("Flow Counter is %d for Go Routine %d on pinFlowSensor %d\n", flowCounter1, goid(), pinFlowSensor1.Pin())
	case pinFlowSensor2.Pin():
		flowCounter2++
		fmt.Printf("Flow Counter is %d for Go Routine %d on pinFlowSensor %d\n", flowCounter2, goid(), pinFlowSensor2.Pin())
	case pinFlowSensor3.Pin():
		flowCounter3++
		fmt.Printf("Flow Counter is %d for Go Routine %d on pinFlowSensor %d\n", flowCounter3, goid(), pinFlowSensor3.Pin())
	case pinFlowSensor4.Pin():
		flowCounter4++
		fmt.Printf("Flow Counter is %d for Go Routine %d on pinFlowSensor %d\n", flowCounter4, goid(), pinFlowSensor4.Pin())
	case pinFlowSensor5.Pin():
		flowCounter5++
		fmt.Printf("Flow Counter is %d for Go Routine %d on pinFlowSensor %d\n", flowCounter5, goid(), pinFlowSensor5.Pin())
	case pinFlowSensor6.Pin():
		flowCounter6++
		fmt.Printf("Flow Counter is %d for Go Routine %d on pinFlowSensor %d\n", flowCounter6, goid(), pinFlowSensor6.Pin())
	case pinFlowSensor7.Pin():
		flowCounter7++
		fmt.Printf("Flow Counter is %d for Go Routine %d on pinFlowSensor %d\n", flowCounter7, goid(), pinFlowSensor7.Pin())
	case pinFlowSensor8.Pin():
		flowCounter8++
		fmt.Printf("Flow Counter is %d for Go Routine %d on pinFlowSensor %d\n", flowCounter8, goid(), pinFlowSensor8.Pin())
	default:
		fmt.Println("handleFlowEdge Invalid Tap #!!")
	}

	//fmt.Printf("Go Routine %d is referencing flowCounter address %X \n", goid(), &flowCounter)
	//flowCounter++
	//fmt.Printf("Flow Counter is %d for Go Routine %d \n", flowCounter, goid())
}

func Pour(size int, tap int, wg *sync.WaitGroup) {
	// Call Done() using defer as it's be easiest way to guarantee it's called at every exit
	defer wg.Done()

	//Reset flow counter for this tap
	//var flowCounter int = 0

	//Open selected tap and meter flow
	fmt.Printf("Start func Pour on tap %d of GOID: %d\n", tap, goid())

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

//DEBUGGING PURPOSES ONLY!
func goid() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
