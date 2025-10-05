package model

import (
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	invalidIP       = "Некорректный IPv4-адрес"
	invalidPort     = "Порт должен быть числом 1–65535"
	invalidLogin    = "Логин не должен быть пустым"
	invalidPassword = "Пароль не должен быть пустым"
	invalidURI      = "RTSP-URI должен начинаться с rtsp://"
	invalidUriPair  = "Оба RTSP-URI совпадают после нормализации"
)

type Form struct {
	ip       string
	port     string
	login    string
	password string

	uriLow  string
	uriHigh string
}

func (f *Form) UriHigh() string { return f.uriHigh }
func (f *Form) UriLow() string  { return f.uriLow }

func (f *Form) ChangeIP(ip string) {
	f.ip = strings.TrimSpace(ip)
}

func (f *Form) ChangePort(port string) {
	f.port = strings.TrimSpace(port)
}

func (f *Form) ChangeLogin(login string) {
	f.login = strings.ToLower(strings.TrimSpace(login))
}

func (f *Form) ChangePassword(password string) {
	f.password = password
}

func (f *Form) ChangeUriHigh(uri string) {
	f.uriHigh = normaliseUri(uri)
}

func (f *Form) ChangeUriLow(uri string) {
	f.uriLow = normaliseUri(uri)

}

func (f *Form) Validate() error {
	if net.ParseIP(f.ip) == nil {
		return errors.New(invalidIP)
	}

	portInt, err := strconv.Atoi(f.port)
	if err != nil {
		return errors.New(invalidPort)
	}

	if portInt < 1 || portInt > 65535 {
		return errors.New(invalidPort)
	}

	if len(f.login) == 0 {
		return errors.New(invalidLogin)
	}

	if len(f.password) == 0 {
		return errors.New(invalidPassword)
	}

	if !strings.HasPrefix(f.uriHigh, "rtsp://") {
		log.Println(f.uriHigh)
		return errors.New(invalidURI)
	}

	if !strings.HasPrefix(f.uriLow, "rtsp://") {
		log.Println(f.uriLow)
		return errors.New(invalidURI)
	}

	if f.uriHigh == f.uriLow {
		log.Println(f.uriHigh, f.uriLow)
		return errors.New(invalidUriPair)
	}
	return nil
}

func normaliseUri(uri string) string {
	return strings.TrimSuffix(strings.ToLower(strings.TrimSpace(uri)), "/")
}
