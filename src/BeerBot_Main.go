package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/warthog618/gpio"
	gpio_rpi "gpio_rpi"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	//"time"
)

const (
	sizeSixOunce     int = 234
	sizeTwelveOunce  int = 468
	sizeSixteenOunce int = 624
)

var (
	user string
	//drinkSize int
	tapSize = [8]int{}
)

// TODO:
// NOTE: Make sure you have disabled I2C interface in sudo raspi-config - I
// think it might be enabled by default which might cause your pin 2 to be
// changed from input mode to I2C.
// SEE: https://github.com/stianeikeland/go-rpio/issues/35
// Add dtoverlay=gpio-no-irq to /boot/config.txt and restart your pi
//  This disables IRQ which may break some other GPIO libs/drivers

//Verify that the order exists on the API order list
func verifyOrder(uname string) []byte {
	url := "http://96.30.245.134:3000/orders/verify"
	payload := strings.NewReader("{\n\t\"order\": {\n\t\t\"username\": \"" + uname + "\"\n\t}\n}")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "96.30.245.134:3000")
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

	/*
	    Order does exist response:
	    {
	      "id": 591,
	      "username": "test",
	      "created_at": "2019-11-15T16:31:21.321Z",
	      "updated_at": "2019-11-15T16:31:21.321Z",
	      "url": "http://96.30.245.134:3000/orders/591.json"
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
}

// Tells API that order processed and deletes order from API order list
func processOrder(uname string) []byte {
	url := "http://96.30.245.134:3000/orders/processed"

	payload := strings.NewReader("{\n\t\"order\": {\n\t\t\"username\": \"" + uname + "\"\n\t}\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "PostmanRuntime/7.19.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Postman-Token", "aa411b8b-f735-484e-9f7c-2a871763a9dc,7b8a5042-c09d-47b0-a4ea-68ba0a8dda7a")
	req.Header.Add("Host", "96.30.245.134:3000")
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

func togglePour(drinkSize int, tap int) {
	//Solenoid normal state = closed

	fmt.Printf("Begin measuring flow for %d on tap %d\n", drinkSize, tap)
	gpio_rpi.Pour(drinkSize, tap)
	fmt.Println("Pour limit reached!")

}

func main() {

	//Initialize GPIO pins
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

	//Scan the Bar/QR Code
	fmt.Println("Scan Barcode Now!")
	//user = scanCode()

	//TEST VALUES HERE
	user = "test"
	tap = 8;
	tapSize[tap-1] = sizeSixOunce

	//Verify and process the order!
	fmt.Println("Verify Order")

	var verifyResp []byte = verifyOrder(user)
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
			togglePour(tapSize[tap-1], tap)
			fmt.Println("ORDER PROCESSED, LET USER POUR")
		} else {
			fmt.Println("ORDER DOES NOT EXIST, DO NOT LET USER POUR")
		}
	}

	//Close GPIO/clear GPIO memory at end of program
	gpio.Close()

}
