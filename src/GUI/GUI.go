package main

import (
	"image/color"
	"io/ioutil"
	"net/http"
	// "io"
	"os"
	"log"
	// "fmt"
	"strconv"
	"time"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

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

//Called by makeUI() to create a VBox with several labels and an image
func (b *beerbot) makeTapItems(tapNum int) fyne.CanvasObject {
	//Sets tap number and tap label
	b.orders[tapNum-1].tapNum = tapNum
	tapLabel := "Tap " + strconv.Itoa(tapNum)
	b.orders[tapNum-1].label = canvas.NewText(tapLabel, color.Gray{128})
	b.orders[tapNum-1].label.Alignment = fyne.TextAlignCenter
	//Scan tag
	b.orders[tapNum-1].status = canvas.NewText("Scan Tag to Pour", color.Gray{128})
	b.orders[tapNum-1].status.Alignment = fyne.TextAlignCenter
	//name
	b.orders[tapNum-1].userName = canvas.NewText("User Name", color.Gray{128})
	b.orders[tapNum-1].userName.Alignment = fyne.TextAlignCenter
	//dob
	b.orders[tapNum-1].dob = canvas.NewText("1/1/1984", color.Gray{128})
	b.orders[tapNum-1].dob.Alignment = fyne.TextAlignCenter
	//Beer choice
	b.orders[tapNum-1].beer = canvas.NewText("Miller Lite", color.Gray{128})
	b.orders[tapNum-1].beer.Alignment = fyne.TextAlignCenter
	//Pour size
	b.orders[tapNum-1].size = canvas.NewText("12 Ounces", color.Gray{128})
	b.orders[tapNum-1].size.Alignment = fyne.TextAlignCenter

	//Initial image
	myURL := "https://miro.medium.com/max/868/1*Hyd_x4yW3H_wxn_f8tFYLQ.png" //Pic of sam
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
		//img.FillMode = canvas.ImageFillOriginal
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(125,125)) // approx ~1:1.5 (ID picture ratio)

		return img
}//end loadImage

//Change the existing id image shown on the GUI
func changeImage (url string, img *canvas.Image){
	time.Sleep(3*time.Second)

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
	// img.FillMode = canvas.ImageFillOriginal //same issue as above
	img.SetMinSize(fyne.NewSize(125,125)) // approx ~1:1.5 (ID picture ratio)
	img.Refresh()

	file.Close()
}//end changeImage


func main() {
	a := app.New()
	w := a.NewWindow("BeerBot Tap Display")
	b := beerbot{}

	b.c = container.NewPadded(b.makeUI())
	w.SetContent(b.c)

	//changes the ID image
	go changeImage("https://i.kym-cdn.com/entries/icons/original/000/035/432/41rtwpO9McL.jpg", b.orders[0].img)

	w.Resize(fyne.NewSize(1024, 700)) //wouldn't fit on my screen lol
	w.ShowAndRun()
}//end main
