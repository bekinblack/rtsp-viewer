package components

import "fyne.io/fyne/v2/widget"

func NewEntry(placeholder string, onChanged func(s string)) *widget.Entry {
	e := widget.NewEntry()
	e.PlaceHolder = placeholder
	e.OnChanged = onChanged

	return e
}

func NewPwdEntry(placeholder string, onChanged func(s string)) *widget.Entry {
	e := widget.NewPasswordEntry()
	e.PlaceHolder = placeholder
	e.OnChanged = onChanged

	return e
}
