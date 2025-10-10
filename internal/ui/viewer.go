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
	buf    []byte
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
		buf:    make([]byte, width*height*4),
	}
}

func (v *Viewer) View(r io.Reader) {
	for {
		if _, err := io.ReadFull(r, v.buf); err != nil {
			v.clear()
			break
		}

		v.refresh()
	}
}

func (v *Viewer) refresh() {
	copy(v.image.Pix, v.buf)
	fyne.Do(v.viewer.Refresh)
}

func (v *Viewer) clear() {
	v.buf = make([]byte, v.width*v.height*4)
	v.refresh()
}

func (v *Viewer) Content() fyne.CanvasObject {
	return v.viewer
}
