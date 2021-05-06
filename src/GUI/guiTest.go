package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"time"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"log"
)

var(
defaultIDP string = "https://twirpz.files.wordpress.com/2015/06/twitter-avi-gender-balanced-figure.png?w=640"
)

type beerbot struct {
	c fyne.Canvas
	id1 *canvas.Image
	headerL *widget.Label
}

type Order struct {
	uuid string //Tap's UUID
	user int //Order's user/customer
	orderID int
	tapID int
}

type orderLabels struct {
	tapIDL *widget.Label
	beerIDL *widget.Label
	priceL *widget.Label
	sizeL *widget.Label
	FirstLastL *widget.Label
	DOBL *widget.Label
	emailL *widget.Label
	userPic fyne.CanvasObject
}

// This will be the GUI code
func main() {
	a := app.New()
	w := a.NewWindow("Beerbot")
	// myCanvas := w.Canvas()
	bGUI := beerbot{
		c: w.Canvas(),
		// id1: downloadImage("https://i.pinimg.com/originals/64/86/60/648660b8d170ba0540bc1ed50f33de4e.jpg"),
		headerL: widget.NewLabel("test"),
	}


	//Order creation
	o1 := Order{uuid: "testTap", user: 1, orderID: 1, tapID: 1}

	//Label creation
	hL  := orderLabels{
		tapIDL: widget.NewLabel("Tap ID:"),
		beerIDL:widget.NewLabel("Beer ID: "),
		priceL:widget.NewLabel("Price: "),
		sizeL:widget.NewLabel("Size : "),
		FirstLastL:widget.NewLabel("Name: "),
		DOBL:widget.NewLabel("DOB: "),
		emailL:widget.NewLabel("Email: ")}
	oL1 := orderLabels{
		tapIDL: widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		beerIDL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		priceL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		sizeL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		FirstLastL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		DOBL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		emailL:widget.NewLabelWithStyle("-",fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		userPic:setUserPic(defaultIDP),
	}
  oL2 := orderLabels{
		tapIDL: widget.NewLabelWithStyle("-",fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
	  beerIDL:widget.NewLabelWithStyle("-",fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
	  priceL:widget.NewLabelWithStyle("-",fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
	  sizeL:widget.NewLabelWithStyle("-",fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
	  FirstLastL:widget.NewLabelWithStyle("-",fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
	  DOBL:widget.NewLabelWithStyle("-",fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		emailL:widget.NewLabelWithStyle("-",fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		userPic:setUserPic(defaultIDP),
	}
  oL3 := orderLabels{
		tapIDL: widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		beerIDL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		priceL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		sizeL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		FirstLastL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		DOBL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		emailL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		userPic:setUserPic(defaultIDP),
	}
  oL4 := orderLabels{
		tapIDL: widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		beerIDL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		priceL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		sizeL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		FirstLastL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		DOBL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		emailL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		userPic:setUserPic(defaultIDP),
	}
  oL5 := orderLabels{
		tapIDL: widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		beerIDL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		priceL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		sizeL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		FirstLastL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		DOBL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		emailL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		userPic:setUserPic(defaultIDP),
	}
  oL6 := orderLabels{
		tapIDL: widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		beerIDL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		priceL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		sizeL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		FirstLastL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		DOBL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		emailL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		userPic:setUserPic(defaultIDP),
	}
	oL7 := orderLabels{
		tapIDL: widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		beerIDL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		priceL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		sizeL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		FirstLastL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		DOBL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		emailL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		userPic:setUserPic(defaultIDP),
	}
	oL8 := orderLabels{
		tapIDL: widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		beerIDL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		priceL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		sizeL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		FirstLastL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		DOBL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		emailL:widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		userPic:setUserPic(defaultIDP),
	}


	headerVBox := container.New(layout.NewVBoxLayout(),
		hL.tapIDL, hL.beerIDL, hL.priceL, hL.sizeL,
	 	hL.FirstLastL, hL.DOBL, hL.emailL)
	userVBox1 := container.New(layout.NewVBoxLayout(),
		oL1.tapIDL, oL1.beerIDL, oL1.priceL, oL1.sizeL,
		oL1.FirstLastL, oL1.DOBL, oL1.emailL, oL1.userPic)
	userVBox2 := container.New(layout.NewVBoxLayout(),
		oL2.tapIDL, oL2.beerIDL, oL2.priceL, oL2.sizeL,
		oL2.FirstLastL, oL2.DOBL, oL2.emailL, oL2.userPic)
	userVBox3 := container.New(layout.NewVBoxLayout(),
		oL3.tapIDL, oL3.beerIDL, oL3.priceL, oL3.sizeL,
		oL3.FirstLastL, oL3.DOBL, oL3.emailL, oL3.userPic)
	userVBox4 := container.New(layout.NewVBoxLayout(),
		oL4.tapIDL, oL4.beerIDL, oL4.priceL, oL4.sizeL,
		oL4.FirstLastL, oL4.DOBL, oL4.emailL, oL4.userPic)
	userVBox5 := container.New(layout.NewVBoxLayout(),
		oL5.tapIDL, oL5.beerIDL, oL5.priceL, oL5.sizeL,
		oL5.FirstLastL, oL5.DOBL, oL5.emailL, oL5.userPic)
	userVBox6 := container.New(layout.NewVBoxLayout(),
		oL6.tapIDL, oL6.beerIDL, oL6.priceL, oL6.sizeL,
		oL6.FirstLastL, oL6.DOBL, oL6.emailL, oL6.userPic)
	userVBox7 := container.New(layout.NewVBoxLayout(),
		oL7.tapIDL, oL7.beerIDL, oL7.priceL, oL7.sizeL,
		oL7.FirstLastL, oL7.DOBL, oL7.emailL, oL7.userPic)
	userVBox8 := container.New(layout.NewVBoxLayout(),
		oL8.tapIDL, oL8.beerIDL, oL8.priceL, oL8.sizeL,
		oL8.FirstLastL, oL8.DOBL, oL8.emailL, oL8.userPic)
	//Layout Config
	top := widget.NewLabelWithStyle(
					"Orders", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	left := headerVBox //headers

	middle := container.New(layout.NewHBoxLayout(),
							userVBox1, userVBox2, userVBox3, userVBox4, userVBox5, userVBox6,
							userVBox7, userVBox8)

	content := fyne.NewContainerWithLayout(
										layout.NewBorderLayout(top, nil, left, nil),
																			top, left, middle)

	bGUI.c.SetContent(content)
	// bGUI.bCanvas.SetContent(content)
	//update guigithub.com/fredbi/uri
  // go changeContent(myCanvas, oL1, oL2, oL3, oL4, oL5, oL6, oL7, oL8, o1)
	go changeContent(bGUI.c, oL1, oL2, oL3, oL4, oL5, oL6, oL7, oL8, o1, content)

	w.Resize(fyne.NewSize(300, 200))
	w.ShowAndRun()
}


// This will be basically everything in main currently
func changeContent(c fyne.Canvas, oL1 orderLabels, oL2 orderLabels,
									oL3 orderLabels, oL4 orderLabels, oL5 orderLabels,
									oL6 orderLabels, oL7 orderLabels, oL8 orderLabels,
									o Order, content fyne.CanvasObject) {
  // for true{
	// // Changes Value
	time.Sleep(time.Second)
	oL1.tapIDL.SetText("1")
	time.Sleep(time.Second)
	oL2.beerIDL.SetText("2")

	oL2.userPic = setUserPic("https://www.auburn.edu/administration/tigercard/images/sample_id_photo.jpg")
	c.SetContent(content)
	oL2.beerIDL.SetText("2")
	c.Refresh(content)

	time.Sleep(time.Second)
	oL3.priceL.SetText("3")
	time.Sleep(time.Second)
	oL4.sizeL.SetText("4")
	time.Sleep(time.Second)
	oL5.FirstLastL.SetText("5")
	time.Sleep(time.Second)
	oL6.DOBL.SetText("6")
	time.Sleep(time.Second)
	oL7.emailL.SetText("7")
	time.Sleep(time.Second)
	oL8.DOBL.SetText("8")
  // }
}//end changeContent

//Add the id face picture with the given parameters
func setUserPic(url string) fyne.CanvasObject {
		//Grabs content from url
 		response, e := http.Get(url)
		if e != nil {
				log.Fatal("Unable to Get URL", e)
		}
		defer response.Body.Close()

		//creates tmp file wth a unique name
		file, err := ioutil.TempFile(os.TempDir(), "idPic.*.jpg")
		if err != nil{
			log.Fatal ("ioutil TempFile error", err)
		}
		defer file.Close()

		// Use io.Copy to just dump the response body to the file. This supports huge files
		_, err = io.Copy(file, response.Body) //copy data from get request into file
		if err != nil {
				log.Fatal(err)
		}
		fmt.Println(file.Name())
		img := canvas.NewImageFromFile(file.Name())
		img.SetMinSize(fyne.NewSize(125,125)) // approx ~1:1.5 (ID picture ratio)

		return img
} //end addFacePic
