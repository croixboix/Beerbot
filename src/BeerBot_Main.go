package main
// NOTE: Make sure you have disabled I2C interface in sudo raspi-config - I
// think it might be enabled by default which might cause your pin 2 to be
// changed from input mode to I2C.
// SEE: https://github.com/stianeikeland/go-rpio/issues/35
// Add dtoverlay=gpio-no-irq to /boot/config.txt and restart your pi
//  This disables IRQ which may break some other GPIO libs/drivers
import (
	_"bufio"
	"encoding/json"
	"time"
	"fmt"
	"github.com/warthog618/gpio"
	gpio_rpi "gpio_rpi"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"log"
	"sync"
	"strconv"
	"math"
	"bytes"
)

const (
	// 10,200 Pulses per Gallon
	// 10,200 Pulses per 128 fluid ounce
	//  # pulses = (size in floz) / 0.0125490196078431372549019607843137254901960784313725490196078431



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

	//Testing variables below ONLY
	//testMessage string = "Tap ID and Order submitted!"
)

type Order struct {
	//Tap's UUID
	uuid string
	//Order's user/customer
	user int
	orderID int
	//Tap(s) to pour on with array value being drink size
	tap [numberOfTaps + 1]int
}

type OrderResponse struct {
	OrderID		int			`json:"id"`
	UserID    int    	`json:"user_id"`
	TapID 		int 		`json:"tap_id"`
	BeerID	 	int 		`json:"beer_id"`
	Price  		string  `json:"price"`
	WasPoured bool		`json:"was_poured"`
	Size      string	`json:"oz"`
}

type CheckResponse struct {
	OrderID		int			`json:"id"`
	WasPoured bool		`json:"was_poured"`
}

type ProcessResponse struct {
	Processed bool `json:"processed"`
}


// ######################## MAIN PROGRAM PROGRAM PROGRAM #######################
func main() {
	//Interrupt to handle command line crtl-c and exit cleanly
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
        <-interrupt
				fmt.Println("ctrl-c interrupt, exit cleanly!")
        endProgram()
    }()



	//Initialize GPIO interfaces
	gpio_rpi.GPIO_INIT()
	fmt.Println("GPIO Initialized!")

	//Main program loop
	for webConnectionAlive == true{

		time.Sleep(1*time.Second)

		//Check order queue for orders to pull
		orderIdToServe := checkOrders(tapUUID)
		fmt.Println("Order IDs to serve Array: ", orderIdToServe)

		//If there are orders to serve then let us fullfill them
		if len(orderIdToServe) >= 1 {

				for i := 0; i < len(orderIdToServe); i++ {
					//Get user orders
					userOrders := getOrders(tapUUID, orderIdToServe[i])

					//############ POUR/FULLFILL ORDER BLOCK #####################################
					//This is just a timeout function so that the program will timeout
					c1 := make(chan string, 1)
					// Run your long running function in it's own goroutine and pass back it's
					 // response into our channel.
					go func() {
						togglePour(*userOrders)
						text := "togglePour Finished!"
						c1 <- text
						}()
					// Listen on our channel AND a timeout channel - which ever happens first.
					select {
						case res := <-c1:
							fmt.Println(res)
						case <-time.After(120 * time.Second):
							fmt.Println("out of time :(")
							tempClearOrder := Order{uuid: tapUUID}
							//Clear pour process
							togglePour(*tempClearOrder)
						}
					//############ END POUR/FULLFILL ORDER BLOCK ######################################

					if processOrder(tapUUID, orderIdToServe[i]) == true{
							orderIdToServe = append(orderIdToServe[:i], orderIdToServe[i+1:]...)
							fmt.Println("Order IDs to server after processOrder update: ", orderIdToServe)
					}
				}
		}

	}

	//Run all the stuff needed to cleanly exit ( IMPORTANT THIS HAPPENS )
	endProgram()

}


//Get orders from the orderqueue
func getOrders(uuid string, orderID int) *Order {
	o := Order{uuid: tapUUID}

	url := "http://96.30.244.56:3000/api/v1/tap_orders/"+ strconv.Itoa(orderID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "96.30.244.56:3000")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	//fmt.Println(res)
	fmt.Println("getOrders body: ", string(body))

	var verifyResp []byte = body
	var verifyData OrderResponse

	err := json.Unmarshal(verifyResp, &verifyData)
	if err != nil {
		fmt.Println("unmarshal error:", err)
	}

	fmt.Println("Verify Order Response Dump:")
	fmt.Println("verifyData: ", verifyData)
	//fmt.Println("userID: ", verifyData.UserID)

	//If data isn't empty then import data into local order struct
	if verifyData.UserID != 0 && verifyData.WasPoured == false{
		o.user = verifyData.UserID
		o.orderID = verifyData.OrderID
		//fmt.Println("Order Username: ", o.user)

		//# pulses = (size in floz) / 0.012549
		//Calculate the number of pulses for the local order struct
		pulses, errPulseConversion := strconv.ParseFloat(verifyData.Size,32)
		if errPulseConversion != nil {
    	fmt.Println("size to pulse conversion error", errPulseConversion)
   	}
		pulses = pulses/0.012549
		//Round our float and store it away in local order struct
		o.tap[verifyData.TapID-1] = int(math.Round(pulses))

	}

	fmt.Println("getOrders o: ", o)

	return &o
}


//Check for orders to be served, returns array of ordersId to be served
func checkOrders(uuid string) []int{
	var orderIDs []int
	fmt.Println("Fetch order ids to fullfill")
	url := "http://96.30.244.56:3000/api/v1/tap_orders"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "96.30.244.56:3000")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	//fmt.Println(res)
	//fmt.Println("body: ", string(body))

	var verifyResp []byte = body
	var verifyData []CheckResponse

	err := json.Unmarshal([]byte(verifyResp), &verifyData)
	if err != nil {
		fmt.Println("unmarshal error:", err)
	}

	//fmt.Println("Verify Order Response Dump:")
	//fmt.Println("verifyData: ", verifyData)

	for i := 0; i < len(verifyData); i++ {
		//Check for orders
		if verifyData[i].OrderID != 0 && verifyData[i].WasPoured == false{
			//Tells main program there is an order to pour
			orderIDs = append(orderIDs, verifyData[i].OrderID)
		}
	}
	return orderIDs
}


// Tells API that order processed and deletes order from API order list
func processOrder(uuid string, orderID int) bool {
	url := "http://96.30.244.56:3000/api/v1/tap_orders/"+ strconv.Itoa(orderID)

	orderResp := CheckResponse{orderID,true}
	//var processData CheckResponse

	payload, err := json.Marshal(orderResp)
	if err != nil {
		fmt.Println("marshal error:", err)
	}

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
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

	fmt.Println("Process Order res: ", res)
	fmt.Println("Process Order body: ", string(body))

	return true
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
			fmt.Printf("(Go Routines)Begin measuring flow for user: %d on tap: %d of size: %d\n", customerOrder.user, i+1, customerOrder.tap[i])
			//fmt.Printf("Pour limit reached for user: %s on tap: %d of size: %d\n", customerOrder.user, i+1, customerOrder.tap[i])
		}
	}
	// Wait for all goroutines to be finished
	wg.Wait()
	fmt.Println("Finished all go routines!")
}


//Run all the stuff needed to cleanly exit ( IMPORTANT THIS HAPPENS )
func endProgram(){
	//Close GPIO/clear GPIO memory at end of program ( IMPORTANT THIS HAPPENS )
	gpio_rpi.CloseSolenoids()
	gpio.Close()
	log.Println("Program ended cleanly!")
	os.Exit(1)
}





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
