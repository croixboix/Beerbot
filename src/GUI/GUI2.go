package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

)

// XKCD is an app to get xkcd images and display them
type XKCD struct {
	ID         int    `json:"num"`
	Title      string `json:"title"`
	Day        string `json:"day"`
	Month      string `json:"month"`
	Year       string `json:"year"`
	Link       string `json:"link"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	News       string `json:"news"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`

	image   *canvas.Image
	iDEntry *widget.Entry
	labels  map[string]*widget.Label
}

func main() {
	a := app.New()
	//a.SetIcon(resourceIconPng)

	w := a.NewWindow("BeerBot")

	w.SetContent(Show(w))
	w.Resize(fyne.NewSize(480, 360))
	w.ShowAndRun()
}
