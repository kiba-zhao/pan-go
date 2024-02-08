package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Desktop interface {
	Destroy()
	Show()
}
type desktopImpl struct {
	app fyne.App
}

func New() Desktop {

	desktop := new(desktopImpl)
	desktop.app = app.New()
	return desktop
}

func (d *desktopImpl) Show() {
	w := d.app.NewWindow("Hello")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	w.Resize(fyne.NewSize(600, 400))
	w.Show()
	d.app.Run()
}

func (d *desktopImpl) Destroy() {

}
