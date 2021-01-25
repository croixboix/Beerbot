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
	"runtime"
	"strings"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

const (
	// 10,200 Pulses per Gallon
	// 10,200 Pulses per 128 fluid ounce
	//  # pulses = (size in floz) / 0.0125490196078431372549019607843137254901960784313725490196078431



	//Define number of taps on system (# of physical taps -1)
	//Ex: A 4 tap system would be = 3
	numberOfTaps int = 7

)

var (
	//This is how we will set the tap system's ID
	tapUUID string = "d95c9dcc-3b99-4ad6-9a4e-e60418fcd359"
	tapControlID = 1
	authToken string


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
	orderID int
	user int
	tapID int
	beerID int
	price string
	size string

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

type AuthPOST struct {
	TapControl struct {
		TapControlID int 	`json:"id"`
		TapUUID string		`json:"uuid"`
	} `json:"tap_control"`
}

type AuthResponse struct {
	TapControlID int 		`json:"id"`
	TapUUID string			`json:"uuid"`
	AuthenToken string 	`json:"authentication_token"`
}

type orderLabels struct {
	orderIDL *widget.Label
	userIDL *widget.Label
	tapIDL *widget.Label
	beerIDL *widget.Label
	priceL *widget.Label
	sizeL *widget.Label
}


// ######################## MAIN PROGRAM PROGRAM PROGRAM #######################
func main() {


	a := app.New()
	w := a.NewWindow("Canvas")
	myCanvas := w.Canvas()


	oL1 := orderLabels{orderIDL: widget.NewLabel("-"),
										 userIDL:widget.NewLabel("-"),
										 tapIDL:widget.NewLabel("-"),
										 beerIDL:widget.NewLabel("-"),
										 priceL:widget.NewLabel("-"),
										 sizeL:widget.NewLabel("-")}
  oL2 := orderLabels{orderIDL: widget.NewLabel("-"),
 										 userIDL:widget.NewLabel("-"),
 										 tapIDL:widget.NewLabel("-"),
 										 beerIDL:widget.NewLabel("-"),
 										 priceL:widget.NewLabel("-"),
 										 sizeL:widget.NewLabel("-")}
  oL3 := orderLabels{orderIDL: widget.NewLabel("-"),
										 userIDL:widget.NewLabel("-"),
									   tapIDL:widget.NewLabel("-"),
										 beerIDL:widget.NewLabel("-"),
										 priceL:widget.NewLabel("-"),
										 sizeL:widget.NewLabel("-")}
  oL4 := orderLabels{orderIDL: widget.NewLabel("-"),
										 userIDL:widget.NewLabel("-"),
										 tapIDL:widget.NewLabel("-"),
										 beerIDL:widget.NewLabel("-"),
										 priceL:widget.NewLabel("-"),
										 sizeL:widget.NewLabel("-")}
  oL5 := orderLabels{orderIDL: widget.NewLabel("-"),
										 userIDL:widget.NewLabel("-"),
										 tapIDL:widget.NewLabel("-"),
										 beerIDL:widget.NewLabel("-"),
										 priceL:widget.NewLabel("-"),
										 sizeL:widget.NewLabel("-")}
  oL6 := orderLabels{orderIDL: widget.NewLabel("-"),
 										userIDL:widget.NewLabel("-"),
 										tapIDL:widget.NewLabel("-"),
 										beerIDL:widget.NewLabel("-"),
 										priceL:widget.NewLabel("-"),
 										sizeL:widget.NewLabel("-")}
	oL7 := orderLabels{orderIDL: widget.NewLabel("-"),
										 userIDL:widget.NewLabel("-"),
										 tapIDL:widget.NewLabel("-"),
										 beerIDL:widget.NewLabel("-"),
										 priceL:widget.NewLabel("-"),
										 sizeL:widget.NewLabel("-")}
	oL8 := orderLabels{orderIDL: widget.NewLabel("-"),
										 userIDL:widget.NewLabel("-"),
										 tapIDL:widget.NewLabel("-"),
										 beerIDL:widget.NewLabel("-"),
										 priceL:widget.NewLabel("-"),
										 sizeL:widget.NewLabel("-")}

	orderL  := widget.NewLabel("Order ID: ")
  userL   := widget.NewLabel("User ID: ")
  tapL    := widget.NewLabel("Tap ID: ")
  beerL   := widget.NewLabel("Beer ID: ")
  priceL  := widget.NewLabel("Price: ")
  sizeL   := widget.NewLabel("Size: ")


//Box frenzy
	myCanvas.SetContent(
    widget.NewHBox(
      widget.NewVBox(
        orderL, userL, tapL, beerL, priceL, sizeL,),//end Heading VBox
      widget.NewVBox(
        oL1.orderIDL, oL1.userIDL, oL1.tapIDL, oL1.beerIDL,
				oL1.priceL, oL1.sizeL,),//end user1 Vbox
			widget.NewVBox(
        oL2.orderIDL, oL2.userIDL, oL2.tapIDL, oL2.beerIDL,
				oL2.priceL, oL2.sizeL,),//end user2 Vbox
			widget.NewVBox(
        oL3.orderIDL, oL3.userIDL, oL3.tapIDL, oL3.beerIDL,
				oL3.priceL, oL3.sizeL,),//end user3 Vbox
			widget.NewVBox(
        oL4.orderIDL, oL4.userIDL, oL4.tapIDL, oL4.beerIDL,
				oL4.priceL, oL4.sizeL,),//end user4 Vbox
			widget.NewVBox(
				oL5.orderIDL, oL5.userIDL, oL5.tapIDL, oL5.beerIDL,
				oL5.priceL, oL5.sizeL,),//end user5 Vbox
			widget.NewVBox(
				oL6.orderIDL, oL6.userIDL, oL6.tapIDL, oL6.beerIDL,
				oL6.priceL, oL6.sizeL,),//end user6 Vbox
			widget.NewVBox(
				oL7.orderIDL, oL7.userIDL, oL7.tapIDL, oL7.beerIDL,
				oL7.priceL, oL7.sizeL,),//end user7 Vbox
			widget.NewVBox(
				oL8.orderIDL, oL8.userIDL, oL8.tapIDL, oL8.beerIDL,
				oL8.priceL, oL8.sizeL,),//end user8 Vbox
    ),//end Hbox
  ) //adding label widget to window


	go runProgram(myCanvas, oL1, oL2, oL3, oL4, oL5, oL6, oL7, oL8)

	w.Resize(fyne.NewSize(500, 230))
	w.ShowAndRun()



	endProgram()
}

func runProgram(c fyne.Canvas, oL1 orderLabels, oL2 orderLabels, oL3 orderLabels, oL4 orderLabels, oL5 orderLabels, oL6 orderLabels, oL7 orderLabels, oL8 orderLabels) {
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

	//Authenticate with API
	authToken = authTapController(tapUUID, tapControlID)


	//Main program loop
	for webConnectionAlive == true{

		//fmt.Println("Created togglePour goroutine wait groups!")

		time.Sleep(1*time.Second)

		//Check order queue for orders to pull
		orderIdToServe := checkOrders(tapUUID)
		fmt.Println("Order IDs to serve: ", orderIdToServe)

		//If there are orders to serve then let us fullfill them
		if len(orderIdToServe) >= 1 {

				for i := 0; i < len(orderIdToServe); i++ {
					//Get user orders
					userOrders := getOrders(tapUUID, orderIdToServe[i])

					fmt.Println("In change Content")
				  oL1.orderIDL.SetText(strconv.Itoa(userOrders.orderID))
					oL1.userIDL.SetText(strconv.Itoa(userOrders.user))
					oL1.tapIDL.SetText(strconv.Itoa(userOrders.tapID))
					oL1.beerIDL.SetText(strconv.Itoa(userOrders.beerID))
					oL1.priceL.SetText(userOrders.price)
					oL1.sizeL.SetText(userOrders.size)


					go togglePour(*userOrders)

				}

				//fmt.Println("Order ID Array before processOrder: ", orderIdToServe)
				//fmt.Println("len(orderIdToServe): ", len(orderIdToServe))
				// Mark the orders we just fullfilled/poured as poured on the orders API
				for i := len(orderIdToServe) - 1; i >= 0; i-- {
					// Call to process order
					if processOrder(tapUUID, orderIdToServe[i]) == true{
							orderIdToServe = append(orderIdToServe[:i], orderIdToServe[i+1:]...)
						}
				}
		}
	}
	//Run all the stuff needed to cleanly exit ( IMPORTANT THIS HAPPENS )
	endProgram()
}


//Tap Controller Authentication with API
func authTapController(uuid string, tapControlID int) string{
	url := "http://96.30.244.56:3000/api/v1/tap_sessions"

	authPost := AuthPOST{tapControlID,uuid}
	//var processData CheckResponse

	payload, err := json.Marshal(authPost)
	if err != nil {
		fmt.Println("marshal error:", err)
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))

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
	fmt.Println("body: ", string(body))

	var authResp []byte = body
	var verifyAuth AuthResponse

	errR := json.Unmarshal(authResp, &verifyAuth)
	if errR != nil {
		fmt.Println("unmarshal error:", errR)
	}

	fmt.Println(verifyAuth.AuthenToken)
	return verifyAuth.AuthenToken
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
		o.tapID = verifyData.TapID
		o.beerID = verifyData.BeerID
		o.price = verifyData.Price
		o.size = verifyData.Size
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
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	//body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println("Process Order body: ", string(body))

	//fmt.Println("Process Order res: ", res.Status)
	if res.Status == "204 No Content" {
		fmt.Println("Processed orderID: ", orderID)
	 	return true
 	} else {
		fmt.Println("FAILED processing orderID: ", orderID)
		return false
	}
}


