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
	"gpio_rpi"
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
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
)

const (
	// 10,200 Pulses per Gallon
	// 10,200 Pulses per 128 fluid ounce
	//  # pulses = (size in floz) / 0.01254901960784313725490196078431372549019607
	pulseDivConst float64 = 0.01254902


	//Define number of taps on system (# of physical taps -1)
	//Ex: A 4 tap system would be = 3
	numberOfTaps int = 7
)

var (
	//This is how we will set the tap system's ID
	tapUUID string = "d95c9dcc-3b99-4ad6-9a4e-e60418fcd359"
	tapControlID = 1
	authToken string

	//Keeps track of whether connection is alive
	webConnectionAlive bool = true

	defaultIDP string = "https://twirpz.files.wordpress.com/2015/06/twitter-avi-gender-balanced-figure.png?w=640"
)

type Order struct {
	//Tap_Orders Data
	uuid string
	orderID int
	user int
	tapID int
	beerID int
	price string
	size string
	//Tap(s) to pour on with array value being drink size
	tap [numberOfTaps + 1]int

	//Tap_Users Data
	email string
	firstName string
	lastName string
	dob string
	mobilePhone string
	pictureURL string
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
	TapControlPOST struct {
		TapControlID int 	`json:"id"`
		TapUUID string		`json:"uuid"`
	} `json:"tap_control"`
}

type AuthResponse struct {
	TapControlID int 		`json:"id"`
	TapUUID string			`json:"uuid"`
	AuthenToken string 	`json:"authentication_token"`
}

type TapUserResponse struct{
	UserID int							`json:"id"`
	UserEmail string				`json:"email"`
	FirstName string				`json:"first_name"`
	LastName	string				`json:"last_name"`
	DoB	string							`json:"dob"`
	MobilePhone string			`json:"string"`
	PhotoURL string					`json:"photo"`
	DriverLFrontURL string	`json:"drivers_license_front"`
	DriverLBackURL string		`json:"drivers_license_back"`
}

//Holds orders (uses labels from makeItems()) and container for window
type beerbot struct{
	orders [8]orderInfo
	c *fyne.Container
}

//Holds the order information and image for organization and ease for passing around the program
type orderInfo struct {
	tapNum 		int
	label 	 *canvas.Text
	status 	 *canvas.Text
	userName *canvas.Text
	dob 		 *canvas.Text
	beer 		 *canvas.Text
	size 		 *canvas.Text
	img 		 *canvas.Image
	}

// ######################## MAIN PROGRAM PROGRAM PROGRAM #######################
func main() {

	a := app.New()
	w := a.NewWindow("BeerBot Tap Display")
	b := beerbot{}

	b.c = container.NewPadded(b.makeUI())
	w.SetContent(b.c)

	go runProgram(b)

	w.Resize(fyne.NewSize(1920, 1080)) //wouldn't fit on my screen lol
	w.SetFixedSize(true) //the weird stuff doesn't happen when I put this line in
	w.ShowAndRun()

	endProgram()
}


func runProgram(b beerbot) {
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

	//Authenticate with API to get our Authentication Token for API communication
	authToken = authTapController(tapUUID, tapControlID)

	//Main program loop
	for webConnectionAlive == true{


		//Sleep for 100ms to save CPU time
		time.Sleep(100*time.Millisecond)

		//Check order queue for orders to pull
		orderIdToServe := checkOrders(tapUUID, authToken)
		//fmt.Println("Order IDs to serve: ", orderIdToServe)

		//If there are orders to serve then let us fullfill them
		if len(orderIdToServe) >= 1 {

				for i := 0; i < len(orderIdToServe); i++ {
					//Get user orders
					userOrders := getOrderData(tapUUID, orderIdToServe[i], authToken)

					go togglePour(*userOrders, b)

				}
				//fmt.Println("Order ID Array before processOrder: ", orderIdToServe)
				//fmt.Println("len(orderIdToServe): ", len(orderIdToServe))

				// Mark the orders we just fullfilled/poured as poured on the orders API
				for i := len(orderIdToServe) - 1; i >= 0; i-- {
					// Call to process order
					if processOrder(tapUUID, orderIdToServe[i], authToken) == true{
							// Remove order from local queue
							orderIdToServe = append(orderIdToServe[:i], orderIdToServe[i+1:]...)
						}
				}
		}
	}
	//Run all the stuff needed to cleanly exit ( IMPORTANT THIS HAPPENS )
	endProgram()
}

