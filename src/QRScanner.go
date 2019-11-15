package main

import(
  "fmt"
  "encoding/json"
  //"bufio"
  //"os"
  "strings"
  "net/http"
  "io/ioutil"
)

var (
)

//Verify that the order exists on the API order list
func verifyOrder(uname string) []byte{
  url := "http://96.30.245.134:3000/orders/verify"
  payload := strings.NewReader("{\n\t\"order\": {\n\t\t\"username\": \""+uname+"\"\n\t}\n}")
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
func processOrder(uname string) []byte{
  url := "http://96.30.245.134:3000/orders/processed"

	payload := strings.NewReader("{\n\t\"order\": {\n\t\t\"username\": \""+uname+"\"\n\t}\n}")

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


func main() {

  type verifyResponse struct {
    ID int `json:"id"`
    Username string `json:"username"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string  `json:"updated_at"`
    URL string  `json:"url"`
  }

  type processResponse struct {
    Processed bool `json:"processed"`
  }

  var user string = "test"

  fmt.Println("Verify Order")

  var verifyResponse []byte = verifyOrder(user)
  var verifyData verifyResponse

  err := json.Unmarshal(verifyResponse, &verifyData)
  if err != nil {
        fmt.Println("error:", err)
    }

  fmt.Println("Verify Order Response Dump:")
  fmt.Println("id: ", verifyData.ID)
  fmt.Println("username: ", verifyData.Username)
  fmt.Println("created_at: ", verifyData.CreatedAt)
  fmt.Println("updated_at: ", verifyData.UpdatedAt)
  fmt.Println("url: ", verifyData.URL)


  if(verifyData.Username != "null"){

    fmt.Println("Process/Delete Order")

    var processData processResponse

    var processResponse []byte = processOrder(verifyData.Username)

    err := json.Unmarshal(processResponse, &processData)
    if err != nil {
          fmt.Println("error:", err)
      }

    fmt.Println("Process Order Response Dump:")
    fmt.Println("processed: ", processData.Processed)

    if(processData.Processed == true){
      //Call pour
      fmt.Println("ORDER PROCESSED, LET USER POUR")
    } else{
      fmt.Println("ORDER DOES NOT EXIST, DO NOT LET USER POUR")
    }
  }




  /*
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
  */


}
