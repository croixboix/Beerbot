package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)


func makeImageItem() fyne.CanvasObject {
	label := canvas.NewText("label", color.Gray{128})
	label.Alignment = fyne.TextAlignCenter

	img := canvas.NewRectangle(color.Black)
	return container.NewBorder(nil, label, nil, nil, img)
}


func makeImageGrid() fyne.CanvasObject {
	items := []fyne.CanvasObject{}

	for range []int{1,2,3} {
		img := makeImageItem()
		items = append(items, img)
	}

	cellSize := fyne.NewSize(160, 120)
	return container.NewGridWrap(cellSize, items...)
}


func makeStatus()fyne.CanvasObject{
	return canvas.NewText("status", color.Gray{128})
}


func makeUI() fyne.CanvasObject {
	status := makeStatus()
	content := makeImageGrid()
	return container.NewBorder(nil, status, nil, nil, content)
}



func main() {
	a := app.New()
	w := a.NewWindow("Image Browser")
	w.SetContent(makeUI())
	w.Resize(fyne.NewSize(480,360))
	w.ShowAndRun()


}