//Called by makeUI() to create a VBox with several labels and an image
func (b *beerbot) makeTapItems(tapNum int) fyne.CanvasObject {
	//Sets tap number and tap label
	b.orders[tapNum-1].tapNum = tapNum
	tapLabel := "Tap " + strconv.Itoa(tapNum)
	b.orders[tapNum-1].label = canvas.NewText(tapLabel, color.Gray{200})
	b.orders[tapNum-1].label.Alignment = fyne.TextAlignCenter
	b.orders[tapNum-1].label.TextSize = 18
	//Scan tag
	b.orders[tapNum-1].status = canvas.NewText("Scan Tag to Pour", color.Gray{200})
	b.orders[tapNum-1].status.Alignment = fyne.TextAlignCenter
	b.orders[tapNum-1].status.TextSize = 18
	//name
	b.orders[tapNum-1].userName = canvas.NewText("-", color.Gray{128})
	b.orders[tapNum-1].userName.Alignment = fyne.TextAlignCenter
	b.orders[tapNum-1].userName.TextSize = 12
	//dob
	b.orders[tapNum-1].dob = canvas.NewText("-", color.Gray{128})
	b.orders[tapNum-1].dob.Alignment = fyne.TextAlignCenter
	b.orders[tapNum-1].userName.TextSize = 10
	//Beer choice
	b.orders[tapNum-1].beer = canvas.NewText("Miller Lite", color.Gray{128})
	b.orders[tapNum-1].beer.Alignment = fyne.TextAlignCenter
	b.orders[tapNum-1].userName.TextSize = 10
	//Pour size
	b.orders[tapNum-1].size = canvas.NewText("-", color.Gray{128})
	b.orders[tapNum-1].size.Alignment = fyne.TextAlignCenter
	b.orders[tapNum-1].userName.TextSize = 10

	//Initial image
	myURL := "https://twirpz.files.wordpress.com/2015/06/twitter-avi-gender-balanced-figure.png?w=640" //Pic of sam
	img := loadImage(myURL)

	b.orders[tapNum-1].img = img


	return container.NewVBox(b.orders[tapNum-1].label,
												b.orders[tapNum-1].status, b.orders[tapNum-1].userName,
												b.orders[tapNum-1].dob, b.orders[tapNum-1].beer,
												b.orders[tapNum-1].size, b.orders[tapNum-1].img)
}//end makeTapIems


//creates the 8 orders and put them in a 4x2 grid
func (b *beerbot) makeUI() fyne.CanvasObject {
	items := []fyne.CanvasObject{}

	//For each order, create a set of labels and a canvas image
	for _, v := range []int{1,2,3,4,5,6,7,8} {
		orderContainer := b.makeTapItems(v) //creates label and image and put them into VBOX
		items = append(items, orderContainer) //adds the VBox to an array
	}

	return container.NewGridWithRows(2, items...)
}//end makeUI


//Add the id face picture with the given parameters
func loadImage(url string) *canvas.Image {
		req, _ := http.NewRequest("GET", url, nil)
		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		imgLoc := os.TempDir() + "/idImg.jpg" //Location for any device
		err := ioutil.WriteFile(imgLoc, body, 0644)
		if err != nil{
			log.Fatal ("ioutil TempFile error", err)
		}

		img := canvas.NewImageFromFile(imgLoc)
		//For some reason when this is called, all the images are updated instead of just 1
		// img.FillMode = canvas.ImageFillContain //try commenting this out
		img.SetMinSize(fyne.NewSize(300,300)) // approx ~1:1.5 (ID picture ratio)

		return img
}//end loadImage

