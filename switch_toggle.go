package main

import(
	"fmt"
	"os"

	"github.com/stianeikeland/go-rpio"
)

var (
	pin = rpio.Pin(SAM REPLACE THIS WITH THE PIN NUMBER AND THE THREE DIGITS ON THE BACK)
)

func main() {

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()

	pin.Output()
	
	document.addEventListener('keydown',function(e) {
		var key = e.keyCode || e.which;
		if(key === 81) {
			pin.Toggle()
		}
 }
 
}

