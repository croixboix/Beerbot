package main
// NOTE: Make sure you have disabled I2C interface in sudo raspi-config - I
// think it might be enabled by default which might cause your pin 2 to be
// changed from input mode to I2C.
// SEE: https://github.com/stianeikeland/go-rpio/issues/35
// Add dtoverlay=gpio-no-irq to /boot/config.txt and restart your pi
//  This disables IRQ which may break some other GPIO libs/drivers
import (
	_ "bufio"
	_ "encoding/json"
	"time"
	"fmt"
	"github.com/warthog618/gpio"
	gpio_rpi "gpio_rpi"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"log"
	"strings"
	"sync"
)

const (
	// 10,200 Pulses per Gallon
	// 10,200 Pulses per 128 fluid ounce
	//  # pulses = (size in floz) / 0.0125490196078431372549019607843137254901960784313725490196078431
	//Constants for drink size integers for flow sensor
	sizeFourOunce    int = 318
	sizeSixOunce     int = 478
	sizeTwelveOunce  int = 956
	sizeSixteenOunce int = 1275
	/* 	TODO:				Create helper function to calucate # pulses from size (maybe) */

	//Define number of taps on system (# of physical taps -1)
	//Ex: A 4 tap system would be = 3
	numberOfTaps int = 7
	/*	TODO:				Get this info from the API CALL!!! 						*/
)


var (
	//This is how we will set the tap system's ID
	tapUUID string = "TestTap"
	tapToken string = "a"

	//Size of order queue
	orderQueueSize int = 0

	//Keeps track of whether connection is alive
	webConnectionAlive bool = true
	failedPingCounter int = 0

	//Testing variables below ONLY
	//testMessage string = "Tap ID and Order submitted!"
)


type Order struct {
	//Tap's UUID
	uuid string
	//Order's user/customer
	user string
	//Tap(s) to pour on with array value being drink size
	tap [numberOfTaps + 1]int
}


//Initiates pour routine (this should be the last thing called, serves order)
func togglePour(customerOrder Order) {
	//Create a wait group for goroutines
	var wg sync.WaitGroup
	//wg.Add(numberOfTaps + 1)
	fmt.Println("Created goroutine wait groups!")
	//Solenoid normal state = closed
	for i := 0; i <= numberOfTaps; i++ {
		//fmt.Printf("(Go Routines)Begin measuring flow for user: %s on tap: %d of size: %d\n", customerOrder.user, i+1, customerOrder.tap[i])
		if customerOrder.tap[i] != 0 {
			wg.Add(1)
			go gpio_rpi.Pour(customerOrder.tap[i], i+1, &wg)
			fmt.Printf("(Go Routines)Begin measuring flow for user: %s on tap: %d of size: %d\n", customerOrder.user, i+1, customerOrder.tap[i])
			//fmt.Printf("Pour limit reached for user: %s on tap: %d of size: %d\n", customerOrder.user, i+1, customerOrder.tap[i])
		}
	}
	// Wait for all goroutines to be finished
	wg.Wait()
	fmt.Println("Finished all go routines!")
}


// ######################## MAIN PROGRAM PROGRAM PROGRAM #######################
func main() {
	//Interrupt to handle command line crtl-c and exit cleanly
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	//Initialize GPIO interfaces
	gpio_rpi.GPIO_INIT()
	fmt.Println("GPIO Initialized!")

	type verifyResponse struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		URL       string `json:"url"`
	}

	type processResponse struct {
		Processed bool `json:"processed"`
	}

	//Main program loop
	for webConnectionAlive == true{
		//Check for ctrl-c CLI input to end program cleanly
		select {
			case <-interrupt:
				log.Println("Ctrl-c input detected, exiting cleanly")
				endProgram()
				return
			}

		time.Sleep(1*time.Second)
		order := getOrder(tapUUID)

		if orderQueueSize > 1 {
			//############ TEST/DEMO CODE BLOCK ######################################
				//Get test order from API
				testOrder := getOrder(tapUUID)

				//This is just a timeout function so that the program will timeout
				c1 := make(chan string, 1)
				// Run your long running function in it's own goroutine and pass back it's
				 // response into our channel.
				go func() {
					togglePour(*testOrder)
					text := "togglePour Finished!"
					c1 <- text
					}()
				// Listen on our channel AND a timeout channel - which ever happens first.
				select {
					case res := <-c1:
						fmt.Println(res)
					case <-time.After(120 * time.Second):
						fmt.Println("out of time :(")
					}
			//############ END TEST/DEMO CODE BLOCK ######################################
		}

	}


	//Run all the stuff needed to cleanly exit ( IMPORTANT THIS HAPPENS )
	endProgram()

}

