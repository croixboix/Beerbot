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


func makeTapItems(tapNum int) fyne.CanvasObject {
	//Tap number
	tapLabel := "Tap " + strconv.Itoa(tapNum+1) //I thought values would increase as going across
	label := canvas.NewText(tapLabel, color.Gray{128})
	label.Alignment = fyne.TextAlignCenter
	//Scan tag
	status := canvas.NewText("Scan Tag to Pour", color.Gray{128})
	status.Alignment = fyne.TextAlignCenter
	//name
	userName := canvas.NewText("User Name", color.Gray{128})
	userName.Alignment = fyne.TextAlignCenter
	//dob
	dob := canvas.NewText("1/1/1984", color.Gray{128})
	dob.Alignment = fyne.TextAlignCenter
	//Beer choice
	beer := canvas.NewText("Miller Lite", color.Gray{128})
	beer.Alignment = fyne.TextAlignCenter
	//Pour size
	size := canvas.NewText("12 Ounces", color.Gray{128})
	size.Alignment = fyne.TextAlignCenter

	//Initial image
	myURL := "https://miro.medium.com/max/868/1*Hyd_x4yW3H_wxn_f8tFYLQ.png" //Pic of sam
	img := loadImage(myURL)

	// newURL := "https://i.kym-cdn.com/entries/icons/original/000/035/432/41rtwpO9McL.jpg" //pic of Than
	// changeImage(newURL, img)

	return container.NewVBox(label, status, userName, dob, beer, size, img)
}//end makeTapIems


func makeUI() fyne.CanvasObject {
	items := []fyne.CanvasObject{}

	for i := range []int{1,2,3,4,5,6,7,8} {
		img := makeTapItems(i)
		items = append(items, img)
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
		// err := ioutil.WriteFile("/home/wint/bPics/id1.jpg", body, 0644)
		err := ioutil.WriteFile(imgLoc, body, 0644)
		if err != nil{
			log.Fatal ("ioutil TempFile error", err)
		}

		img := canvas.NewImageFromFile(imgLoc)
		img.SetMinSize(fyne.NewSize(125,125)) // app/rox ~1:1.5 (ID picture ratio)

		return img
}//end loadImage

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
	img.Refresh()

	file.Close()
}//end changeImage

func main() {
	a := app.New()
	w := a.NewWindow("BeerBot Tap Display")
	w.SetContent(container.NewPadded(makeUI()))
	// w.Resize(fyne.NewSize(1920,1080))
	w.Resize(fyne.NewSize(1024, 700)) //wouldn't fit on my screen lol
	w.ShowAndRun()


}//end main
