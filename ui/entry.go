package ui

import "fyne.io/fyne/v2/widget"

func Entry(placeholder string, OnChanged func(s string)) *widget.Entry {
	e := widget.NewEntry()
	e.PlaceHolder = placeholder
	e.OnChanged = OnChanged

	return e
}

func PwdEntry(placeholder string, OnChanged func(s string)) *widget.Entry {
	e := widget.NewPasswordEntry()
	e.PlaceHolder = placeholder
	e.OnChanged = OnChanged

	return e
}
