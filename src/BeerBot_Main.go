package main

import(
  "fmt"
  gpio "./gpio/"
  "time"
  "github.com/stianeikeland/go-rpio"
)

var (

)

// TODO:
// NOTE: Make sure you have disabled I2C interface in sudo raspi-config - I
// think it might be enabled by default which might cause your pin 2 to be
// changed from input mode to I2C.
// SEE: https://github.com/stianeikeland/go-rpio/issues/35
// Add dtoverlay=gpio-no-irq to /boot/config.txt and restart your pi
//  This disables IRQ which may break some other GPIO libs/drivers

func main() {

  //Initialize GPIO pins
  gpio.GPIO_INIT()
  fmt.Println("GPIO Initialized!")

  time.Sleep(time.Second)

  //Solenoid normal state = closed
  //Open solenoid
  gpio.Taggle()
  fmt.Println("Solenoid opened")

 fmt.Println("Begin measuring flow (12oz cutoff)")
  gpio.Pour()
 fmt.Println("Pour limit reached! (12oz)")

  fmt.Println("Closing solenoid")
  gpio.Taggle()


  //Give the Pi some time to catch up, remove this eventually just for debugging
  time.Sleep(time.Second)

  //Close GPIO/clear GPIO memory
  rpio.Close()



}
