package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
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


	headerVBox := widget.NewVBox(
					hL.tapIDL, hL.beerIDL, hL.priceL, hL.sizeL,
				 	hL.FirstLastL, hL.DOBL, hL.emailL,)
	userVBox1 := widget.NewVBox(
		oL1.tapIDL, oL1.beerIDL, oL1.priceL, oL1.sizeL,
		oL1.FirstLastL, oL1.DOBL, oL1.emailL, oL1.userPic)

	//Layout Config
	top := widget.NewLabelWithStyle(
					"Orders", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	left := headerVBox //headers

	middle := widget.NewHBox(userVBox1)
	right := widget.NewVBox(bGUI.headerL, bGUI.id1)
	// bottom := widget.NewHBox(widget.NewVBox(oL1.userPic),
	// 					)


	content := fyne.NewContainerWithLayout(
										layout.NewBorderLayout(top, nil, left, nil), top, left, middle)
	// bGUI.content = fyne.NewContainerWithLayout(
	// 									layout.NewBorderLayout(top, bottom, left, nil),
	// 																		top, left, middle)

	fmt.Println("we here")

	bGUI.c.SetContent(content)


	// bGUI.bCanvas.SetContent(content)
	//update guigithub.com/fredbi/uri
  // go changeContent(myCanvas, oL1, oL2, oL3, oL4, oL5, oL6, oL7, oL8, o1)
	go changeContent(bGUI.c, oL1, o1, content)

	w.Resize(fyne.NewSize(300, 200))
	w.ShowAndRun()
}


// This will be basically everything in main currently
func changeContent(c fyne.Canvas, oL1 orderLabels,
									o Order, content fyne.CanvasObject) {
  // for true{

	// // Changes Value
	time.Sleep(time.Second)
	oL1.tapIDL.SetText("1")
	oL1.userPic = setUserPic("https://www.auburn.edu/administration/tigercard/images/sample_id_photo.jpg")
	oL1.beerIDL.SetText("2")
	c.SetContent(content)
	c.Refresh(content)

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
		img.SetMinSize(fyne.NewSize(100,125)) // approx ~1:1.5 (ID picture ratio)


		return img
} //end addFacePic

func (bGUI *beerbot) downloadImage(url string){
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
	// fmt.Println(file.Name())
	bGUI.id1.File = file.Name()
	canvas.Refresh(bGUI.id1)
	// img := canvas.NewImageFromFile(file.Name())
	// img.SetMinSize(fyne.NewSize(100,125)) // approx ~1:1.5 (ID picture ratio)

} //end downloadImage

// func (bGUI *beerbot) setUserPic(url string) {
// 		//Grabs content from url
//  		response, e := http.Get(url)
// 		if e != nil {
// 				log.Fatal("Unable to Get URL", e)
// 		}
// 		defer response.Body.Close()
//
// 		//creates tmp file wth a unique name
// 		file, err := ioutil.TempFile(os.TempDir(), "idPic.*.jpg")
// 		if err != nil{
// 			log.Fatal ("ioutil TempFile error", err)
// 		}
// 		defer file.Close()
//
// 		// Use io.Copy to just dump the response body to the file. This supports huge files
// 		_, err = io.Copy(file, response.Body) //copy data from get request into file
// 		if err != nil {
// 				log.Fatal(err)
// 		}
//
// 		bGUI.id1.File = file.Name()
// 		canvas.Refresh(bGUI.id1)
// 		fmt.Println(file.Name())
// 		// img := canvas.NewImageFromFile(file.Name())
// 		// img.SetMinSize(fyne.NewSize(100,125)) // approx ~1:1.5 (ID picture ratio)
// 		//
// 		// return img
// } //end addFacePic

// func newbGUI() *beerbot{
// 	bGUI :
// 	return &beerbot{
// 		id1: setUserPic("https://image.shutterstock.com/image-photo/man-posing-police-mugshot-260nw-637218115.jpg"),
// 		id2: setUserPic("https://www.gocivilairpatrol.com/media/cms/Membership_ID_photo_FA67888970A73.jpg"),
// 		id3: setUserPic("https://static.wikia.nocookie.net/darling-in-the-franxx/images/b/b3/Zero_Two_appearance.jpg/revision/latest/scale-to-width-down/340?cb=20180807204943"),
// 		id4: setUserPic("https://i.pinimg.com/originals/64/86/60/648660b8d170ba0540bc1ed50f33de4e.jpg"),
// 		id5: setUserPic("https://i.pinimg.com/originals/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg"),
// 		id6: setUserPic("https://www.auburn.edu/administration/tigercard/images/sample_id_photo.jpg"),
// 		id7: setUserPic("https://cdn.cloudflare.steamstatic.com/steamcommunity/public/images/avatars/4a/4a5a8309b3ec29a8e3e1cd3f64704ab54427bb4b_full.jpg"),
// 		id8: setUserPic("https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTlb9odaRFcDKI8LF4x3lOiajYk2CXs_PcGRg&usqp=CAU"),
// 	}
// } //end newBGUI

	// func (bGUI *beerbot) DataToScreen() {
	// 	myType := reflect.TypeOf(bGUI).Elem()
	// 	myValue := reflect.ValueOf(bGUI).Elem()
	// 	for i := 0; i < myType.NumField(); i++ {
	// 			tag := myType.Field(i).Tag.Get("json")
	// 			switch tag {
	// 			case "": // not a dipackage GUI
// 			case "img": // special field for images
	// 				url := myValue.Field(i).String()
	//
	// 				go x.downloadImage(url)
	// 			case "num":
	// 				v := myValue.Field(i).Int()
	// 				bGUI.iDEntry.SetText(fmt.Sprintf("%d", v))
	// 			default:
	// 				v := myValue.Field(i).String()
	// 				if newline := strings.IndexAny(v, "\n.-,"); newline > -1 {
	// 					v = v[:newline] + "..."
	// 				}
	// 				bGUI.labels[tag].SetText(v)
	// 			}
	// 		}
	// } //end bGUI