//Run all the stuff needed to cleanly exit ( IMPORTANT THIS HAPPENS )
func endProgram(){
	//Close GPIO/clear GPIO memory at end of program ( IMPORTANT THIS HAPPENS )
	gpio_rpi.CloseSolenoids()
	gpio.Close()
	log.Println("Program ended cleanly!")
	os.Exit(1)
}


func connectionAliveTest(failedPingCounter int){
	if failedPingCounter >= 100{
		webConnectionAlive = false
	} else {
		webConnectionAlive = true
	}
}

// Tells API that order processed and deletes order from API order list
func processOrder(uname string) []byte {
	url := "http://96.30.244.56:3000/orders/processed"
	payload := strings.NewReader("{\n\t\"order\": {\n\t\t\"username\": \"" + uname + "\"\n\t}\n}")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "96.30.244.56:3000")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Length", "39")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	//fmt.Println(res)
	//fmt.Println(string(body))

	return body
}


//Get orders from the orderqueue
func getOrder(uuid string) *Order {
	o := Order{uuid: tapUUID}
	fmt.Println("Fetch orders")
	url := "http://96.30.244.56:3000/orders"
	//payload := strings.NewReader("{\n\t\"order\": {\n\t\t\"username\": \"" + uname + "\"\n\t}\n}")
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "96.30.244.56:3000")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Length", "39")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var verifyResp []byte = body
	var verifyData verifyResponse

	err := json.Unmarshal(verifyResp, &verifyData)
	if err != nil {
		fmt.Println("error:", err)
	}

	if body.user != "null"{
		orderQueueSize++
		o.user = body.user
		o.tap = body.tap
		fmt.Printf("Username: %s\n", o.user)
		for i := 0; i <= numberOfTaps; i++ {
			//o.tap[i] = tap[i]
			//fmt.Printf("numberOfTaps = %d | i = %d | tap[i] = %d | o.tap[i] = %d\n", numberOfTaps, i, tap[i], o.tap[i])
			fmt.Printf("Tap # %d value(drink size) is %d\n", i+1, o.tap[i])
		}
	}

	return &o
}

/*
		Order does exist response:
		{
			"id": 591,
			"username": "test",
			"created_at": "2019-11-15T16:31:21.321Z",
			"updated_at": "2019-11-15T16:31:21.321Z",
			"url": "http://96.30.244.56:3000/orders/591.json"
		}

		Order does not exist response:
		{
		"id": null,
		"username": null,
		"created_at": null,
		"updated_at": null,
		"url": null
		}
*/






/*#############################DEPRECATED/FOR REFERENCE ONLY##############################################################*/



/*
//Create struct for verify response
type verifyResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	URL       string `json:"url"`
}

type processResponse struct {
	Processed bool `json:"processed"`
}
*/

/*
//Test code for reading from USB (STD-IN) QR scanner
func scanCode() string {
	var userCode string

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		fmt.Println("Scanned barcode: ", scanner.Text())
		userCode = scanner.Text()
		if scanner.Text() != "" {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return userCode
}
*/


/*
func main() {

//API test code below

//Verify and process the order!
fmt.Println("Verify Order")

var verifyResp []byte = verifyOrder(testOrder.user)
var verifyData verifyResponse

err := json.Unmarshal(verifyResp, &verifyData)
if err != nil {
	fmt.Println("error:", err)
}

fmt.Println("Verify Order Response Dump:")
fmt.Println("id: ", verifyData.ID)
fmt.Println("username: ", verifyData.Username)
fmt.Println("created_at: ", verifyData.CreatedAt)
fmt.Println("updated_at: ", verifyData.UpdatedAt)
fmt.Println("url: ", verifyData.URL)

if verifyData.Username != "null" {

	fmt.Println("Process/Delete Order")

	var processData processResponse

	var processResponse []byte = processOrder(verifyData.Username)

	err := json.Unmarshal(processResponse, &processData)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println("Process Order Response Dump:")
	fmt.Println("processed: ", processData.Processed)

	if processData.Processed == true {
		//Let user pour the drink!
		//Call pour!
		//togglePour(drinkSize[tap-1], tap)
		togglePour(*testOrder)
		fmt.Println("ORDER PROCESSED, LET USER POUR")
	} else {
		fmt.Println("ORDER DOES NOT EXIST, DO NOT LET USER POUR")
	}
}
}

*/