//Initiates pour routine (this should be the last thing called, serves order)
func togglePour(customerOrder Order) {
	//This is just a timeout function so that the program will timeout
	c1 := make(chan string, 1)
	// Run your long running function in it's own goroutine and pass back it's
	// response into our channel.

	solenoidToClose := 9

	go func() {
		var wg1 sync.WaitGroup

		//Solenoid normal state = closed
		for i := 0; i <= numberOfTaps; i++ {
			if customerOrder.tap[i] != 0 {
				wg1.Add(1)
				go gpio_rpi.Pour(customerOrder.tap[i], i+1, &wg1)
				solenoidToClose = i+1
			}
		}
		wg1.Wait()
		text := "Pour Finished!"
		c1 <- text
	}()
	// Listen on our channel AND a timeout channel - which ever happens first.
	select {
		case res := <-c1:
			fmt.Println(res)
		case <-time.After(60 * time.Second):
			fmt.Println("out of time :(")
			gpio_rpi.CloseSolenoids(solenoidToClose)
	}
}




//DEBUGGING PURPOSES ONLY! // USE goid() to return thread id
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


//Run all the stuff needed to cleanly exit ( IMPORTANT THIS HAPPENS )
func endProgram(){
	//Close GPIO/clear GPIO memory at end of program ( IMPORTANT THIS HAPPENS )
	gpio_rpi.CloseSolenoids(0)
	gpio.Close()
	log.Println("Program ended cleanly!")
	os.Exit(1)
}





/*#############################DEPRECATED/FOR REFERENCE ONLY##############################################################*/

/*
 OLD MULTI-THREADED POUR FUNCTION
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
*/


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
