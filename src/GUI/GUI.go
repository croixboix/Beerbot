package main

import (
	"image/color"
	"io/ioutil"
	"net/http"
	"io"
	"os"
	"log"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)


func makeTapItems() fyne.CanvasObject {
	label := canvas.NewText("Tap 1", color.Gray{128})
	label.Alignment = fyne.TextAlignCenter

	status := canvas.NewText("Scan Tag to Pour", color.Gray{128})
	status.Alignment = fyne.TextAlignCenter

	userName := canvas.NewText("User Name", color.Gray{128})
	userName.Alignment = fyne.TextAlignCenter

	dob := canvas.NewText("1/1/1984", color.Gray{128})
	dob.Alignment = fyne.TextAlignCenter

	beer := canvas.NewText("Miller Lite", color.Gray{128})
	beer.Alignment = fyne.TextAlignCenter

	size := canvas.NewText("12 Ounces", color.Gray{128})
	size.Alignment = fyne.TextAlignCenter


	myURL := "https://miro.medium.com/max/868/1*Hyd_x4yW3H_wxn_f8tFYLQ.png"
	img := loadImage(myURL)


	return container.NewVBox(label, status, userName, dob, beer, size, img)
}


func makeUI() fyne.CanvasObject {
	items := []fyne.CanvasObject{}

	for range []int{1,2,3,4,5,6,7,8} {
		img := makeTapItems()
		items = append(items, img)
	}

	return container.NewGridWithRows(2, items...)
}


//Add the id face picture with the given parameters
func loadImage(url string) fyne.CanvasObject {
		//Grabs content from url
 		response, e := http.Get(url)
		if e != nil {
				log.Fatal("Unable to Get URL", e)
				img := canvas.NewRectangle(color.Black)
				return img
		}
		defer response.Body.Close()

		//creates tmp file wth a unique name
		file, err := ioutil.TempFile(os.TempDir(), "userPic.*.jpg")
		if err != nil{
			log.Fatal("ioutil TempFile error", err)
			img := canvas.NewRectangle(color.Black)
			return img
		}
		defer file.Close()

		// Use io.Copy to just dump the response body to the file. This supports huge files
		_, err = io.Copy(file, response.Body) //copy data from get request into file
		if err != nil {
			log.Fatal(err)
			img := canvas.NewRectangle(color.Black)
			return img
		}
		fmt.Println(file.Name())
		img := canvas.NewImageFromFile(file.Name())
		img.SetMinSize(fyne.NewSize(255,340)) // 4:3 aspect ratio
		img.FillMode = canvas.ImageFillContain

		return img
}



func main() {
	a := app.New()
	w := a.NewWindow("BeerBot Tap Display")
	w.SetContent(container.NewPadded(makeUI()))
	w.Resize(fyne.NewSize(1920,1080))
	w.ShowAndRun()


}
