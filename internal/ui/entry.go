package ui

import "fyne.io/fyne/v2/widget"

func Entry(text, placeholder string, OnChanged func(s string)) *widget.Entry {
	e := widget.NewEntry()
	e.Text = text
	e.PlaceHolder = placeholder
	e.OnChanged = OnChanged

	return e
}

func PwdEntry(text, placeholder string, OnChanged func(s string)) *widget.Entry {
	e := widget.NewPasswordEntry()
	e.Text = text
	e.PlaceHolder = placeholder
	e.OnChanged = OnChanged

	return e
}