//Change the existing id image shown on the GUI
func changeImage (url string, img *canvas.Image){
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	imgLoc := os.TempDir() + "/idImg.jpg"
	err := ioutil.WriteFile(imgLoc, body, 0644)
	if err != nil{
		log.Fatal ("ioutil TempFile error", err)
	}

	file, err := os.Open(imgLoc)
	if err != nil {
		panic(err)
	}

	img.File = file.Name()
	// img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(300,300)) // approx ~1:1.5 (ID picture ratio)
	img.Refresh()

	file.Close()
}//end changeImage



func (b *beerbot) changeLabel (customerOrder Order, tap int) {
	switch tap {
		case 1:
			t := b.orders[0].status
			t.Text = "Pour Now!"
			t.Refresh()

			a := b.orders[0].userName
			a.Text = customerOrder.firstName + " " + customerOrder.lastName
			a.Refresh()

			c := b.orders[0].dob
			c.Text = customerOrder.dob
			c.Refresh()

			d := b.orders[0].beer
			d.Text = strconv.Itoa(customerOrder.beerID)
			d.Refresh()

			e := b.orders[0].size
			e.Text = customerOrder.size + " Oz"
			e.Refresh()

		case 2:
			//b.orders[1].label.Text = "Changed label"
			t := b.orders[1].label
			t.Text = "Changed label"
			t.Refresh()

			fmt.Println("Changed label case 2")
		case 3:
		case 4:
		case 5:
		case 6:
		case 7:
		case 8:
		default:
	}
}

func (b *beerbot) clearLabel (tap int) {
	switch tap {
		case 1:
			t := b.orders[0].status
			t.Text = "Scan Tag to Pour"
			t.Refresh()

			a := b.orders[0].userName
			a.Text = "-"
			a.Refresh()

			c := b.orders[0].dob
			c.Text = "-"
			c.Refresh()

			d := b.orders[0].beer
			d.Text = "-"
			d.Refresh()

			e := b.orders[0].size
			e.Text = "-"
			e.Refresh()


		case 2:
			//b.orders[1].label.NewText = "Changed label"
			t := b.orders[1].label
			t.Text = "Changed label"
			t.Refresh()

			fmt.Println("Clear label case 2")

		case 3:
		case 4:
		case 5:
		case 6:
		case 7:
		case 8:
		default:
	}
}


//Get user data for given order
func getUserData(customerOrder *Order, authToken string) {
	url := "http://96.30.244.56:3000/api/v1/tap_users/"+ strconv.Itoa(customerOrder.user)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "96.30.244.56:3000")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Authorization", authToken)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	//fmt.Println(res)
	//fmt.Println("getUserData body: ", string(body))

	var verifyResp []byte = body
	var verifyData TapUserResponse

	err := json.Unmarshal(verifyResp, &verifyData)
	if err != nil {
		fmt.Println("unmarshal error:", err)
	}


	customerOrder.email = verifyData.UserEmail
	customerOrder.firstName = verifyData.FirstName
	customerOrder.lastName = verifyData.LastName
	customerOrder.dob = verifyData.DoB
	customerOrder.mobilePhone = verifyData.MobilePhone
	customerOrder.pictureURL = verifyData.PhotoURL

}


//Get orders from the orderqueue
func getOrderData(uuid string, orderID int, authToken string) *Order {
	o := Order{uuid: tapUUID}

	url := "http://96.30.244.56:3000/api/v1/tap_orders/"+ strconv.Itoa(orderID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "96.30.244.56:3000")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Authorization", authToken)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	//fmt.Println(res)
	//fmt.Println("getOrders body: ", string(body))

	var verifyResp []byte = body
	var verifyData OrderResponse

	err := json.Unmarshal(verifyResp, &verifyData)
	if err != nil {
		fmt.Println("unmarshal error:", err)
	}

	//fmt.Println("Verify Order Response Dump:")
	//fmt.Println("verifyData: ", verifyData)

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
		//Calculate the number of pulses for flow sensor for the local order struct
		pulses, errPulseConversion := strconv.ParseFloat(verifyData.Size,32)
		if errPulseConversion != nil {
    	fmt.Println("size to pulse conversion error", errPulseConversion)
   	}
		pulses = pulses/pulseDivConst
		//Round our float and store it away in local order struct
		o.tap[verifyData.TapID-1] = int(math.Round(pulses))

		//Get user data filled into the order struct
		getUserData(&o, authToken)
	}

	fmt.Println("getOrders o: ", o)

	return &o
}


