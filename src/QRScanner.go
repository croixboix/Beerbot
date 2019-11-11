package main

import(
  "fmt"
  "bufio"
  "os"
  "strings"
  "net/http"
  "io/ioutil"
)

var (

)

func requestOrders() string {

  url := "http://96.30.245.134:3000/orders/index"
  payload := strings.NewReader("username=test")
  req, _ := http.NewRequest("GET", url, payload)
  req.Header.Add("Authorization", "eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE1NzM1ODI1MTh9.NbhB_MFF9pV7dBLBp8o5dKtY0EgMR3sK621hN34G1Tw")
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  req.Header.Add("User-Agent", "PostmanRuntime/7.19.0")
  req.Header.Add("Accept", "/")
  req.Header.Add("Cache-Control", "no-cache")
  req.Header.Add("Postman-Token", "711bb03c-4ad6-4809-a89f-c372f1eaab81,af06cae4-ff57-4998-a9b5-26d004bb32d5")
  req.Header.Add("Accept-Encoding", "gzip, deflate")
  req.Header.Add("Cookie", "_beer_bot_api_session=AJ03Hd2KhYHVoMe%2BiRmykVhp%2BdGRQ%2Fi%2BgElgQHFRdJUKEoST6pGtRJM%2Ft4wd9vbr8ZcK%2BlPLNPz8MtS3p3Er9AqCKjZS7By%2Bmx6PsjbxPyhUtQ3QwRHZ7i7eQsvmLs8o3IIVDBtw5PJVi57CN4FJZ5lCjqmXnQ08VfiyFCh4rMa7ixGRXxtIF7nINv7et78s7f9gTwve132c%2FPsYYPQ2gaLfHZ8oJKOjuyNTqoTqz23tssWqCPeiigUGZU6yMBE50J%2FcsTYKTZDVDopqHppc61MkXdzjEAnAubwvp5R%2Fljr0eRa%2BPbnBJNKyptIKIJOWRhd0oM8%3D--ELRctXZO7nN6Z8bp--LMiKintBHTmbRTV7R76NHg%3D%3D")
  req.Header.Add("Referer", "http://96.30.245.134:3000/orders/create")
  req.Header.Add("Connection", "keep-alive")
  req.Header.Add("cache-control", "no-cache")
  res, _ := http.DefaultClient.Do(req)
  defer res.Body.Close()
  body, _ := ioutil.ReadAll(res.Body)
  fmt.Println(res)
  fmt.Println(string(body))


  return string(body)
}


func processOrders() {
  //If sucessful return 0, else 1
  
  url := "http://96.30.245.134:3000/orders/processed"
  payload := strings.NewReader("username=test")
  req, _ := http.NewRequest("POST", url, payload)
  req.Header.Add("Authorization", "eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE1NzM1ODI1MTh9.NbhB_MFF9pV7dBLBp8o5dKtY0EgMR3sK621hN34G1Tw")
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  req.Header.Add("User-Agent", "PostmanRuntime/7.19.0")
  req.Header.Add("Accept", "/")
  req.Header.Add("Cache-Control", "no-cache")
  req.Header.Add("Postman-Token", "711bb03c-4ad6-4809-a89f-c372f1eaab81,606f5a7f-3ce3-4f71-8161-4b387d14df1e")
  req.Header.Add("Accept-Encoding", "gzip, deflate")
  req.Header.Add("Cookie", "_beer_bot_api_session=AJ03Hd2KhYHVoMe%2BiRmykVhp%2BdGRQ%2Fi%2BgElgQHFRdJUKEoST6pGtRJM%2Ft4wd9vbr8ZcK%2BlPLNPz8MtS3p3Er9AqCKjZS7By%2Bmx6PsjbxPyhUtQ3QwRHZ7i7eQsvmLs8o3IIVDBtw5PJVi57CN4FJZ5lCjqmXnQ08VfiyFCh4rMa7ixGRXxtIF7nINv7et78s7f9gTwve132c%2FPsYYPQ2gaLfHZ8oJKOjuyNTqoTqz23tssWqCPeiigUGZU6yMBE50J%2FcsTYKTZDVDopqHppc61MkXdzjEAnAubwvp5R%2Fljr0eRa%2BPbnBJNKyptIKIJOWRhd0oM8%3D--ELRctXZO7nN6Z8bp--LMiKintBHTmbRTV7R76NHg%3D%3D")
  req.Header.Add("Referer", "http://96.30.245.134:3000/orders/create")
  req.Header.Add("Connection", "keep-alive")
  req.Header.Add("cache-control", "no-cache")
  res, _ := http.DefaultClient.Do(req)
  defer res.Body.Close()
  body, _ := ioutil.ReadAll(res.Body)
  fmt.Println(res)
  fmt.Println(string(body))
  
}


func main() {

  orders := requestOrders()

  fmt.Println("requestOrders Return: ", string(orders))

  scanner := bufio.NewScanner(os.Stdin)

  for scanner.Scan() {
    fmt.Println("Scanned barcode: ", scanner.Text())
	if scanner.Text() == string(orders) {
		fmt.Println("QR CODE MATCH!")
		break
	}
  }

  if err := scanner.Err(); err != nil {
    fmt.Fprintln(os.Stderr, "reading standard input:", err)
  }
  

}