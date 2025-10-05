package ui

import (
	"image"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Viewer struct {
	width  int
	height int
	image  *image.RGBA
	viewer *canvas.Image
}

func NewViewer(width, height int) *Viewer {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	viewer := canvas.NewImageFromImage(img)
	viewer.SetMinSize(fyne.NewSize(float32(width), float32(height)))

	return &Viewer{
		image:  img,
		width:  width,
		height: height,
		viewer: viewer,
	}
}

func (v *Viewer) View(r io.Reader) {
	buf := make([]byte, v.width*v.height*4)
	for {
		if _, err := io.ReadFull(r, buf); err != nil {
			copy(v.image.Pix, make([]byte, v.width*v.height*4))
			fyne.Do(v.viewer.Refresh)
			break
		}

		copy(v.image.Pix, buf)
		fyne.Do(v.viewer.Refresh)
	}
}

func (v *Viewer) Content() fyne.CanvasObject { return v.viewer }