//Check for orders to be served, returns array of ordersId to be served
func checkOrders(uuid string, authToken string) []int{
	var orderIDs []int
	//fmt.Println("Fetch order ids to fullfill")
	url := "http://96.30.244.56:3000/api/v1/tap_orders"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "96.30.244.56:3000")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Authorization", authToken)

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
func processOrder(uuid string, orderID int, authToken string) bool{
	url := "http://96.30.244.56:3000/api/v1/tap_orders/"+ strconv.Itoa(orderID)

	orderResp := CheckResponse{orderID,true}

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
	req.Header.Add("Authorization", authToken)

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
func togglePour(customerOrder Order, b beerbot) {
	//This is just a timeout function so that the program will timeout
	c1 := make(chan string, 1)
	// Run your long running function in it's own goroutine and pass back it's
	// response into our channel.

	//Set this outside the switch range as a fail safe
	tapToClose := 9

	go func() {
		var wg1 sync.WaitGroup

		//Solenoid normal state = closed
		for i := 0; i <= numberOfTaps; i++ {
			if customerOrder.tap[i] != 0 {
				//Update GUI with retreived user order
				changeImage(customerOrder.pictureURL, b.orders[i].img)
				b.changeLabel(customerOrder, i+1)
				//Add to our waitgroup
				wg1.Add(1)
				//Shoot off a thread to pour customer's order on specific tap
				go gpio_rpi.Pour(customerOrder.tap[i], i+1, &wg1)
				tapToClose = i+1
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

			//Clear GUI after finished pouring order
			//clearGUIOrder(tapToClose, oL1, oL2, oL3, oL4, oL5, oL6, oL7, oL8)
			go changeImage("https://twirpz.files.wordpress.com/2015/06/twitter-avi-gender-balanced-figure.png?w=640", b.orders[tapToClose-1].img)

		case <-time.After(20 * time.Second):
			fmt.Println("out of time :(")
			// Close solenoids incase timeout
			gpio_rpi.CloseSolenoids(tapToClose)

			//Clear GUI after finished pouring order
			//clearGUIOrder(tapToClose, oL1, oL2, oL3, oL4, oL5, oL6, oL7, oL8)
			go changeImage("https://twirpz.files.wordpress.com/2015/06/twitter-avi-gender-balanced-figure.png?w=640", b.orders[tapToClose-1].img)

	}
}


//Tap Controller Authentication with API
func authTapController(uuid string, tapControlID int) string{
	url := "http://96.30.244.56:3000/api/v1/tap_sessions"

	authPost := AuthPOST{}
	authPost.TapControlPOST.TapControlID = tapControlID
	authPost.TapControlPOST.TapUUID = uuid

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

	var authResp []byte = body
	var verifyAuth AuthResponse

	errR := json.Unmarshal(authResp, &verifyAuth)
	if errR != nil {
		fmt.Println("unmarshal error:", errR)
	}

	fmt.Println(verifyAuth.AuthenToken)
	return verifyAuth.AuthenToken
}


//Run all the stuff needed to cleanly exit ( IMPORTANT THIS HAPPENS )
func endProgram(){
	//Close GPIO/clear GPIO memory at end of program ( IMPORTANT THIS HAPPENS )
	gpio_rpi.CloseSolenoids(0)
	gpio.Close()
	log.Println("Program ended cleanly!")
	os.Exit(1)
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
