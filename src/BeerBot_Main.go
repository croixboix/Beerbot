package main

import(
  "fmt"
  gpio "./gpio/"
  "time"
)

var (

)


func main() {

  //Solenoid normal state = closed
  //Open solenoid
  gpio.Taggle()
  fmt.Println("Solenoid opened")

 fmt.Println("Begin measuring flow (12oz cutoff)")
  gpio.Pour()
 fmt.Println("Pour limit reached! (12oz)")

  fmt.Println("Closing solenoid")
  gpio.Taggle()
  
  time.Sleep(time.Second / 5)

}
