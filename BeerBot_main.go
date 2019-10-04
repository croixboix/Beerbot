package main

import (
  "fmt"
  gpio "pi/gpio"
)

var (

)


func main() {

  //Solenoid normal state = closed
  //Open solenoid
  gpio.toggle()
  fmt.Println("Solenoid opened")

  fmt.Println("Begin measuring flow (12oz cutoff)")
  gpio.pour()
  fmt.Println("Pour limit reached! (12oz)")

  fmt.Println("Closing solenoid")
  gpio.toggle()

}
